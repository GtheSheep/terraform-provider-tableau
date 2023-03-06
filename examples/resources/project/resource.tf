resource "tableau_project" "test" {
  name = "test"
  description = "Moo"
  content_permissions = "LockedToProject"
  parent_project_id = tableau_project.test_parent.id
}
