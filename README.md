# go-todo
simple todo app using golang and postgresql

## How to run
1. Clone this repository
2. Run `go mod tidy`
3. Run `make run`
4. Open `http://localhost:8080` in your browser
5. Done

## Database
1. Create a database
2. set the database credentials in the `.env` file
3. Run the migration using `make migrate`
4. Done