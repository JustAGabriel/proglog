package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	Log *Log
}

func newHTTPServer() *server {
	return &server{
		Log: NewLog(),
	}
}

func (s *server) CreateNewRecord(w http.ResponseWriter, r *http.Request) {
	var req CreateRecordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := CreateRecordResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *server) GetRecord(w http.ResponseWriter, r *http.Request) {
	var req GetRecordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	record, err := s.Log.Read(req.Offset)
	if err == ErrOffsetNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := GetRecordResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func NewHTTPServer(addr string) *http.Server {
	server := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/records/", server.CreateNewRecord).Methods("POST")
	r.HandleFunc("/api/v1/records/", server.GetRecord).Methods("GET")
	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}
