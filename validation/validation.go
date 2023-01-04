package validation

import (
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Init() {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		// This is a workaround for FieldError.Field() bug
		// in validation v10, that returns the actual struct field name
		// instead of the json name, which is needed for custom error messages.
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		//  Register custom validations
		v.RegisterValidation("unik", unik)
	}
}
