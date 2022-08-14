package myDir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	var s Store
	testcases := []struct {
		id             int
		expectedOutput employee
	}{
		{2, employee{2, "jason", "roy", 23, 89000}},
		{3, employee{3, "shikhar", "dhawan", 34, 96000}},
		{4, employee{4, "virat", "kohli", 33, 1259000}},
	}

	db, mock, err := sqlmock.New()
	s.Db = db
	if err != nil {
		fmt.Println(err)
	}
	for i, tc := range testcases {
		row := mock.NewRows([]string{"Id", "Fname", "Lname", "Age", "Salary"}).
			AddRow(tc.expectedOutput.Id, tc.expectedOutput.Fname, tc.expectedOutput.Lname, tc.expectedOutput.Age, tc.expectedOutput.Salary)
		mock.ExpectQuery("select (.+) from employee where Id=?").WithArgs(tc.id).WillReturnRows(row).WillReturnError(err)

		u := fmt.Sprintf("/employee?Id=%v", tc.id)
		req := httptest.NewRequest(http.MethodGet, u, nil)
		res := httptest.NewRecorder()
		s.Handler(res, req)
		temp := res.Result()

		defer temp.Body.Close()
		body, _ := io.ReadAll(temp.Body)
		var actualOutput employee
		json.Unmarshal(body, &actualOutput)

		if actualOutput != tc.expectedOutput {
			t.Errorf("Test Case Failed %v", i+1)
		}
	}
}

func TestPost(t *testing.T) {
	testcases := []struct {
		desc           string
		input          employee
		expectedOutput int
	}{
		{"valid request 1", employee{5, "joe", "root", 35, 76000}, 200},
		{"valid request 2", employee{6, "mark", "wood", 28, 56800}, 200},
	}

	var s Store
	var err error
	db, mock, err := sqlmock.New()
	s.Db = db

	for i, tc := range testcases {
		mock.ExpectExec("insert into employee values").WithArgs(tc.input.Id, tc.input.Fname, tc.input.Lname, tc.input.Age, tc.input.Salary).
			WillReturnResult(sqlmock.NewResult(1, 1)).WillReturnError(err)

		json_data, err := json.Marshal(tc.input)
		if err != nil {
			log.Fatal(err)
		}
		reader := bytes.NewReader(json_data)
		req := httptest.NewRequest(http.MethodPost, "/post", reader)
		res := httptest.NewRecorder()
		s.PostHandler(res, req)
		actualOutput := res.Result().StatusCode

		if actualOutput != tc.expectedOutput {
			t.Errorf("test case failed %v", i+1)
		}
	}
}
