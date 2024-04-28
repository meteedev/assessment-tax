package handler

import (
	"io"

	"github.com/xeipuuv/gojsonschema"
	"github.com/labstack/echo/v4"
	"github.com/meteedev/assessment-tax/apperrs"
	"github.com/meteedev/assessment-tax/constant"
)





func (h *TaxHandler) validateSchema(c echo.Context, schema string) ([]byte,error) {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return nil,apperrs.NewInternalServerError(err.Error())
	}

	// Load JSON schema
	schemaLoader := gojsonschema.NewStringLoader(schema)

	// Load JSON request
	requestLoader := gojsonschema.NewStringLoader(string(body))


	// Validate JSON request against JSON schema
	result, err := gojsonschema.Validate(schemaLoader, requestLoader)
	if err != nil {
		return nil,apperrs.NewBadRequestError(err.Error())
	}

	// Check validation result
	if !result.Valid() {
		return nil,apperrs.NewBadRequestError(constant.MSG_HANDLER_ERR_INVALID_PAYLOAD)
	}

	return body,nil
}
