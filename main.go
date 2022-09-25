package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type Repos struct {
	Repos []Repo `json:"repos"`
}

type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
		})
	})

	r.GET("/github/*language", func(c *gin.Context) {

		// Set fallback values
		var timeFrame string
		if timeFrame = c.Query("since"); timeFrame == "" {
			timeFrame = "daily"
		}

		repos := getRepos(c.Param("language"), timeFrame)

		c.JSON(http.StatusOK, repos)
	})

	r.Run()
}

func getRepos(language, timeFrame string) Repos {
	c := colly.NewCollector()

	var Repos Repos
	c.OnHTML(".Box-row", func(e *colly.HTMLElement) {
		name := e.ChildAttr(".lh-condensed a", "href")
		name, _ = url.PathUnescape(name)
		name = strings.Replace(name, "/", "", 1)

		description := e.ChildText("p")

		Repos.Repos = append(Repos.Repos, Repo{
			Name:        name,
			Description: description,
			URL:         "https://github.com/" + name,
		})
	})

	c.Visit("https://github.com/trending/" + language + "?since=" + timeFrame)

	return Repos
}
