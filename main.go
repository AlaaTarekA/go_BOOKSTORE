package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var connectionString string = "root:myDatabase_123@tcp(docker.for.mac.localhost:3306)/books_database?charset=utf8&parseTime=True&loc=Local"

// a function that returns all the books in the database
func getAllBooks(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connectionString)
	//defer db.Close()

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		return
	}

	var entries []Book

	rows, err := db.Query("SELECT * from book_data;")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	}

	defer rows.Close()
	for rows.Next() {

		var book Book

		var title sql.NullString
		var author sql.NullString
		var publisher sql.NullString
		var publishDate sql.NullTime
		var rating int
		var status sql.NullString

		rows.Scan(&title, &author, &publisher, &publishDate, &rating, &status)
		book.Title = title.String
		book.Author = author.String
		book.Publisher = publisher.String

		book.PublishDate = Date{Day: publishDate.Time.Day(), Month: int(publishDate.Time.Month()), Year: publishDate.Time.Year()}
		book.Rating = rating
		book.Status = status.String
		entries = append(entries, book)
	}
	respondWithJSON(w, http.StatusOK, entries)

}

// a function that returns a certain book by its title

func getBookByTitle(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		return
	}
	defer db.Close()
	id := r.URL.Query().Get("title")
	if id == "" {
		params := mux.Vars(r)
		id = params["title"]
	}
	var title sql.NullString
	var author sql.NullString
	var publisher sql.NullString
	var publishDate sql.NullTime
	var rating int
	var status sql.NullString
	err = db.QueryRow("SELECT * from book_data where title= ?", id).Scan(&title, &author, &publisher, &publishDate, &rating, &status)

	switch {
	case err == sql.ErrNoRows:
		respondWithError(w, http.StatusBadRequest, "No entry found with the title="+id)
		return
	case err != nil:

		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return

	default:

		var book Book
		book.Title = title.String
		book.Author = author.String
		book.Publisher = publisher.String

		day := publishDate.Time.Day()
		month := int(publishDate.Time.Month())
		year := publishDate.Time.Year()
		book.PublishDate = Date{Day: day, Month: month, Year: year}
		book.Rating = rating
		book.Status = status.String
		respondWithJSON(w, http.StatusOK, book)
	}

}

// a function that checks if a specific book already exits in the database

func rowExists(t string) bool {
	db, err := sql.Open("mysql", connectionString)
	var title sql.NullString
	var author string
	var publisher string
	var publishDate string
	var rating int
	var status string
	err = db.QueryRow("SELECT Title, Author, Publisher, PublishDate, Rating, Status from book_data where title= ?", t).Scan(&title, &author, &publisher, &publishDate, &rating, &status)

	return err != sql.ErrNoRows

}

// a function that creates a new book in the database
func createBook(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connectionString)
	//defer db.Close()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		return
	}
	decoder := json.NewDecoder(r.Body)
	var entry Book
	decoder.Decode(&entry)
	if rowExists(entry.Title) {
		respondWithError(w, http.StatusOK, "Already exists.")
		return
	}
	if entry.Status != "checked" && entry.Status != "unchecked" {
		respondWithError(w, http.StatusInternalServerError, "There was problem entering the entry.")
		return
	}
	if entry.Rating < 1 || entry.Rating > 3 {
		respondWithError(w, http.StatusInternalServerError, "There was problem entering the entry.")
		return
	}

	var day int = entry.PublishDate.Day
	var month int = entry.PublishDate.Month
	var year int = entry.PublishDate.Year

	stringDate := strconv.Itoa(day) + "-" + strconv.Itoa(month) + "-" + strconv.Itoa(year)
	re := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])-(0?[1-9]|1[012])-((19|20)\\d\\d)")
	match := re.MatchString(stringDate)
	if !match {
		respondWithError(w, http.StatusInternalServerError, "There was problem entering the entry.")
		return
	}
	dateRes := strconv.Itoa(year) + "-" + strconv.Itoa(month) + "-" + strconv.Itoa(day)

	statement, err := db.Prepare("insert into book_data (Title, Author, Publisher, PublishDate, Rating, Status) values(?,?,?,?,?,?)")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	}

	res, err := statement.Exec(entry.Title, entry.Author, entry.Publisher, dateRes, entry.Rating, entry.Status)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was problem entering the entry.")
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {

		respondWithJSON(w, http.StatusOK, entry)
		return

	}

}

// a function  deletes a book from the database

func deleteBook(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		return
	}
	defer db.Close()

	id := r.URL.Query().Get("title")
	if id == "" {
		params := mux.Vars(r)
		id = params["title"]
	}
	var title sql.NullString
	var author sql.NullString
	var publisher sql.NullString
	var publishDate sql.NullTime
	var rating int
	var status sql.NullString
	err = db.QueryRow("SELECT title, author, publisher, publishDate, rating, status from book_data where title=?", id).Scan(&title, &author, &publisher, &publishDate, &rating, &status)
	switch {
	case err == sql.ErrNoRows:
		respondWithError(w, http.StatusBadRequest, "No entry found with the title="+id)
		return
	case err != nil:
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	default:

		res, err := db.Exec("DELETE from book_data where title=?", id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
			return
		}
		if count == 1 {
			var book Book
			book.Title = title.String
			book.Author = author.String
			book.Publisher = publisher.String

			day := publishDate.Time.Day()
			month := int(publishDate.Time.Month())
			year := publishDate.Time.Year()
			book.PublishDate = Date{Day: day, Month: month, Year: year}
			book.Rating = rating
			book.Status = status.String

			respondWithJSON(w, http.StatusOK, book)
			return
		}

	}
}

// a function updates a certain book either updating the status from checked to unchecked or vice versa or the rating in the rangge form 1 to 3

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		return
	}
	defer db.Close()
	decoder := json.NewDecoder(r.Body)
	var input string
	err = decoder.Decode(&input)
	s := strings.Split(input, ":")
	id := r.URL.Query().Get("title")
	if id == "" {
		params := mux.Vars(r)
		id = params["title"]
	}
	if s[0] == "status" && (s[1] != "checked" && s[1] != "unchecked") {
		respondWithError(w, http.StatusInternalServerError, "There was problem entering the entry.")
		return
	}
	if s[0] == "rating" {
		i, _ := strconv.Atoi(s[1])
		if i < 1 || i > 3 {
			respondWithError(w, http.StatusInternalServerError, "There was problem entering the entry.")
			return
		}
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		return
	}
	if s[0] == "status" {

		statement, err := db.Prepare("update book_data set status=? where title=?")

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
			return
		}

		res, err := statement.Exec(s[1], id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "There was problem entering the entry.")
			return
		}
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
			respondWithJSON(w, http.StatusOK, "Status Updated!")
		}
		defer statement.Close()
	}
	if s[0] == "rating" {
		i, _ := strconv.Atoi(s[1])
		statement, err := db.Prepare("update book_data set rating=? where title=?")

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
			return
		}

		res, err := statement.Exec(i, id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "There was problem entering the entry.")
			return
		}
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
			respondWithJSON(w, http.StatusOK, "Rating Updated!")
		}
		defer statement.Close()
	}

}

func main() {
	r := mux.NewRouter().StrictSlash(true)
	router := r.PathPrefix("/api").Subrouter() // /api will give access to all the API endpoints
	router.HandleFunc("/books/", getAllBooks).Methods("GET")
	router.HandleFunc("/books/{title}", getBookByTitle).Methods("GET")
	router.HandleFunc("/books/", createBook).Methods("POST")
	router.HandleFunc("/books/{title}", updateBook).Methods("POST")
	router.HandleFunc("/books/{title}", deleteBook).Methods("DELETE")
	http.Handle("/", r)
	fmt.Println("Listening on port :8000")
	http.ListenAndServe(":8000", r)
}

// a function that responds with error if an error exists
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// a function Called for responses to encode and send json data
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	//encode payload to json
	response, _ := json.Marshal(payload)

	// set headers and write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
