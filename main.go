package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-librato/appoptics"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: appoptics.Provider})
}
