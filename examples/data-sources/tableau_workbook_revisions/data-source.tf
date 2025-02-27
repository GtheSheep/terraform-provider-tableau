data "tableau_workbook_revisions" "example" {
    id = data.tableau_workbooks.wb[0].id
}
