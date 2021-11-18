package resources_test

import (
	"os"
	"testing"

	"github.com/gthesheep/terraform-provider-tableau/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func providers() map[string]*schema.Provider {
	p := provider.Provider()
	return map[string]*schema.Provider{
		"tableau": p,
	}
}

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = provider.Provider()
	testAccProviders = map[string]*schema.Provider{
		"tableau": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TABLEAU_SERVER_URL"); v == "" {
		t.Fatal("TABLEAU_SERVER_URL must be set for acceptance tests")
	}
	if v := os.Getenv("TABLEAU_SERVER_VERSION"); v == "" {
		t.Fatal("TABLEAU_SERVER_VERSION must be set for acceptance tests")
	}
	if v := os.Getenv("TABLEAU_SITE_NAME"); v == "" {
		t.Fatal("TABLEAU_SITE_NAME must be set for acceptance tests")
	}
}
