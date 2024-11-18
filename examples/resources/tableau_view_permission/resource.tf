resource "tableau_view_permission" "test_permission" {
  view_id = "xxxxx-xxxxx-xxxxx"
	user_id = "xxxxx-xxxxx-xxxxx"
  capability_name = "Write"
	capability_mode = "Deny"
}
