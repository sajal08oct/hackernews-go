package hn

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseUrl       = "https://hacker-news.firebaseio.com/v0/"
	topStoriesUrl = "topstories.json"
	itemUrl       = "/item"
)

type Client struct {
	apiBase string
}

func (c *Client) defaultify() {
	if c.apiBase == "" {
		c.apiBase = baseUrl
	}

}

func (c *Client) GetTopStories() ([]int, error) {
	c.defaultify()
	resp, err := http.Get(fmt.Sprintf("%s%s", c.apiBase, topStoriesUrl))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ids []int
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&ids)
	if err != nil {

		return nil, err
	}
	return ids, nil
}

func (c *Client) GetItem(id int) (Item, error) {
	c.defaultify()
	var item Item
	resp, err := http.Get(fmt.Sprintf("%s%s/%d.json", c.apiBase, itemUrl, id))
	if err != nil {
		fmt.Println(err)
		return item, err
	}
	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&item)

	if err != nil {
		fmt.Println(err)
		return item, err
	}

	return item, nil
}

type Item struct {
	ID          int    `json:"id"`
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	Deleted     bool   `json:"deleted"`
	Type        string `json:"type"`
	Kids        []int  `json:"kids"`
	Parent      string `json:"parent"`
	Text        string `json:"text"`
	Title       string `json:"title"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	URL         string `json:"url"`
}
