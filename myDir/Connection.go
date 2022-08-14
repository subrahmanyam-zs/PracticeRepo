package myDir

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlconfig struct {
	User     string
	Host     string
	Password string
	Port     string
	Dbname   string
}

//user:password@tcp(host:port)/dbname
func Connection(mysql MySqlconfig) (*sql.DB, error) {
	constr := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v", mysql.User, mysql.Password, mysql.Host, mysql.Port, mysql.Dbname)

	db, err := sql.Open("mysql", constr)
	if err != nil {
		fmt.Println(err)
	}
	return db, err
}
