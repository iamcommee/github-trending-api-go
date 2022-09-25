Github Trending API with Go
---

This repository is created by Go with [go colly](http://go-colly.org/) for scraping github trending as API.

URL Endpoint : 

```
https://iamcommee-trending-api.herokuapp.com/github/{programing_language}?since={daily|weekly|monthly}
```

Example : 

- https://iamcommee-trending-api.herokuapp.com/github/ (If `since` param is empty, it will use daily as a default)
- https://iamcommee-trending-api.herokuapp.com/github/?since=weekly
- https://iamcommee-trending-api.herokuapp.com/github/go?since=monthly

Response :

```
{
    "repos": [
        {
            "owner": string,
            "repository_name": string,
            "repository ": string,
            "description": string,
            "programing_language": string,
            "url": string,
            "stars": number,
            "forks": number,
            "time_frame_stars": number
        },
        ...
    ]
}
```