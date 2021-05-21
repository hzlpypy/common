package v1

import (
	"bou.ke/monkey"
	"errors"
	"github.com/hzlpypy/common/databases/mysql"
	"github.com/jinzhu/gorm"
	"reflect"
	"testing"
)

type TestChildren struct {
	Common string `gorm:"type:varchar(32);not null"`
}

type TestBase struct {
	Children TestChildren
	Id       string  `gorm:"type:varchar(32);not null"`
	Name     string  `gorm:"type:varchar(32);not null"`
	Gender   uint    `gorm:"type:tinyint(1);not null"`
	Sort     int     `gorm:"type:int(3);not null"`
	PlayLol  bool    `gorm:"type:bool;not null"`
	Money    float32 `gorm:"type:float(3,2);not null"`
}

type TestEnd struct {
	Id       string `gorm:"type:varchar(32);not null"`
	Name     string `gorm:"type:varchar(32);not null"`
	Children TestChildren
}

func TestGetBatchInsertSql(t *testing.T) {
	type args struct {
		objs      []interface{}
		tableName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "len(objs) == 0",
			args: args{},
			want: "",
		},
		{
			name: "all is ok default a == fieldNum-1",
			args: args{
				objs: []interface{}{
					TestBase{
						Children: TestChildren{
							Common: "123",
						},
						Id:      "1",
						Name:    "test",
						Gender:  1,
						Sort:    2,
						PlayLol: true,
						Money:   2.32,
					},
				},
				tableName: "test_base",
			},
			want: "insert into `test_base` (`common`,`id`,`name`,`gender`,`sort`,`play_lol`,`money`) values ('123','1','test',1,2,1,2.320);",
		},
		{
			name: "all is ok Struct a == fieldNum-1",
			args: args{
				objs: []interface{}{
					TestEnd{
						Id:   "1",
						Name: "test",
						Children: TestChildren{
							Common: "123",
						},
					},
				},
				tableName: "test_base",
			},
			want: "insert into `test_base` (`id`,`name`,`common`) values ('1','test','123');",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetBatchInsertSql(tt.args.objs, tt.args.tableName); got != tt.want {
				t.Errorf("GetBatchInsertSql() = %v, want %v", got, tt.want)
			}
		})
	}
}

func commonDB() *gorm.DB {
	return &gorm.DB{}
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
			monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Transaction", func(_ *gorm.DB, fc func(tx *gorm.DB) error) error {
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
			monkey.PatchInstanceMethod(reflect.TypeOf(db), "Updates", func(_ *gorm.DB, values interface{}, ignoreProtectedAttrs ...bool) *gorm.DB {
				db.Error = nil
				return db
			})
			monkey.PatchInstanceMethod(reflect.TypeOf(db), "Delete", func(_ *gorm.DB, value interface{}, where ...interface{}) *gorm.DB {
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
