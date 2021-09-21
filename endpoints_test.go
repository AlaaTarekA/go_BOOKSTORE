package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// testscreate a new book
func TestCreateBook(t *testing.T) {

	var jsonStr = []byte(`{"title":"11","author":"xyz","publisher":"pqr","publishdate":{"day":1,"month":1,"year":2020},"rating":1,"status":"unchecked"}`)
	req, err := http.NewRequest("POST", "http://localhost:8000/api/books", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"title":"11","author":"xyz","publisher":"pqr","publishdate":{"day":1,"month":1,"year":2020},"rating":1,"status":"unchecked"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//tests creating an book whish already exists, so the output should be the book alraedy exits
func TestCreateExistingBook(t *testing.T) {

	var jsonStr = []byte(`{"title":"11","author":"xyz","publisher":"pqr","publishdate":{"day":1,"month":1,"year":2020},"rating":1,"status":"unchecked"}`)
	req, err := http.NewRequest("POST", "http://localhost:8000/api/books", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"error":"Already exists."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//tests updating the status of the book from checked to unchecked or vice versa
func TestUpdateBookStatus(t *testing.T) {

	var jsonStr = []byte(`"status:checked"`)
	req, err := http.NewRequest("POST", "http://localhost:8000/api/books?title=11", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `"Status Updated!"`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//tests updating book rate within the range from 1 to 3
func TestUpdateBookRating(t *testing.T) {

	var jsonStr = []byte(`"rating:2"`)
	req, err := http.NewRequest("POST", "http://localhost:8000/api/books?title=11", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `"Rating Updated!"`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//tests get All the books from the database table
func TestGetAllBooks(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/api/books", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllBooks)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `[{"title":"ketab","author":"roba","publisher":"solly","publishdate":{"day":7,"month":10,"year":2019},"rating":2,"status":"unchecked"},{"title":"abc","author":"mohannad","publisher":"pooBoo","publishdate":{"day":1,"month":11,"year":2019},"rating":1,"status":"checked"},{"title":"11","author":"xyz","publisher":"pqr","publishdate":{"day":1,"month":1,"year":2020},"rating":2,"status":"checked"}]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

//tests getting a certain book by its title
func TestGetBookByTitle(t *testing.T) {

	req, err := http.NewRequest("GET", "http://localhost:8000/api/books?title=ketab", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getBookByTitle)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"title":"ketab","author":"roba","publisher":"solly","publishdate":{"day":7,"month":10,"year":2019},"rating":2,"status":"unchecked"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//tests that a book doesnt exit in the table
func TestGetBookNotFound(t *testing.T) {
	req, err := http.NewRequest("GET", "http://localhost:8000/api/books?title=nono", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getBookByTitle)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status == http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

//tests delete a book by its title
func TestDeleteBook(t *testing.T) {
	req, err := http.NewRequest("DELETE", "http://localhost:8000/api/books?title=11", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := `{"title":"11","author":"xyz","publisher":"pqr","publishdate":{"day":1,"month":1,"year":2020},"rating":2,"status":"checked"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

//tests invalid date input when creating a new book
func TestInvalidDate(t *testing.T) {
	var jsonStr = []byte(`{"title":"newww","author":"xyz","publisher":"pqr","publishdate":{"day":0,"month":0,"year":0},"rating":1,"status":"unchecked"}`)
	req, err := http.NewRequest("POST", "http://localhost:8000/api/books", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf(" status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	expected := `{"error":"There was problem entering the entry."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

//tests invalid status input when creating a new book

func TestInvalidStatus(t *testing.T) {
	var jsonStr = []byte(`{"title":"newww","author":"xyz","publisher":"pqr","publishdate":{"day":1,"month":1,"year":2000},"rating":1,"status":"none"}`)
	req, err := http.NewRequest("POST", "http://localhost:8000/api/books", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf(" status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	expected := `{"error":"There was problem entering the entry."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

//tests invalid rating input when creating a new book

func TestInvalidRating(t *testing.T) {
	var jsonStr = []byte(`{"title":"newww","author":"xyz","publisher":"pqr","publishdate":{"day":1,"month":1,"year":2000},"rating":5,"status":"checked"}`)
	req, err := http.NewRequest("POST", "http://localhost:8000/api/books", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf(" status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	expected := `{"error":"There was problem entering the entry."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

//tests invalid status input when updating a new book

func TestInvalidUpdateStatus(t *testing.T) {

	var jsonStr = []byte(`"status:none"`)
	req, err := http.NewRequest("POST", "http://localhost:8000/api/books?title=11", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf(" status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	expected := `{"error":"There was problem entering the entry."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

//tests invalid rating input when updating a new book

func TestInvalidUpdateRating(t *testing.T) {

	var jsonStr = []byte(`"rating:none"`)
	req, err := http.NewRequest("POST", "http://localhost:8000/api/books?title=11", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(updateBook)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf(" status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	expected := `{"error":"There was problem entering the entry."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}
