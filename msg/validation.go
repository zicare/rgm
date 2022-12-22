package msg

import "github.com/go-playground/validator/v10"

// Return a validation errors list
func ValidationErrors(err error) (ml MessageList) {

	switch err.(type) {
	/*
		case *time.ParseError:
			//Time %s has a wrong format, required format is %s
			e := err.(*time.ParseError)
			m := msg.Get("22").SetArgs(lib.TrimQuotes(e.Value), "2006-01-02T15:04:05-07:00")
			eml = append(eml, m)
		case *json.UnmarshalTypeError:
			//Value is a %s, required type is %s
			e := err.(*json.UnmarshalTypeError)
			m := msg.Get("23").SetArgs(e.Value, e.Type.String()).SetField(e.Field)
			eml = append(eml, m)
	*/
	case validator.ValidationErrors:
		for _, v := range err.(validator.ValidationErrors) {
			//typ := v.Type().String()
			m := Get("24").SetArgs(v.Value(), v.Tag(), v.Param()).SetField(v.Field())
			ml = append(ml, m)
		}
	}
	return ml
}
