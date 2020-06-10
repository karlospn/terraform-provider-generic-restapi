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

	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaffolding_resource.test", "id", "1234"),
				),
			},
			{
				Config: testResourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaffolding_resource.test", "id", "1234"),
				),
			},
		},
	})
}

const testResource = `
resource "scaffolding_resource" "test" {
  id_attribute = "id"
  data = jsonencode(
    {
        "id": "1234",
		"name" : "john",
		"age" : 23
    })
}`

const testResourceUpdate = `
resource "scaffolding_resource" "test" {
  id_attribute = "id"
  data = jsonencode(
    {
        "id": "1234",
		"name" : "john",
		"age" : 24
    })
}`
