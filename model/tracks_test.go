package model

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	username := os.Getenv("MYSQL_USERNAME")
	dbname := os.Getenv("MYSQL_DB")
	db, err = sql.Open("mysql", username+":@/"+dbname)
	if err != nil {
		panic(err)
	}
}

func TestCRUD(t *testing.T) {
	track := &Track{
		Name:   "track name",
		Artist: "Artist name",
		Length: 120,
		Bpm:    128,
	}
	tx := getTx()
	defer tx.Rollback()
	err := track.Create(tx)
	if err != nil {
		t.Fatal("Could not create a track")
	}
	if track.ID == "" {
		t.Fatal("ID was not set")
	}
	err = track.Find(tx, track.ID)
	if err != nil {
		t.Fatal("Could not find persisted track")
	}
	track.Bpm = 112
	err = track.Update(tx)
	if err != nil {
		t.Fatal("Could not create a track")
	}
	err = track.Find(tx, track.ID)
	if err != nil {
		t.Fatal("Could not create a track")
	}
	if track.Bpm != 112 {
		t.Fatal("Track was not updated")
	}
	err = track.Delete(tx)
	if err != nil {
		t.Fatal("Could not delete a track")
	}
	err = track.Find(tx, track.ID)
	if err == nil {
		t.Fatal("Track was not deleted")
	}
}

func getTx() *sql.Tx {
	tx, err := db.Begin()
	if err != nil {
		panic("Could not start tx")
	}
	return tx
}
