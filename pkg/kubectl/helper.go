package kubectl

import "reflect"

// GetNameNamespace return name and namespace parsing kubectl objects.
// We assume　argument have name property and no empty value.
func GetNameNamespace(i interface{}) (name, ns string) {
	const unknown = "UnknownName"
	rv := reflect.Indirect(reflect.ValueOf(i))
	rt := rv.Type()
	if _, ok := rt.FieldByName("Name"); !ok {
		return unknown, ""
	}
	name = rv.FieldByName("Name").String()
	if name == "" {
		return unknown, ""
	}

	if _, ok := rt.FieldByName("Namespace"); !ok {
		return name, ""
	}
	ns = rv.FieldByName("Namespace").String()
	return name, ns
}
