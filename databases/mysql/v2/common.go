package v2

import (
	"errors"
	"github.com/hzlpypy/common/databases/mysql"
	_ "github.com/go-sql-driver/mysql"
	libMyql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"reflect"
	"time"
)

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



func NewDbConnection(connStr *mysql.Config) (*gorm.DB, error) {
	dbConfiguration := connStr.Username + ":" + connStr.Password + "@(" + connStr.Host + ":" + connStr.Port + ")/" + connStr.DBName + "?charset=" + connStr.Charset + "&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			Colorful:      true,        // Disable color
			LogLevel:      logger.Error,
		},
	)
	if connStr.Debug {
		newLogger = newLogger.LogMode(logger.Info)
	}
	db, err := gorm.Open(libMyql.Open(dbConfiguration), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		Logger:                                   newLogger,
		SkipDefaultTransaction:                   true,
		PrepareStmt:                              true,  // 在执行任何 SQL 时都会创建一个 prepared statement 并将其缓存，以提高后续的效率
		DisableAutomaticPing:                     false, // 在完成初始化后，GORM 会自动 ping 数据库以检查数据库的可用性，若要禁用该特性，可将其设置为 true
		DisableForeignKeyConstraintWhenMigrating: connStr.DisableForeignKeyConstraintWhenMigrating, // 在 AutoMigrate 或 CreateTable 时，GORM 会自动创建外键约束，若要禁用该特性，可将其设置为 true
	})
	if err != nil {
		log.Println(dbConfiguration)
		log.Println(err)
		log.Printf("dbConfiguration=%v\nerr=%v", dbConfiguration, err)
		return nil, err
	}
	db.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("db.Set err=%v", err)
		return nil, err
	}
	// connect pool
	sqlDB.SetConnMaxLifetime(connStr.ConnMaxLifetime)
	sqlDB.SetMaxIdleConns(connStr.MaxIdleConns)
	sqlDB.SetMaxOpenConns(connStr.MaxOpenConns)
	_ = sqlDB.Ping()
	log.Println("database configuration is loaded.")
	return db, nil
}
