package influxdbv2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccReadBucket tests the read Bucket data source
func TestAccReadBucket(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceBucketConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.influxdb-v2_bucket.by_name", "name", "BucketAcctestName"),
					resource.TestCheckResourceAttr("data.influxdb-v2_bucket.by_name", "description", "Desc Bucket Acctest"),
				),
			},
		},
	})
}

func testDataSourceBucketConfig() string {
	return `resource "influxdb-v2_organization" "org" {
                name = "OrgAcctestName"
                description = "Desc Org Acctest"
        }

	    resource "influxdb-v2_bucket" "bucket" {
			name = "BucketAcctestName"
			description = "Desc Bucket Acctest"
			org_id = influxdb-v2_organization.org.id
            retention_rules {
                every_seconds = "3630"
            }
		}

		data "influxdb-v2_bucket" "by_name" {
			name = influxdb-v2_bucket.bucket.name
			//requirement of terraform v0.13
			depends_on = [influxdb-v2_bucket.bucket]
		}
`
}
