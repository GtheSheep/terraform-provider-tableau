---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tableau_virtual_connection Data Source - terraform-provider-tableau"
subcategory: ""
description: |-
  Retrieve virtual Connection details
---

# tableau_virtual_connection (Data Source)

Retrieve virtual Connection details

## Example Usage

```terraform
data "tableau_virtual_connections" "vcs" {}

data "tableau_virtual_connection" "vc0" {
  id = data.tableau_virtual_connections.vcs.virtual_connections[0].id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) ID of the virtual Connection

### Read-Only

- `content` (String) Definition of the virtual connection as JSON
- `name` (String) Name of the virtual connection
- `owner_id` (String) Owner ID of the virtual connection
- `project_id` (String) Project ID of the virtual connection
