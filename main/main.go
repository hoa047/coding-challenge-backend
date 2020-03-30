package main

import (
    "src"    
    "fmt"
    "log"
    "net/http"
    "os"
    _ "github.com/lib/pq"
)

/*
    Main class called by docker-compose
*/

func main() {
    restPort := os.Args[1]
    host := os.Args[2]
    db := src.SetupDb(host)
    defer db.Close()
    src.CreatePoliceOfficersTable(db)
    src.CreateBikeTheftsTable(db)
    http.HandleFunc("/bike-thefts", src.BikeTheftsHandler)
    http.HandleFunc("/officers", src.OfficersHandler)
    http.HandleFunc("/", src.RootHandler)

    fmt.Println("Listening on :",restPort)
    log.Fatal(http.ListenAndServe(":"+restPort, nil))
}