resource "tableau_group_user" "example" {
  group_id = tableau_group.test_users.id
  user_id  = tableau_user.test_user.id
}
