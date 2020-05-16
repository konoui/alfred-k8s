package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/goleak"
)

var testAllServices = []*Service{
	{
		Namespace:  "test3-namespace",
		Name:       "service-test1",
		Type:       "ClusterIP",
		ClusterIP:  "10.100.49.33",
		ExternalIP: "<none>",
		Ports:      "8080/TCP",
		Age:        "11d",
	},
	{
		Namespace:  "test3-namespace",
		Name:       "service-test2",
		Type:       "ClusterIP",
		ClusterIP:  "10.100.75.11",
		ExternalIP: "<none>",
		Ports:      "8080/TCP,9901/TCP",
		Age:        "11d",
	},
}

var testServices = []*Service{
	{
		Name:       "service-test1",
		Type:       "ClusterIP",
		ClusterIP:  "10.100.49.33",
		ExternalIP: "<none>",
		Ports:      "8080/TCP",
		Age:        "11d",
	},
	{
		Name:       "service-test2",
		Type:       "ClusterIP",
		ClusterIP:  "10.100.75.11",
		ExternalIP: "<none>",
		Ports:      "8080/TCP,9901/TCP",
		Age:        "11d",
	},
}

func TestGetServices(t *testing.T) {
	tests := []struct {
		name string
		all  bool
		want []*Service
	}{
		{
			name: "list services",
			want: testServices,
		},
		{
			name: "list services in all namespaces",
			all:  true,
			want: testAllServices,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
			got, err := k.GetServices(tt.all)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
