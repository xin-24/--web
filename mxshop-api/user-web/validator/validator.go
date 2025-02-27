package validator//验证器

import (
	"regexp"
	"github.com/go-playground/validator/v10"
)

func ValidateMobile(fl validator.FieldLevel)bool{
	mobile:=fl.Field().String()
	ok,_:=regexp.MatchString(`^1([38][0-9]|4[579]|5[^4]|6[6]|7[1-35-8]|9[189])\d{8}$`,mobile)
	if !ok{
		return false
	}
	return true

}