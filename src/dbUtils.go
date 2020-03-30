package src

import (
    "database/sql"
    "fmt"
)

var db *sql.DB

const (
  dbname   = "dev"
  user     = "user"
  password = "password123"
  port     = 5432
  sslmode  = "disable"
)

/*
    Class which implements database functionality
*/

func SetupDb(host string) *sql.DB {
    fmt.Println("Opening postgres connection - START")
    dbConnectInfo := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
        user, password, host, port, dbname, sslmode)
    var err error
    db, err = sql.Open("postgres", dbConnectInfo)
    if err != nil {
        panic(err)
    }
    err = db.Ping()
    if err != nil {
        panic(err)
    }
    fmt.Println("Opening postgres connection - FINISHED")
    return db
}
