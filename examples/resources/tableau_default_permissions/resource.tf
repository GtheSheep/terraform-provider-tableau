data "tableau_projects" "example" {}
data "tableau_groups" "example" {}

resource "tableau_default_permissions" "example" {
  project_id  = data.tableau_projects.example.projects[0].id
  target_type = "datasources"
  grantee_capabilities = [
    {
      group_id = data.tableau_groups.example.groups[0].id
      capabilities = [
        {
          name = "PulseMetricDefine"
          mode = "Allow"
        },
      ]
    },
    {
      group_id = data.tableau_groups.example.groups[1].id
      capabilities = [
        {
          name = "SaveAs"
          mode = "Allow"
        },
        {
          name = "ExportXml"
          mode = "Deny"
        },
        {
          name = "VizqlDataApiAccess"
          mode = "Allow"
        },
        {
          name = "Delete"
          mode = "Allow"
        },
        {
          name = "ChangePermissions"
          mode = "Allow"
        },
        {
          name = "Write"
          mode = "Allow"
        },
        {
          name = "Connect"
          mode = "Allow"
        },
        {
          name = "Read" # == View
          mode = "Allow"
        },
        {
          name = "ChangeHierarchy" # == Move
          mode = "Allow"
        },
      ]
    },
  ]
}
