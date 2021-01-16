package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/keremakillioglu/simplebank/util"
)

// Func accepts a FieldLevel interface for all validation needs. The return
// value should be true when validation succeeds.
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	// extract the currency from field level
	// fieldLevel.Field() is a reflection value/ call it with .Interface() to get the value and convert it to string
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		//check if currency is supported
		return util.IsSupportedCurrency(currency)
	}

	return false
}
