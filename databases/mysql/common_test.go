package mysql

import "testing"

func TestFilteredSQLInject(t *testing.T) {
	type args struct {
		toMatchStrs []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "MatchString false",
			args: args{
				toMatchStrs: []string{"id"},
			},
			want: false,
		},

		{
			name: "MatchString true",
			args: args{
				toMatchStrs: []string{"select id from test"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilteredSQLInject(tt.args.toMatchStrs); got != tt.want {
				t.Errorf("FilteredSQLInject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetColumnName(t *testing.T) {
	type args struct {
		jsonName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "all is ok",
			args: args{
				jsonName: "default:1;column:id",
			},
			want: "id",
		},
		{
			name: "return ''",
			args: args{
				jsonName: "default:1",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetColumnName(tt.args.jsonName); got != tt.want {
				t.Errorf("GetColumnName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLimitOffset(t *testing.T) {
	type args struct {
		page     int32
		pageSize int32
	}
	tests := []struct {
		name  string
		args  args
		want  uint32
		want1 uint32
	}{
		{
			name: "all is ok",
			args: args{
				page:     1,
				pageSize: 10,
			},
			want:  10,
			want1: 0,
		},
		{
			name: "all is ok",
			args: args{
				page:     0,
				pageSize: 10,
			},
			want:  10,
			want1: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetLimitOffset(tt.args.page, tt.args.pageSize)
			if got != tt.want {
				t.Errorf("GetLimitOffset() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetLimitOffset() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
