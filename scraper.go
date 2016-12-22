package main

import (
  "fmt"
  "os"
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

func extractUrlsAndNames(html string) (wikipediaUrlsAndNames map[string]string) {
  wikipediaUrlsAndNames = make(map[string]string)
  wikipediaUrlPrefix := "https://en.wikipedia.org"

  untrimmedAnchors := strings.Split(html, "<a href=\"")[1:]
  for _, untrimmedAnchor := range untrimmedAnchors {
    url := untrimmedAnchor[:strings.Index(untrimmedAnchor, "\"")]

    if (strings.HasPrefix(url, "/wiki") && !strings.Contains(url, ":")) {
      fullUrl := wikipediaUrlPrefix + url
      name := untrimmedAnchor[strings.Index(untrimmedAnchor, ">") + 1:strings.Index(untrimmedAnchor, "<")]

      if (!strings.Contains(name, "<") || name == "Read") {
        wikipediaUrlsAndNames[fullUrl] = name
      }
    }
  }

  return wikipediaUrlsAndNames
}

func findPageWithName(startingPage, endingPage string, path []map[string]string, pageFound chan []map[string]string) {
  pageHTML := getUrl(startingPage)
  urls := extractUrlsAndNames(pageHTML)

  temp := make([]map[string]string, len(path))
  copy(temp, path)

  for url, name := range urls {
    select {
      case <- pageFound:
        return
      default:
        temp = append(path, map[string]string {startingPage: name})
        if (url == endingPage) {
          pageFound <- temp
          close(pageFound)
          return
        } else {
          go findPageWithName(url, endingPage, temp, pageFound)
        }
    }
  }
}

func main() {
  startingPage := "https://en.wikipedia.org/wiki/Pikachu"
  endingPage := "https://en.wikipedia.org/wiki/Central_processing_unit"

  if (len(os.Args) == 3) { // Command line arguments
    startingPage = os.Args[1]
    endingPage = os.Args[2]
  }

  pageFound := make(chan []map[string]string)

  var path []map[string]string

  go findPageWithName(startingPage, endingPage, path, pageFound)
  found := <- pageFound

  fmt.Println("===== Begin  =====")
  fmt.Println("- Begin at:", startingPage)
  for index, urlAndNextStep := range found {
    for url, nextStep := range urlAndNextStep {
      fmt.Println("===== Step", (index + 1), "=====")
      fmt.Println("- Start at:", url)
      if (nextStep == "") {
        fmt.Println("- Click on: [Error - Unknown Next Step]")
      } else {
        fmt.Println("- Click on:", nextStep)
      }
    }
  }
  fmt.Println("===== End    =====")
  fmt.Println("- End at:", endingPage)

}
