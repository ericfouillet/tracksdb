package model

import (
	"database/sql"
	"time"
)

// Track represents a track
type Track struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	Artist string        `json:"artist"`
	Length time.Duration `json:"length"`
	Bpm    int           `json:"bpm"`
}

// Tracker is an interface for all handlers of tracks
type Tracker interface {
	Find(*sql.Tx, string) error
	FindAll(*sql.Tx) ([]*Track, error)
	Create(*sql.Tx) error
	Update(*sql.Tx) error
	Delete(*sql.Tx) error
}
