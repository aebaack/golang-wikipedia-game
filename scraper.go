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

func extractUrlsAndNames(html string) (wikipediaUrlsAndNames map[string]string) {
  wikipediaUrlsAndNames = make(map[string]string)
  wikipediaUrlPrefix := "https://en.wikipedia.org"

  untrimmedAnchors := strings.Split(html, "<a href=\"")[1:]
  for _, untrimmedAnchor := range untrimmedAnchors {
    url := untrimmedAnchor[:strings.Index(untrimmedAnchor, "\"")]

    if (strings.HasPrefix(url, "/wiki") && !strings.Contains(url, ":")) {
      fullUrl := wikipediaUrlPrefix + url
      name := untrimmedAnchor[strings.Index(untrimmedAnchor, ">") + 1:strings.Index(untrimmedAnchor, "<")]

      wikipediaUrlsAndNames[fullUrl] = name
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
    temp = append(path, map[string]string {startingPage: name})
    select {
      case <- pageFound:
        return
      default:
        if (url == endingPage) {
          pageFound <- temp
        } else {
          go findPageWithName(url, endingPage, temp, pageFound)
        }
    }
  }
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

  // pageFound := make(chan []string)
  //
  // go findPage(startingPage, endingPage, []string{}, pageFound)
  // found := <- pageFound
  //
  // found = append(found, endingPage)
  //
  // for index, url := range found {
  //   fmt.Println("Step", index, ":", url)
  // }

  pageFound := make(chan []map[string]string)

  var path []map[string]string

  go findPageWithName(startingPage, endingPage, path, pageFound)
  found := <- pageFound

  for index, urlAndNextStep := range found {
    for url, nextStep := range urlAndNextStep {
      fmt.Println("===== Step", (index + 1), "=====")
      fmt.Println("Start at", url, "and click", nextStep)
    }
  }

}
