package main

import (
  "crypto/tls"
  "flag"
  "fmt"
  "github.com/jackdanger/collectlinks"
  "github.com/mediocregopher/radix.v2/pool"
  // Import the Radix.v2 redis package (we need access to its Nil type).
  "github.com/mediocregopher/radix.v2/redis"
  "net/http"
  "net/url"
  "os"
)
var visited = make(map[string]bool)

func main() {
  flag.Parse()

  args := flag.Args()
  fmt.Println(args)
  if len(args) < 1 {
    fmt.Println("Please specify start page")
    os.Exit(1)
  }

  queue := make(chan string)

  go func() { queue <- args[0] }()

  for uri := range queue {
    enqueue(uri, queue)
  }
}

func enqueue(uri string, queue chan string) {
  fmt.Println("fetching", uri)
  visited[uri] = true
  transport := &http.Transport{
    TLSClientConfig: &tls.Config{
      InsecureSkipVerify: true,
    },
  }
  client := http.Client{Transport: transport}
  resp, err := client.Get(uri)
  if err != nil {
    return
  }
  defer resp.Body.Close()

  links := collectlinks.All(resp.Body)

  for _, link := range links {
    absolute := fixUrl(link, uri)
    if uri != "" {
      if !visited[absolute] {
        go func() { queue <- absolute }()
      }
    }
  }
}

func fixUrl(href, base string) (string) {
  uri, err := url.Parse(href)
  if err != nil {
    return ""
  }
  baseUrl, err := url.Parse(base)
  if err != nil {
    return ""
  }
  //fmt.Printf("uri: %s \nbaseURL: %s\nResolveReference: %s\n", uri,baseUrl,baseUrl.ResolveReference(uri))
  uri = baseUrl.ResolveReference(uri)
  return uri.String()
}
