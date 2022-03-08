
package cmd

import (
	"fmt"
    "github.com/gocolly/colly/v2"
    "encoding/json"
	"io/ioutil"
    "os"
    "regexp"
	"github.com/spf13/cobra"
    "github.com/spf13/viper"
    "strings"
)

var addSiteUrl          string
var outDir             string
var selectors           []string
var options             Options
var sites               []Site
var site                Site
var siteExists          bool
var cnt                 int
var curr                Site

type Config struct {
    OutDir string `mapstructure:"OUT_DIR"`
}

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

type Options struct {
    Sites       []Site      `json:"sites"`
    Output      string      `json:"output"`
}

func init() {
    rootCmd.AddCommand(scrapeCmd)
	// Here are defined flags and configuration settings for commands
    scrapeCmd.Flags().StringVarP(&outDir, "output", "o", "assets/news_site_articles.json", "set your json output file")
	scrapeCmd.Flags().StringVarP(&addSiteUrl, "add-url", "a", "", "add new news source site url")
}

var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Command used for running news scraping scripts. The command will scrape all known sites if no flags are provided!",
    Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

        config, err := LoadConfig(".")
        outDir = config.OutDir
        if err != nil {
            fmt.Println("cannot load config:", err)
        }

        if addSiteUrl != "" {
            ReadFile(outDir)
            curr = Site{len(sites), addSiteUrl, make([]string, 0), make([]Article, 0)}
            sites = append(sites, site)
            Scrape(Options{sites, outDir})
        } else {
            ReadFile(outDir)
            Scrape(Options{sites, outDir})
        }
	},
}

func Scrape(options Options) {
    siteExists = false
    sites = options.Sites
    for _, site := range sites {
        curr = site
        articleNum := len(site.Articles)
        c := colly.NewCollector()
        
        selectors := curr.Selectors
        selector := FilterSelector(selectors[0])

        c.OnRequest(func(r *colly.Request) {
                CheckForValidSourceEntry(outDir)
                fmt.Printf("Visiting %s\n", r.URL)
            })
        
        c.OnHTML(selector, func(h *colly.HTMLElement) {
                selection := h.DOM
                newArticle := &Article{
                    ID: articleNum,
                    Title: selection.Find("a").Text(),
                    Buzzwords: make([]string, 0),
                }

                //fmt.Println(newArticle.Title)

                if !contains(curr.Articles, newArticle.Title) {
                    curr.Articles = append(curr.Articles, *newArticle);
                    articleNum++
                }
            })

        c.OnError(func(r *colly.Response, e error) {
            fmt.Printf("Error while scraping %s\n", e.Error())
        })

        options.Sites = sites
        c.Visit(curr.Url)
        WriteFile(options, curr)    
    }
}


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

func FilterUrl(url string) string { 
  re := regexp.MustCompile(`^(.+?)\/{2}`)
  url = strings.ReplaceAll(url, re.FindStringSubmatch(url)[0], "")
  if !strings.Contains(url, "www.") {
      return "www."+url
  }
  return url
}

func FilterSelector(selector string) string {
    selector = strings.ReplaceAll(selector, "\u003e", ">")
    arr := strings.Split(selector, ">")
    return strings.Join(arr[:len(arr)-1], ">")
}

func CheckForValidSourceEntry(source string) {
    ReadFile(source)
    for _, page := range sites {
        if page.Url == curr.Url {
            siteExists = true
            curr = page
            cnt = len(curr.Articles)
        } else {
            fmt.Println("Adding new site: ", curr.Url)
            curr.ID = len(sites)
        }
    }
}

func WriteFile(options Options, site Site) {
    if !siteExists { 
        options.Sites = append(options.Sites, site) 
    } else { 
        fmt.Println(options.Sites)
        fmt.Println(options.Sites[site.ID])
        options.Sites[site.ID] = site
    }
    dataBytes, err := json.MarshalIndent(sites, "", " ")

    if err != nil {
        fmt.Println(err)
    }

    err = ioutil.WriteFile(options.Output, dataBytes, 0644)
    if err != nil {
        fmt.Println(err)
    }
}

func ReadFile(fileName string) {
    err := CheckIfFileExists(fileName)

    if err != nil {
        fmt.Println(err)
    }

    file, err := ioutil.ReadFile(fileName)
    if err != nil {
        fmt.Println(err)
    }

    json.Unmarshal(file, &sites)
}

func CheckIfFileExists(filename string) error {
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

