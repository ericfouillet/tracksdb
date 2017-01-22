package tracksdb

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"

	"github.com/ericfouillet/tracksdb/model"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// Start registers handlers and starts listening on port 8080
func Start() {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/tracks", tracksHandler).Methods("GET")
	r.HandleFunc("/tracks/", newTrackHandler).Methods("POST")
	r.HandleFunc("/tracks/{id}", trackHandler).Methods("GET", "PUT", "DELETE")
	http.ListenAndServe(":8080", r)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tracks", http.StatusTemporaryRedirect)
}

func tracksHandler(w http.ResponseWriter, r *http.Request) {
	tx, err := getTx()
	if err != nil {
		http.Error(w, "Issue with database connection: "+err.Error(), http.StatusBadRequest)
		return
	}
	track := &model.Track{}
	tracks, err := track.FindAll(tx)
	if err != nil {
		http.Error(w, "Could not find tracks "+err.Error(), http.StatusBadRequest)
		tx.Rollback()
		return
	}
	tx.Commit()
	json.NewEncoder(w).Encode(tracks)
}

func newTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		track := &model.Track{}
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&track); err != nil {
			http.Error(w, "Could not read track details "+err.Error(), http.StatusBadRequest)
			return
		}
		tx, err := getTx()
		if err != nil {
			http.Error(w, "Issue with database connection: "+err.Error(), http.StatusBadRequest)
			tx.Rollback()
			return
		}
		if err := track.Create(tx); err != nil {
			http.Error(w, "Could not create track "+err.Error(), http.StatusBadRequest)
			tx.Rollback()
			return
		}
		if err := json.NewEncoder(w).Encode(track); err != nil {
			http.Error(w, "Could not encode track "+err.Error(), http.StatusBadRequest)
			tx.Rollback()
			return
		}
		tx.Commit()
	}
}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		deleteTrack(w, r)
	case http.MethodPut:
		updateTrack(w, r)
	case http.MethodGet:
		findTrack(w, r)
	default:
		http.Error(w, "Only GET, POST and DELETE are supported", http.StatusMethodNotAllowed)
	}
}

func deleteTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tx, err := getTx()
	if err != nil {
		http.Error(w, "Issue with database connection: "+err.Error(), http.StatusBadRequest)
		tx.Rollback()
		return
	}
	track := &model.Track{ID: id}
	if err := track.Delete(tx); err != nil {
		http.Error(w, "Could not delete track"+err.Error(), http.StatusBadRequest)
		tx.Rollback()
		return
	}
	tx.Commit()
	tx, err = getTx()
	if err != nil {
		http.Error(w, "Issue with database connection: "+err.Error(), http.StatusBadRequest)
		tx.Rollback()
		return
	}
	track = &model.Track{}
	err = track.Find(tx, id)
	if err == nil { // Expect an error during Find
		http.Error(w, "Could not delete track "+err.Error(), http.StatusBadRequest)
		return
	}
	if track.ID != "" {
		http.Error(w, "Track was not deleted: "+track.ID, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func updateTrack(w http.ResponseWriter, r *http.Request) {
	tx, err := getTx()
	if err != nil {
		http.Error(w, "Issue with database connection: "+err.Error(), http.StatusBadRequest)
		tx.Rollback()
		return
	}
	track := &model.Track{}
	if err := json.NewDecoder(r.Body).Decode(track); err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		tx.Rollback()
		return
	}
	err = track.Update(tx)
	if err != nil {
		http.Error(w, "Could not update track "+err.Error(), http.StatusBadRequest)
		tx.Rollback()
		return
	}
	if err := json.NewEncoder(w).Encode(track); err != nil {
		http.Error(w, "Could not encode track "+err.Error(), http.StatusBadRequest)
		tx.Rollback()
		return
	}
	tx.Commit()
}

func findTrack(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tx, err := getTx()
	if err != nil {
		http.Error(w, "Issue with database connection: "+err.Error(), http.StatusBadRequest)
		tx.Rollback()
		return
	}
	track := &model.Track{ID: id}
	err = track.Find(tx, id)
	if err != nil {
		http.Error(w, "Could not find track "+err.Error(), http.StatusBadRequest)
		tx.Rollback()
		return
	}
	if err := json.NewEncoder(w).Encode(track); err != nil {
		http.Error(w, "Could not encode track "+err.Error(), http.StatusBadRequest)
		tx.Rollback()
		return
	}
	tx.Commit()
}

func getTx() (*sql.Tx, error) {
	username := os.Getenv("MYSQL_USERNAME")
	dbname := os.Getenv("MYSQL_DB")
	db, err := sql.Open("mysql", username+":@/"+dbname)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open connection to the database")
	}
	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "Could not open connection to the database")
	}
	tx, err := db.Begin()
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "Could not start transaction")
	}
	return tx, nil
}
