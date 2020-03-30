package src

import (	    
	"bytes"
    "database/sql"
    "fmt"
    "io"
	"io/ioutil"
	"net/http"
    "strconv"
)

/*
    Class which implements general functionality used in other classes 
*/

func listUnassignedBikeTheft( db *sql.DB) []int {
    return getIdFromRows(db.Query(`
        SELECT id 
        FROM bike_thefts b
        WHERE b.solved=false and b.officer_id=-1;
        `))
}

func downloadImage(db *sql.DB, id string, path string) bool {
    dbFileName, imgBytes := getImageFromDb(db, id)
    if imgBytes != nil {
    	err := ioutil.WriteFile(path + id + "-" + dbFileName, imgBytes, 0644)
    	if err != nil {
        	return false
    	}
	}
    return true
}

func getImageFromDb (db *sql.DB, id string) (string, []byte) {
    row := db.QueryRow(`
        SELECT IMAGE_NAME, IMAGE FROM bike_thefts
        WHERE id = $1;
        `, id)
    var imageName string
    var image []byte
    switch err := row.Scan(&imageName, &image); err {
    case sql.ErrNoRows:
        return "", nil
    case nil:
    default:
        panic(err)
    }
    return imageName, image
}

func getImageAsBytes(w http.ResponseWriter, r *http.Request, imageKey string, idKey string, imageNameKey string) (int, string, []byte) {
	err := r.ParseMultipartForm(64 << 20) // maxMemory 64MB
    if err != nil {
        fmt.Fprintf(w, "No image specified in form-data")
        return -1, "", nil
    }
    
    var id int
    var imageName string
    if len(r.MultipartForm.File[imageKey]) == 0 || 
        len(r.MultipartForm.Value[idKey]) == 0 || 
        len(r.MultipartForm.Value[imageNameKey]) == 0 {
        fmt.Fprintf(w, "image, bikeTheft id and imageName must be specified")
        return -1, "", nil   
    }
    id, err = strconv.Atoi(r.MultipartForm.Value["id"][0])
    if err != nil {
        fmt.Fprintf(w, "id must be an integer")
        return -1 ,"", nil
    }

    imageName = r.MultipartForm.Value["imageName"][0]
    fileHeader := r.MultipartForm.File["image"][0]
    
    img, err := fileHeader.Open() 
    defer img.Close()        
    if err != nil {
        fmt.Fprintf(w, "Could not open image: " + imageName)
        return -1 ,"", nil
    } 
    //multipart.File to bytes 
    buf := bytes.NewBuffer(nil)
    if _, err := io.Copy(buf, img); err != nil {
        fmt.Fprintf(w, "Could not convert following image to bytes: " + imageName)
        return -1, "", nil
    }
    return id, imageName, buf.Bytes()
}
