package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

// NewHTTPServer for API
func NewHTTPServer(address string) *http.Server {
	s := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/", s.handleProducer).Methods("POST")
	r.HandleFunc("/", s.handelConsumer).Methods("GET")
	return &http.Server{
		Addr:    address,
		Handler: r,
	}
}

// ProducerRequest JSON in API
type ProducerRequest struct {
	Record Record `json:"record"`
}

// ProducerResponse JSON in API
type ProducerResponse struct {
	Offset uint64 `json:"offset"`
}

// ConsumerRequest JSON in API
type ConsumerRequest struct {
	Offset uint64 `json:"offset"`
}

// ConsumerResponse JSON in API
type ConsumerResponse struct {
	Record Record `json:"record"`
}

func (s *httpServer) handleProducer(w http.ResponseWriter, r *http.Request) {
	req := ProducerRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	offset, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := ProducerResponse{
		Offset: offset,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *httpServer) handelConsumer(w http.ResponseWriter, r *http.Request) {
	req := ConsumerRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	record, err := s.Log.Read(req.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res := ConsumerResponse{
		Record: record,
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
