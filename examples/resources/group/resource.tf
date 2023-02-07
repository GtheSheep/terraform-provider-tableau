resource "tableau_group" "example" {
  grant_license_mode = "onLogin"
  minimum_site_role  = "Explorer"
  name               = "Test Users"
}
