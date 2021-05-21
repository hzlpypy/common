package utils

import "testing"

func TestRegexpEmail(t *testing.T) {
	type args struct {
		m string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true",
			args: args{
				m: "123yet@163.com",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RegexpEmail(tt.args.m); got != tt.want {
				t.Errorf("RegexpEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpIPV4(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true",
			args: args{
				ip: "127.0.0.1",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RegexpIPV4(tt.args.ip); got != tt.want {
				t.Errorf("RegexpIPV4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegexpMobile(t *testing.T) {
	type args struct {
		m string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "true",
			args: args{
				m: "13235253635",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RegexpMobile(tt.args.m); got != tt.want {
				t.Errorf("RegexpMobile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateCardNumber(t *testing.T) {
	type args struct {
		cardNumber string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "false",
			args: args{
				cardNumber: "123yet@163.com",
			},
			want: false,
		},
		{
			name: "true",
			args: args{
				cardNumber: "421087199402118784",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateCardNumber(tt.args.cardNumber); got != tt.want {
				t.Errorf("ValidateCardNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
