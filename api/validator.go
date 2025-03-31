package api

import (
	"github.com/definitely-unique-username/simple_bank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
	if currncy, ok := fl.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currncy)
	}

	return false
}
