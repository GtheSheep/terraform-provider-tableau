resource "tableau_datasource_permission" "test_permission" {
  datasource_id = "xxxxx-xxxxx-xxxxx"
	user_id = "xxxxx-xxxxx-xxxxx"
  capability_name = "Write"
	capability_mode = "Deny"
}
