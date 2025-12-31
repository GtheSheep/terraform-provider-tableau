data "tableau_projects" "example" {}
data "tableau_groups" "example" {}

resource "tableau_project_permissions" "example" {
  project_id  = data.tableau_projects.example.projects[0].id
  grantee_capabilities = [
    {
      group_id = data.tableau_groups.example.groups[0].id
      capabilities = [
        {
          name = "Read"
          mode = "Allow"
        },
        {
          name = "Write"
          mode = "Deny"
        },
      ]
    },
    {
      group_id = data.tableau_groups.example.groups[1].id
      capabilities = [
        {
          name = "Read"
          mode = "Allow"
        },
        {
          name = "Write"
          mode = "Allow"
        },
      ]
    },
  ]
}