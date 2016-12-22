package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "strings"
)

func getUrl(url string) (pageData string){
  res, err := http.Get(url)

  if (err != nil) {
    return
  }

  bytes, _ := ioutil.ReadAll(res.Body)
  defer res.Body.Close()

  return (string(bytes))
}

func extractUrls(html string) (wikipediaUrls []string) {
  wikipediaUrlPrefix := "https://en.wikipedia.org"

  untrimmedAnchors := strings.Split(html, "<a href=\"")[1:]
  for _, untrimmedAnchor := range untrimmedAnchors {
    url := untrimmedAnchor[:strings.Index(untrimmedAnchor, "\"")]

    if (strings.HasPrefix(url, "/wiki") && !strings.Contains(url, ":")) {
      wikipediaUrls = append(wikipediaUrls, wikipediaUrlPrefix + url)
    }
  }

  return wikipediaUrls
}

func findPage(startingPage, endingPage string, depth int, pageFound chan int) {
  pageHTML := getUrl(startingPage)
  urls := extractUrls(pageHTML)

  for _, url := range urls {
    select {
      case <- pageFound:
        return
      default:
        if (url == endingPage) {
          fmt.Println(url)
          pageFound <- depth
        } else {
          go findPage(url, endingPage, depth + 1, pageFound)
        }
    }
  }
}

func main() {
  startingPage := "https://en.wikipedia.org/wiki/Pikachu"
  endingPage := "https://en.wikipedia.org/wiki/Central_processing_unit"

  pageFound := make(chan int)

  go findPage(startingPage, endingPage, 0, pageFound)
  found := <- pageFound

  fmt.Println(found)
}
