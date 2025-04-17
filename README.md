# terraform-provider-googlesheets

[![CI](https://github.com/winebarrel/terraform-provider-googlesheets/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/terraform-provider-googlesheets/actions/workflows/ci.yml)
[![terraform docs](https://img.shields.io/badge/terraform-docs-%35835CC?logo=terraform)](https://registry.terraform.io/providers/winebarrel/googlesheets/latest/docs)

Terraform provider for retrieving data from Google Sheets.

## Usage

```tf
terraform {
  required_providers {
    lambdaconfig = {
      source  = "winebarrel/googlesheets"
      version = ">= 0.2.0"
    }
  }
}

provider "googlesheets" {
  credentials_json = file("credentials.json")
}

data "googlesheets_sheet" "my_sheet" {
  sheet_id = "..."
  range    = "sheet1!A2:B2"
}

output "values" {
  value = jsondecode(data.googlesheets_sheet.my_sheet.json)
}
# values = [
#   [
#     "A1 TEXT",
#     "B1 TEXT",
#   ],
#   [
#     "A2 TEXT",
#     "B2 TEXT",
#   ],
# ]

output "sensitive_values" {
  value     = jsondecode(data.googlesheets_sheet.my_sheet.sensitive_json)
  sensitive = true
}
```

## Run locally for development

```sh
# TODO: Create "credentials.json".
#       see https://cloud.google.com/iam/docs/keys-create-delete
cp googlesheets.tf.sample googlesheets.tf
make
make tf-plan
make tf-apply
```
