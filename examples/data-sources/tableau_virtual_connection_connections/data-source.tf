data "tableau_virtual_connection_connections" "example" {
    id = data.tableau_virtual_connections.vc[0].id
}