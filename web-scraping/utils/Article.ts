export class Article {
  private news: string;
  private author: undefined | string;
  constructor(news: string) {
    this.news = news
    this.author = undefined
  }
}

export default Article;