resource "tableau_project_default_permission" "test_permission" {
  project_id = "xxxxx-xxxxx-xxxxx"
  user_id = "xxxxx-xxxxx-xxxxx"
  target_type = "workbooks"
  capability_name = "Write"
  capability_mode = "Deny"
}
