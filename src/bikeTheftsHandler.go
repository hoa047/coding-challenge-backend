package src

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strings"
    "time"
)

/*
    Class which implements the bike theft API endpoint
*/

func BikeTheftsHandler(w http.ResponseWriter, r *http.Request) {
    if r.Header.Get("Content-Type") != "" && r.Method != "GET" {
        bikeTheft := getBikeTheftFromRequest(w, r)
        switch r.Method {
        case "POST":   
            createBikeTheft(w, r, bikeTheft, db)
        case "PUT":   
            contentType := r.Header.Get("Content-Type")
            if contentType == "application/json" {
                updateBikeTheftSolved(w, r, bikeTheft, db)
            } else if strings.Contains(contentType , "multipart/form-data") {
                updateBikeTheftImage(w, r, db)
            }
        case "DELETE":    
            deleteBikeTheft(w, r, bikeTheft, db)
        default:
            sendAPINotFound(w)
        }    
    } else if r.Method == "GET" {    
        params := r.URL.Query()
        if len(params["fetch"]) > 0 && params["fetch"][0] == "unassigned"  {
            sendJsonResponse(w, 
                listUnassignedBikeTheft(db))
        } else if len(params["id"]) > 0 {
            listBikeTheftsById(w,r, params["id"][0], db)
        } else {
            listBikeThefts(w, r, db)
        }
    } else {
        sendAPINotFound(w)
    }
}

func getBikeTheftFromRequest(w http.ResponseWriter, r *http.Request) bikeTheft {
    var bikeTheft bikeTheft
    err := json.NewDecoder(r.Body).Decode(&bikeTheft)
    if err != nil {
        http.Error(w, err.Error(), 400)
    }
    return bikeTheft
}

func CreateBikeTheftsTable(db *sql.DB) {
    fmt.Println("Creating table bike_thefts - START")
    processExecResult(
        db.Exec(`
            CREATE TABLE IF NOT EXISTS bike_thefts (
            id SERIAL PRIMARY KEY,
            title TEXT,
            brand TEXT,
            city TEXT,
            description TEXT,
            reported TIMESTAMP,
            updated TIMESTAMP,
            solved BOOLEAN,
            officer_id INT NOT NULL,
            image_name TEXT,
            image BYTEA)
            `))
    fmt.Println("Creating table bike_thefts - FINISHED")
}

func DeleteBikeTheftsTable(db *sql.DB) {
    fmt.Println("Deleting table bike_thefts - START")
    processExecResult(
        db.Exec(`
            DROP TABLE IF EXISTS bike_thefts
            `))
    fmt.Println("Deleting table bike_thefts - FINISHED")
}

func createBikeTheft(w http.ResponseWriter, r *http.Request, bikeTheft bikeTheft, db *sql.DB) {
    //assign vacant officer to this case
    vacantOfficerId := -1
    if(len(listVacantOfficers(db)) > 0) {
        vacantOfficerId = listVacantOfficers(db)[0]
    } 
    sendBikeTheftJsonResponseRow(w, db.QueryRow(`
        INSERT INTO bike_thefts (title, brand, city, description, reported, updated, solved, officer_id, image_name, image)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING *;
    `,bikeTheft.Title, bikeTheft.Brand, bikeTheft.City, bikeTheft.Description, time.Now(), time.Now(), false, vacantOfficerId, 
    bikeTheft.ImageName, ""))
}


func listBikeTheftsById(w http.ResponseWriter, r *http.Request, id string, db *sql.DB) {   
    downloadImage(db, id, "/go/src/app/")
    sendBikeTheftJsonResponseRow(w, db.QueryRow(`
        SELECT * FROM bike_thefts
        WHERE id = $1;
        `, id))
}

func updateBikeTheftImage(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    id, imageName, imgBytes := getImageAsBytes(w, r, "image", "id", "imageName")
    db.QueryRow(`
        UPDATE bike_thefts
        SET image_name = $1,
        image = $2
        WHERE id = $3
        RETURNING *;
    `, imageName, imgBytes, id)
}

func updateBikeTheftSolved(w http.ResponseWriter, r *http.Request, bikeTheft bikeTheft, db *sql.DB) {
    //case solved, one officer made vacant
    vacantOfficerId := sendBikeTheftJsonResponseRow(w, db.QueryRow(`
        UPDATE bike_thefts
        SET solved = true
        WHERE id = $1
        RETURNING *;
    `, bikeTheft.Id))
    //assign vacant officer
    unassigned := listUnassignedBikeTheft(db)
    if(len(unassigned) > 0) {
        db.Exec(`
            UPDATE bike_thefts 
            SET officer_id = $2
            WHERE id = $1;
        `, unassigned[0], vacantOfficerId)
    }  
}

func deleteBikeTheft(w http.ResponseWriter, r *http.Request, bikeTheft bikeTheft, db *sql.DB) {
    sendBikeTheftJsonResponseRow(w, db.QueryRow(`
        DELETE FROM bike_thefts
        WHERE id = $1
        RETURNING *;`, 
    bikeTheft.Id))        
}

func listBikeThefts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    rows, err := db.Query(`SELECT * FROM bike_thefts;`)
    sendBikeTheftJsonResponseRows(w, rows, err)
}

func sendBikeTheftJsonResponseRow(w http.ResponseWriter, row *sql.Row) int {
    var id          int         
    var title       string     
    var brand       string      
    var city        string       
    var description string
    var reported    time.Time
    var updated     time.Time
    var solved      bool     
    var officerId   int    
    var imageName   string
    var image       []byte
    switch err := row.Scan(&id, &title, &brand, &city, &description, &reported, &updated, &solved, &officerId, &imageName, &image); err {
    case sql.ErrNoRows:
        sendJsonResponse(w, 
            map[string]string{"status":"bike theft not found"})
        return -1    
    case nil:
        sendJsonResponse(w, 
            bikeTheft {
                Id:             id, 
                Title:          title,
                Brand:          brand,
                City:           city,       
                Description:    description,
                Reported:       reported,    
                Updated:        updated,      
                Solved:         solved,    
                OfficerId:      officerId,  
                ImageName:      imageName,
            } )
        return officerId          
    default:
        panic(err)
    }
}

func sendBikeTheftJsonResponseRows(w http.ResponseWriter, rows *sql.Rows,err error) {
    var id          int         
    var title       string     
    var brand       string      
    var city        string       
    var description string
    var reported    time.Time
    var updated     time.Time
    var solved      bool     
    var officerId   int    
    var imageName  string
    var image       []byte

    bikeThefts :=[]bikeTheft{}
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    for rows.Next() {
        if err := rows.Scan(&id, &title, &brand, &city, &description, &reported, &updated, &solved, &officerId, &imageName, &image); err != nil {
                log.Fatal(err)
        }
        bikeThefts = append(bikeThefts,  
            bikeTheft {
                Id:             id, 
                Title:          title,
                Brand:          brand,
                City:           city,       
                Description:    description,
                Reported:       reported,    
                Updated:        updated,      
                Solved:         solved,    
                OfficerId:      officerId,  
                ImageName:      imageName,
            } )
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
    sendJsonResponse(w, bikeThefts)
}

func listVacantOfficers( db *sql.DB) []int {
    return getIdFromRows(db.Query(`
        SELECT id 
        FROM police_officers p
        WHERE p.id not in(
        SELECT b.officer_id FROM bike_thefts b WHERE b.solved=false);
        `))
}

func getIdFromRows(rows *sql.Rows, err error) []int {
    id      := 0
    ids :=[]int{}
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    for rows.Next() {
        if err := rows.Scan(&id); err != nil {
                log.Fatal(err)
        }
        ids = append(ids, id)
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
    return ids
}