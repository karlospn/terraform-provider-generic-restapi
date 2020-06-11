package generic

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccScaffoldingObject_Basic(t *testing.T) {

	os.Setenv("REST_API_URI", "http://127.0.0.1:5000")
	os.Setenv("REST_API_READ_METHOD", "/api/builds/projectA/{id}")
	os.Setenv("REST_API_CREATE_METHOD", "/api/builds/projectA")
	os.Setenv("REST_API_UPDATE_METHOD", "/api/builds/projectA/{id}")
	os.Setenv("REST_API_DESTROY_METHOD", "/api/builds/projectA/{id}")
	os.Setenv("TF_LOG", "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaffolding_resource.test", "id", "some_more_apps"),
				),
			},
			{
				Config: testResourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaffolding_resource.test", "id", "some_more_apps"),
				),
			},
		},
	})
}

const testResource = `
resource "scaffolding_resource" "test" {
  id_attribute = "applicationName"
  data = jsonencode(
    {
		"project": "projectA",
        "applicationName" : "some_more_apps",
        "buildTemplate": "template_11",
        "pool": "MyPool",
        "repository" : "App1",
        "branch" : "master"
    })
}`

const testResourceUpdate = `
resource "scaffolding_resource" "test" {
  id_attribute = "applicationName"
  data = jsonencode(
    {
		"project": "projectA",
        "applicationName" : "some_more_apps",
        "buildTemplate": "template_12",
        "pool": "MyPool",
        "repository" : "App1",
        "branch" : "master"
    })
}`
