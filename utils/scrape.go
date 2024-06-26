package utils

import (
	"autojob/models"
	"log"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

func scrapeJobData(scrapeData *[]models.ScrapeData, searchTerm string) {
	var pageUrls []string
	getPageUrls(searchTerm, &pageUrls)

	var jobUrls []string
	getJobUrls(&pageUrls, &jobUrls)

	getJobDetails(&jobUrls, scrapeData)
}

func getPageUrls(searchTerm string, pages *[]string) {
	c := colly.NewCollector()

	c.OnHTML("a[aria-label='Next']", func(e *colly.HTMLElement) {
		*pages = append(*pages, e.Request.URL.String())

		nextPage := e.Attr("href")
		if nextPage != "" {
			e.Request.Visit("https://www.seek.com.au" + nextPage)
		}
	})

	url := "https://www.seek.com.au/" + strings.ReplaceAll(searchTerm, " ", "-") + "-jobs/full-time?daterange=1"

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting to get page urls", r.URL.String())
	})

	c.Visit(url)
}

func getJobUrls(pageUrls *[]string, jobUrls *[]string) {

	var wg sync.WaitGroup

	for _, url := range *pageUrls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			c := colly.NewCollector()

			c.OnHTML("[data-automation='normalJob']", func(e *colly.HTMLElement) {
				route := e.ChildAttr("a[data-automation='jobTitle']", "href")
				link := "https://www.seek.com.au" + route

				*jobUrls = append(*jobUrls, link)
			})

			c.OnRequest(func(r *colly.Request) {
				log.Println("Visiting to get job urls", r.URL.String())
			})

			c.Visit(url)
		}(url)

		wg.Wait()
	}
}

func getJobDetails(jobUrls *[]string, jobDetails *[]models.ScrapeData) {

	var wg sync.WaitGroup

	for _, url := range *jobUrls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			c := colly.NewCollector()

			c.OnHTML("div[data-automation='jobDetailsPage']", func(e *colly.HTMLElement) {
				title := e.ChildText("[data-automation='job-detail-title']")
				company := e.ChildText("[data-automation='advertiser-name']")
				location := e.ChildText("[data-automation='job-detail-location']")
				description := e.ChildText("[data-automation='jobAdDetails']")
				formattedDescription := strings.ReplaceAll(description, "\n", "\\n")

				job := models.ScrapeData{
					Title:       title,
					Company:     company,
					Location:    location,
					Description: formattedDescription,
					Link:        e.Request.URL.String(),
				}

				*jobDetails = append(*jobDetails, job)
			})

			c.OnRequest(func(r *colly.Request) {
				log.Println("Visiting to get job details", r.URL.String())
			})

			c.Visit(url)
		}(url)

		wg.Wait()
	}
}
