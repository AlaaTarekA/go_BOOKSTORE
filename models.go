package main

type Book struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Publisher   string `json:"publisher"`
	PublishDate Date   `json:"publishdate"`
	Rating      int    `json:"rating"`
	Status      string `json:"status"`
}

type Date struct {
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}
