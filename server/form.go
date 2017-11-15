package server

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// parseForm parses the supplied form, validates it, and stores the values in
// the supplied struct. The return value contains a map of all errors
// encountered during parsing.
func parseForm(form url.Values, v interface{}) map[string][]string {
	var (
		vType  = reflect.TypeOf(v).Elem()
		vVal   = reflect.ValueOf(v).Elem()
		errors = map[string][]string{}
	)
	for i := 0; i < vType.NumField(); i++ {
		var (
			fType = vType.Field(i)
			fVal  = vVal.Field(i)
			s     = form.Get(vType.Field(i).Name)
		)
		for _, tag := range strings.Split(fType.Tag.Get("form"), ",") {
			switch tag {
			case "required":
				if len(s) == 0 {
					errors[fType.Name] = append(
						errors[fType.Name],
						"this field is required",
					)
				}
			}
		}
		switch fVal.Kind() {
		case reflect.String:
			fVal.SetString(s)
		case reflect.Int64:
			iVal, _ := strconv.ParseInt(s, 10, 64)
			fVal.SetInt(iVal)
		case reflect.Bool:
			fVal.SetBool(len(s) > 0)
		}
	}
	return errors
}
