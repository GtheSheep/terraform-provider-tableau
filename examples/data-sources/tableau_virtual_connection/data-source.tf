data "tableau_virtual_connections" "vcs" {}

data "tableau_virtual_connection" "vc0" {
  id = data.tableau_virtual_connections.vcs.virtual_connections[0].id
}