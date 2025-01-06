package dns

import (
	"reflect"
	"testing"
)

func Test_serializeFlags(t *testing.T) {
	type args struct {
		flags Flags
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		{
			name: "qr rd ra",
			args: args{
				flags: Flags{
					Query:               1,
					Opcode:              0,
					AuthoritativeAnswer: 0,
					Truncation:          0,
					RecursionDesired:    1,
					RecursionAvailable:  1,
					Z:                   0,
					ResponseCode:        0,
				},
			},
			want: 0b1000000110000000,
		},
		{
			name: "qr aa tr rd ra",
			args: args{
				flags: Flags{
					Query:               1,
					Opcode:              0,
					AuthoritativeAnswer: 1,
					Truncation:          1,
					RecursionDesired:    1,
					RecursionAvailable:  1,
					Z:                   0,
					ResponseCode:        0,
				},
			},
			want: 0b1000011110000000,
		},
		{
			name: "qr oc=1 aa tr rd ra",
			args: args{
				flags: Flags{
					Query:               1,
					Opcode:              1,
					AuthoritativeAnswer: 1,
					Truncation:          1,
					RecursionDesired:    1,
					RecursionAvailable:  1,
					Z:                   0,
					ResponseCode:        0,
				},
			},
			want: 0b1000111110000000,
		},
		{
			name: "qr oc=12 aa tr rd ra",
			args: args{
				flags: Flags{
					Query:               1,
					Opcode:              12,
					AuthoritativeAnswer: 1,
					Truncation:          1,
					RecursionDesired:    1,
					RecursionAvailable:  1,
					Z:                   0,
					ResponseCode:        0,
				},
			},
			want: 0b1110011110000000,
		},
		{
			name: "qr oc=12 aa tr rd ra rc=5",
			args: args{
				flags: Flags{
					Query:               1,
					Opcode:              12,
					AuthoritativeAnswer: 1,
					Truncation:          1,
					RecursionDesired:    1,
					RecursionAvailable:  1,
					Z:                   0,
					ResponseCode:        5,
				},
			},
			want: 0b1110011110000101,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serializeFlags(tt.args.flags); got != tt.want {
				t.Errorf("serializeFlags() = %b, want %b", got, tt.want)
			}
		})
	}
}

func Test_serializeQuery(t *testing.T) {
	type args struct {
		q MessageQuestion
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "[OK] google.com",
			args: args{
				q: MessageQuestion{
					Qname:  []string{"google", "com"},
					Qclass: 1,
					Qtype:  1,
				},
			},
			want: []byte{0x6, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x3, 0x63, 0x6f, 0x6d, 0x0, 0x0, 0x1, 0x0, 0x1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serializeQuery(tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("serializeQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serializeAnswer(t *testing.T) {
	type args struct {
		a []MessageResourceRecord
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "[OK] google.com -> 172.217.23.206",
			args: args{
				a: []MessageResourceRecord{
					{
						Name:  []string{"google", "com"},
						Type:  1,
						Class: 1,
						TTL:   20,
						Data:  []int{172, 217, 23, 206},
					},
				},
			},
			want: []byte{0xc0, 0xc, 0x0, 0x1, 0x0, 0x1, 0x0, 0x0, 0x0, 0x14, 0x0, 0x4, 0xac, 0xd9, 0x17, 0xce},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serializeAnswer(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("serializeAnswer() = %v, want %v", got, tt.want)
			}
		})
	}
}
