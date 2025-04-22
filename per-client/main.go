package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Message struct{
  Status string `json:"status"`
  Body string `json:"body"`

}

func endpointHandler(writer http.ResponseWriter, r *http.Request){
  writer.Header().Set("Content-Type", "application/json")
  writer.WriteHeader(http.StatusOK)
  message := Message{
    Status: "Successful",
    Body: "Hi, you have reached the API",
  }

  err := json.NewEncoder(writer).Encode(&message)
  if err != nil {
    return
  }
}

func main(){
  http.Handle("/ping", perClientRateLimiter(endpointHandler))
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    log.Println("There was an error on port :8080", err)
  }
}
