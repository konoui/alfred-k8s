package kubectl

import (
	"reflect"
	"testing"

	"go.uber.org/goleak"
)

func Test_startJob(t *testing.T) {
	type args struct {
		name string
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		want    []string
	}{
		{
			name: "check stdout",
			args: args{
				name: "echo",
				args: []string{
					"line1\nline2",
				},
			},
			want: []string{"line1", "line2"},
		},
		{
			name: "wait the command",
			args: args{
				name: "bash",
				args: []string{
					"-c",
					"echo line1; echo line2; sleep 5",
				},
			},
			want: []string{"line1", "line2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			status, errChan := startJob(tt.args.name, tt.args.args...)
			var got []string
			for line := range status.ReadLine() {
				got = append(got, line)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want %v, got %v", tt.want, got)
			}

			// wait for finish command
			err := <-errChan
			if err != nil {
				t.Fatal(err)
			}

			pid1 := status.Pid()
			pid2 := status.Pid()
			pid3 := status.Pid()
			if pid1 != pid2 || pid2 != pid3 {
				t.Fatal(pid1, pid2)
			}

		})
	}
}
