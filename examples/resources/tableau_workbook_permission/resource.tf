resource "tableau_workbook_permission" "test_permission" {
  workbook_id = "xxxxx-xxxxx-xxxxx"
	user_id = "xxxxx-xxxxx-xxxxx"
  capability_name = "Write"
	capability_mode = "Deny"
}
