resource "tableau_group_user" "data_science_gary_james" {
  group_id = tableau_group.test_users.id
  user_id  = tableau_user.test_user.id
}
