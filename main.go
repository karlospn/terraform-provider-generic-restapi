package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/karlospn/terraform-provider-generic-restapi/generic"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return generic.Provider()
		},
	})
}
