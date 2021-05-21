package mysql

import "testing"

func TestOp_Validate(t *testing.T) {
	type args struct {
		opInter *OpInter
	}
	tests := []struct {
		name    string
		o       Op
		args    args
		wantErr bool
	}{
		{
			name:    "o == Op_CREATE ",
			o:       Op_CREATE,
			args:    args{},
			wantErr: false,
		},
		{
			name: "len(opInter.Where) == 0",
			o:    Op_UPDATE,
			args: args{
				opInter: &OpInter{
					Where: "",
				},
			},
			wantErr: true,
		},
		{
			name: "all is ok",
			o:    Op_UPDATE,
			args: args{
				opInter: &OpInter{
					Where: "id = 1",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.o.Validate(tt.args.opInter); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
