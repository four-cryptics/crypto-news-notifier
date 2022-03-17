package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// Article serves as a single news article in a news source
type Article struct {
	ID                  int             `json:"id"`
    Title               string          `json:"title"`
}

// Site serves as a struct of each news source (wepage)
type Site struct {
    ID                  int             `json:"id"`
    Url                 string          `json:"url"`
	Selectors	        []string	    `json:"selectors"`
    Articles            []Article       `json:"articles"`
}

// ScrapeData is a struct of all the needed read/saved/generated site information for scraping and data checking
type ScrapeData struct {
    Sites               []Site          `json:"all_known_sites"`
    OutputFileName      string          `json:"output_file_name"`
    CurrentSite         Site            `json:"current_site"`
}

// ContainsArticle checks if current site already contains the exact same article
func (s *Site) ContainsArticle(title string) bool {
	for _, article := range s.Articles  {
		if article.Title == title {
            return true
		}
	}
	return false
}

// ScrapeAllSites scrapes all sites of data and stores it
func (sd *ScrapeData) ScrapeAllSites() {
    for i, site := range sd.Sites {
        sd.CurrentSite = site
        sd.CurrentSite.Scrape()
        sd.Sites[i] = sd.CurrentSite
    }
    sd.ExportSites()
}


// Scrape traverses the current site and scrapes the content of provided elements.
func (s *Site) Scrape() {
    articleNum := len(s.Articles)
    c := colly.NewCollector(
        colly.MaxDepth(1),
    )
    c.Limit(&colly.LimitRule{
        Delay: 2 * time.Second,
        RandomDelay: 2 * time.Second,
    })

    c.OnRequest(func(r *colly.Request) {
        // fmt.Printf("Visiting %s\n", r.URL)
    })

    c.OnResponse(func(r *colly.Response) {
        fmt.Println("Visited", r.Request.URL)
    })

    c.OnError(func(_ *colly.Response, err error) {
        fmt.Println("Something went wrong:", err)
    })

    for _, selector := range s.Selectors {
        selector, el := FilterSelector(selector)
        c.OnHTML(selector, func(e *colly.HTMLElement) {
            selection := e.DOM
            re := regexp.MustCompile(`(\s|\n|\r|\t)+`)
            title := re.ReplaceAllString(selection.Find(el).Text(), " ") 

            newArticle := &Article{
                ID: articleNum,
                Title: title,
            }

            if len(newArticle.Title) > 5 {
                if !s.ContainsArticle(newArticle.Title) {
                    s.Articles = append(s.Articles, *newArticle)
                    articleNum++
                    fmt.Println("ADDED! -> ", newArticle.Title)
                }
            }
        })
    }
    c.Visit(s.Url)
    if len(s.Articles) < 1 {
        fmt.Println("There seems to be something wrong with the provided selectors")
    }
    fmt.Println(s.Url, "now has: ", len(s.Articles), "articles.")
}

// EnsureOutputFile checks if the output file exists, otherwise creates it
func (sd *ScrapeData) EnsureOutputFile() error {
    _, err := os.Stat(sd.OutputFileName)
    if os.IsNotExist(err) {
        _, err := os.Create(sd.OutputFileName)
            if err != nil {
                fmt.Println(err)
            }
    }
    return nil
}

// ExportSites writes all article data to the pre-defined JSON output file
func (sd ScrapeData) ExportSites() {
    dataBytes, err := json.MarshalIndent(sd.Sites, "", " ")
    if err != nil {
        fmt.Println(err)
    }

    err = ioutil.WriteFile(sd.OutputFileName, dataBytes, 0644)
    if err != nil {
        fmt.Println(err)
    }
}

// Export outputs single site to file
func (s *Site) Export(sites []Site, OutputFile string) {
    sites = append(sites, *s)
    WriteToJsonFile(sites, OutputFile)
}

func WriteToJsonFile(sites []Site, OutputFile string) {
    dataBytes, err := json.MarshalIndent(sites, "", " ")
    if err != nil {
        fmt.Println(err)
    }

    err = ioutil.WriteFile(OutputFile, dataBytes, 0644)
    if err != nil {
        fmt.Println(err)
    }
}

// CreateSite creates a new valid news source
func (sd *ScrapeData) CreateSite(x int, url string, selector []string, articles []Article) Site {
    sd.CurrentSite = Site{x, url, selector, articles}
    sd.CurrentSite.Selectors = TakeUserElementSelectorInput()
    return sd.CurrentSite
}

// EnsurePageExists makes sure that the news webpage is added if it doesn't exist
func (sd ScrapeData) EnsurePageExists(url string) {
    if len(sd.Sites) > 0 {
        for _, page := range sd.Sites { 
            if page.Url == url {
                fmt.Println("Provided news source already exists!")
                sd.ScrapeAllSites()
            }
        }
        fmt.Println("New news source added!")
        page := sd.CreateSite(len(sd.Sites), url, []string{}, []Article{})
        page.Scrape()
        page.Export(sd.Sites, sd.OutputFileName)
    } else {
        fmt.Println("New news source added!")
        page := sd.CreateSite(0, url, []string{}, []Article{})
        page.Scrape()
        page.Export(sd.Sites, sd.OutputFileName)
    }
}

// ReadFile calls EnsureOutputFile to check for output file validity and reads site articles data
func (sd *ScrapeData) ReadFile() []Site{
    sd.EnsureOutputFile()

    file, err := ioutil.ReadFile(sd.OutputFileName)
    if err != nil {
        fmt.Println(err)
    }

    json.Unmarshal(file, &sd.Sites)
    return sd.Sites
}

// FilterUrl modifies URL address to fit the specified requirements (remove 'http://' etc.)
func FilterUrl(url string) string { 
  re := regexp.MustCompile(`^(.+?)\/{2}`)
  url = strings.ReplaceAll(url, re.FindStringSubmatch(url)[0], "")
  if !strings.Contains(url, "www.") {
      return "www."+url
  }
  return url
}

// FilterSelector converts characters and returns last element separate from the full selector path
func FilterSelector(selector string) (string, string) {
    selector = strings.ReplaceAll(selector, "\u003e", ">")
    arr := strings.Split(selector, ">")
    return strings.Join(arr[:len(arr)-1], ">"), strings.TrimSpace(arr[len(arr)-1])
}

// TakeUserElementSelectorInput takes user input element selectors one by one and adds them to the site (use with invoking -a or -add flag on scrape)
func TakeUserElementSelectorInput() []string {
    selectors := []string{}
    stop := false
    i := 1
    for !stop {
        scanner := bufio.NewScanner(os.Stdin)
        fmt.Println("Enter DOM selector! (press q to quit)")
        fmt.Print("Enter element number ",i," or press q to quit:")
        scanner.Scan()

        if(scanner.Text() == "q") {
            stop = true
            break;
        } else {
            selectors = append(selectors, scanner.Text())
            fmt.Println("Selector added!")
            i++
        }

        if scanner.Err() != nil {
            fmt.Println("Error: ", scanner.Err())
        }
    }
    if len(selectors) > 0 {
        return selectors
    } else {
        return []string{}
    }
}
