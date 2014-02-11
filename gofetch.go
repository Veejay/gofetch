package main

import (
  "code.google.com/p/go.net/html"
  "github.com/gorilla/websocket"
  "github.com/gorilla/mux"
  "fmt"
  "encoding/json"
  "net/http"
  "net/url"
)

type QueryURL struct {
  Url string
}

type HttpResponse struct {
  Url        string
  StatusCode int
}

func getHypertextReference(tag html.Token) (href string) {
  for _, attr := range tag.Attr {
    if attr.Key == "href" {
      href = attr.Val
      break
    }
  }
  return href
}

func extractLinksFromPage(address string, c chan<- string) {
  response, err := http.Get(address)
  if err != nil {
    fmt.Printf("An error occurred while issuing a HTTP GET request to %s\n", address)
    return
  }
  i := 0
  defer response.Body.Close()
  tokenizer := html.NewTokenizer(response.Body)
  for {
    tokenType := tokenizer.Next()
    if tokenType == html.ErrorToken {
      return
    }
    token := tokenizer.Token()
    switch tokenType {
    case html.StartTagToken:
      if token.Data == "html" {
        i += 1
      }
      if token.Data == "a" {
        href := getHypertextReference(token)
        location, err := url.Parse(href)
        if err != nil {
          // FIXME: That's wrong, handle error
          panic(err)
        }
        location.RawQuery = url.QueryEscape(location.RawQuery)
        if location.Scheme == "http" || location.Scheme == "https" {
          c <- location.String()
        }
      }
    case html.EndTagToken:
      if token.Data == "html" {
        i -= 1
        if i == 0 {
          close(c)
        }
      }
    }
  }
}

func checkLink(href string, responses chan<- HttpResponse) {
  response, err := http.Get(href)
  // FIXME: This retry is awful, we might want to 
  // send that to a retry channel or something
  if err != nil {
    response, err := http.Get(href)
    if err != nil {
      responses <- HttpResponse{href, 999}
      return
    }
    defer response.Body.Close()
    responses <- HttpResponse{href, response.StatusCode}
  }
  defer response.Body.Close()
  responses <- HttpResponse{href, response.StatusCode}
}

func WebSocketHandler(rw http.ResponseWriter, request *http.Request) {
  ws, err := websocket.Upgrade(rw, request, nil, 1024, 1024)
  if err != nil {
    return
  }
  for {
    messageType, p, err := ws.ReadMessage()
    if err != nil {
        return
    } else {
      var urlData QueryURL
      err := json.Unmarshal(p, &urlData)
      if err != nil {
        panic(err)
      }
      hrefs := make(chan string)
      httpResponses := make(chan HttpResponse)

      go extractLinksFromPage(urlData.Url, hrefs)
      numberOfLinks := 0
      for href := range hrefs {
        numberOfLinks += 1
        go checkLink(href, httpResponses)
      }
      for i := 0; i < numberOfLinks; i += 1 {
        response := <-httpResponses
        responseJSON, err := json.Marshal(response)
        if err != nil {
          panic(err)
        }
        ws.WriteMessage(messageType, responseJSON)
      }
    }
  }
}

// Parameters are simply the request and the response writer
func RootHandler(rw http.ResponseWriter, request *http.Request) {
  fmt.Println("GET /")
  fmt.Println("Processing by root handler")
  http.ServeFile(rw, request, "public/gofetch.html")
}

func QueryHandler(rw http.ResponseWriter, request *http.Request) {
  // Process the POST request
}

func main() {
  r := mux.NewRouter()

  s := r.Schemes("http").Host("localhost").Subrouter()
  s.HandleFunc("/", RootHandler)

  s.HandleFunc("/websocket", WebSocketHandler)

  http.Handle("/", s)
  fmt.Println("Up and listening on port 12345")
  http.ListenAndServe(":12345", nil)

}
