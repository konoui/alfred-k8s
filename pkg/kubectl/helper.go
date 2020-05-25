package kubectl

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// GetNameNamespace return name and namespace parsing kubectl objects.
// We assume　argument have name property and no empty value.
func GetNameNamespace(i interface{}) (name, ns string) {
	const unknown = "UnknownName"
	rv := reflect.Indirect(reflect.ValueOf(i))
	rt := rv.Type()
	nameField, ok := rt.FieldByName("Name")
	if !ok {
		return unknown, ""
	}
	name = rv.FieldByName(nameField.Name).String()
	if name == "" {
		return unknown, ""
	}

	nsField, ok := rt.FieldByName("Namespace")
	if !ok {
		return name, ""
	}
	ns = rv.FieldByName(nsField.Name).String()
	return name, ns
}

// MakeResourceStruct sets struct fields of `res` to corresponding line value with header location.
func MakeResourceStruct(line, rawHeaders string, res interface{}) error {
	rv := reflect.Indirect(reflect.ValueOf(res))
	rt := rv.Type()
	if rt.Kind() != reflect.Struct {
		return errors.New("argument is not struct")
	}

	indexMap := makeIndexMap(rawHeaders)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.Type.Kind() != reflect.String {
			continue
		}
		value := getFieldNameValueFromRawLineWithIndexMap(f.Name, line, indexMap)
		rv.Field(i).SetString(value)
	}
	return nil
}

// getFieldNameValueFromRawLineWithIndexMap confirms `indexMap` has `fiedName` and return corresponding `value` from `line`
func getFieldNameValueFromRawLineWithIndexMap(fieldName, line string, indexMap map[string]int) string {
	key := strings.ToLower(fieldName)
	begin, ok := indexMap[key]
	if !ok || begin == -1 {
		return ""
	}

	if begin > len(line) {
		return ""
	}
	return strings.Fields(line[begin:])[0]
}

// makeIndexMap parses rawHeaders and saves header name as key and header location(index) as value
func makeIndexMap(rawHeaders string) (indexMap map[string]int) {
	headers := strings.Fields(rawHeaders)
	indexMap = make(map[string]int, len(headers))
	for _, h := range headers {
		key := replaceHeader(h)
		// Note: search "NAME " with space as NAME will match NAMESPACE
		begin := strings.Index(rawHeaders, h+" ")
		if begin == -1 {
			// Note: try without space as last column will have no space
			begin = strings.Index(rawHeaders, h)
		}
		indexMap[key] = begin
	}
	return
}

func replaceHeader(h string) string {
	key := strings.ToLower(h)
	// Note replace deployment UP-TO-DATE with UPTODATE
	key = strings.Replace(key, "-", "", -1)
	// Note replace service PORT(S) with PORTS
	key = strings.Replace(key, "(", "", -1)
	key = strings.Replace(key, ")", "", -1)
	return key
}
