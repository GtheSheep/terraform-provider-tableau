data "tableau_projects" "all" {
}

data "tableau_project_permissions" "project_permissions" {
    id = data.tableau_projects.all.projects[0].id
}
