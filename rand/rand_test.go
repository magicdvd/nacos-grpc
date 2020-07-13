package rand

import "testing"

func TestUint32n(t *testing.T) {
	type args struct {
		n uint32
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "zero",
			args: args{
				n: 0,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Uint32n(tt.args.n); got != tt.want {
				t.Errorf("Uint32n() = %v, want %v", got, tt.want)
			}
		})
	}
}
