package main_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/vilisseranen/castellers"
)

var a main.App

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize("test_database.db")

	ensureTablesExist()

	code := m.Run()

	clearTables()

	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTables()

	req, _ := http.NewRequest("GET", "/events", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentEvent(t *testing.T) {
	clearTables()

	req, _ := http.NewRequest("GET", "/events/deadbeef", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Event not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Event not found'. Got '%s'", m["error"])
	}
}

func TestCreateEvent(t *testing.T) {
	clearTables()
	addAdmin("deadbeef")

	payload := []byte(`{"name":"diada","startDate":"2018-06-01 23:16", "endDate":"2018-06-03 17:14"}`)

	req, _ := http.NewRequest("POST", "/admins/deadbeef/events", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "diada" {
		t.Errorf("Expected event name to be 'diada'. Got '%v'", m["name"])
	}

	if m["startDate"] != "2018-06-01 23:16" {
		t.Errorf("Expected event start date to be '2018-06-01 23:16'. Got '%v'", m["date"])
	}

	if m["endDate"] != "2018-06-03 17:14" {
		t.Errorf("Expected event end date to be '2018-06-03 17:14'. Got '%v'", m["date"])
	}
}

func TestCreateEventNonAdmin(t *testing.T) {
	clearTables()
	addAdmin("deadbeef")

	payload := []byte(`{"name":"diada","startDate":"2018-06-01 23:16", "endDate":"2018-06-03 17:14"}`)

	req, _ := http.NewRequest("POST", "/admins/4b1d/events", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusUnauthorized, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != nil {
		t.Errorf("Expected event name to be ''. Got '%v'", m["name"])
	}

	if m["startDate"] != nil {
		t.Errorf("Expected event start date to be ''. Got '%v'", m["date"])
	}

	if m["endDate"] != nil {
		t.Errorf("Expected event end date to be '2018-06-03 17:14'. Got '%v'", m["date"])
	}
}

func TestCreateMember(t *testing.T) {
	clearTables()
	addAdmin("deadbeef")

	payload := []byte(`{"name":"clement","roles": ["baix", "second"], "extra":"Santi"}`)

	req, _ := http.NewRequest("POST", "/admins/deadbeef/members", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "clement" {
		t.Errorf("Expected member name to be 'clement'. Got '%v'", m["name"])
	}

	if m["extra"] != "Santi" {
		t.Errorf("Expected extra to be 'Santi'. Got '%v'", m["extra"])
	}
}

func TestCreateMemberNoExtra(t *testing.T) {
	clearTables()
	addAdmin("deadbeef")

	payload := []byte(`{"name":"clement","roles": ["baix", "second"]}`)

	req, _ := http.NewRequest("POST", "/admins/deadbeef/members", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "clement" {
		t.Errorf("Expected member name to be 'clement'. Got '%v'", m["name"])
	}

	if m["extra"] != "" {
		t.Errorf("Expected extra to be ''. Got '%v'", m["extra"])
	}
}

func TestGetEvent(t *testing.T) {
	clearTables()
	addEvent("deadbeef", "An event", "2018-06-03 18:00", "2018-06-03 21:00")

	req, _ := http.NewRequest("GET", "/events/deadbeef", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "An event" {
		t.Errorf("Expected event name to be 'An event'. Got '%v'", m["name"])
	}

	if m["startDate"] != "2018-06-03 18:00" {
		t.Errorf("Expected event start date to be '2018-06-01 23:16'. Got '%v'", m["date"])
	}

	if m["endDate"] != "2018-06-03 21:00" {
		t.Errorf("Expected event end date to be '2018-06-03 17:14'. Got '%v'", m["date"])
	}
}

func TestGetMember(t *testing.T) {
	clearTables()
	addMember("deadbeef", "Clément")

	req, _ := http.NewRequest("GET", "/members/deadbeef", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "Clément" {
		t.Errorf("Expected member name to be 'Clément'. Got '%v'", m["name"])
	}
}

func TestGetEvents(t *testing.T) {
	clearTables()
	addEvent("deadbeef", "An event", "2018-06-03 18:00", "2018-06-03 21:00")
	addEvent("deadfeed", "Another event", "2018-06-04 18:00", "2018-06-04 21:00")

	req, _ := http.NewRequest("GET", "/events?count=2&start=0", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m [2]map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m[0]["name"] != "An event" {
		t.Errorf("Expected event name to be 'An event'. Got '%v'", m[0]["name"])
	}

	if m[0]["startDate"] != "2018-06-03 18:00" {
		t.Errorf("Expected event start date to be '2018-06-03 18:00'. Got '%v'", m[0]["date"])
	}

	if m[0]["endDate"] != "2018-06-03 21:00" {
		t.Errorf("Expected event end date to be '2018-06-03 21:00'. Got '%v'", m[0]["date"])
	}

	if m[1]["name"] != "Another event" {
		t.Errorf("Expected event name to be 'Another event'. Got '%v'", m[1]["name"])
	}

	if m[1]["startDate"] != "2018-06-04 18:00" {
		t.Errorf("Expected event start date to be '2018-06-04 18:00'. Got '%v'", m[1]["date"])
	}

	if m[1]["endDate"] != "2018-06-04 21:00" {
		t.Errorf("Expected event end date to be '2018-06-04 21:00'. Got '%v'", m[1]["date"])
	}
}

func TestUpdateEvent(t *testing.T) {
	clearTables()
	addEvent("deadbeef", "An event", "2018-06-03 18:00", "2018-06-03 21:00")

	req, _ := http.NewRequest("GET", "/events/deadbeef", nil)

	response := executeRequest(req)
	var originalEvent map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalEvent)

	payload := []byte(`{"name":"test event - updated name","startDate":"2018-06-03 19:00", "endDate":"2018-06-03 22:00"}`)

	req, _ = http.NewRequest("PUT", "/events/deadbeef", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] == originalEvent["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalEvent["name"], "test event - updated name", m["name"])
	}

	if m["startDate"] == originalEvent["startDate"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalEvent["date"], "2018-06-03 19:00", m["startDate"])
	}
	if m["endDate"] == originalEvent["endDate"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalEvent["date"], "2018-06-03 22:00", m["enDate"])
	}
}

func TestDeleteEvent(t *testing.T) {
	clearTables()
	addEvent("deadbeef", "An event", "2018-06-03 18:00", "2018-06-03 21:00")

	req, _ := http.NewRequest("GET", "/events/deadbeef", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/events/deadbeef", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/events/deadbeef", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestParticipateEvent(t *testing.T) {
	clearTables()
	addMember("deadbeef", "toto")
	addEvent("deadbeef", "diada", "2018-06-05 22:55", "2018-06-05 23:55")

	payload := []byte(`{"answer":"yes"}`)

	req, _ := http.NewRequest("POST", "/events/deadbeef/members/deadbeef", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["answer"] != "yes" {
		t.Errorf("Expected answer to be 'yes'. Got '%v'", m["name"])
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func addEvent(uuid, name, startDate, endDate string) {
	tx, err := a.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO events(uuid, name, startDate, endDate) VALUES(?, ?, ?, ?);")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(uuid, name, startDate, endDate)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func addAdmin(uuid string) {
	tx, err := a.DB.Begin()

	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO admins(uuid) VALUES(?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(uuid)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func addMember(uuid, name string) {
	tx, err := a.DB.Begin()

	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO members(uuid, name, extra) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(uuid, name, "")
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func ensureTablesExist() {
	a.DB.Exec("DROP TABLE events")
	a.DB.Exec("DROP TABLE admins")
	a.DB.Exec("DROP TABLE members")
	a.DB.Exec("DROP TABLE presences")
	a.DB.Exec(main.EventsTableCreationQuery)
	a.DB.Exec(main.AdminsTableCreationQuery)
	a.DB.Exec(main.MembersTableCreationQuery)
	a.DB.Exec(main.PresencesTableCreationQuery)
}

func clearTables() {
	a.DB.Exec("DELETE FROM events")
	a.DB.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'events'")
	a.DB.Exec("DELETE FROM admins")
	a.DB.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'admins'")
	a.DB.Exec("DELETE FROM members")
	a.DB.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'members'")
	a.DB.Exec("DELETE FROM presences")
	a.DB.Exec("UPDATE sqlite_sequence SET seq = 0 WHERE name = 'presences'")
}
