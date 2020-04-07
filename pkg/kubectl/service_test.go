package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
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

var FakeServiceFunc = func(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllServices := GetByteFromTestFile(t, "testdata/raw-services-in-all-namespaces.txt")
	rawDataServices := GetByteFromTestFile(t, "testdata/raw-services.txt")
	if len(args) >= 4 {
		if args[1] == "service" && args[2] == allNamespaceFlag {
			return &executor.Response{
				Stdout: rawDataAllServices,
			}, nil
		}
	}
	if len(args) >= 2 {
		if args[1] == "service" {
			return &executor.Response{
				Stdout: rawDataServices,
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
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
