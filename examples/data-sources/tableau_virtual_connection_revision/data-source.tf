data "tableau_virtual_connections" "vc" {}
data "tableau_virtual_connection_revisions" "example" {
    id       = data.tableau_virtual_connections.vc.revisions[0].id
    revision = 1
}

/*
 * Revisions numbering starts from 1 and if you give latest revision number, it will complain
 * com.tableausoftware.domain.content.publishedconnections.exceptions.WrongPublishedConnectionRevisionType: Revision c9b09d76-3824-4ef2-befe-ba3f775a8da3 was expected to be HISTORICAL but was ACTIVE (errorCode=360010))","code":"400200"}}
 */