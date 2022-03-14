package main

import (
	"github.com/appoptics/terraform-provider-appoptics/appoptics"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: appoptics.Provider,
	})
}
