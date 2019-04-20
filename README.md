# Go GraphQL API Boilerplate

## Stacks

- Go
- GraphQL : [graphql-go](https://github.com/graph-gophers/graphql-go)
- ORM : [gorm](https://github.com/jinzhu/gorm)

## How to Run

### Initialize DB

1. Create a database

```shell
postgres=# CREATE DATABASE go;
```

2. Create a user as owner of database

```shell
postgres=# CREATE USER go WITH ENCRYPTED PASSWORD 'go';

postgres=# ALTER DATABASE go OWNER TO go;
```

3. Grant all privileges to user for the database

```shell
postgres=# GRANT ALL PRIVILEGES ON DATABASE go TO go;
```

4. Configure the db in `db.go`

```go
// ConnectDB : connecting DB
func ConnectDB() (*DB, error) {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=go dbname=go password=go sslmode=disable")

	if err != nil {
		panic(err)
	}

	return &DB{db}, nil
}
```

### Initial Migration

```shell
$ go run ./migrations/init.go
```

This will generate the `users` table in the database as per the User Model declared in `./model/user.go` Model

### Run the server

```shell
$ go run server.go
```

### GraphQL Playground

Connect to http://localhost:8080/playground

## Usage

### Sign Up

```graphql
mutation {
  signUp(
    email: "test@test.com"
    password: "12345678"
    firstName: "graphql"
    lastName: "go"
  ) {
    ok
    error
    user {
      id
      email
      firstName
      lastName
      createdAt
      updatedAt
    }
  }
}
```

## Next to do

- [x] Sign-Up
- [ ] Query the profile with implementing `context.Context`
- [ ] Sign-In with JWT
- [ ] Change the password
- [ ] Change the profile
