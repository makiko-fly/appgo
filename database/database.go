package database

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/oxfeeefeee/appgo"
)

func MysqlConnStr() string {
	c := &appgo.Conf.Mysql
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.DbName, c.Charset)
}

func Open(driver string) (*gorm.DB, error) {
	switch driver {
	case "mysql":
		return gorm.Open(driver, MysqlConnStr())
	default:
		return nil, fmt.Errorf("database: unknown driver %q", driver)
	}
}

func SqlStr(str string) sql.NullString {
	return sql.NullString{str, true}
}
