resource "tableau_group" "test_users" {
  grant_license_mode = "onLogin"
  minimum_site_role  = "Explorer"
  name               = "Test Users"
}
