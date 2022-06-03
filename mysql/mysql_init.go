/*
mysql_init.go
mysql配置及连接
*/

package mysql

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zituocn/esme/logx"
)

var (
	dbs map[string]*gorm.DB

	defaultDBName string
)

const (
	dbType = "mysql"
)

type DBConfig struct {
	Name            string
	User            string
	Password        string
	Host            string
	Port            int
	Debug           bool
	DisablePrepared bool
}

func InitDefaultDB(db *DBConfig) (err error) {
	if db == nil {
		err = errors.New("[mysql] no database to initialize")
		return
	}
	defaultDBName = db.Name
	dbs = make(map[string]*gorm.DB, 1)
	newORM(db)
	return
}

func InitDB(list []*DBConfig) (err error) {
	if len(list) == 0 {
		err = errors.New("[mysql] no database to initialize")
		return
	}
	dbs = make(map[string]*gorm.DB, len(list))
	for _, item := range list {
		newORM(item)
	}

	return
}

// GetORM return default *gorm.DB
func GetORM() *gorm.DB {
	m, ok := dbs[defaultDBName]
	if !ok {
		logx.Panic("[DB] the database is not initialized, please refer to the instructions for use")
	}
	return m
}

// GetORMByName get orm by name
func GetORMByName(name string) *gorm.DB {
	m, ok := dbs[name]
	if !ok {
		logx.Panic("[DB] the database is not initialized, please refer to the instructions for use")
	}
	return m
}

// newORM a new ORM
func newORM(db *DBConfig) {
	var (
		orm *gorm.DB
		err error
	)
	if db.User == "" || db.Password == "" || db.Host == "" || db.Port == 0 {
		panic(fmt.Sprintf("[DB]-[%s] failed to get database configuration", db.Name))
	}

	str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", db.User, db.Password, db.Host, db.Port, db.Name) + "?charset=utf8mb4&parseTime=true&loc=Local"
	if db.DisablePrepared {
		str = str + "&interpolateParams=true"
	}
	for orm, err = gorm.Open(dbType, str); err != nil; {
		logx.Errorf("[DB]-[%v] connection exception: %retrying: %v", db.Name, err, str)
		time.Sleep(5 * time.Second)
		orm, err = gorm.Open(dbType, str)
	}
	orm.LogMode(db.Debug)
	orm.CommonDB()
	orm.SingularTable(true)
	dbs[db.Name] = orm
}
