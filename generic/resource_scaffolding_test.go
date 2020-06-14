package generic

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccScaffoldingObject_Basic(t *testing.T) {

	os.Setenv("REST_API_URI", "http://127.0.0.1:5000")
	os.Setenv("REST_API_READ_METHOD", "/api/builds/projectB/{id}")
	os.Setenv("REST_API_CREATE_METHOD", "/api/builds/projectB")
	os.Setenv("REST_API_UPDATE_METHOD", "/api/builds/{id}")
	os.Setenv("REST_API_DESTROY_METHOD", "/api/builds/{id}")
	os.Setenv("TF_LOG", "1")

	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaffolding_resource.test", "id", "99"),
				),
			},
			{
				Config: testResourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaffolding_resource.test", "id", "99"),
				),
			},
		},
	})
}

const testResource = `
resource "scaffolding_resource" "test" {
  create_method  = "/api/builds"
  read_method = "/api/builds/projectA/{id}"
  payload = jsonencode(
    {
		"project": "projectA",
        "applicationName" : "some_more_apps_21",
        "buildTemplate": "template_11",
        "pool": "MyPool",
        "repository" : "App1",
        "branch" : "master"
    })
}`

const testResourceUpdate = `
resource "scaffolding_resource" "test" {
	payload = jsonencode(
	  {
		  "project": "projectA",
		  "applicationName" : "some_more_apps_10",
		  "buildTemplate": "template_111",
		  "pool": "MyPool",
		  "repository" : "App1",
		  "branch" : "master"
	  })
  }`
