package main

import (
	"context"

	"github.com/gthesheep/terraform-provider-tableau/tableau"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name tableau

func main() {
	providerserver.Serve(context.Background(), tableau.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/gthesheep/tableau",
	})
}
