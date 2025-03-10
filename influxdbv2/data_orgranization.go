package influxdbv2

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/domain"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceOrganization() *schema.Resource {
	return &schema.Resource{
		Description: "Lookup an Organization in InfluxDB2.",
		ReadContext: DataSourceOrganizationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				Description: "Name of the Organization.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the Organization.",
			},
			// Computed outputs
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the Organization.",
			},
		},
	}
}

func DataSourceOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// use the meta value to retrieve your client from the provider configure method
	influx := meta.(influxdb2.Client)
	orgAPI := influx.OrganizationsAPI()

	// Warning or errors can be collected in a slice type
	var (
		diags diag.Diagnostics
		org   *domain.Organization
		err   error
	)

	if v, ok := d.GetOk("name"); ok {
		orgName := v.(string)
		if org, err = orgAPI.FindOrganizationByName(ctx, orgName); err != nil {
			diags = append(diags, diag.FromErr(err)...)
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Can't find Organization with name: %s", orgName),
			})
			return diags
		}
	}

	id := org.Id
	if id == nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Organization not found",
		})
		return diags
	}

	d.SetId(*id)
	err = d.Set("id", *id)
	if err != nil {
		return nil
	}
	err = d.Set("name", org.Name)
	if err != nil {
		return nil
	}
	if org.Description != nil {
		err := d.Set("description", *org.Description)
		if err != nil {
			return nil
		}
	}

	return diags
}
