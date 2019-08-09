package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
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
		stories, err := getTopStories(numStories)
		if err != nil {
			http.Error(w, "Failed to process the template ::"+err.Error(), http.StatusInternalServerError)
			return
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

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.GetTopStories()
	if err != nil {
		return nil, err
	}
	type result struct {
		item  item
		index int
		err   error
	}
	resultCh := make(chan result)

	for index, id := range ids {
		go func(ind, id int) {
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultCh <- result{index: ind, err: err}
			}
			resultCh <- result{index: ind, item: parseHNItem(hnItem)}
		}(index, id)
	}
	var results []result
	counter := 0
	for res := range resultCh {
		counter++
		if counter == len(ids) {
			break
		}
		if isStoryLink(res.item) {
			results = append(results, res)
		}
	}
	fmt.Println("Got results")
	sort.Slice(results, func(i int, j int) bool {
		return results[i].index < results[j].index
	})
	var stories []item

	for _, res := range results {
		stories = append(stories, res.item)
		if len(stories) >= numStories {
			break
		}
	}

	return stories, nil
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
