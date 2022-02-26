"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Site = void 0;
const Article_1 = __importDefault(require("./Article"));
class Site {
    constructor(link) {
        this.url = link;
        this.type = '';
        this.validity = 0;
        this.articles = [];
    }
    addArticle(title, news) {
        this.articles[title] = new Article_1.default(news);
    }
}
exports.Site = Site;
exports.default = Site;
