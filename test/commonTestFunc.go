package test

import (
    "../src"
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "strconv"
    "testing"
)

var emptyArray = "[]"
var noFreeResourceLeft = "[]"
var ids [8]string
var nonExistentId = "8888"

/*
    Class which implements test functionality used by other classes
*/

func provisionCommonData() {
    fmt.Println("Provison commonTestData - START")
    for i := 1; i < len(ids)+1; i++ {
        ids[i-1] = "{\"id\":"+strconv.Itoa(i)+"}"
    }
    fmt.Println("Provison commonTestData - START")
}

func compareValue(t *testing.T, v1 interface {}, v2 interface {}) {
    if v1 != v2 {
        t.Errorf("compareValue: Expected %v got %v",
            v1, v2)
    }
}

func verifyIntArray(t *testing.T, resRec *httptest.ResponseRecorder, body string) {
    verifyStatusCode(t, resRec)
    var resRecJsonData []interface{}
    resRecJson := json.Unmarshal([]byte(resRec.Body.String()), &resRecJsonData)
    if resRecJson != nil  {
        t.Errorf("Expected json int array. Got %v",
            resRec.Body.String())
    }

    var reqBodyJsonData []interface{}
    reqBodyJson := json.Unmarshal([]byte(body), &reqBodyJsonData)
    if reqBodyJson != nil {
        t.Errorf("Could not parse input json. Got %v",
            body)
    }
    if len(resRecJsonData) != len(reqBodyJsonData) {
        t.Errorf("Expected %v json array elements, got %v",
            len(reqBodyJsonData), len(resRecJsonData))

    }

    for i := range(reqBodyJsonData) { 
        if reqBodyJsonData[i] != resRecJsonData[i]{
            t.Errorf("int array diff at [%v]: expected %v got %v",
                i, reqBodyJsonData[i], resRecJsonData[i])
        }
    }
}

func verifyJsonObject(t *testing.T, jsonBody string) map[string]interface{} {
    var jsonObjectData map[string]interface{}
    jsonObject := json.Unmarshal([]byte(jsonBody), &jsonObjectData)
    if jsonObject != nil  {
        t.Errorf("Expected json object. Got %v",
            jsonBody)
    }
    return jsonObjectData
}

func verifyJsonArray(t *testing.T, jsonBody string) []map[string]interface{} {
    var jsonArrayData []map[string]interface{}
    jsonArray := json.Unmarshal([]byte(jsonBody), &jsonArrayData)
    if jsonArray != nil  {
        t.Errorf("Expected json array. Got %v",
            jsonBody)
    }
    return jsonArrayData
}


func verifyStatusCode(t *testing.T, resRec *httptest.ResponseRecorder) {
    if status := resRec.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
}

// send & record
func sendBodyAndRecord(t *testing.T, requestMethod string, endpoint string, body string) *http.Request {
    req, err := http.NewRequest(requestMethod, endpoint, strings.NewReader(body))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")
    return req
}

func sendNoBodyAndRecord(t *testing.T, requestMethod string, endpoint string) *http.Request {
    req, err := http.NewRequest(requestMethod, endpoint, nil)
    if err != nil {
        t.Fatal(err)
    }
    return req
}

// serveHttp*
func serveHttpRecordOfficer(request *http.Request) *httptest.ResponseRecorder {
    resRec := httptest.NewRecorder()
    handler := http.HandlerFunc(src.OfficersHandler)
    handler.ServeHTTP(resRec, request)
    return resRec
}

func serveHttpRecordRoot(request *http.Request) *httptest.ResponseRecorder {
    resRec := httptest.NewRecorder()
    handler := http.HandlerFunc(src.RootHandler)
    handler.ServeHTTP(resRec, request)
    return resRec
}

func serveHttpRecordBikeTheft(request *http.Request) *httptest.ResponseRecorder {
    resRec := httptest.NewRecorder()
    handler := http.HandlerFunc(src.BikeTheftsHandler)
    handler.ServeHTTP(resRec, request)
    return resRec
}