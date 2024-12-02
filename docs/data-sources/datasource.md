---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "tableau_datasource Data Source - terraform-provider-tableau"
subcategory: ""
description: |-
  Retrieve datasource details
---

# tableau_datasource (Data Source)

Retrieve datasource details

## Example Usage

```terraform
data "tableau_datasource" "example" {
    name = "moo"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) ID of the datasource
- `name` (String) Name for the datasource

### Read-Only

- `certification_note` (String) Certification note
- `content_url` (String) URL of the datasource content
- `description` (String) Datasource description
- `encrypt_extracts` (String) Whether or not this datasource encrypts extracts
- `has_extracts` (Boolean) Whether or not this datasource has extracts
- `is_certified` (Boolean) Whether or not this datasource is certified
- `owner_id` (String) Datasource Owner ID
- `project_id` (String) Datasource Project ID
- `tags` (List of String) List of tags on the datasource
- `type` (String) Type of datasource
- `use_remote_query_agent` (Boolean) Whether or not this datasource uses a remote query agent
- `web_page_url` (String) Web page URL for the datasource