package form

import (
	// "fmt"
	"reflect"
	// "strings"
	"github.com/go-macaron/binding"
	
	"github.com/midoks/dztasks/internal/tools"
)

type Empty struct {
	No bool
}

func init() {
	binding.SetNameMapper(tools.ToSnakeCase)
}

// Assign assign form values back to the template data.
func Assign(form interface{}, data map[string]interface{}) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		} else if len(fieldName) == 0 {
			fieldName = tools.ToSnakeCase(field.Name)
		}

		data[fieldName] = val.Field(i).Interface()
	}
}
