package generic

import (
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProvider terraform.ResourceProvider
var testAccProviders map[string]terraform.ResourceProvider

func init() {
	testAccProvider = Provider().(terraform.ResourceProvider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"sccafolding_resource": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func TestResourceProvider_RequireBasic(t *testing.T) {
	rp := Provider()

	raw := map[string]interface{}{}

	rawConfig, err := config.NewRawConfig(raw)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	err = rp.Configure(terraform.NewResourceConfig(rawConfig))
	if err != nil {
		t.Fatalf("Provider failed with error: %s", err)
	}
}
