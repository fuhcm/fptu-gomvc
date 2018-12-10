package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gosu-team/cfapp-api/lib"
)

// Feed ...
type Feed struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

// Item ...
type Item struct {
	Title       string   `json:"title"`
	PubDate     string   `json:"pubDate"`
	Link        string   `json:"link"`
	GUID        string   `json:"guid"`
	Author      string   `json:"author"`
	Thumbnail   string   `json:"thumbnail"`
	Description string   `json:"description"`
	Content     string   `json:"content"`
	Categories  []string `json:"categories"`
}

// FeedReponse ...
type FeedReponse struct {
	Status string `json:"status"`
	Feed   Feed   `json:"feed"`
	Items  []Item `json:"items"`
}

func getFeeds(body []byte) (*FeedReponse, error) {
	var f = new(FeedReponse)
	err := json.Unmarshal(body, &f)
	if err != nil {
		fmt.Println("Whoops:", err)
	}

	return f, err
}

// GetPostsByURLHandler ...
func GetPostsByURLHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")

	res := lib.Response{ResponseWriter: w}

	apiKey := os.Getenv("API_KEY")
	resp, err := http.Get("https://api.rss2json.com/v1/api.json?rss_url=" + url + "/feed" + "&api_key=" + apiKey + "&count=10&order_by=pubDate")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	f, err := getFeeds([]byte(body))

	res.SendOK(f)
}