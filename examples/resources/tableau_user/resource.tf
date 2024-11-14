resource "tableau_user" "example" {
  auth_setting = "SAML"
  email        = "test.user@email.com"
  full_name    = "test.user@email.com"
  name         = "test.user@email.com"
  site_role    = "SiteAdministratorCreator"
}
