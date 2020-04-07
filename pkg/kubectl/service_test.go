package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.uber.org/goleak"
)

var testAllServices = []*Service{
	&Service{
		Namespace:  "test3-namespace",
		Name:       "service-test1",
		Type:       "ClusterIP",
		ClusterIP:  "10.100.49.33",
		ExternalIP: "<none>",
		Ports:      "8080/TCP",
		Age:        "11d",
	},
	&Service{
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
	&Service{
		Name:       "service-test1",
		Type:       "ClusterIP",
		ClusterIP:  "10.100.49.33",
		ExternalIP: "<none>",
		Ports:      "8080/TCP",
		Age:        "11d",
	},
	&Service{
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
		name     string
		fakeFunc FakeFunc
		all      bool
		want     []*Service
	}{
		{
			name:     "list services",
			fakeFunc: FakeServiceFunc,
			want:     testServices,
		},
		{
			name:     "list services in all namespaces",
			fakeFunc: FakeServiceFunc,
			all:      true,
			want:     testAllServices,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeFunc)
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
