package kubectl

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_parsePortsFromResponse(t *testing.T) {
	type args struct {
		resp *Response
	}
	tests := []struct {
		name      string
		args      args
		wantPorts []string
	}{
		{
			name:      "tcp/udp port",
			wantPorts: []string{"80", "8080"},
			args: args{
				resp: &Response{
					stdout: bytes.NewBufferString("80/TCP\n8080/TCP\n9090/UDP"),
				},
			},
		},
		{
			name: "invalid port",
			args: args{
				resp: &Response{
					stdout: bytes.NewBufferString("AAA/AAA"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotPorts := parsePortsFromResponse(tt.args.resp); !reflect.DeepEqual(gotPorts, tt.wantPorts) {
				t.Errorf("parsePortsFromResponse() = %v, want %v", gotPorts, tt.wantPorts)
			}
		})
	}
}
