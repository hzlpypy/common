package v1

import (
	"errors"
	"fmt"
	"github.com/hzlpypy/common/databases/mysql"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
)

// -- gorm v1 batch insert --
// GetBatchInsertSql:获取批量添加数据sql语句,支持内嵌struct(暂时只支持一个内嵌的struct)
func GetBatchInsertSql(objs []interface{}, tableName string) string {
	if len(objs) == 0 {
		return ""
	}
	fieldName := ""
	var valueTypeList []string
	fieldNum := reflect.TypeOf(objs[0]).NumField()
	fieldT := reflect.TypeOf(objs[0])
	superpositionFieldNum := fieldNum
	jilu := []int{}
	for a := 0; a < fieldNum; a++ {
		structFieldType := fieldT.Field(a).Type
		switch structFieldType.Kind() {
		case reflect.Struct:
			fieldBaseNum := structFieldType.NumField()
			// 记录内嵌struct，开始和结束的索引
			if len(jilu) == 0 {
				jilu = []int{a, a + fieldBaseNum - 1}
			}
			superpositionFieldNum += fieldBaseNum
			for b := 0; b < fieldBaseNum; b++ {
				name := gorm.ToColumnName(structFieldType.Field(b).Name)
				//name := GetColumnName(structFieldType.Field(b).Tag.Get("gorm"))
				// 添加字段名
				if a == fieldNum-1 {
					fieldName += fmt.Sprintf("`%s`", name)
				} else {
					fieldName += fmt.Sprintf("`%s`,", name)
				}
				structFieldBaseType := structFieldType.Field(b).Type
				valueTypeList = GetValueTypeList(structFieldBaseType, valueTypeList, b)
			}
		default:
			name := gorm.ToColumnName(fieldT.Field(a).Name)
			//name := GetColumnName(fieldT.Field(a).Tag.Get("gorm"))
			// 添加字段名
			if a == fieldNum-1 {
				fieldName += fmt.Sprintf("`%s`", name)
			} else {
				fieldName += fmt.Sprintf("`%s`,", name)
			}
			valueTypeList = GetValueTypeList(structFieldType, valueTypeList, a)

		}
	}
	var valueList []string
	for _, obj := range objs {
		objV := reflect.ValueOf(obj)
		v := "("
		for index, i := range valueTypeList {
			if index == len(valueTypeList)-1 {
				v += GetFormatField(objV, index, i, "", jilu)
			} else {
				v += GetFormatField(objV, index, i, ",", jilu)
			}
		}
		v += ")"
		valueList = append(valueList, v)
	}
	insertSql := fmt.Sprintf("insert into `%s` (%s) values %s", tableName, fieldName, strings.Join(valueList, ",")+";")
	return insertSql
}

// getValueTypeList:获取字段类型
func GetValueTypeList(structFieldType reflect.Type, valueTypeList []string, a int) []string {
	// 获取字段类型,暂支持以下类型,缺少请在下面添加
	if structFieldType.Name() == "string" {
		valueTypeList = append(valueTypeList, "string")
	} else if strings.Index(structFieldType.Name(), "uint") != -1 {
		valueTypeList = append(valueTypeList, "uint")
	} else if strings.Index(structFieldType.Name(), "int") != -1 {
		valueTypeList = append(valueTypeList, "int")
	} else if strings.Index(structFieldType.Name(), "bool") != -1 {
		valueTypeList = append(valueTypeList, "bool")
	} else if strings.Index(structFieldType.Name(), "float") != -1 {
		valueTypeList = append(valueTypeList, "float")
	}
	return valueTypeList
}

//  GetFormatField:获取字段类型值转为字符串
func GetFormatField(objV reflect.Value, index int, t string, sep string, jilu []int) string {
	one := jilu[0]
	two := jilu[1]
	if index >= one && index <= two {
		objKind := objV.Field(one).Kind()
		if objKind == reflect.Struct {
			objV = reflect.ValueOf(objV.Field(one).Interface())
		}
		index = index - one
	} else {
		if index >= two {
			index = one + (index - two)
		}
	}
	// 依赖getValueTypeList,补充了什么类型，添加什么类型
	v := ""
	switch t {
	case "string":
		newStr := strings.Replace(objV.Field(index).String(), "'", "`", -1)
		v += fmt.Sprintf("'%s'%s", newStr, sep)
	case "uint":
		v += fmt.Sprintf("%d%s", objV.Field(index).Uint(), sep)
	case "int":
		v += fmt.Sprintf("%d%s", objV.Field(index).Int(), sep)
	case "bool":
		var boolInt int
		if objV.Field(index).Bool() {
			boolInt = 1
		}
		v += fmt.Sprintf("%d%s", boolInt, sep)
	case "float":
		v += fmt.Sprintf("%.3f%s", objV.Field(index).Float(), sep)
	}
	return v

}

func OpTables(tx *gorm.DB, opInter *mysql.OpInter, errChan chan error, i int, stopChan chan bool) {
	select {
	case stop := <-stopChan:
		if stop {
			//fmt.Print(fmt.Sprintf("error count:%d", i) + "\n")
			return
		}
		//fmt.Print(fmt.Sprintf("success count:%d", i) + "\n")
		err := opInter.Op.Validate(opInter)
		if err != nil {
			errChan <- err
			return
		}
		object := opInter.Object
		switch opInter.Op {
		case mysql.Op_CREATE:
			err := tx.Create(object).Error
			errChan <- err
		case mysql.Op_UPDATE:
			err := tx.Updates(object).Where(opInter.Where).Error
			errChan <- err
		case mysql.Op_DELETE:
			err := tx.Where(opInter.Where).Delete(object).Error
			errChan <- err
		default:
			errChan <- errors.New("op is invalidate")
		}
		return
	}
}

// OpObjectRelation:op指定模型及其关联数据，支持多个模型op(create，update，delete)
func OpObjectRelation(db *gorm.DB, opInters []*mysql.OpInter) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var errChan = make(chan error, 5)
		// 是否执行关闭
		var stopChan = make(chan bool)
		var validateOpInter []*mysql.OpInter
		for _, opInter := range opInters {
			if reflect.ValueOf(opInter.Object).IsNil() {
				continue
			}
			validateOpInter = append(validateOpInter, opInter)
		}
		// 处于等待状态的协程数量
		goWaitNum := len(validateOpInter)
		for i, info := range validateOpInter {
			go OpTables(tx, info, errChan, i, stopChan)
		}
		stopChan <- false
		goWaitNum -= 1
		var count int
		for {
			select {
			case err := <-errChan:
				count += 1
				if err != nil {
					// 关闭等待状态协程，否则会导致产生更多协程阻塞，且无法被消费/回收,导致内存泄漏
					for i := 0; i < goWaitNum; i++ {
						stopChan <- true
					}
					close(errChan)
					close(stopChan)
					return err
				} else {
					if count == len(validateOpInter) {
						close(errChan)
						close(stopChan)
						return nil
					} else {
						stopChan <- false
						goWaitNum -= 1
					}
				}
			}
		}
	})

}
