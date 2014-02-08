package main

import (
  "code.google.com/p/go.net/html"
  "github.com/gorilla/websocket"
  "github.com/gorilla/mux"
  "fmt"
  "net/http"
  "net/url"
  /* "os" */
)

type HttpResponse struct {
  url        string
  statusCode int
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
  // Used to keep track of the opening and closing <html> tags
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
        if location.Scheme == "http" || location.Scheme == "https" {
          c <- href
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
  if err != nil {
    // FIXME: This is absolutely not a 999. The HttpResponse should
    // actually be named something that embeds the URL, the response and
    // any potential errors that occurred
    responses <- HttpResponse{href, 999}
    return
  }
  defer response.Body.Close()
  responses <- HttpResponse{href, response.StatusCode}
}

func WebSocketHandler(rw http.ResponseWriter, request *http.Request) {
  ws, err := websocket.Upgrade(rw, request, nil, 1024, 1024)
  // TODO: Handle error
  if err != nil {
    return
  }
  for {
    messageType, p, err := ws.ReadMessage()
    if err != nil {
        return
    } else {
      fmt.Printf("%s\n", p)
    }
    err = ws.WriteMessage(messageType, p)
    if err != nil {
        return
    }
  }
}

// Parameters are simply the request and the response writer
func RootHandler(rw http.ResponseWriter, request *http.Request) {
  fmt.Printf("GET /\nProcessing by root handler\n")

  http.ServeFile(rw, request, "public/gofetch.html")
}

func QueryHandler(rw http.ResponseWriter, request *http.Request) {
  // Process the POST request
}

func main() {

  r := mux.NewRouter()

  s := r.Schemes("http").Host("localhost").Subrouter()
  s.HandleFunc("/", RootHandler)

  s.HandleFunc("/check", QueryHandler).
    Methods("POST").
    Headers("X-Requested-With", "XMLHttpRequest")

  s.HandleFunc("/websocket", WebSocketHandler)

  http.Handle("/", s)
  fmt.Println("Up and listening on port 12345")
  http.ListenAndServe(":12345", nil)

  /* hrefs := make(chan string) */
  /* httpResponses := make(chan HttpResponse) */
  /* // Here we go */
  /* go extractLinksFromPage(os.Args[1], hrefs) */
  /* numberOfLinks := 0 */
  /* for href := range hrefs { */
  /*   numberOfLinks += 1 */
  /*   go checkLink(href, httpResponses) */
  /* } */
  /* for i := 0; i < numberOfLinks; i += 1 { */
  /*   response := <-httpResponses */
  /*   fmt.Printf("Status %d for URL %s\n", response.statusCode, response.url) */
  /* } */
}
