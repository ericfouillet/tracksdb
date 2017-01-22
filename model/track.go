package model

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// Create creates a track
func (t *Track) Create(tx *sql.Tx) error {
	stmt, err := tx.Prepare("INSERT INTO Tracks(id, name, artist, length, bpm) values(?, ?, ?, ?, ?)")
	if err != nil {
		return errors.Wrap(err, "Could not prepare insert statement")
	}
	defer stmt.Close()
	id, err := getNextID("track")
	if err != nil {
		return errors.Wrap(err, "Could not generate a new ID")
	}
	_, err = stmt.Exec(id, t.Name, t.Artist, t.Length, t.Bpm)
	if err != nil {
		return errors.Wrap(err, "Could not persist track")
	}
	t.ID = id
	return nil
}

// Update updates a track
func (t *Track) Update(tx *sql.Tx) error {
	stmt, err := tx.Prepare("UPDATE Tracks SET Name=?, Artist=?, Length=?, Bpm=? WHERE ID=?")
	if err != nil {
		return errors.Wrap(err, "Could not prepare insert statement")
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.Name, t.Artist, t.Length, t.Bpm, t.ID)
	if err != nil {
		return errors.Wrap(err, "Could not persist track")
	}
	return nil
}

// Delete ccc
func (t *Track) Delete(tx *sql.Tx) error {
	stmt, err := tx.Prepare("DELETE FROM Tracks where ID = ?")
	if err != nil {
		return errors.Wrap(err, "Could not prepare delete statement")
	}
	defer stmt.Close()
	res, err := stmt.Exec(t.ID)
	if err != nil {
		return errors.Wrap(err, "Could not delete track")
	}
	n, err := res.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Could not delete track")
	}
	if n == 0 {
		return errors.New("0 rows affected")
	}
	return nil
}

// Find finds a single track
func (t *Track) Find(tx *sql.Tx, idQuery string) error {
	row := tx.QueryRow("SELECT * FROM Tracks WHERE ID = ?", idQuery)
	var id, name, artist string
	var length, bpm int
	if err := row.Scan(&id, &name, &artist, &length, &bpm); err != nil {
		return errors.Wrap(err, "Could not read row from DB")
	}
	t.ID, t.Name, t.Artist, t.Length, t.Bpm = id, name, artist, time.Duration(length), bpm
	return nil
}

// FindAll finds all persisted tracks
func (t *Track) FindAll(tx *sql.Tx) ([]*Track, error) {
	rows, err := tx.Query("SELECT * FROM Tracks ORDER BY ID DESC")
	if err != nil {
		return nil, errors.Wrap(err, "Could not find tracks")
	}
	defer rows.Close()
	var tracks []*Track
	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "Error while parsing result rows")
	}
	for rows.Next() {
		var id, name, artist string
		var length, bpm int
		rows.Scan(&id, &name, &artist, &length, &bpm)
		tracks = append(tracks, &Track{ID: id, Name: name, Artist: artist, Length: time.Duration(length), Bpm: bpm})
	}
	return tracks, nil
}
