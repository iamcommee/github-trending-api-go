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
	Owner              string `json:"owner"`
	RepositoryName     string `json:"repository_name"`
	Repository         string `json:"repository "`
	Description        string `json:"description"`
	ProgramingLanguage string `json:"programing_language"`
	URL                string `json:"url"`
	Stars              int    `json:"stars"`
	Forks              int    `json:"forks"`
	TimeFrameStars     int    `json:"time_frame_stars"`
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
		description := e.ChildText("p")
		programingLanguage := e.ChildText("span[itemprop='programmingLanguage']")

		starHtmlAttr := fmt.Sprintf(`a[href='%s']`, "/"+repo+"/stargazers")
		starHtmlValue := e.ChildText(starHtmlAttr)
		stars, err := strconv.Atoi(strings.Replace(starHtmlValue, ",", "", -1))

		if err != nil {
			log.Println(err)
		}

		forkHtmlAttr := fmt.Sprintf(`a[href='%s']`, "/"+repo+"/network/members"+"."+repoName)
		forkHtmlValue := e.ChildText(forkHtmlAttr)
		forks, err := strconv.Atoi(strings.Replace(forkHtmlValue, ",", "", -1))

		if err != nil {
			log.Println(err)
		}

		// example : 9,999 stars today
		timeFrameStarHtmlValue := e.ChildText("span[class='d-inline-block float-sm-right']")
		timeFrameStarSplitter := strings.Split(timeFrameStarHtmlValue, " ")
		timeFrameStars, err := strconv.Atoi(strings.Replace(timeFrameStarSplitter[0], ",", "", -1))

		if err != nil {
			log.Println(err)
		}

		Repos.Repos = append(Repos.Repos, Repo{
			Owner:              owner,
			RepositoryName:     repoName,
			Repository:         repo,
			Description:        description,
			ProgramingLanguage: programingLanguage,
			URL:                "https://github.com/" + repo,
			Stars:              stars,
			Forks:              forks,
			TimeFrameStars:     timeFrameStars,
		})
	})

	c.Visit("https://github.com/trending/" + language + "?since=" + timeFrame)

	return Repos
}
