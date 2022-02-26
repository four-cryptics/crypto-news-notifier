const puppeteer = require("puppeteer");
import fs from 'fs'
const timeout = 5000;
import { readIfFileExists, fixNumbering } from '../utils/file_utils'
import { setup } from '../utils/puppeteer_utils'
import { Site } from '../utils/Site'
import { Article } from '../utils/Article'

//let siteList = await readIfFileExists();
const sites: Site[] = []
const site1 = new Site('https://cryptonews.com/')
site1.addElement("body > main > section:nth-child(3) > div > div > div.col-12.col-lg-9 > section > article > div > div.col-12.col-md-7.column-45__right.d-flex.flex-column.justify-content-center > a")
site1.type = 'crypto-news'
site1.validity = 2

scrape(site1)

async function scrape(site: Site) {
  // console.log(site)
  let browser = await setup();
  const page = await browser.newPage();
  await page.setViewport({
    width: 1920,
    height: 1080,
    deviceScaleFactor: 1,
  });
  
  // Get webpage elements
  await page.goto(site.url);
  const articles = await page.evaluate((site: Site) => Array.from(document.querySelectorAll(site.elements[0]), element => element.innerHTML), site);
  console.log(articles);
  for(let article of articles) {
    site.addArticle(article);
  }
  sites.push(site);
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
  fs.writeFile('./data/news_sites.json', JSON.stringify(sites, null, 4), function (err) {
    if (err) throw err;
    console.log('File is created successfully.');
  });
  await page.close();
  await browser.close();
}
