package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
	"time"

	"./hn"
)

func main() {
	var port, numStories int

	flag.IntVar(&port, "port", 3000, "the port to start web server on")
	flag.IntVar(&numStories, "num_stories", 30, "number of top stories to display")
	flag.Parse()
	tpl := template.Must(template.ParseFiles("./index.gohtml"))
	http.HandleFunc("/", handler(numStories, tpl))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))

}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var client hn.Client
		ids, err := client.GetTopStories()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		var stories []item

		for _, id := range ids {
			hnItem, err := client.GetItem(id)
			if err != nil {
				http.Error(w, "Failed to load top stories here", http.StatusInternalServerError)
				continue
			}

			item := parseHNItem(hnItem)
			if isStoryLink(item) {
				stories = append(stories, item)
				if len(stories) >= numStories {
					break
				}
			}
		}

		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}

	})
}

type templateData struct {
	Stories []item
	Time    time.Duration
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}
func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

type item struct {
	hn.Item
	Host string
}
