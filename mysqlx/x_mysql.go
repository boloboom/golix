package mysqlx

import (
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlSetting struct {
	RunMode string

	User     string
	Password string
	Host     string

	MySqlMaxIdleConn int
	MySqlMaxOpenConn int
	MySqlConnMaxLife time.Duration
}

var (
	mysqlSetting *MysqlSetting
	mysqlConnMap sync.Map
)

// Setup initializes the MySQLSetting.
//
// The function takes a pointer to a MysqlSetting struct as a parameter.
// It sets the global variable mysqlSetting to the provided setting.
func Setup(setting *MysqlSetting) {
	mysqlSetting = setting
}

// GetDb returns a *gorm.DB and an error.
//
// It takes a string parameter named "name" which represents the name of the database.
// It returns a *gorm.DB and an error.
func GetDb(name string) (*gorm.DB, error) {
	if db, ok := mysqlConnMap.Load(name); ok {
		return db.(*gorm.DB), nil
	}
	if db, err := connectDb(name); err == nil {
		mysqlConnMap.Store(name, db)
		return db, nil
	} else {
		log.Printf("community database init failed :%s", name)
		return nil, err
	}
}

// connectDb connects to a MySQL database using the provided name.
//
// Parameters:
//   - name: the name of the database to connect to.
//
// Returns:
//   - *gorm.DB: the connected database instance.
//   - error: an error if the connection fails.
func connectDb(name string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlSetting.User,
			mysqlSetting.Password,
			mysqlSetting.Host,
			name,
		)))
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if mysqlSetting.MySqlMaxIdleConn == 0 {
		mysqlSetting.MySqlMaxIdleConn = 10
	}
	sqlDB.SetMaxIdleConns(mysqlSetting.MySqlMaxIdleConn)

	if mysqlSetting.MySqlMaxOpenConn == 0 {
		mysqlSetting.MySqlMaxOpenConn = 50
	}
	sqlDB.SetMaxOpenConns(mysqlSetting.MySqlMaxOpenConn)

	if mysqlSetting.MySqlConnMaxLife == 0 {
		mysqlSetting.MySqlConnMaxLife = 3600
	}
	sqlDB.SetConnMaxLifetime(mysqlSetting.MySqlConnMaxLife)

	return db, nil
}
