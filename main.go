package main

import (
	"context"

<<<<<<< Updated upstream
	"github.com/gthesheep/terraform-provider-tableau/tableau"
=======
>>>>>>> Stashed changes
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"gitlab.com/tailormed/devops/terraform-provider-tableau/tableau"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name tableau

func main() {
	providerserver.Serve(context.Background(), tableau.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/gthesheep/tableau",
	})
}
