
package cmd

import (
	"fmt"
    "github.com/gocolly/colly/v2"
    "encoding/json"
	"io/ioutil"
    "os"
	"github.com/spf13/cobra"
)

type Article struct {
	ID          int         `json:"id"`
    Title       string      `json:"title"`
    Buzzwords   []string    `json:"buzzwords"`
}

type Site struct {
    ID          int         `json:"id"`
    Url         string      `json:"url"`
	Selectors	[]string	`json:"selectors"`
    Articles    []Article   `json:"articles"`
}

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Run news scraping scripts",
	Run: func(cmd *cobra.Command, args []string) {
		scrape()
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scraperCmd.PersistentFlags().String("foo", "", "A help for foo")
	scrapeCmd.Flags().String("add", "a", "Add new news source site")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scraperCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func scrape() {
    scrapeUrl := "https://cryptonews.com/"
    articleDataFilePath := "assets/news_site_articles.json"
	selector1 := "section:nth-child(3) > div > div > div > section > article > div > div.justify-content-center"

    c := colly.NewCollector(colly.AllowedDomains("cryptonews.com", "www.cryptonews.com", ))

    sites := []Site{}
    siteExists := false
    var site = Site{0, scrapeUrl, []Article{}, []string{}}
    cnt := 0

    c.OnHTML(, func(h *colly.HTMLElement) {
        selection := h.DOM
        newArticle := &Article{
            ID: cnt,
            Title: selection.Find("a").Text(),
            Buzzwords: []string{"a", "e", "i", "o", "u"},
        }

        if !contains(site.Articles, newArticle.Title) {
            site.Articles = append(site.Articles, *newArticle);
            cnt++
        }
    })

    c.OnRequest(func(r *colly.Request) {
        err := checkFile(articleDataFilePath)

        if err != nil {
            fmt.Println(err)
        }

        file, err := ioutil.ReadFile(articleDataFilePath)
        if err != nil {
            fmt.Println(err)
        }

        json.Unmarshal(file, &sites)
        for _, x := range sites {
            if x.Url == site.Url {
                siteExists = true
                site = x
                cnt = len(site.Articles)
            } else {
                fmt.Println("Adding new site: ", site.Url)
                site.ID = len(sites)
            }
        }
        fmt.Printf("Visiting %s\n", r.URL)
    })

    c.OnError(func(r *colly.Response, e error) {
        fmt.Printf("Error while scraping %s\n", e.Error())
    })

    c.Visit(scrapeUrl)

    if !siteExists { 
        sites = append(sites, site) 
    } else { 
        sites[site.ID] = site 
    }
    dataBytes, err := json.MarshalIndent(sites, "", " ")

    if err != nil {
        fmt.Println(err)
    }

    err = ioutil.WriteFile(articleDataFilePath, dataBytes, 0644)
    if err != nil {
        fmt.Println(err)
    }
}

func checkFile(filename string) error {
    _, err := os.Stat(filename)
        if os.IsNotExist(err) {
            _, err := os.Create(filename)
                if err != nil {
                    return err
                }
        }
        return nil
}

func contains(a []Article, title string) bool {
	for _, v := range a {
		if v.Title == title {
			return true
		}
	}

	return false
}

