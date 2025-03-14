resource "tableau_default_permissions" "org-admin_group-datasources" {
  project_id  = resource.tableau_project.org.id
  target_type = "datasources"
  grantee_capabilities = [
    {
      group_id = var.all_users_group_id
      capabilities = [
        {
          name = "PulseMetricDefine"
          mode = "Allow"
        },
      ]
    },
    {
      group_id = var.admin_group_id
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
