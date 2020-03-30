package src

import (
    "database/sql"
    "fmt"
    "encoding/json"
    "log"
    "net/http"
    "strconv"
)

/*
    Class which implements the officer API endpoint
*/

func OfficersHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Content-Type") != "" && r.Method != "GET" {
        if r.Header.Get("Content-Type") != "application/json" {
            http.Error(w, 
                "Content-Type header is not application/json, cannot handle request", 
                    http.StatusUnsupportedMediaType)
            return
        }
        officer := getPoliceOfficerFromRequest(w, r)
        switch r.Method {
        case "POST":    
            createPoliceOfficer(w, r, officer, db)
        case "PUT":    
            updatePoliceOfficer(w, r, officer, db)
        case "DELETE":    
            deletePoliceOfficer(w, r, officer, db)
        default:
            sendAPINotFound(w)
        }
    } else if r.Method == "GET" {     
        params := r.URL.Query()
        if len(params["fetch"]) > 0 && params["fetch"][0] == "vacant"  {
            sendJsonResponse(w, 
                listVacantOfficers(db))
        } else if len(params["id"]) > 0 {
            listPoliceOfficerById(w,r, params["id"][0], db)
        } else {
            listPoliceOfficers(w, r, db)
        }
    } else {
        sendAPINotFound(w)
    }
}

func getPoliceOfficerFromRequest(w http.ResponseWriter, r *http.Request) officer {
    var officer officer
    err := json.NewDecoder(r.Body).Decode(&officer)
    if err != nil {
        http.Error(w, err.Error(), 400)
    }
    return officer
}

func CreatePoliceOfficersTable(db *sql.DB) {
    fmt.Println("Creating table police_officers - START")
    processExecResult(
        db.Exec(`
            CREATE TABLE IF NOT EXISTS police_officers (
            id SERIAL PRIMARY KEY,
            name TEXT)
            `))
    fmt.Println("Creating table police_officers - FINISHED")    
}

func DeletePoliceOfficersTable(db *sql.DB) {
    fmt.Println("Deleting table police_officers - START")
    processExecResult(
        db.Exec(`
            DROP TABLE IF EXISTS police_officers
            `))
    fmt.Println("Deleting table police_officers - FINISHED")    
}

func createPoliceOfficer(w http.ResponseWriter, r *http.Request, officer officer, db *sql.DB) {
    newOfficerId := sendOfficerJsonResponseRow(w, db.QueryRow(`
        INSERT INTO police_officers (name)
        VALUES ($1)
        RETURNING *;
    `, officer.Name))
    //Assign case to officer
    unassigned := listUnassignedBikeTheft(db)
    if(len(unassigned) > 0) {
        db.Exec(`
            UPDATE bike_thefts 
            SET officer_id = $2
            WHERE id = $1;
        `, unassigned[0], newOfficerId)
    }
}

func updatePoliceOfficer(w http.ResponseWriter, r *http.Request, officer officer, db *sql.DB) {
    sendOfficerJsonResponseRow(w, db.QueryRow(`
        UPDATE police_officers
        SET name = $2
        WHERE id = $1
        RETURNING *;
    `, officer.Id, officer.Name))
}

func deletePoliceOfficer(w http.ResponseWriter, r *http.Request, officer officer, db *sql.DB) {
    var id int
    row := db.QueryRow(`
       SELECT id FROM bike_thefts
       WHERE officer_id = $1 
       AND solved = false;
    `, officer.Id)
    switch err := row.Scan(&id); err {
    case sql.ErrNoRows:
        sendOfficerJsonResponseRow(w, db.QueryRow(`
            DELETE FROM police_officers
            WHERE id = $1
            RETURNING *;
        `, officer.Id))                        
    case nil:
        sendJsonResponse(w, 
            map[string]string{"status":"officer cannot be deleted. Currently assigned to bikeTheft: " +  strconv.Itoa(id)})
            //row != 0 i.e officer is already assigned    
    default:
        panic(err)
    }
}

func listPoliceOfficerById(w http.ResponseWriter, r *http.Request, id string, db *sql.DB) {
    sendOfficerJsonResponseRow(w, db.QueryRow(`
        SELECT * FROM police_officers
        WHERE id = $1;
    `, id))
}

func listPoliceOfficers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    rows, err := db.Query(`
        SELECT * FROM police_officers;
    `)
    sendOfficerJsonResponseRows(w, rows, err)
}

func sendOfficerJsonResponseRow(w http.ResponseWriter, row *sql.Row) int{
    var id int
    var name string
    switch err := row.Scan(&id, &name); err {
    case sql.ErrNoRows:
        sendJsonResponse(w, 
            map[string]string{"status":"officer not found"})
        return -1    
    case nil:
        sendJsonResponse(w, 
            officer {
                Id: id, 
                Name: name, 
            } )
        return id        
    default:
        panic(err)
    }
}

func sendOfficerJsonResponseRows(w http.ResponseWriter, rows *sql.Rows, err error) {
    var id int
    var name string
    officers :=[]officer{}
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    for rows.Next() {
        if err := rows.Scan(&id, &name); err != nil {
                log.Fatal(err)
        }
        officers = append(officers,  
            officer {
                Id:     id, 
                Name:   name, 
            } )
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
    sendJsonResponse(w, officers)
}