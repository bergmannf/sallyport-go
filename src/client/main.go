package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Event struct {
	Data string
}

type EventSource struct {
	Url    string
	Events chan Event
}

func (e EventSource) ParseData(data []byte) []Event {
	events := []Event{}
	rawEvents := strings.Split(string(data), "\n\n")
	for _, event := range rawEvents {
		trimmed := strings.TrimLeft(event, "data: ")
		events = append(events, Event{trimmed})
	}
	return events
}

func ReadLine(reader *bufio.Reader) ([]byte, error) {
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return line, nil
}

func ReadEvent(buffer *bytes.Buffer, reader *bufio.Reader) {
	line, err := ReadLine(reader)
	res := bytes.Compare(line, []byte{'\n'})
	if res == 0 {
		return
	}
	buffer.Write(line)
	line, err = ReadLine(reader)
	if err != nil {
		buffer.Reset()
		return
	}
	res = bytes.Compare(line, []byte{'\n'})
	if res == 0 {
		return
	} else {
		buffer.Write(line)
		ReadEvent(buffer, reader)
	}
}

func (e EventSource) readResponse(response *http.Response) {
	for {
		log.Println("Reading body")
		buffer := bytes.Buffer{}
		reader := bufio.NewReader(response.Body)
		ReadEvent(&buffer, reader)
		log.Println("Reading body ...")
		events := e.ParseData(buffer.Bytes())
		log.Println("Parsed body", events)
		for _, event := range events {
			e.Events <- event
		}
	}
}

func listen(address string) {
	log.Println("Performing Request")
	response, err := http.Get(address)
	log.Println("Request done")
	if err != nil {
		log.Println("Error while executing GET request: ", err)
	}
	es := EventSource{Url: address, Events: make(chan Event, 100)}
	es.readResponse(response)
}

func forward() {

}

func main() {
	endpoint := flag.String("endpoint-address", "", "Address of the endpoint that will sent out events.")
	forward := flag.String("forward-address", "", "Address to forward the webhook received events to.")
	flag.Parse()
	if *endpoint == "" || *forward == "" {
		fmt.Println("endpoint-address must be set")
		os.Exit(1)
	}
	if *forward == "" {
		fmt.Println("forward-address must be set")
		os.Exit(2)
	}
	listen(*endpoint)
}
