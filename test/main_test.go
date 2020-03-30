package test

import (
    "../src"
    "fmt"
    "os"
    "testing"
    _ "github.com/lib/pq"
)

/*
    Class which starts the tests, called by go test
*/

func TestMain(m *testing.M) {
    os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
    fmt.Println("   testMain START")
    host := "localhost"
    db := src.SetupDb(host)
    defer db.Close()
    src.DeletePoliceOfficersTable(db)
    src.DeleteBikeTheftsTable(db)
    src.CreatePoliceOfficersTable(db)
    src.CreateBikeTheftsTable(db)
    provisionCommonData()
    provisionOfficerData()
    provisionBikeTheftData()
    fmt.Println("   testMain END")
    return m.Run()
}