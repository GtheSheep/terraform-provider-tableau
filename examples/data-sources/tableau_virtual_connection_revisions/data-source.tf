data "tableau_virtual_connections" "vc" {}
data "tableau_virtual_connection_revisions" "example" {
    id = data.tableau_virtual_connections.vc.revisions[0].id
}