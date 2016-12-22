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

// This function has an annoying amount of parameters...
func findPage(startingPage, endingPage string, path []string, pageFound chan []string) {
  pageHTML := getUrl(startingPage)
  urls := extractUrls(pageHTML)

  path = append(path, startingPage)

  for _, url := range urls {
    select {
      case <- pageFound:
        return
      default:
        if (url == endingPage) {
          pageFound <- path
        } else {
          go findPage(url, endingPage, path, pageFound)
        }
    }
  }
}

func main() {
  startingPage := "https://en.wikipedia.org/wiki/Pikachu"
  endingPage := "https://en.wikipedia.org/wiki/Central_processing_unit"

  pageFound := make(chan []string)

  go findPage(startingPage, endingPage, []string{}, pageFound)
  found := <- pageFound

  found = append(found, endingPage)

  for index, url := range found {
    fmt.Println("Step", index, ":", url)
  }
}
