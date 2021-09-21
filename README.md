# go_BOOKSTORE

# MySQL Setup
1. Create a books_database, create the tables and insert data using sql-dump.sql.
2. In main.go change the password in the connection string with your database password.


# Build Docker Image Locally

docker build . -t go-sample-app


# Run Locally

docker run -p 8000:8000 go-sample-app

# Postman Collection Documentation

https://documenter.getpostman.com/view/17574060/UUxuhpeT


# Run Tests

docker run -p 8000:8000 go-sample-app sh -c "go test -v"
