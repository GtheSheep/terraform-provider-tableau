resource "tableau_user" "example" {
  auth_setting = "SAML"
  email        = "test.user@email.com"
  full_name    = "Test User"
  name         = "test.user"
  site_role    = "SiteAdministratorCreator"
}
