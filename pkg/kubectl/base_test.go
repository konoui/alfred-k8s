package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
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

var FakePodResourceFunc = func(args ...string) (*executor.Response, error) {
	if len(args) >= 1 {
		if args[0] == "get" && args[1] == testResourceName {
			if len(args) == 3 && args[2] == allNamespaceFlag {
				return &executor.Response{
					Stdout: []byte(testAllPodsRawData),
				}, nil
			}

			return &executor.Response{
				Stdout: []byte(testPodsRawData),
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

func TestPodResource(t *testing.T) {
	tests := []struct {
		name         string
		fakeExecutor executor.Executor
		all          bool
		want         []*BaseResource
	}{
		{
			name:         "list pods for base resource",
			fakeExecutor: NewFakeExecutor(FakePodResourceFunc),
			want:         testBasePods,
		},
		{
			name:         "list pods in all namespaces for base resource",
			fakeExecutor: NewFakeExecutor(FakePodResourceFunc),
			all:          true,
			want:         testBaseAllPods,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := SetupKubectl(t, tt.fakeExecutor)
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
