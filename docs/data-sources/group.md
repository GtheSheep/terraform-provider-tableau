---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tableau_group Data Source - terraform-provider-tableau"
subcategory: ""
description: |-
  Retrieve group details
---

# tableau_group (Data Source)

Retrieve group details

## Example Usage

```terraform
data "tableau_group" "example" {
    id = "abc"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) ID of the group

### Read-Only

- `minimum_site_role` (String) Minimum site role for the group
- `name` (String) Name for the group
