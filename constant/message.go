package constant

//server message
const (
	MSG_SERVER_GRACEFUL_SHUTDOWN = "shutting down the server"
)

// app errors message 
const (
	MSG_APP_ERR_INTERNAL_SERVER_ERROR = "Internal Server Error"
	MSG_APP_ERR_UNEXPECTED_ERROR      = "Unexpected error occurred"
)

// business logic message
const (
	MSG_BU_GENERAL_ERROR = "sorry for the inconvenience Unavailable at this time"
	
	MSG_BU_VALIDATE_CSV_DIGIT_ONLY = "column data in csv must only digits"
	MSG_BU_VALIDATE_CSV_GREATER_EQUAL_ZERO = "column data in csv greater than or equal 0"
	MSG_BU_INVALID_CSV_RECORD_COLUMN_NUMBERS = "column data in csv must 3 columns"

	MSG_BU_INVALID_TOTAL_INCOME_LESS_THAN_OR_EQUAL_ZERO = "totalIncome must greater than 0 "
	
	MSG_BU_INVALID_WHT_LESS_THAN_ZERO = "wht must not be less than 0 "
	MSG_BU_INVALID_WHT_GREATER_THAN_TOTALINCOME = "wht can not greater than Total income"
	
	MSG_BU_INVALID_PERSONAL_ALLOW_LESS_THAN_ZERO = "personal allowance must not be less than 0 "
	MSG_BU_INVALID_PERSONAL_ALLOW_MININUM = "personal deductibles start at "

	MSG_BU_INVALID_K_RECEIPT_ALLOW_LESS_THAN_ZERO = "k-receipt allowance must not be less than 0 "
	MSG_BU_INVALID_K_RECEIPT_ALLOW_MININUM = "k-receipt deductibles start at "

	MSG_BU_DEDUCT_UPD_PERSONAL_FAILED = "update personal allowance failed" 
	MSG_BU_DEDUCT_PERSONAL_CONFIG_NOT_FOUND = "personal allowance config not found in database"

	MSG_BU_DEDUCT_UPD_K_RECEIPT_FAILED = "update k-receipt allowance failed" 
	MSG_BU_DEDUCT_K_RECEIPT_CONFIG_NOT_FOUND = "k-receipt allowance config not found in database"

)

const(
	MSG_HANDLER_ERR_LOADING_SCHEMA = "error loading schema"
	MSG_HANDLER_ERR_VALIDATE_SCHEMA = "error validating schema"
	MSG_HANDLER_ERR_INVALID_PAYLOAD = "invalid payload"
)


const(
	MSG_UPLOAD_CSV_WRONG_FORMAT  = "csv wrong format"
)