# Crypto-News-Notifier

![crypto](assets/crypto_thumbnail.png)

## About
Crypto-News-Notifier, is a platform made by a crafty team of friends and colleagues:
- **@MatejBasa** :ok_man:
- **@CraftyIvan** :person_in_tuxedo:
- **@Leamoonmoon** :vampire_woman:
- **@YuriTheRed** :mage_man:

The platform was developed solely in the spirit of education. It's ment for gathering crypto data and making predictions based on the gathered data and selected prediction models.

It currently provides below mentioned contents:
- Easy to use scripts for large web data scraping/collecting,

## How to use

Any kind of CLI interaction is done trough the ./cmd/ folder, using the Cobra package which provides easy-to-use CLI.

For now it provides a couple of data gathering commands, for example:
- Add a new news source for scraping by calling `go run main.go scrape -a/-add <news_url>`,
- Scrape all known sites of their articles and add them to the known ones by calling `go run main.go scrape`

**Note** When you're adding a new site to the known news sources, you'll be asked to provide the DOM element selectors of the news articles (you can provide multiple if needed)