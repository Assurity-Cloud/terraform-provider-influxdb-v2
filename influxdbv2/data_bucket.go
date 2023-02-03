package influxdbv2

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceBucket() *schema.Resource {
	return &schema.Resource{
		Description: "Lookup a Bucket in InfluxDB2.",
		ReadContext: DataSourceBucketRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				Description: "Name of the Bucket.",
			},
			// Computed outputs
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the Bucket.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the Bucket.",
			},
		},
	}
}

func DataSourceBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	influx := meta.(influxdb2.Client)
	bucketsAPI := influx.BucketsAPI()

	var (
		diags  diag.Diagnostics
		bucket *domain.Bucket
		err    error
	)

	if v, ok := d.GetOk("name"); ok {
		bucketName := v.(string)
		if bucket, err = bucketsAPI.FindBucketByName(ctx, bucketName); err != nil {
			diags = append(diags, diag.FromErr(err)...)
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Can't find Bucket with name: %s", bucketName),
			})
			return diags
		}
	}

	id := bucket.Id
	if id == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Bucket not found",
		})
		return diags
	}

	d.SetId(*id)
	err = d.Set("id", *id)
	if err != nil {
		return nil
	}
	err = d.Set("name", bucket.Name)
	if err != nil {
		return nil
	}
	if bucket.Description != nil {
		err := d.Set("description", *bucket.Description)
		if err != nil {
			return nil
		}
	}

	return diags
}
