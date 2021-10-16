# terraform-provider-tableau
Terraform Provider for Tableau

Originally aimed at Tableau Online, but can expand.

```terraform
terraform {
  required_providers {
    tableau = {
      source  = "GtheSheep/tableau"
      version = "<version>"
    }
  }
}
```

## Authentication

## Examples
Check out the `examples/` folder for some usage options, these are intended to
simply showcase what this module can do rather than be best practices for any
given use case.
