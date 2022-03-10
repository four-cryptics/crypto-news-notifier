package cmd

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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var options             Options

type Article struct {
	ID                  int             `json:"id"`
    Title               string          `json:"title"`
}

type Site struct {
    ID                  int             `json:"id"`
    Url                 string          `json:"url"`
	Selectors	        []string	    `json:"selectors"`
    Articles            []Article       `json:"articles"`
}

// Check if current site already contains the exact same article
func (s Site) containsArticle(title string) bool {
	for _, article := range s.Articles  {
		if article.Title == title {
            return true
		}
	}
	return false
}

// Prints all present site data
func (s Site) Print() {
    fmt.Println("----------------------------------")
    fmt.Println(s.Url, "ALL SITE ARTICLES:")
    fmt.Println("")
    for i, article := range s.Articles {
        fmt.Println(i, " --> ", article.Title)
    }
    fmt.Println("----------------------------------")
}

type Options struct {
    Sites               []Site          `json:"sites"`
    Output              string          `json:"output"`
    CurrentSite         Site            `json:"current_site"`
    IdNum               int             `json:"id_num"`
    AddedSite           Site            `json:"added_site`
}

// Check if file exists, otherwise create it
func (o Options) CheckIfFileExists() error {
    _, err := os.Stat(o.Output)
    if os.IsNotExist(err) {
        _, err := os.Create(o.Output)
            if err != nil {
                fmt.Println(err)
            }
    }
    return nil
}

// This gets executed automatically when rootCmd.Execute() is called from main.go, due to it being a subcommand
func init() {
    rootCmd.AddCommand(scrapeCmd)
	// Here are defined flags and configuration settings for commands
    scrapeCmd.Flags().StringVarP(&options.Output, "output", "o", "assets/news_site_articles.json", "set your json output file")
	scrapeCmd.Flags().StringVarP(&options.AddedSite.Url, "add-url", "a", "", "add new news source site url")
}

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Command used for running news scraping scripts. The command will scrape all known sites if no flags are provided!",
    Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
        config, err := LoadConfig(".")
        if err != nil {
            fmt.Println("cannot load config:", err)
        }

        options = Options{options.ReadFile(), config.OutDir, Site{}, 0, Site{0, options.AddedSite.Url, []string{}, []Article{}}}

        if options.AddedSite.Url != "" {
            exists, page := options.DoesPageExist(options.AddedSite.Url)
            if exists {
                fmt.Println("The news source already exists!")
                options.CurrentSite = page
                page = page.Scrape()
                options.Sites[page.ID] = page
            } else {
                fmt.Println("A new news source!")
                options.CurrentSite = page
                page = page.Scrape()
                options.Sites = append(options.Sites, page)
            }
        } else {
            for _, site := range options.Sites {
                site = site.Scrape()
                options.Sites[site.ID] = site
            }
        }
        options.ExportSites()
	},
}

// Write to the JSON output file for news article data (predefined but an be changed using the CLI)
func (o Options) ExportSites() {
    dataBytes, err := json.MarshalIndent(o.Sites, "", " ")
    if err != nil {
        fmt.Println(err)
    }

    err = ioutil.WriteFile(o.Output, dataBytes, 0644)
    if err != nil {
        fmt.Println(err)
    }
}

func (o Options) CreateSite(x int, url string, selector []string, articles []Article) Site {
    o.CurrentSite = Site{x, url, selector, articles}
    return o.CurrentSite
}

// Makes sure the news webpage is added if it doesn't exist
func (o Options) DoesPageExist(url string) (bool, Site) {
    site := Site{}
    if len(o.Sites) > 0 {
        for _, page := range o.Sites { 
            if page.Url == url {
                return true, page
            }
        }
        site = o.CreateSite(len(o.Sites), url, []string{}, []Article{})
        site.Selectors = TakeUserElementSelectorInput()
        return false, site
    }
    site = o.CreateSite(0, url, []string{}, []Article{})
    site.Selectors = TakeUserElementSelectorInput()
    return false, site
}

func (o Options) ScrapeAll() {
    for _, site := range o.Sites {
        site = site.Scrape()
    }
}

func (s Site) Scrape() Site {
    articleNum := len(s.Articles)
    c := colly.NewCollector(
        colly.MaxDepth(1),
    )
    c.Limit(&colly.LimitRule{
        Delay: 2 * time.Second,
        RandomDelay: 2 * time.Second,
    })

    c.OnRequest(func(r *colly.Request) {
        fmt.Printf("Visiting %s\n", r.URL)
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
                if !s.containsArticle(newArticle.Title) {
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
    return s
}

// Checks for news article output file validity and reads site articles data
func (o Options) ReadFile() []Site{
    o.CheckIfFileExists()

    file, err := ioutil.ReadFile(o.Output)
    if err != nil {
        fmt.Println(err)
    }

    json.Unmarshal(file, &o.Sites)
    return o.Sites
}

// Filters URL address to fit the specified requirements (remove 'http://' etc.)
func FilterUrl(url string) string { 
  re := regexp.MustCompile(`^(.+?)\/{2}`)
  url = strings.ReplaceAll(url, re.FindStringSubmatch(url)[0], "")
  if !strings.Contains(url, "www.") {
      return "www."+url
  }
  return url
}

// Converts characters and return last element separate from the full selector path
func FilterSelector(selector string) (string, string) {
    selector = strings.ReplaceAll(selector, "\u003e", ">")
    arr := strings.Split(selector, ">")
    return strings.Join(arr[:len(arr)-1], ">"), strings.TrimSpace(arr[len(arr)-1])
}

// Reads app.env environmental config
func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("app")
    viper.SetConfigType("env")
    viper.AutomaticEnv()

    err = viper.ReadInConfig()
    if err != nil {
        return
    }

    err = viper.Unmarshal(&config)
    return
}

// Let's user input selectors one by one with invoking -a or -add flag on scrape
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

// Logs all known site articles data
func (o Options) PrintAllSites() {   
    for _, site := range o.Sites {
        site.Print()
    }
}

// Logs all articles that were scraped and added when this script is ran
func (o Options) PrintAllAddedArticles() {
    fmt.Println("----------------------------------")
    if len(o.AddedSite.Articles) > 0 {
        fmt.Println(o.AddedSite.Url, "ALL ", len(o.AddedSite.Articles)," ADDED ARTICLES:")
        for _, article := range o.AddedSite.Articles {
            fmt.Println(article.Title)
        }
    } else {
        fmt.Println("0 ARTICLES WERE ADDED")
    }
    fmt.Println("----------------------------------")
}