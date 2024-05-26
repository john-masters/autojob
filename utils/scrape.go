package utils

import (
	"autojob/models"
	"strings"

	"github.com/gocolly/colly"
)

func scrapeJobData(jobs *[]models.Job, searchTerm string) {
	c := colly.NewCollector()

	c.OnHTML("[data-card-type='JobCard']", func(e *colly.HTMLElement) {
		title := e.ChildText("a[data-automation='jobTitle']")
		company := e.ChildText("a[data-automation='jobCompany']")
		link := e.ChildAttr("a[data-automation='jobTitle']", "href")

		fullLink := "https://www.seek.com.au" + link

		job := models.Job{
			Title:   title,
			Company: company,
			Link:    fullLink,
		}

		*jobs = append(*jobs, job)

		e.Request.Visit(fullLink)
	})

	c.OnHTML("div[data-automation='jobAdDetails']", func(e *colly.HTMLElement) {
		description := e.Text

		for i := range *jobs {
			if (*jobs)[i].Link == e.Request.URL.String() {
				(*jobs)[i].Description = description
				break
			}
		}
	})

	c.OnHTML("a[aria-label='Next']", func(e *colly.HTMLElement) {
		nextPage := e.Attr("href")
		if nextPage != "" {
			e.Request.Visit("https://www.seek.com.au" + nextPage)
		}
	})

	url := "https://www.seek.com.au/" + strings.ReplaceAll(searchTerm, " ", "-") + "-jobs/full-time?daterange=1"

	c.Visit(url)
}
