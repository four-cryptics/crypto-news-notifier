"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const puppeteer = require("puppeteer");
const timeout = 5000;
const file_utils_1 = require("../utils/file_utils");
const puppeteer_utils_1 = require("../utils/puppeteer_utils");
const Site_1 = require("../utils/Site");
const sites = [];
const site = new Site_1.Site('https://cryptonews.com/');
site.type = 'crypto-news';
site.validity = 2;
sites.push(site);
scrape(site);
function scrape(site) {
    return __awaiter(this, void 0, void 0, function* () {
        let siteList = yield (0, file_utils_1.readIfFileExists)();
        // console.log(site)
        let browser = yield (0, puppeteer_utils_1.setup)();
        const page = yield browser.newPage();
        yield page.setViewport({
            width: 1920,
            height: 1080,
            deviceScaleFactor: 1,
        });
        // Get webpage elements
        yield page.goto(site.url);
        console.log(site.url);
        const articles = yield page.evaluate(() => Array.from(document.querySelectorAll('body > main > section:nth-child(3) > div > div > div > section > article > div > div > a'), element => element.innerHTML));
        console.log(articles);
        // const quotes = await page.evaluate(() => Array.from(document.querySelectorAll('#article_content > p > strong'), element => element.innerHTML));
        // // const authors = await page.evaluate(() => Array.from(document.querySelectorAll('div > div > div.single-post-wrap.entry-content > p > strong'), element => element.innerHTML));
        // for(let i = 0; i < quotes.length; i++) {
        //   // console.log("Element: "+ elements[i]);
        //   let quote = quotes[i].replace(/“/, "").replace(/”/g, "");
        //   // Check if quote already exists
        //   let exists = false;
        //   for(let j = 0; j < numOfExistingArticles; j++) {
        //     if(siteList[j].quote.includes(quote)) { exists = true; }
        //   }
        //   if(!exists) {
        //     siteList[numOfExistingArticles+i] = new Article(quote, 'Undefined', 'unknown');
        //   }
        // }
        // articleList = 
        // fs.writeFile('../data/articles.json', JSON.stringify(siteList, null, 4), function (err) {
        //   if (err) throw err;
        //   console.log('File is created successfully.');
        // });
        yield page.close();
        yield browser.close();
    });
}
