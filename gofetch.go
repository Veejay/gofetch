package main

import (
  "fmt"
  "net/http"
  "os"
  "net/url"
  "code.google.com/p/go.net/html"
)

type HttpResponse struct {
  url string
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

func extractLinksFromPage (address string, c chan<- string) {
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
        fmt.Printf("Encounted a <html> tag. The value of i is %d\n", i)
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
        fmt.Printf("Encounted a </html> tag. The value of i is %d\n", i)
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
  if response.StatusCode == http.StatusNotFound {
    responses <- HttpResponse{href, response.StatusCode}
  } else {
    responses <- HttpResponse{href, response.StatusCode}
  }
}

func main() {
  hrefs := make(chan string)
  httpResponses := make(chan HttpResponse)
  // Here we go
  go extractLinksFromPage(os.Args[1], hrefs)
  numberOfLinks := 0
  for href := range hrefs {
    numberOfLinks++
    go checkLink(href, httpResponses)
  }
  for i := 0; i < numberOfLinks; i++ {
    response := <-httpResponses
    fmt.Printf("Status %d for URL %s\n", response.statusCode, response.url)
  }
}
