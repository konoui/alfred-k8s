package kubectl

import "testing"

func TestGetNameNamespace(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantNs   string
	}{
		{
			name: "parse non-namespaced pod",
			args: args{
				i: testPods[0],
			},
			wantName: "test1-pod",
			wantNs:   "",
		},
		{
			name: "parse namespaced pod",
			args: args{
				i: testAllPods[0],
			},
			wantName: "test1-pod",
			wantNs:   "test1-namespace",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotNs := GetNameNamespace(tt.args.i)
			if gotName != tt.wantName {
				t.Errorf("GetNameNamespace() gotName = %v, want %v", gotName, tt.wantName)
			}
			if gotNs != tt.wantNs {
				t.Errorf("GetNameNamespace() gotNs = %v, want %v", gotNs, tt.wantNs)
			}
		})
	}
}
