terraform {
  required_providers {
    dbt = {
      source  = "GtheSheep/tableau"
      version = "0.0.12"
    }
  }
}

provider "tableau" {
  server_url = "https://my.tableau.server.com"
  server_version = "3.13"
  site           = "my_site"
}
