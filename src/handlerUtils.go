package src

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
)

/*
    Class which implements common utility functions for API hanlders
*/

func RootHandler(w http.ResponseWriter, r *http.Request) {
    sendAPINotFound(w)
}

func sendAPINotFound(w http.ResponseWriter){
    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, "REST API Summary: /officers, /bike-thefts") 
}

func sendJsonResponse(w http.ResponseWriter, v interface{}){
    w.Header().Set("Content-Type", "application/json") 
    json.NewEncoder(w).Encode(v)    
}

func processExecResult(res sql.Result,err error) {
    if err != nil {
        panic(err)
    }
}
