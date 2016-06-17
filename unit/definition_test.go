package unit

import (
	"reflect"
	"strings"
	"testing"
)

func TestDefinition(t *testing.T) {
	var def definition

	contents := strings.NewReader(`[Unit]
Description=Description
Documentation=Documentation

Wants=Wants
Requires=Requires
Conflicts=Conflicts
Before=Before
After=After

[Install]
WantedBy=WantedBy
RequiredBy=RequiredBy`)

	if err := parseDefinition(contents, &def); err != nil {
		t.Errorf("Error parsing definition: %s", err)
	}

	defVal := reflect.ValueOf(&def).Elem()
	for i := 0; i < defVal.NumField(); i++ {

		section := defVal.Field(i)
		sectionType := section.Type()

		for j := 0; j < section.NumField(); j++ {
			option := struct {
				reflect.Value
				Name string
			}{
				section.Field(j),
				sectionType.Field(j).Name,
			}

			switch option.Kind() {
			case reflect.String:
				if option.String() != option.Name {
					t.Errorf("%s != %s", option, option.Name)
				}
			case reflect.SliceOf(reflect.TypeOf(reflect.String)).Kind():
				slice := option.Interface().([]string)

				if !reflect.DeepEqual(slice, []string{option.Name}) {
					t.Errorf("%s != %s", slice, []string{option.Name})
				}
			}
		}
	}
}