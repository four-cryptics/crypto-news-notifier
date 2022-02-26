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
exports.getActivePage = exports.setup = void 0;
const puppeteer = require("puppeteer");
// const puppeteer_extra = require("puppeteer-extra");
// const pluginStealth = require('puppeteer-extra-plugin-stealth');
// puppeteer_extra.use(pluginStealth())
// const iPhonex = puppeteer.devices['Pixel 4'];
function setup() {
    return __awaiter(this, void 0, void 0, function* () {
        const browser = yield puppeteer.launch({
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
    });
}
exports.setup = setup;
function getActivePage(browser, timeout) {
    return __awaiter(this, void 0, void 0, function* () {
        var start = new Date().getTime();
        while (new Date().getTime() - start < timeout) {
            var pages = yield browser.pages();
            var arr = [];
            for (const p of pages) {
                if (yield p.evaluate(() => { return document.visibilityState == 'visible'; })) {
                    arr.push(p);
                }
            }
            if (arr.length == 1)
                return arr[0];
        }
        throw "Unable to get active page";
    });
}
exports.getActivePage = getActivePage;
