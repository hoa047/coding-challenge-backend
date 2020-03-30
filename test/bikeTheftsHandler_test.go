package test

import (
    "fmt"
    "net/http/httptest"
    "strconv"
    "testing"
)

const nbrOfbikeThefts = 8
var bikeTheftNonExistent = "{\"id\":"+nonExistentId+"}"
var bikeThefts [nbrOfbikeThefts]string
var bikeTheftsSolved [nbrOfbikeThefts]string
var threeBikeTheftsArray string
var threeBikeTheftsArraySolved string
var twoUnassignedBikeThefts = "[5,6]"
var bikeTheftNotFound = "bike theft not found"
var bikeTheftsUrl = "/bike-thefts"
var bikeTheftSolvedUrl = "/bike-theft-solved"
var bikeTheftsFetchIdUrl = "/bike-thefts?id="
var bikeTheftsUnassignedUrl = "/bike-thefts?fetch=unassigned"

/*
    Class which tests the bike theft API
*/

func provisionBikeTheftData() {
    fmt.Println("Provision bikeTheftTestData - START")
    for i := 1; i < len(bikeThefts)+1; i++ {
        x := strconv.Itoa(i)
        bikeThefts[i-1] = "{\"title\":\"TLE"+x+"\"," + 
                            "\"brand\":\"BRD"+x+"\","+
                            "\"city\":\"CTY"+x+"\","+
                            "\"description\":\"DSC"+x+"\","+
                            "\"solved\":false,"+
                            "\"image_name\":\"IMG_N"+x+"\","+
                            "\"image\":\"IMG"+x+"\"}"
        bikeTheftsSolved[i-1] = "{\"title\":\"TLE"+x+"\"," + 
                            "\"brand\":\"BRD"+x+"\","+
                            "\"city\":\"CTY"+x+"\","+
                            "\"description\":\"DSC"+x+"\","+
                            "\"solved\":true,"+
                            "\"image_name\":\"IMG_N"+x+"\","+
                            "\"image\":\"IMG"+x+"\"}"
    }
    threeBikeTheftsArray = "["+bikeThefts[0]+","+bikeThefts[1]+","+bikeThefts[2]+"]"
    threeBikeTheftsArraySolved = "["+bikeTheftsSolved[0]+","+bikeTheftsSolved[1]+","+bikeTheftsSolved[2]+"]"
    fmt.Println("Provision bikeTheftTestData - FINISHED")
}

// BikeTheft Tests
func TestCRUDBikeThefts(t *testing.T) {
    resRec := serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUrl))
    verifyBikeTheftRespArray(t, resRec, emptyArray)
    for i := 0; i < 3; i++ {
        resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "POST", bikeTheftsUrl, bikeThefts[i]))
        verifyBikeTheftRespObject(t, resRec, bikeThefts[i])
    }
    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUrl))
    verifyBikeTheftRespArray(t, resRec, threeBikeTheftsArray)
    for i := 0; i < 3; i++ {
        resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "PUT", bikeTheftSolvedUrl, ids[i]))
        verifyBikeTheftRespObject(t, resRec, bikeTheftsSolved[i])
    }
    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUrl))
    verifyBikeTheftRespArray(t, resRec, threeBikeTheftsArraySolved)
    for i := 0; i < 3; i++ {
        resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "DELETE", bikeTheftsUrl, ids[i]))
        verifyBikeTheftRespObject(t, resRec, bikeTheftsSolved[i])
    }
    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUrl))
    verifyBikeTheftRespArray(t, resRec, emptyArray)
}

func TestGetBikeTheftById(t *testing.T) {
    resRec := serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsFetchIdUrl + "4"))
    verifyNonExistentBikeTheftRespObject(t, resRec)

    resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "POST", bikeTheftsUrl, bikeThefts[3]))
    verifyBikeTheftRespObject(t, resRec, bikeThefts[3])

    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsFetchIdUrl + "4"))
    verifyBikeTheftRespObject(t, resRec, bikeThefts[3])

    resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "DELETE", bikeTheftsUrl, ids[3]))
    verifyBikeTheftRespObject(t, resRec, bikeThefts[3])

    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsFetchIdUrl + "4"))
    verifyNonExistentBikeTheftRespObject(t, resRec)
}

func TestGetBikeTheftsUnassigned(t *testing.T) {
    resRec := serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUnassignedUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)
    for i := 4; i < 6; i++ {
        resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "POST", bikeTheftsUrl, bikeThefts[i]))
        verifyBikeTheftRespObject(t, resRec, bikeThefts[i])
    }
    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUnassignedUrl))
    verifyIntArray(t, resRec, twoUnassignedBikeThefts)
    for i := 4; i < 6; i++ {
        resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "DELETE", bikeTheftsUrl, ids[i]))
        verifyBikeTheftRespObject(t, resRec, bikeThefts[i])
    }
    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUnassignedUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)
}

func TestGetNonExistentBikeTheftById(t *testing.T) {
    resRec := serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsFetchIdUrl + nonExistentId))
    verifyNonExistentBikeTheftRespObject(t, resRec)
}

func TestPutNonExistentBikeTheftById(t *testing.T) {
    resRec := serveHttpRecordBikeTheft(sendBodyAndRecord(t, "PUT", bikeTheftsUrl, bikeTheftNonExistent))
    verifyNonExistentBikeTheftRespObject(t, resRec)
}

func TestDeleteNonExistentBikeTheftById(t *testing.T) {
    resRec := serveHttpRecordBikeTheft(sendBodyAndRecord(t, "DELETE", bikeTheftsUrl, bikeTheftNonExistent))
    verifyNonExistentBikeTheftRespObject(t, resRec)
}

// bikeTheftsHandler_test support functions
func verifyNonExistentBikeTheftRespObject(t *testing.T, resRec *httptest.ResponseRecorder) {
    verifyStatusCode(t, resRec)
    resRecJsonData := verifyJsonObject(t, resRec.Body.String())
    if resRecJsonData["status"] != bikeTheftNotFound {
        t.Errorf("Expected %v got %v",
            bikeTheftNotFound, resRecJsonData["status"])
    }
}

func verifyBikeTheftRespObject(t *testing.T, resRec *httptest.ResponseRecorder, body string) {
    verifyStatusCode(t, resRec)
    resRecJsonData := verifyJsonObject(t, resRec.Body.String())
    reqBodyJsonData := verifyJsonObject(t, body)
    compareBikeTheftValues(t, reqBodyJsonData, resRecJsonData)
}


func verifyBikeTheftRespArray(t *testing.T, resRec *httptest.ResponseRecorder, body string) {
    verifyStatusCode(t, resRec)
    resRecJsonData := verifyJsonArray(t, resRec.Body.String())
    reqBodyJsonData := verifyJsonArray(t, body)
    if len(resRecJsonData) != len(reqBodyJsonData) {
        t.Errorf("Expected %v json objects, got %v",
            len(reqBodyJsonData), len(resRecJsonData))
    }
    for i := range(reqBodyJsonData) { 
        compareBikeTheftValues(t, reqBodyJsonData[i], resRecJsonData[i])
    }
}

func compareBikeTheftValues(t *testing.T, o1 map[string]interface{}, o2 map[string]interface{}) {
    compareValue(t, o1["title"], o2["title"])
    compareValue(t, o1["brand"], o2["brand"])
    compareValue(t, o1["city"], o2["city"])
    compareValue(t, o1["description"], o2["description"])
    compareValue(t, o1["solved"], o2["solved"])
    compareValue(t, o1["image_name"], o2["image_name"])
}