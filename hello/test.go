package main

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "encoding/json"
  "log"
  "net/http"
)

const (
  defaultServiceProviderPath string = "http://localhost:3010/api/v1/services"
  defaultHost string = "localhost"
  defaultPort int = 3090
  defaultName string = "hello-service"
)

func main() {
  http.HandleFunc("/", handler)
  http.HandleFunc("/register", register_service)
  http.HandleFunc("/online", make_online)
  http.HandleFunc("/deregister", deregister_service)
  http.HandleFunc("/chuck", get_joke)
  log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, you've requested: %q\n", r.URL.Path)
}

func register_service(w http.ResponseWriter, r *http.Request) {
  url := fmt.Sprintf("%s/%s:%d", defaultServiceProviderPath, defaultName, defaultPort)
  fmt.Println("URL:>", url)

  var jsonStr = []byte(fmt.Sprintf(`{
    "__v": 0,
    "serviceType": "dummy",
    "name": %q,
    "description": "Dummy service for InfoHub presentation",
    "host": %q,
    "port": %d,
    "path": "/api/dummy/v1",
    "status": "online",
    "_id": "%s:%d",
    "buildInfo": {
        "version": "local",
        "commit": "development",
        "buildDate": "now"
    }
}`, defaultName, defaultHost, defaultPort, defaultHost, defaultPort))

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  fmt.Println("response Status:", resp.Status)
  if (resp.Status == "200") {
    fmt.Fprintf(w, "The service is registered!")
  } else {
    fmt.Fprintf(w, "Service registration failed...")
  }
}

func make_online(w http.ResponseWriter, r *http.Request) {
  url := fmt.Sprintf("%s/%s:%d/online", defaultServiceProviderPath, defaultName, defaultPort)
  fmt.Println("URL:>", url)

  req, err := http.NewRequest("POST", url, nil)
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  fmt.Println("response Status:", resp.Status)
  if (resp.Status == "200") {
    fmt.Fprintf(w, "The service is online!")
  } else {
    fmt.Fprintf(w, "Setting service to online failed...")
  }
}

func deregister_service(w http.ResponseWriter, r *http.Request) {
  url := fmt.Sprintf("%s/%s:%d", defaultServiceProviderPath, defaultName, defaultPort)
  fmt.Println("URL:>", url)

  req, err := http.NewRequest("DELETE", url, nil)

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  fmt.Println("response Status:", resp.Status)
  if (resp.Status == "200") {
    fmt.Fprintf(w, "The service is deregistered!")
  } else {
    fmt.Fprintf(w, "Service deregistration failed...")
  }
}

func get_joke(w http.ResponseWriter, r *http.Request) {
  url := "http://api.icndb.com/jokes/random"
  fmt.Println("URL:>", url)

  req, err := http.NewRequest("GET", url, nil)

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }
  defer resp.Body.Close()

  fmt.Println("response Status:", resp.Status)
  b, err := ioutil.ReadAll(resp.Body)
  if err != nil {
      panic(err)
  }

  var m Message
  err := json.Unmarshal(b, &m)
  fmt.Printf("%s", m)
}