const puppeteer = require("puppeteer");
import fs from 'fs'
const timeout = 5000;
import { readIfFileExists, fixNumbering } from '../utils/file_utils'
import { setup } from '../utils/puppeteer_utils'
import { Site } from '../utils/Site'
import { Article } from '../utils/Article'

const sites: Site[] = []
const site = new Site('https://cryptonews.com/')
site.type = 'crypto-news'
site.validity = 2
sites.push(site)

scrape(site)

async function scrape(site: Site) {
  let siteList = await readIfFileExists();
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
  console.log(site.url)
  const articles = await page.evaluate(() => Array.from(document.querySelectorAll('body > main > section:nth-child(3) > div > div > div > section > article > div > div > a'), element => element.innerHTML));
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
  await page.close();
  await browser.close();
}
