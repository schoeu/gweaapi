package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func ErrHandle(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// 创建数据库链接
func OpenDb(dbTyepe string, dbStr string) *sql.DB {
	if dbTyepe == "" {
		dbTyepe = "mysql"
	}
	db, err := sql.Open(dbTyepe, dbStr)
	ErrHandle(err)

	err = db.Ping()
	ErrHandle(err)
	return db
}
