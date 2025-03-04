data "tableau_users" "users" {}
data "tableau_projects" "projects" {}

resource tableau_workbook "example" {
  name               = "example"
  description        = "example about workbook resource"
  encrypt_extracts   = "true"
  thumbnails_user_id = data.tableau_users.users[0].id
  project_id         = data.tableau_projects.projects[0].id
  show_tabs          = "false"
  workbook_filename  = "wbfile.twb"
  workbook_content   = file("wbfile.twb")
}

/*
Example of minimalistic workbook content of wbfile.twb:
<?xml version='1.0' encoding='utf-8' ?>

<!-- build 20251.25.0219.1921                               -->
<workbook original-version='18.1' source-build='2024.3.0 (20243.25.0110.1701)' version='18.1' xml:base='https://SITE_NAME.tableau.com' xmlns:user='http://www.tableausoftware.com/xml/user'>
  <document-format-change-manifest />
  <preferences />
  <datasources />
  <worksheets>
    <worksheet name='Sheet 1' />
  </worksheets>
  <windows />
</workbook>
*/