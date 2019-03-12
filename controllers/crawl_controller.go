package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gosu-team/fptu-api/config"
	"github.com/gosu-team/fptu-api/lib"
)

// MediumURLs ...
var MediumURLs = []string{"https://codeburst.io",
	"https://medium.freecodecamp.org",
	"https://hackernoon.com",
	"https://medium.com/javascript-scene",
	"https://blog.logrocket.com/",
	"https://medium.com/tag/react",
	"https://medium.com/tag/golang"}

// FPTURLs ...
var FPTURLs = []string{"https://daihoc.fpt.edu.vn"}

// CodeDaoURLS ...
var CodeDaoURLS = []string{"https://toidicodedao.com", "https://codeaholicguy.com/"}

// Feed ...
type Feed struct {
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

// MiniItem ...
type MiniItem struct {
	Title       string `json:"title"`
	PubDate     string `json:"pubDate"`
	GUID        string `json:"guid"`
	Thumbnail   string `json:"thumbnail"`
	Description string `json:"description"`
}

// FeedReponse ...
type FeedReponse struct {
	Status string `json:"status"`
	Feed   Feed   `json:"feed"`
	Items  []Item `json:"items"`
}

// MiniFeedReponse ...
type MiniFeedReponse struct {
	Items []MiniItem `json:"items"`
}

func getFeeds(body []byte) (*FeedReponse, error) {
	var f = new(FeedReponse)
	err := json.Unmarshal(body, &f)
	if err != nil {
		fmt.Println("Whoops:", err)
	}

	return f, err
}

func minimizeItems(items []Item) []MiniItem {
	var miniItems []MiniItem
	for _, v := range items {
		miniItems = append(miniItems, MiniItem{
			Title:       v.Title,
			PubDate:     v.PubDate,
			GUID:        v.GUID,
			Thumbnail:   v.Thumbnail,
			Description: v.Description,
		})
	}

	return miniItems
}

func resolveMediumURL(url string) string {
	urlParts := strings.Split(url, "/")
	string mediumChannel
	if strings.Contains(url, "tag") {
		mediumChannel = urlParts[4]
	} else {
		mediumChannel = urlParts[3]
	}

	return "https://medium.com/feed/" + mediumChannel
}

func getFeedFromURL(url string) *FeedReponse {
	if strings.Contains(url, "https://medium.com/") {
		url = resolveMediumURL(url)
	} else {
		url = url + "/feed"
	}

	// Get and parse API
	apiKey := os.Getenv("API_KEY")
	resp, err := http.Get("http://api.rss2json.com/v1/api.json?rss_url=" + url + "&api_key=" + apiKey + "&count=10&order_by=pubDate")
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}

	feed, err := getFeeds([]byte(body))

	return feed
}

func findArticleByID(name string, items []Item, id string) int {
	for index, v := range items {
		switch name {
		case "medium":
			if v.GUID == "https://medium.com/p/"+id {
				return index
			}
		case "codedao":
			if v.GUID == "http://toidicodedao.com/?p="+id || v.GUID == "http://codeaholicguy.com/?p="+id {
				return index
			}
		case "fpt":
			if v.GUID == "https://daihoc.fpt.edu.vn/?p="+id+"/" {
				return index
			}
		default:
			if v.GUID == id {
				return index
			}
		}
	}

	return -1
}

func getDataFromURLs(urlArr []string) []Item {
	var wg sync.WaitGroup
	var itemCrawl []Item

	wg.Add(len(urlArr))

	for index := range urlArr {
		go func(index int) {
			defer wg.Done()
			crawl := getFeedFromURL(urlArr[index])
			itemCrawl = append(itemCrawl, crawl.Items...)
		}(index)
	}

	wg.Wait()

	return itemCrawl
}

func getDataFromSite(name string) []Item {
	// Reference to system cache
	c := config.GetCache()
	defaultExpiration := config.GetDefaultExpiration()

	var urls []string
	switch name {
	case "medium":
		urls = MediumURLs
	case "codedao":
		urls = CodeDaoURLS
	case "fpt":
		urls = FPTURLs
	default:
		urls = MediumURLs
	}

	// Check cache and use data from cache
	var articles []Item
	articleGot, found := c.Get(name)
	if found {
		articles, _ = articleGot.([]Item)
	} else {
		articlesCache := getDataFromURLs(urls)
		if len(articlesCache) > 0 {
			c.Set(name, articlesCache, defaultExpiration)
		}
		articles = articlesCache
	}

	return articles
}

// GetHomeFeedHandler ...
func GetHomeFeedHandler(w http.ResponseWriter, r *http.Request) {
	res := lib.Response{ResponseWriter: w}
	vars := mux.Vars(r)
	name := vars["name"]

	if len(name) == 0 {
		res.SendNotFound()
		return
	}

	articles := getDataFromSite(name)

	if len(articles) == 0 {
		res.SendBadRequest("Cannot crawl this page")
		return
	}

	res.SendOK(minimizeItems(articles))
}

// GetPostFeedHandler ...
func GetPostFeedHandler(w http.ResponseWriter, r *http.Request) {
	res := lib.Response{ResponseWriter: w}
	vars := mux.Vars(r)
	name := vars["name"]
	id := vars["id"]

	if len(name) == 0 || len(id) == 0 {
		res.SendNotFound()
		return
	}

	articles := getDataFromSite(name)

	index := findArticleByID(name, articles, id)

	if index == -1 {
		res.SendBadRequest("Cannot found that post")
		return
	}

	res.SendOK(articles[index])
}

// GetResolveGithubGist ...
func GetResolveGithubGist(w http.ResponseWriter, r *http.Request) {
	res := lib.Response{ResponseWriter: w}
	url := r.URL.Query().Get("url")

	resp, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}

	gistURL := resp.Request.URL.Scheme + "://" + resp.Request.URL.Host + resp.Request.URL.Path

	res.SendOK(gistURL)
}
