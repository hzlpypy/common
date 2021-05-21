package v2

import (
	"bou.ke/monkey"
	"database/sql"
	"errors"
	"github.com/hzlpypy/common/databases/mysql"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func commonDB() *gorm.DB {
	return &gorm.DB{
		Statement: &gorm.Statement{},
	}
}

type TestBase struct {
}

func TestOpObjectRelation(t *testing.T) {
	var base *TestBase
	db := commonDB()
	object := new(string)
	type args struct {
		db       *gorm.DB
		opInters []*mysql.OpInter
	}
	tests := []struct {
		name         string
		args         args
		wantErr      bool
		wantCreatErr bool
	}{
		{
			name: "Op.Validate error",
			args: args{
				db: db,
				opInters: []*mysql.OpInter{
					&mysql.OpInter{
						Op:     mysql.Op_UPDATE,
						Object: object,
						Where:  "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Create error",
			args: args{
				db: db,
				opInters: []*mysql.OpInter{
					&mysql.OpInter{
						Op:     mysql.Op_CREATE,
						Object: base,
						Where:  "",
					},
					&mysql.OpInter{
						Op:     mysql.Op_CREATE,
						Object: object,
						Where:  "",
					},
					&mysql.OpInter{
						Op:     mysql.Op_CREATE,
						Object: object,
						Where:  "",
					},
					&mysql.OpInter{
						Op:     mysql.Op_UPDATE,
						Object: object,
						Where:  "ojbk",
					},
				},
			},
			wantErr:      true,
			wantCreatErr: true,
		},
		{
			name: "update all is ok",
			args: args{
				db: db,
				opInters: []*mysql.OpInter{
					&mysql.OpInter{
						Op:     mysql.Op_UPDATE,
						Object: object,
						Where:  "ojbk",
					},
					&mysql.OpInter{
						Op:     mysql.Op_UPDATE,
						Object: object,
						Where:  "ojbk",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "delete all is ok",
			args: args{
				db: db,
				opInters: []*mysql.OpInter{
					&mysql.OpInter{
						Op:     mysql.Op_DELETE,
						Object: object,
						Where:  "ojbk",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "op is invalidate",
			args: args{
				db: db,
				opInters: []*mysql.OpInter{
					&mysql.OpInter{
						Op:     "error",
						Object: object,
						Where:  "ojbk",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Transaction", func(_ *gorm.DB, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
				err := fc(db)
				return err
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Create", func(_ *gorm.DB, value interface{}) (tx *gorm.DB) {
				if tt.wantCreatErr {
					db.Error = errors.New("test")
					return db
				}
				db.Error = nil
				return db
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(db), "Updates", func(_ *gorm.DB, values interface{}) (tx *gorm.DB) {
				db.Error = nil
				return db
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(db), "Delete", func(_ *gorm.DB, value interface{}, conds ...interface{}) (tx *gorm.DB) {
				db.Error = nil
				return db
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(db), "Where", func(_ *gorm.DB, query interface{}, args ...interface{}) (tx *gorm.DB) {
				return db
			})
			if err := OpObjectRelation(tt.args.db, tt.args.opInters); (err != nil) != tt.wantErr {
				t.Errorf("OpObjectRelation() error = %v, wantErr %v", err, tt.wantErr)
			}
			monkey.UnpatchAll()
		})
	}
}

func TestOpTables(t *testing.T) {
	type args struct {
		tx       *gorm.DB
		opInter  *mysql.OpInter
		errChan  chan error
		i        int
		stopChan chan bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
