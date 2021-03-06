package server

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// copyStruct enumerates the exported fields in src and copies their values to
// the equivalent fields in dst. It is assumed that fields with the same name
// share the same type, otherwise a panic may result.
func copyStruct(src interface{}, dst interface{}) {
	var (
		srcType = reflect.TypeOf(src).Elem()
		srcVal  = reflect.ValueOf(src).Elem()
		dstVal  = reflect.ValueOf(dst).Elem()
	)
	for i := 0; i < srcType.NumField(); i++ {
		fDstVal := dstVal.FieldByName(srcType.Field(i).Name)
		if fDstVal.Kind() != reflect.Invalid {
			fDstVal.Set(srcVal.Field(i))
		}
	}
}

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
