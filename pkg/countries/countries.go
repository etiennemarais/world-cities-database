package countries

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"go.uber.org/zap"
)

type Country struct {
	Name    string
	Code    string
	Regions []Region
}

type Region struct {
	ID   string
	Name string
}

func GetAll(logger *zap.Logger) []Country {
	listView := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.AllowedDomains("www.worldcitiesdb.com"),
		colly.Async(true),
	)

	listView.Limit(&colly.LimitRule{
		Parallelism: 1,
		RandomDelay: 30 * time.Second,
	})

	singleView := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.AllowedDomains("www.worldcitiesdb.com"),
	)
	singleView.Limit(&colly.LimitRule{
		Parallelism: 5,
		Delay:       1 * time.Second / 3,
		RandomDelay: 2 * time.Second,
	})

	stateView := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.AllowedDomains("www.worldcitiesdb.com"),
	)
	stateView.Limit(&colly.LimitRule{
		Parallelism: 5,
		Delay:       1 * time.Second / 2,
		RandomDelay: 2 * time.Second,
	})

	countries := []Country{}

	listView.OnHTML("div.maincontent", func(e *colly.HTMLElement) {
		e.DOM.Find("table tbody tr td:nth-child(2) a").Each(func(i int, s *goquery.Selection) {
			url, exists := s.Attr("href")
			if exists {
				singleView.Visit(fmt.Sprintf("http://www.worldcitiesdb.com%s", url))
			}
		})
	})

	// Get the country name and code
	singleView.OnHTML("div.maincontent.alignc", func(c *colly.HTMLElement) {
		// Country Code
		countryCode := trim(c.ChildText("table tr:nth-child(1) td:nth-child(2)"))
		countryName := c.ChildText("table tr:nth-child(3) td:nth-child(2) a")
		country := Country{
			Code: countryCode,
			Name: escape(countryName),
		}

		countries = append(countries, country)
	})

	singleView.OnHTML("div.infolink.alignc", func(i *colly.HTMLElement) {
		// Find the link to the provinces/state list
		url := i.ChildAttr("a.minfo[href*=\"state\"]", "href")
		stateView.Visit(fmt.Sprintf("http://www.worldcitiesdb.com%s", url))
	})

	// Get the state/province name and id
	stateView.OnHTML("#content", func(p *colly.HTMLElement) {
		// Find the country from the breadcrumbs
		countryName := escape(p.DOM.Find("ul li:nth-child(3) a").Text())
		regions := []Region{}
		p.DOM.Find("table tbody tr").Each(func(i int, s *goquery.Selection) {
			region := Region{
				ID:   trim(s.Find("td:nth-child(2)").Text()),
				Name: escape(s.Find("td:nth-child(3) a").Text()),
			}
			regions = append(regions, region)
		})

		// Ignore regions that are empty like islands
		if len(regions) > 0 {
			// Write the regions to the country data as it comes in
			for ic, country := range countries {
				if countryName == country.Name {
					countries[ic].Regions = regions
				}
			}
		}
	})

	// Before making a request print "Visiting ..."
	listView.OnRequest(func(r *colly.Request) {
		logger.Info(fmt.Sprintf("List View Visiting %s", r.URL.String()))
	})

	singleView.OnRequest(func(r *colly.Request) {
		logger.Info(fmt.Sprintf("Single View Visiting %s", r.URL.String()))
	})

	stateView.OnRequest(func(r *colly.Request) {
		logger.Info(fmt.Sprintf("State View Visiting %s", r.URL.String()))
	})

	listView.OnResponse(func(r *colly.Response) {
		logger.Info(fmt.Sprintf("Response received %d", r.StatusCode))
	})

	listView.OnError(func(r *colly.Response, err error) {
		fmt.Println("List View Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	singleView.OnError(func(r *colly.Response, err error) {
		fmt.Println("Single View Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	listView.Visit("http://www.worldcitiesdb.com/country/list")

	listView.Wait()

	return countries
}

func trim(str string) string {
	return strings.Replace(str, " ", "", -1)
}

func escape(str string) string {
	return strings.Replace(str, "'", "’", -1)
}
