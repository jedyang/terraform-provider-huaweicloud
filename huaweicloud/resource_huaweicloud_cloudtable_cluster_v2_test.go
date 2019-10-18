// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file at
//     https://www.github.com/huaweicloud/magic-modules
//
// ----------------------------------------------------------------------------

package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
)

func TestAccCloudtableClusterV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudtableClusterV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCloudtableClusterV2_basic(acctest.RandString(10)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudtableClusterV2Exists(),
				),
			},
		},
	})
}

func testAccCloudtableClusterV2_basic(val string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup" {
  name = "terraform_test_security_group%s"
  description = "terraform security group acceptance test"
  timeouts {
    delete = "20m"
  }
}

resource "huaweicloud_cloudtable_cluster_v2" "cluster" {
  availability_zone = "%s"
  name = "terraform-test-cluster%s"
  rs_num = 2
  security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup.id}"
  subnet_id = "%s"
  vpc_id = "%s"
  storage_type = "COMMON"
}
	`, val, OS_AVAILABILITY_ZONE, val, OS_NETWORK_ID, OS_VPC_ID)
}

func testAccCheckCloudtableClusterV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.sdkClient(OS_REGION_NAME, "cloudtable", serviceProjectLevel)
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cloudtable_cluster_v2" {
			continue
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return err
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"X-Language":   "en-us",
			}})
		if err == nil {
			return fmt.Errorf("huaweicloud_cloudtable_cluster_v2 still exists at %s", url)
		}
	}

	return nil
}

func testAccCheckCloudtableClusterV2Exists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.sdkClient(OS_REGION_NAME, "cloudtable", serviceProjectLevel)
		if err != nil {
			return fmt.Errorf("Error creating sdk client, err=%s", err)
		}

		rs, ok := s.RootModule().Resources["huaweicloud_cloudtable_cluster_v2.cluster"]
		if !ok {
			return fmt.Errorf("Error checking huaweicloud_cloudtable_cluster_v2.cluster exist, err=not found this resource")
		}

		url, err := replaceVarsForTest(rs, "clusters/{id}")
		if err != nil {
			return fmt.Errorf("Error checking huaweicloud_cloudtable_cluster_v2.cluster exist, err=building url failed: %s", err)
		}
		url = client.ServiceURL(url)

		_, err = client.Get(url, nil, &golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
				"X-Language":   "en-us",
			}})
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return fmt.Errorf("huaweicloud_cloudtable_cluster_v2.cluster is not exist")
			}
			return fmt.Errorf("Error checking huaweicloud_cloudtable_cluster_v2.cluster exist, err=send request failed: %s", err)
		}
		return nil
	}
}
