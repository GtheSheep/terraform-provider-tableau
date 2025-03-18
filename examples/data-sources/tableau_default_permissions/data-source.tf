data "tableau_projects" "all" {
}

data "tableau_default_permissions" "default_permissions" {
    project_id  = data.tableau_projects.all.projects[0].id
    target_type = "workbooks"
}