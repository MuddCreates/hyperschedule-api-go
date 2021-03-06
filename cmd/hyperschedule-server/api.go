package main

import (
  "github.com/MuddCreates/hyperschedule-api-go/internal/api"
  "log"
  "net/http"
  "encoding/json"
)

func getIp(req *http.Request) string {
  f := req.Header.Get("X-Forwarded-For")
  if len(f) != 0 {
    return f
  }
  return req.RemoteAddr
}

func apiV3Handler(resp http.ResponseWriter, req *http.Request) {
  log.Printf("API: received request from %s", getIp(req))

  data := state.GetData()
  if data == nil {
    log.Printf("received api request before loaded")
    resp.WriteHeader(http.StatusServiceUnavailable)
    return
  }
  output, err := json.Marshal(api.MakeV3(state.GetData()))
  if err != nil {
    log.Printf("api: failed to jsonify, %s", err)
    resp.WriteHeader(http.StatusInternalServerError)
    return
  }

  resp.Header().Add("Content-Type", "application/json")
  resp.Header().Add("Access-Control-Allow-Origin", "*")
  resp.Write(output)
}
