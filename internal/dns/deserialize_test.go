package dns

import (
	"reflect"
	"testing"
)

func TestMessage_parseDNSNameInQuery(t *testing.T) {
	tests := []struct {
		name    string
		want    MessageQuestion
		data    []byte
		wantErr bool
	}{
		{
			name: "Qname=google.com, Qtype = 1, Qclass = 1",
			want: MessageQuestion{
				Qname:  []string{"google", "com"},
				Qtype:  1,
				Qclass: 1,
			},
			data:    []byte{0x6, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x3, 0x63, 0x6f, 0x6d, 0x0, 0x0, 0x1, 0x0, 0x1},
			wantErr: false,
		},
		{
			name: "Qname=example.com, Qtype = 0, Qclass = 10",
			want: MessageQuestion{
				Qname:  []string{"example", "com"},
				Qtype:  0,
				Qclass: 10,
			},
			data:    []byte{7, 'e', 'x', 'a', 'm', 'p', 'l', 'e', 3, 'c', 'o', 'm', 0x0, 0x0, 0x0, 0x0, 0xA},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Message{}

			err := m.parseDNSNameInQuery(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDNSNameInQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(m.Question, tt.want) {
				t.Errorf("parseDNSNameInQuery() = %v, want %v", m.Question, tt.want)
			}
		})
	}
}

func Test_parseFlags(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Flags
		wantErr bool
	}{
		{
			name: "qr oc=12 aa tr rd ra rc=5",
			args: args{
				data: []byte{0b11100111, 0b10000101},
			},
			want: Flags{
				Query:               1,
				Opcode:              12,
				AuthoritativeAnswer: 1,
				Truncation:          1,
				RecursionDesired:    1,
				RecursionAvailable:  1,
				Z:                   0,
				ResponseCode:        5,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFlags(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFlags() got = %v, want %v", got, tt.want)
			}
		})
	}
}
