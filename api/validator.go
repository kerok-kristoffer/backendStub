package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/kerok-kristoffer/backendStub/util"
)

// sample validator to demonstrate custom validators for api calls seen in #14 in tut
var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency) // never implemented util.IsSupportedCurrency
	}
	return false
}
