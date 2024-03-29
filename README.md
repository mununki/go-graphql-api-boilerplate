# Go GraphQL API Boilerplate

## Stacks

- Go
- GraphQL : [graphql-go](https://github.com/graph-gophers/graphql-go)
- Querybuilder : [goqu](https://github.com/doug-martin/goqu)

## Features

- User Sign Up & Sign In with OAuth (Google, Kakao)
- Change the Profile

## How to Run

### Initialize DB

1. Create a database using `sql/create.sql` (MYSQL)

2. Configure the db in `db/db.go`

```go
// ConnectDB : connecting DB
func ConnectDB() (*DB, error) {
	db, err := sql.Open("mysql", "api:your_password$@/database_name?parseTime=true")

	if err != nil {
		panic(err)
	}

	// https://github.com/go-sql-driver/mysql/#important-settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10)

	errPing := db.Ping()
	if errPing != nil {
		panic(err.Error())
	}

	qb := goqu.New("mysql", db)

	return &DB{qb}, nil
}
```

### Make `.env` file
```env
STAGE=PROD // PROD for production or else for debug
JWT_SECRET=my_secret
GOOGLE_CLIENT_ID=your_google_web_client_id
KAKAO_REST_API_KEY=your_kakao_rest_api_key
KAKAO_REDIRECT_URI=http://localhost:8280/oauth/kakao/redirect
```

### Run the server

```shell
$ go run server.go
```

### GraphQL Playground

Connect to http://localhost:8080

### Authentication : JWT

You need to set the Http request headers `Authorization`: `Bearer {JWT_token}`

## Usage

### Sign In

```graphql
mutation {
  signInGoogle(idToken: "12345678") {
    ok
    error
    token
  }
}
```

```graphql
mutation {
  signInKakao(code: "12345678") {
    ok
    error
    token
  }
}
```

### Change a Profile

```graphql
mutation {
  changeProfile(nickname: "Go developer") {
    ok
    error
    user {
      id
      email
      nickname
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
      nickname
      createdAt
      updatedAt
    }
  }
}
```
