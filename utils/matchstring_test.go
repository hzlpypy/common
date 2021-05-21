package utils

import (
	"bou.ke/monkey"
	"github.com/google/uuid"
	"reflect"
	"testing"
	"time"
)

func TestNewMD5(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "len(str) == 0",
			args: args{},
			want: "",
		},
		{
			name: "all is ok",
			args: args{
				str: "test",
			},
			want: "098f6bcd4621d373cade4e832627b4f6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMD5(tt.args.str); got != tt.want {
				t.Errorf("NewMD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUUID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "all is ok",
			want: "ba44421c147949268adb7c0e817824c0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(uuid.UUID{}), "String", func(_ uuid.UUID) string {
				return "ba44421c-1479-4926-8adb-7c0e817824c0"
			})
			if got := NewUUID(); got != tt.want {
				t.Errorf("NewUUID() = %v, want %v", got, tt.want)
			}
			monkey.UnpatchAll()
		})
	}
}

func TestNewUnixtime(t *testing.T) {
	tests := []struct {
		name string
		want uint
	}{
		{
			name: "all is ok",
			want: 1617849131,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(time.Time{}), "Unix", func(_ time.Time) int64 {
				return int64(1617849131)
			})
			if got := NewUnixtime(); got != tt.want {
				t.Errorf("NewUnixtime() = %v, want %v", got, tt.want)
			}
			monkey.UnpatchAll()
		})
	}
}

func TestRandStringRunes(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "all is ok",
			args: args{
				n: 3,
			},
			want: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RandStringRunes(tt.args.n); reflect.TypeOf(got) != reflect.TypeOf(tt.want) {
				t.Errorf("RandStringRunes() = %v, want %v", got, tt.want)
			}
		})
	}
}
