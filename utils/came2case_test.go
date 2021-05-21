package utils

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBuffer_Append(t *testing.T) {
	b := new(bytes.Buffer)
	type fields struct {
		Buffer *bytes.Buffer
	}
	type args struct {
		i interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Buffer
	}{
		{
			name: "int",
			fields: fields{
				Buffer: b,
			},
			args: args{
				i: 1,
			},
			want: &Buffer{
				b,
			},
		},
		{
			name: "int64",
			fields: fields{
				Buffer: b,
			},
			args: args{
				i: int64(1),
			},
			want: &Buffer{
				b,
			},
		},
		{
			name: "uint",
			fields: fields{
				Buffer: b,
			},
			args: args{
				i: uint(1),
			},
			want: &Buffer{
				b,
			},
		},
		{
			name: "uint64",
			fields: fields{
				Buffer: b,
			},
			args: args{
				i: uint64(1),
			},
			want: &Buffer{
				b,
			},
		},
		{
			name: "string",
			fields: fields{
				Buffer: b,
			},
			args: args{
				i: "1",
			},
			want: &Buffer{
				b,
			},
		},
		{
			name: "[]byte",
			fields: fields{
				Buffer: b,
			},
			args: args{
				i: []byte("1"),
			},
			want: &Buffer{
				b,
			},
		},
		{
			name: "rune",
			fields: fields{
				Buffer: b,
			},
			args: args{
				i: rune(1),
			},
			want: &Buffer{
				b,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				Buffer: tt.fields.Buffer,
			}
			if got := b.Append(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Append() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuffer_append(t *testing.T) {
	b := &bytes.Buffer{}
	b = nil
	type fields struct {
		Buffer *bytes.Buffer
	}
	type args struct {
		i string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Buffer
	}{
		{
			name: "panic",
			fields: fields{
				Buffer: b,
			},
			args: args{
				i: "123",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Buffer{
				Buffer: tt.fields.Buffer,
			}
			if got := b.append(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Append() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCamel2Case(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{

		{
			name: "all is ok",
			args: args{
				name: "ASDFig",
			},
			want: "a_s_d_fig",
		},
		{
			name: "all is ok",
			args: args{
				name: "aSDFig",
			},
			want: "a_s_d_fig",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Camel2Case(tt.args.name); got != tt.want {
				t.Errorf("Camel2Case() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBuffer(t *testing.T) {
	tests := []struct {
		name string
		want *Buffer
	}{
		{
			name: "all is ok",
			want: &Buffer{
				Buffer: new(bytes.Buffer),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBuffer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}
