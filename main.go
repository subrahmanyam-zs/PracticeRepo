package main

import (
	"fmt"
	"net/http"
	"proHTTPDB/myDir"
)

func main() {
	var s myDir.Store
	mysql := myDir.MySqlconfig{"root", "localhost", "Jason@470", "3306", "go"}
	db, err := myDir.Connection(mysql)
	if err != nil {
		fmt.Println("error", err)
	}
	s.Db = db
	http.HandleFunc("/", s.PostHandler)
	http.ListenAndServe(":8080", nil)

}
