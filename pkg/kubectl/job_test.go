package kubectl

import (
	"context"
	"reflect"
	"testing"
	"time"

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
		errMsg  string
		want    []string
		timeout time.Duration
	}{
		{
			name: "check stdout",
			args: args{
				name: "echo",
				args: []string{
					"line1\nline2",
				},
			},
			want:    []string{"line1", "line2"},
			timeout: 10 * time.Second,
		},
		{
			name: "wait the command",
			args: args{
				name: "bash",
				args: []string{
					"-c",
					"echo line1; sleep 5; echo line2;",
				},
			},
			want:    []string{"line1", "line2"},
			timeout: 10 * time.Second,
		},
		{
			name: "timeout",
			args: args{
				name: "bash",
				args: []string{
					"-c",
					"echo line1; sleep 10; echo line2",
				},
			},
			want:    []string{"line1"},
			timeout: 4 * time.Second,
			errMsg:  "job is canceled",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer goleak.VerifyNone(t)
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()
			status, errCh := startJob(ctx, tt.args.name, tt.args.args...)

			var got []string
			for line := range status.ReadLine() {
				got = append(got, line)
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("want %v, got %v", tt.want, got)
			}

			// wait for finish command
			err := <-errCh
			if err != nil {
				// TODO expected error or not
				if err.Error() != tt.errMsg {
					t.Errorf("want %v, got %v", tt.errMsg, err.Error())
				}
			}

			pid1 := status.Pid
			pid2 := status.Pid
			pid3 := status.Pid
			if pid1 != pid2 || pid2 != pid3 {
				t.Fatal(pid1, pid2)
			}
		})
	}
}
