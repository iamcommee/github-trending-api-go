package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type Repos struct {
	Repos []Repo `json:"repos"`
}

type Repo struct {
	Owner          string `json:"owner"`
	RepositoryName string `json:"repository_name"`
	Repository     string `json:"repository "`
	Description    string `json:"description"`
	URL            string `json:"url"`
	Star           int    `json:"star"`
	Fork           int    `json:"fork"`
	TimeFrameStar  int    `json:"time_frame_star"`
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

		repoHtmlAttr := e.ChildAttr(".lh-condensed a", "href")
		repoHtmlValue, _ := url.PathUnescape(repoHtmlAttr)

		// 0 = white space
		// 1 = owner
		// 2 = repository name
		repoSplitter := strings.Split(repoHtmlValue, "/")

		owner := repoSplitter[1]
		repoName := repoSplitter[2]
		repo := owner + "/" + repoName

		starHtmlAttr := fmt.Sprintf(`a[href='%s']`, "/"+repo+"/stargazers")
		forkHtmlAttr := fmt.Sprintf(`a[href='%s']`, "/"+repo+"/network/members"+"."+repoName)

		description := e.ChildText("p")

		starHtmlValue := e.ChildText(starHtmlAttr)
		star, err := strconv.Atoi(strings.Replace(starHtmlValue, ",", "", -1))

		if err != nil {
			log.Println(err)
		}

		forkHtmlValue := e.ChildText(forkHtmlAttr)
		fork, err := strconv.Atoi(strings.Replace(forkHtmlValue, ",", "", -1))

		if err != nil {
			log.Println(err)
		}

		timeFrameStarHtmlValue := e.ChildText("span[class='d-inline-block float-sm-right']")

		timeFrameStarSplitter := strings.Split(timeFrameStarHtmlValue, " ")

		timeFrameStar, err := strconv.Atoi(strings.Replace(timeFrameStarSplitter[0], ",", "", -1))

		if err != nil {
			log.Println(err)
		}

		Repos.Repos = append(Repos.Repos, Repo{
			Owner:          owner,
			RepositoryName: repoName,
			Repository:     repo,
			Description:    description,
			URL:            "https://github.com/" + repo,
			Star:           star,
			Fork:           fork,
			TimeFrameStar:  timeFrameStar,
		})

	})

	c.Visit("https://github.com/trending/" + language + "?since=" + timeFrame)

	return Repos
}
