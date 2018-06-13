package model

import (
	"fmt"
	"log"

	"github.com/vilisseranen/castellers/common"
)

const EVENTS_TABLE = "events"

// Tables creation queries
const EventsTableCreationQuery = `CREATE TABLE IF NOT EXISTS events
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	startDate INTEGER NOT NULL,
	endDate INTEGER NOT NULL,
	description TEXT,
	uuid TEXT NOT NULL,
	CONSTRAINT uuid_unique UNIQUE (uuid)
);`

type Recurring struct {
	Interval string `json:"interval"`
	Until    uint   `json:"until"`
}

type Event struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	StartDate uint      `json:"startDate"`
	EndDate   uint      `json:"endDate"`
	Recurring Recurring `json:"recurring"`
}

func (e *Event) Get() error {
	stmt, err := db.Prepare(fmt.Sprintf("SELECT name, startDate, endDate FROM %s WHERE uuid= ?", EVENTS_TABLE))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	err = stmt.QueryRow(e.UUID).Scan(&e.Name, &e.StartDate, &e.EndDate)
	return err
}

func (e *Event) GetAll(start, count int) ([]Event, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT uuid, name, startDate, endDate FROM %s LIMIT ? OFFSET ?", EVENTS_TABLE), count, start)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	Events := []Event{}

	for rows.Next() {
		var e Event
		if err = rows.Scan(&e.UUID, &e.Name, &e.StartDate, &e.EndDate); err != nil {
			return nil, err
		}
		Events = append(Events, e)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return Events, nil
}

func (e *Event) UpdateEvent() error {
	stmt, err := db.Prepare(fmt.Sprintf("Update %s SET name = ?, startDate = ?, endDate = ? WHERE uuid= ?", EVENTS_TABLE))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.Name, e.StartDate, e.EndDate, e.UUID)
	return err
}

func (e *Event) DeleteEvent() error {
	stmt, err := db.Prepare(fmt.Sprintf("DELETE FROM %s WHERE uuid= ?", EVENTS_TABLE))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.UUID)
	return err
}

func (e *Event) CreateEvent() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(fmt.Sprintf("INSERT INTO %s (uuid, name, startDate, endDate) VALUES (?, ?, ?, ?)", EVENTS_TABLE))
	if err != nil {
		return err
	}
	defer stmt.Close()
	e.UUID = common.GenerateUUID()
	_, err = stmt.Exec(e.UUID, e.Name, e.StartDate, e.EndDate)
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}
