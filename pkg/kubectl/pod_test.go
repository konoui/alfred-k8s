package kubectl

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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

func TestGetPods(t *testing.T) {
	tests := []struct {
		name string
		all  bool
		want []*Pod
	}{
		{
			name: "list pods",
			want: testPods,
		},
		{
			name: "list pods in all namespaces",
			all:  true,
			want: testAllPods,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			k := SetupKubectl(t, nil)
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
