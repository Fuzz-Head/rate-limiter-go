package main

import (
	"log"
	"net/http"
  "encoding/json"
)

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func endpointHandler(writer http.ResponseWriter, request *http.Request){
  writer.Header().Set("Content-Type", "application/json")
  writer.WriteHeader(http.StatusOK)
  message := Message{
    Status: "Successful",
    Body: "Hi, You have reached the API, please proceed.",
  }

  err := json.NewEncoder(writer).Encode(&message)
  if err != nil {
    return 
  }
}

func main(){
  http.Handle("/ping", rateLimiter(endpointHandler))
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Println("There was an error on port :8080", err)
  }
}
