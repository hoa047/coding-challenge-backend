package test

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "strconv"
    "testing"
)

const nbrOfOfficers = 8
var officerNonExistent = "{\"id\":"+nonExistentId+", \"name\":\"NON_EXISTENT\"}"
var officerNameChangeId = "{\"id\":6, \"name\":\"1st_CHANGE\"}"
var officerNames [nbrOfOfficers]string
var officerNameIds [nbrOfOfficers]string
var threeOfficersArray string
var twoVacantOfficers = "[4,5]"
var threeVacantOfficers = "[1,2,3]"
var officerNotFound = "officer not found"
var apiSummary = "REST API Summary: /officers, /bike-thefts"
var officersUrl = "/officers"
var officersFetchIdUrl = "/officers?id="
var officersFetchVacantUrl = "/officers?fetch=vacant"
var officerCannotBeDeleted7 = "{\"status\":\"officer cannot be deleted. Currently assigned to bikeTheft: 7\"}"
var officerCannotBeDeleted8 = "{\"status\":\"officer cannot be deleted. Currently assigned to bikeTheft: 8\"}"

/*
    Class which tests the officer API
*/

func provisionOfficerData() {
    fmt.Println("Provison officerTestData - START")
    for i := 1; i < len(officerNames)+1; i++ {
        x := strconv.Itoa(i)
        officerNames[i-1] = "{\"name\":\"nbr "+x+"\"}"
        officerNameIds[i-1] = "{\"id\":"+strconv.Itoa(i)+","+
                                "\"name\":\"nbr "+x+"\"}"
    }
    threeOfficersArray = "["+officerNameIds[0]+","+officerNameIds[1]+","+officerNameIds[2]+"]"
    fmt.Println("Provison officerTestData - FINISHED")
}

// Officer Tests
func TestGetOfficers(t *testing.T) {
    resRec := serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersUrl))
    verifyOfficerRespArray(t, resRec, emptyArray)
    for i := 0; i < 3; i++ {
        resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "POST", officersUrl, officerNames[i]))
        verifyOfficerRespObject(t, resRec, officerNameIds[i])
    }
    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersUrl))
    verifyOfficerRespArray(t, resRec, threeOfficersArray)

    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, threeVacantOfficers)
    for i := 0; i < 3; i++ {
        resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "DELETE", officersUrl, ids[i]))
        verifyOfficerRespObject(t, resRec, officerNameIds[i])
    }
    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersUrl))
    verifyOfficerRespArray(t, resRec, emptyArray)
}

func TestGetOfficerVacant(t *testing.T) {
    resRec := serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)
    for i := 3; i < 5; i++ {
        resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "POST", officersUrl, officerNames[i]))
        verifyOfficerRespObject(t, resRec, officerNameIds[i])
    }
    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, twoVacantOfficers)
    for i := 3; i < 5; i++ {
        resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "DELETE", officersUrl, ids[i]))
        verifyOfficerRespObject(t, resRec, officerNameIds[i])
    }
    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)  
}

func TestCRUDOfficer(t *testing.T) {
    resRec := serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersUrl))
    verifyOfficerRespArray(t, resRec, emptyArray)

    resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "POST", officersUrl, officerNames[5]))
    verifyOfficerRespObject(t, resRec, officerNameIds[5])

    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchIdUrl + "6"))
    verifyOfficerRespObject(t, resRec, officerNameIds[5])

    resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "PUT", officersUrl, officerNameChangeId))
    verifyOfficerRespObject(t, resRec, officerNameChangeId)

    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchIdUrl + "6"))
    verifyOfficerRespObject(t, resRec, officerNameChangeId)

    resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "DELETE", officersUrl, ids[5]))
    verifyOfficerRespObject(t, resRec, officerNameChangeId)

    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersUrl))
    verifyOfficerRespArray(t, resRec, emptyArray)
}

func TestGetNonExistentOfficerById(t *testing.T) {
    resRec := serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchIdUrl + nonExistentId))
    verifyNonExistentOfficerRespObject(t, resRec)
}

func TestPutNonExistentOfficer(t *testing.T) {
    resRec := serveHttpRecordOfficer(sendBodyAndRecord(t, "PUT", officersUrl, officerNonExistent))
    verifyNonExistentOfficerRespObject(t, resRec)
}

func TestDeleteNonExistentOfficer(t *testing.T) {
    resRec := serveHttpRecordOfficer(sendBodyAndRecord(t, "DELETE", officersUrl, officerNonExistent))
    verifyNonExistentOfficerRespObject(t, resRec)
}

func TestRootPage(t *testing.T) {
    resRec := serveHttpRecordRoot(sendNoBodyAndRecord(t, "GET", officersUrl+"nil"))
    if status := resRec.Code; status != http.StatusNotFound {
        t.Errorf("handler returned wrong status code: expected %v got %v",
            http.StatusNotFound, status)
    }

    if resRec.Body.String() != apiSummary  {
        t.Errorf("wrong summary message:\nexpected\n%v\ngot\n%v",
            apiSummary, resRec.Body.String())
    }
}

//Test officer and bikeTheft
func TestOfficerAssignedtoBikeTheft(t *testing.T) {
    //Ensure that there are no officers and bikeThefts
    resRec := serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)

    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUnassignedUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)

    //create officer
    resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "POST", officersUrl, officerNames[6]))
    verifyOfficerRespObject(t, resRec, officerNameIds[6])
    
    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, "[7]")

    //create bikeTheft, vacant officer gets assigned
    resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "POST", bikeTheftsUrl, bikeThefts[6]))
    verifyBikeTheftRespObject(t, resRec, bikeThefts[6])

    //verify that officer is assigned
    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)  

    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUnassignedUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)

    //cannot delete officer assigned to a case
    resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "DELETE", officersUrl, ids[6]))
    verifyOfficerCannotBeDeleted(t, resRec, officerCannotBeDeleted7)

    //bikeTheft solved -> officer vacant
    resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "PUT", bikeTheftsUrl, ids[6]))
    verifyBikeTheftRespObject(t, resRec, bikeTheftsSolved[6])

    //verify that officer is vacant
    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, "[7]")  

    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUnassignedUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)

    //delete officer and bikeTheft
    resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "DELETE", officersUrl, ids[6]))
    verifyOfficerRespObject(t, resRec, officerNameIds[6])

    resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "DELETE", bikeTheftsUrl, ids[6]))
    verifyBikeTheftRespObject(t, resRec, bikeTheftsSolved[6])

    //verify that officer and bikeTheft are deleted
    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)  

    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUnassignedUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)
}

//Test bikeTheft and officer
func TestBikeTheftTakenByOfficer(t *testing.T) {
    //Ensure that there are no officers and bikeThefts
    resRec := serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)

    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUnassignedUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)

    //Create bikeTheft, it becomes unassigned
    resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "POST", bikeTheftsUrl, bikeThefts[7]))
    verifyBikeTheftRespObject(t, resRec, bikeThefts[7])

    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUnassignedUrl))
    verifyIntArray(t, resRec, "[8]")

    //create officer, bikeTheft gets assigned to officer
    resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "POST", officersUrl, officerNames[7]))
    verifyOfficerRespObject(t, resRec, officerNameIds[7])
    
    //verify that officer is assigned to bikeTheft
    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)  

    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUnassignedUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)

    //cannot delete officer assigned to a case
    resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "DELETE", officersUrl, ids[7]))
    verifyOfficerCannotBeDeleted(t, resRec, officerCannotBeDeleted8)

    //bikeTheft solved -> officer vacant
    resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "PUT", bikeTheftsUrl, ids[7]))
    verifyBikeTheftRespObject(t, resRec, bikeTheftsSolved[7])

    //verify that officer is vacant
    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, "[8]")  

    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUnassignedUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)

    //delete officer and bikeTheft
    resRec = serveHttpRecordOfficer(sendBodyAndRecord(t, "DELETE", officersUrl, ids[7]))
    verifyOfficerRespObject(t, resRec, officerNameIds[7])

    resRec = serveHttpRecordBikeTheft(sendBodyAndRecord(t, "DELETE", bikeTheftsUrl, ids[7]))
    verifyBikeTheftRespObject(t, resRec, bikeTheftsSolved[7])

    //verify that officer and bikeTheft are deleted
    resRec = serveHttpRecordOfficer(sendNoBodyAndRecord(t, "GET", officersFetchVacantUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)  

    resRec = serveHttpRecordBikeTheft(sendNoBodyAndRecord(t, "GET", bikeTheftsUnassignedUrl))
    verifyIntArray(t, resRec, noFreeResourceLeft)
}

// offcersHandler_test support functions
func verifyOfficerCannotBeDeleted(t *testing.T, resRec *httptest.ResponseRecorder, body string) {
    verifyStatusCode(t, resRec)
    resRecJsonData := verifyJsonObject(t, resRec.Body.String())
    reqBodyJsonData := verifyJsonObject(t, body)
    if len(reqBodyJsonData) != len(resRecJsonData) {
        t.Errorf("Map size expected %v , got %v",
            len(reqBodyJsonData), len(resRecJsonData))
    }
    if reqBodyJsonData["status"] != resRecJsonData["status"] {
        t.Errorf("Map size expected %v , got %v",
            reqBodyJsonData["status"], resRecJsonData["status"])
    }
}

func verifyNonExistentOfficerRespObject(t *testing.T, resRec *httptest.ResponseRecorder) {
    verifyStatusCode(t, resRec)
    resRecJsonData := verifyJsonObject(t, resRec.Body.String())
    if resRecJsonData["status"] != officerNotFound {
        t.Errorf("Expected %v got %v",
            officerNotFound, resRecJsonData["status"])
    }
}

func verifyOfficerRespObject(t *testing.T, resRec *httptest.ResponseRecorder, body string) {
    verifyStatusCode(t, resRec)
    resRecJsonData := verifyJsonObject(t, resRec.Body.String())
    reqBodyJsonData := verifyJsonObject(t, body)
    compareOfficerValues(t, reqBodyJsonData, resRecJsonData)
}

func verifyOfficerRespArray(t *testing.T, resRec *httptest.ResponseRecorder, body string) {
    verifyStatusCode(t, resRec)
    resRecJsonData := verifyJsonArray(t, resRec.Body.String())
    reqBodyJsonData := verifyJsonArray(t, body)
    if len(resRecJsonData) != len(reqBodyJsonData) {
        t.Errorf("Expected %v json objects, got %v",
            len(reqBodyJsonData), len(resRecJsonData))
    }
    for i := range(reqBodyJsonData) { 
        compareOfficerValues(t, reqBodyJsonData[i], resRecJsonData[i])
    }
}

func compareOfficerValues(t *testing.T, o1 map[string]interface{}, o2 map[string]interface{}) {    
    compareValue(t, o1["name"], o2["name"])
}