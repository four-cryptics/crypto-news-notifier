import Article from './Article';

export class Site {
    public url: string;
    public type: string;
    public validity: number;
    public articles: Article[];

    constructor(link: string) {
        this.url = link;
        this.type = '';
        this.validity = 0;
        this.articles = []
    } 

    addArticle(title: any, news: string) {
        this.articles[title] = new Article(news);
    }
}

export default Site;