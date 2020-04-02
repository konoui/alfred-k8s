package kubectl

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/konoui/alfred-k8s/pkg/executor"
)

var (
	testAllPodsRawData = `NAMESPACE	NAME	READY	STATUS	RESTARTS	AGE
	test1-namespace	test1-pod	1/1	Running	0	11d
	test2-namespace	test2-pod	2/2	Running	1	11d`
	testPodsRawData = `NAME	READY	STATUS	RESTARTS	AGE
	test1-pod	1/1	Running	0	11d
	test2-pod	2/2	Running	1	11d`
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

var FakePodFunc = func(args ...string) (*executor.Response, error) {
	if len(args) >= 3 {
		if args[1] == "pod" && args[2] == allNamespaceFlag {
			return &executor.Response{
				Stdout: []byte(testAllPodsRawData),
			}, nil
		}
	}
	if len(args) >= 2 {
		if args[1] == "pod" {
			return &executor.Response{
				Stdout: []byte(testPodsRawData),
			}, nil
		}
	}
	return &executor.Response{}, fmt.Errorf("match no command args")
}

func TestGetPods(t *testing.T) {
	tests := []struct {
		name         string
		fakeExecutor executor.Executor
		all          bool
		want         []*Pod
	}{
		{
			name:         "list pods",
			fakeExecutor: NewFakeExecutor(FakePodFunc),
			want:         testPods,
		},
		{
			name:         "list pods in all namespaces",
			fakeExecutor: NewFakeExecutor(FakePodFunc),
			all:          true,
			want:         testAllPods,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := SetupKubectl(t, tt.fakeExecutor)
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
