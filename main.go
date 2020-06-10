package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/karlospn/terraform-provider-generic-restapi/generic"
)

func main() {

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: generic.Provider})

}
