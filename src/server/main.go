package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const DB_PATH = "./sallyport.db"
// Map endpoints to the backing Queue
var endpoints map[string]Queue = make(map[string]Queue)

type Request struct {
	Uri     string
	Method  string
	Headers map[string][]string
	Body    string
}

func (r Request) json() string {
	serialized, err := json.Marshal(r)
	if err != nil {
		fmt.Println("Could not serialize the HTTP request")
	}
	fmt.Println("Serialized: ", string(serialized))
	return string(serialized)
}

func createNewEndpoint(path string) string {
	newEndpoint := fmt.Sprintf("/endpoint/%s", path)
	fmt.Println("Register new endpoint at ", newEndpoint)
	endpoints[newEndpoint] = NewMemoryQueue()
	return newEndpoint
}

func newEndpoint(w http.ResponseWriter, r *http.Request) {
	uuid := uuid.New()
	endpoint := createNewEndpoint(uuid.String())
	http.HandleFunc(endpoint, handleEndpoint)
}

func serverSideEventsConnection(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Can not use server sent events", http.StatusInternalServerError)
		return
	}
	fmt.Println("ENDPOINT GET: ", r.RequestURI)
	clientQueue := endpoints[r.RequestURI]
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	for {
		request := <- clientQueue.NotificationChannel()
		w.Write([]byte(request.json()))
		flusher.Flush()
	}
}

func storeRequest(r *http.Request) {
	clientQueue := endpoints[r.RequestURI]
	buf := new(bytes.Buffer)
	buf.ReadFrom(r.Body)
	body := buf.String()
	request := Request{
		Uri:     r.RequestURI,
		Method:  r.Method,
		Headers: r.Header,
		Body:    body,
	}
	clientQueue.Put(request)
}

func handleEndpoint(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		storeRequest(r)
	case http.MethodGet:
		serverSideEventsConnection(w, r)
	default:
		w.Write([]byte("Use Post as the webhook and Get to retrieved stored requests."))
	}
}

func main() {
	http.HandleFunc("/new", newEndpoint)
	http.HandleFunc(createNewEndpoint("static"), handleEndpoint)
	http.ListenAndServe("localhost:8080", nil)
}
