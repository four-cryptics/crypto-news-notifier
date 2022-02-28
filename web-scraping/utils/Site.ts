import Article from './Article';

export class Site {
    public url: string;
    public type: string;
    public validity: number;
    public articles: Article[];
    public elements: string[];

    constructor(link: string) {
        this.url = link;
        this.type = '';
        this.validity = 0;
        this.articles = []
        this.elements = [];
    } 

    public addArticle(title: string) {
        this.articles.push(new Article(title));
    }

    public addElement(elementSelector: string) {
        this.elements.push(elementSelector);
    }

    public getElements() {
        return this.elements;
    }

    public getElementsLength() {
        return this.elements.length;
    }
}

export default Site;