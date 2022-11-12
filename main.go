package main

import (
	"github.com/glovo/terraform-provider-onelogin/onelogin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: onelogin.Provider,
	})
}
