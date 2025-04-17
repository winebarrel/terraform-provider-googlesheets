data "googlesheets_sensitive_sheet" "my_sheet" {
  sheet_id = "..."
  range    = "sheet1!A2:B2"
}

output "values" {
  value     = jsondecode(data.googlesheets_sensitive_sheet.my_sheet.json)
  sensitive = true
}
