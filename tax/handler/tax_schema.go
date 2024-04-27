package handler


const TAX_REQUEST_SCHEMA = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Tax Request Schema",
  "type": "object",
  "properties": {
    "totalIncome": {
      "type": "number",
      "minimum": 0
    },
    "wht": {
      "type": "number",
      "minimum": 0
    },
    "allowances": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "allowanceType": {
            "type": "string"
          },
          "amount": {
            "type": "number",
            "minimum": 0
          }
        },
        "required": ["allowanceType", "amount"]
      }
    }
  },
  "required": ["totalIncome", "wht", "allowances"]
}
`
const UPDATE_DEDUCT_REQUEST = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Update Deduct Request Schema",
  "type": "object",
  "properties": {
    "amount": {
      "type": "number"
    }
  },
  "required": ["amount"]
}
`