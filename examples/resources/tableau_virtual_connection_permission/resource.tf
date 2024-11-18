resource "tableau_virtual_connection_permission" "test_permission" {
  virtual_connection_id = "xxxxx-xxxxx-xxxxx"
	user_id = "xxxxx-xxxxx-xxxxx"
  capability_name = "Write"
	capability_mode = "Deny"
}
