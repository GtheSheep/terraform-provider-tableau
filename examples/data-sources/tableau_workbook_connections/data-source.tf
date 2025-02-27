data "tableau_workbook_connections" "example" {
    id = data.tableau_workbooks.wb[0].id
}
