# Go GraphQL API Boilerplate

## Stacks

- Go
- GraphQL : [graphql-go](https://github.com/graph-gophers/graphql-go)
- ORM : [gorm](https://github.com/jinzhu/gorm)

## Features

- User Sign Up & Sign In
- Change a Password, Profile

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

or with docker
```
$ docker run --rm go-graphql-api migrate
```

This will generate the `users` table in the database as per the User Model declared in `./model/user.go` Model

### Run the server

```shell
$ go run server.go
```
or with docker
```
$ docker build -t go-graphql-api .
$ docker run -p 8080:8080 go-graphql-api
```

### GraphQL Playground

Connect to http://localhost:8080

### Authentication : JWT

You need to set the Http request headers `Authorization`: `{JWT_token}`

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
      bio
      avatar
      createdAt
      updatedAt
    }
  }
}
```

### Sign In

```graphql
mutation {
  signIn(email: "test@test.com", password: "12345678") {
    ok
    error
    token
  }
}
```

### Change a Password

```graphql
mutation {
  changePassword(userID: 1, password: "87654321") {
    ok
    error
    user {
      id
      email
      firstName
      lastName
      bio
      avatar
      createdAt
      updatedAt
    }
  }
}
```

### Change a Profile

```graphql
mutation {
  changeProfile(userID: 1, bio: "Go developer", avatar: "go-developer.png") {
    ok
    error
    user {
      id
      email
      firstName
      lastName
      bio
      avatar
      createdAt
      updatedAt
    }
  }
}
```

### Get my profile

```graphql
query {
  getMyProfile {
    ok
    error
    user {
      id
      email
      firstName
      lastName
      bio
      avatar
      createdAt
      updatedAt
    }
  }
}
```

## Next to do

- [x] Sign-Up
- [x] Query the profile with implementing `context.Context`
- [x] Sign-In with JWT
- [x] Change the password
- [x] Change the profile
- [x] Merging \*.graphql files to a schema with `packr`
- [ ] Using Configuration file for DB & JWT secret_key
