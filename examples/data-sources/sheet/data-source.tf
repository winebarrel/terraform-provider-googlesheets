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
