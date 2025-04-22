package main

import (
	"encoding/json"
	"log"
	"net/http"

	toolbooth "github.com/didip/tollbooth/v7"
)

type Message struct {
  Status string `json:"status"`
  Body string `json:"body"`
}

func endpointHandler(writer http.ResponseWriter, request *http.Request) {
  writer.Header().Set("Content-Type", "application/json")
  writer.WriteHeader(http.StatusOK)
  message := Message{
    Status: "Successful",
    Body: "You have reached the API",
  }

  err := json.NewEncoder(writer).Encode(&message)
  if err != nil {
    return
  }
}

func main() {
  message := Message{
    Status: "Request failed",
    Body: "The API is at capacity",
  }
  jsonMessage, _ := json.Marshal(message)
  tlbthLimiter := toolbooth.NewLimiter(1, nil)
  tlbthLimiter.SetMessageContentType("application/json")
  tlbthLimiter.SetMessage(string(jsonMessage))

  http.Handle("/ping", toolbooth.LimitFuncHandler(tlbthLimiter, endpointHandler))
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Println("Error listening on port :8080", err)
  }
}
