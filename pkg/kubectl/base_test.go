package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
	"go.uber.org/goleak"
)

const testResourceName = "pod"

var (
	testBaseAllPods = []*BaseResource{
		&BaseResource{
			Namespace: "test1-namespace",
			Name:      "test1-pod",
			Age:       "11d",
		},
		&BaseResource{
			Namespace: "test2-namespace",
			Name:      "test2-pod",
			Age:       "11d",
		},
	}
	testBasePods = []*BaseResource{
		&BaseResource{
			Name: "test1-pod",
			Age:  "11d",
		},
		&BaseResource{
			Name: "test2-pod",
			Age:  "11d",
		},
	}
)

var FakePodBaseResourceFunc = func(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllPods := GetByteFromTestFile(t, "testdata/raw-pods-in-all-namespaces.txt")
	rawDataPods := GetByteFromTestFile(t, "testdata/raw-pods.txt")
	if len(args) >= 1 {
		if args[0] == "get" && args[1] == testResourceName {
			if len(args) == 3 && args[2] == allNamespaceFlag {
				return &executor.Response{
					Stdout: []byte(rawDataAllPods),
				}, nil
			}

			return &executor.Response{
				Stdout: []byte(rawDataPods),
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

func TestPodBaseResource(t *testing.T) {
	tests := []struct {
		name     string
		fakeFunc FakeFunc
		all      bool
		want     []*BaseResource
	}{
		{
			name:     "list pods for base resource",
			fakeFunc: FakePodBaseResourceFunc,
			want:     testBasePods,
		},
		{
			name:     "list pods in all namespaces for base resource",
			fakeFunc: FakePodBaseResourceFunc,
			all:      true,
			want:     testBaseAllPods,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeFunc)
			got, err := k.GetBaseResources(testResourceName, tt.all)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
