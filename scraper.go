package main

import (
  "fmt"
  "os"
  "io/ioutil"
  "net/http"
  "strings"
)

// Receives a URL and returns all of the HTML
func getUrl(url string) (pageData string){
  res, err := http.Get(url) // Make the request

  if (err != nil) { // Request has an error
    return
  }

  bytes, _ := ioutil.ReadAll(res.Body) // Read the request
  defer res.Body.Close() // Set the request body to close after function return

  return (string(bytes)) // Return the webpage HTML
}

// Extracts the links and the names of the links from Wikipedia page HTML
func extractUrlsAndNames(html string) (wikipediaUrlsAndNames map[string]string) {
  wikipediaUrlsAndNames = make(map[string]string) // This will hold the links and their name
  wikipediaUrlPrefix := "https://en.wikipedia.org"

  untrimmedAnchors := strings.Split(html, "<a href=\"")[1:] // A list of anchor tags that needs to be trimmed down
  for _, untrimmedAnchor := range untrimmedAnchors {
    url := untrimmedAnchor[:strings.Index(untrimmedAnchor, "\"")] // Grab the link URL

    // Ignore links that do not start with /wiki or are internal links
    if (strings.HasPrefix(url, "/wiki") && !strings.Contains(url, ":")) {
      fullUrl := wikipediaUrlPrefix + url // Create the full URL
      name := untrimmedAnchor[strings.Index(untrimmedAnchor, ">") + 1:strings.Index(untrimmedAnchor, "<")] // Grab the link name

      // Ignore links with internal HTML or the read link at the top of the page
      if (!strings.Contains(name, "<") || name == "Read") {
        wikipediaUrlsAndNames[fullUrl] = name
      }
    }
  }

  return wikipediaUrlsAndNames
}

// Searches through links on a given Wikipedia page until it finds the ending page
func findPage(startingPage, endingPage string, path []map[string]string, pageFound chan []map[string]string) {
  pageHTML := getUrl(startingPage) // Grabs the pages HTML
  urls := extractUrlsAndNames(pageHTML) // Finds all the links in the HTML

  // Temp holds the path that findPage takes, but it is a copy to ensure that the path of each request is not being appended to the same slice
  temp := make([]map[string]string, len(path))
  copy(temp, path)

  // Traverse all URLs on the page and see if any match the ending page URL
  // If not, take the URL and run findPage with the new URL as startingPage
  for url, name := range urls {
    select {
      case <- pageFound: // Checks if the page has already been found to end all searching goroutines
        return
      default:
        temp = append(path, map[string]string {startingPage: name}) // Add the new search to the path
        if (url == endingPage) { // Found the ending page
          pageFound <- temp // Send the path through the channel
          close(pageFound) // Close the channel
          return
        } else {
          go findPage(url, endingPage, temp, pageFound) // Run findPage with the new URL to continue the search
        }
    }
  }
}

func main() {
  // Example starting and ending page if none are supplied through command line
  startingPage := "https://en.wikipedia.org/wiki/Pikachu"
  endingPage := "https://en.wikipedia.org/wiki/Central_processing_unit"

  // Starting and ending pages have been supplied through command line
  if (len(os.Args) == 3) {
    startingPage = os.Args[1]
    endingPage = os.Args[2]
  }

  pageFound := make(chan []map[string]string) // Create channel to listen for the ending page being found

  var path []map[string]string // Create a beginning slice to hold the path that the program takes to reach the ending page

  // Begin finding the page at the starting page
  go findPage(startingPage, endingPage, path, pageFound)
  found := <- pageFound // Holds the path from the starting page to the ending page

  // Display the path that the program took
  fmt.Println("===== Begin  =====")
  fmt.Println("- Begin at:", startingPage)
  for index, urlAndNextStep := range found {
    for url, nextStep := range urlAndNextStep {
      fmt.Println("===== Step", (index + 1), "=====")
      fmt.Println("- Start at:", url)
      if (nextStep == "") { // Occasional error where there is no next step
        fmt.Println("- Click on: [Error - Unknown Next Step]")
      } else {
        fmt.Println("- Click on:", nextStep)
      }
    }
  }
  fmt.Println("===== End    =====")
  fmt.Println("- End at:", endingPage)
}
