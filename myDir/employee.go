package myDir

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Store struct {
	Db *sql.DB
}

type employee struct {
	Id     int
	Fname  string
	Lname  string
	Age    int
	Salary int
}

func (s Store) Handler(res http.ResponseWriter, req *http.Request) {
	var e employee
	q := req.URL.Query().Get("Id")
	if q == "" {
		res.WriteHeader(http.StatusBadRequest)

	} else {
		id, err := strconv.Atoi(q)
		row, err := s.Db.Query("select * from employee where Id=?;", id)
		if err != nil {
			fmt.Println(err)
			return
		}
		row.Next()
		row.Scan(&e.Id, &e.Fname, &e.Lname, &e.Age, &e.Salary)
		encrypt, _ := json.Marshal(e)

		if e.Id == 0 {
			res.WriteHeader(http.StatusNotFound)
			return
		}
		res.Write(encrypt)
	}
}

func (s Store) PostHandler(res http.ResponseWriter, req *http.Request) {

	r := req.Body

	json_data, err := io.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	var data employee
	json.Unmarshal(json_data, &data)

	if data.Id == 0 {
		res.WriteHeader(http.StatusBadRequest)
		return
	} else {
		_, err := s.Db.Exec("insert into employee values (?,?,?,?,?)", data.Id, data.Fname, data.Lname, data.Age, data.Salary)
		if err != nil {
			fmt.Println(err)
			return
		} else {
			res.Write([]byte("succces"))
		}
	}

}
