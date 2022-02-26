const puppeteer = require("puppeteer");
// const puppeteer_extra = require("puppeteer-extra");
// const pluginStealth = require('puppeteer-extra-plugin-stealth');
// puppeteer_extra.use(pluginStealth())
// const iPhonex = puppeteer.devices['Pixel 4'];

export async function setup() {
  const browser = await puppeteer.launch({
    headless: false,
    devtools: true,
    ignoreHTTPSErrors: true,
    defaultViewport: {
      isMobile: true,
      width: 375,
      height: 667,
    }
  });
  //await page.emulate(iPhonex);
  return browser;
  // const pages = await browser.pages();
  // for (const p of pages) {
  //   // await p.emulate(iPhone)
  //   return p; 
  // }
}

export async function getActivePage(browser: any, timeout: number) {
        var start = new Date().getTime();
        while(new Date().getTime() - start < timeout) {
            var pages = await browser.pages();
            var arr = [];
            for (const p of pages) {
                if(await p.evaluate(() => { return document.visibilityState == 'visible' })) {
                    arr.push(p);
                }
            }
            if(arr.length == 1) return arr[0];
        }
        throw "Unable to get active page";
    }
    