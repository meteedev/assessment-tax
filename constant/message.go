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
	MSG_BU_GERNERAL_ERROR = "Sorry for the inconvenience Unavailable at this time"
	MSG_BU_INVALID_WHT_LESS_THAN_ZERO = "WHT must not be less than 0 "
	MSG_BU_INVALID_WHT_GREATER_THAN_TOTALINCOME = "WHT can not greater than Total income"
)