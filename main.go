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
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func main() {
	r := gin.Default()
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
		title := e.ChildAttr("a", "href")
		title, _ = url.PathUnescape(title)
		title = strings.Replace(title, "/login?return_to=", "", -1)
		title = strings.Replace(title, "/", "", 1)

		description := e.ChildText("p")

		Repos.Repos = append(Repos.Repos, Repo{
			Title:       title,
			Description: description,
			URL:         "https://github.com/" + title,
		})
	})

	c.Visit("https://github.com/trending/" + language + "?since=" + timeFrame)

	return Repos
}
