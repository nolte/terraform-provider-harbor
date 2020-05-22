package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/nolte/terraform-provider-harbor/harbor"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: harbor.Provider})
}
