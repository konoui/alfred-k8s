package kubectl

import "reflect"

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
