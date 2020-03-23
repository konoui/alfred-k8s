package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
)

var (
	testAllServicesRawData = `test3-namespace   service-test1             ClusterIP   10.100.49.33     <none>        8080/TCP            11d
	test3-namespace   service-test2           ClusterIP   10.100.75.11     <none>        8080/TCP,9901/TCP   11d`
	testServicesRawData = `service-test1     ClusterIP   10.100.49.33     <none>        8080/TCP            11d
	service-test2   ClusterIP   10.100.75.11     <none>        8080/TCP,9901/TCP   11d`
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

var FakeServiceFunc = func(args ...string) (*executor.Response, error) {
	if len(args) >= 4 {
		if args[1] == "service" && args[2] == allNamespaceFlag {
			return &executor.Response{
				Stdout: []byte(testAllServicesRawData),
			}, nil
		}
	}
	if len(args) >= 2 {
		if args[1] == "service" {
			return &executor.Response{
				Stdout: []byte(testServicesRawData),
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

func TestGetServices(t *testing.T) {
	tests := []struct {
		name         string
		fakeExecutor executor.Executor
		want         []*Service
	}{
		{
			name:         "list services",
			fakeExecutor: NewFakeExecutor(FakeServiceFunc),
			want:         testServices,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := SetupKubectl(t, tt.fakeExecutor)
			got, err := k.GetServices()
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}

func TestGetAllServices(t *testing.T) {
	tests := []struct {
		name         string
		fakeExecutor executor.Executor
		want         []*Service
	}{
		{
			name:         "list services in all namespaces",
			fakeExecutor: NewFakeExecutor(FakeServiceFunc),
			want:         testAllServices,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := SetupKubectl(t, tt.fakeExecutor)
			got, err := k.GetAllServices()
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
