package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewValidator() *CustomValidator {
	v := validator.New()
	// Mengambil nama field dari tag json, bukan nama struct field
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &CustomValidator{Validator: v}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

// FormatError mengubah error validator menjadi string list
func FormatError(err error) []string {
	var errors []string
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, fmt.Sprintf("Field '%s' failed on tag '%s'", e.Field(), e.Tag()))
		}
	}
	return errors
}
