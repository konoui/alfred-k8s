package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
	"go.uber.org/goleak"
)

var (
	testAllPods = []*Pod{
		&Pod{
			Namespace: "test1-namespace",
			Name:      "test1-pod",
			Ready:     "1/1",
			Status:    "Running",
			Restarts:  "0",
			Age:       "11d",
		},
		&Pod{
			Namespace: "test2-namespace",
			Name:      "test2-pod",
			Ready:     "2/2",
			Status:    "Running",
			Restarts:  "1",
			Age:       "11d",
		},
	}
	testPods = []*Pod{
		&Pod{
			Name:     "test1-pod",
			Ready:    "1/1",
			Status:   "Running",
			Restarts: "0",
			Age:      "11d",
		},
		&Pod{
			Name:     "test2-pod",
			Ready:    "2/2",
			Status:   "Running",
			Restarts: "1",
			Age:      "11d",
		},
	}
)

var FakePodFunc = func(t *testing.T, args ...string) (*executor.Response, error) {
	rawDataAllPods := GetByteFromTestFile(t, "testdata/raw-pods-in-all-namespaces.txt")
	rawDataPods := GetByteFromTestFile(t, "testdata/raw-pods.txt")
	if len(args) >= 3 {
		if args[1] == "pod" && args[2] == allNamespaceFlag {
			return &executor.Response{
				Stdout: rawDataAllPods,
			}, nil
		}
	}
	if len(args) >= 2 {
		if args[1] == "pod" {
			return &executor.Response{
				Stdout: rawDataPods,
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

func TestGetPods(t *testing.T) {
	tests := []struct {
		name     string
		fakeFunc FakeFunc
		all      bool
		want     []*Pod
	}{
		{
			name:     "list pods",
			fakeFunc: FakePodFunc,
			want:     testPods,
		},
		{
			name:     "list pods in all namespaces",
			fakeFunc: FakePodFunc,
			all:      true,
			want:     testAllPods,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, tt.fakeFunc)
			got, err := k.GetPods(tt.all)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("+want -got\n%+v", diff)
			}
		})
	}
}
