package main

import (
	"encoding/json"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func perClientRateLimiter(next func(writer http.ResponseWriter, request *http.Request)) http.Handler{

  type client struct {
    limiter *rate.Limiter
    lastSeen time.Time
  }

  var (
    mu sync.Mutex
    clients = make(map[string]*client)
  )

  go func() {
    for {
      time.Sleep(time.Minute)
      mu.Lock()
      for ip, client := range clients {
        if time.Since(client.lastSeen) > 5*time.Minute {
          delete(clients, ip)
        }
      }
      mu.Unlock()
    }
  }()

  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      return
    }

    mu.Lock()
    if _, found := clients[ip]; !found {
      clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
    }

    clients[ip].lastSeen = time.Now()
    if !clients[ip].limiter.Allow(){
      message := Message{
        Status: "Request failed",
        Body: "The API is at capacity",
      }

      w.WriteHeader(http.StatusTooManyRequests)
      json.NewEncoder(w).Encode(&message)
      return
    }

    mu.Unlock()
    next(w, r)
  })
}
