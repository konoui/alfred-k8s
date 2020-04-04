package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
	"go.uber.org/goleak"
)

var (
	testAllIngressesRawData = `test1-namespace	test1-ingress   *       ingress1.hoge.hoge   80      24h
	test2-namespace	test2-ingress   *       ingress2.hoge.hoge   80      24h`
	testIngressesRawData = `test1-ingress   *       ingress1.hoge.hoge   80      24h
	test2-ingress   *       ingress2.hoge.hoge   80      24h`
)

var testAllIngresses = []*Ingress{
	&Ingress{
		Namespace: "test1-namespace",
		Name:      "test1-ingress",
		Hosts:     "*",
		Address:   "ingress1.hoge.hoge",
		Ports:     "80",
		Age:       "24h",
	},
	&Ingress{
		Namespace: "test2-namespace",
		Name:      "test2-ingress",
		Hosts:     "*",
		Address:   "ingress2.hoge.hoge",
		Ports:     "80",
		Age:       "24h",
	},
}
var testIngresses = []*Ingress{
	&Ingress{
		Name:    "test1-ingress",
		Hosts:   "*",
		Address: "ingress1.hoge.hoge",
		Ports:   "80",
		Age:     "24h",
	},
	&Ingress{
		Name:    "test2-ingress",
		Hosts:   "*",
		Address: "ingress2.hoge.hoge",
		Ports:   "80",
		Age:     "24h",
	},
}

var FakeIngressFunc = func(args ...string) (*executor.Response, error) {
	if len(args) >= 4 {
		if args[1] == "ingress" && args[2] == allNamespaceFlag {
			return &executor.Response{
				Stdout: []byte(testAllIngressesRawData),
			}, nil
		}
	}
	if len(args) >= 2 {
		if args[1] == "ingress" {
			return &executor.Response{
				Stdout: []byte(testIngressesRawData),
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

func TestGetIngresses(t *testing.T) {
	tests := []struct {
		name         string
		fakeExecutor executor.Executor
		all          bool
		want         []*Ingress
	}{
		{
			name:         "list ingresses",
			fakeExecutor: NewFakeExecutor(FakeIngressFunc),
			want:         testIngresses,
		},
		{
			name:         "list deployments in all namespaces",
			fakeExecutor: NewFakeExecutor(FakeIngressFunc),
			all:          true,
			want:         testAllIngresses,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeExecutor)
			got, err := k.GetIngresses(tt.all)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
