resource "tableau_project_permission" "test_permission" {
  project_id = "xxxxx-xxxxx-xxxxx"
	user_id = "xxxxx-xxxxx-xxxxx"
  capability_name = "Write"
	capability_mode = "Deny"
}
