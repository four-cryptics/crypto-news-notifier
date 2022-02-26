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
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.fixNumbering = exports.writeFile = exports.readIfFileExists = void 0;
const fs_1 = __importDefault(require("fs"));
function readIfFileExists() {
    return __awaiter(this, void 0, void 0, function* () {
        if (fs_1.default.existsSync('./data/news_sites.json')) {
            console.log("A file for sites exists");
            return JSON.parse(fs_1.default.readFileSync('./data/news_sites.json', 'utf-8'));
        }
        else {
            console.log("Sites file doesn't exist yet.");
            return {};
        }
    });
}
exports.readIfFileExists = readIfFileExists;
function writeFile(filePath, value) {
    return __awaiter(this, void 0, void 0, function* () {
        return new Promise((resolve, reject) => {
            fs_1.default.writeFile(filePath, value, 'utf-8', resolveOrRejectIfError(resolve, reject));
        });
    });
}
exports.writeFile = writeFile;
function fixNumbering(articleList) {
    let length = Object.keys(articleList).length;
    console.log(length);
    for (let i = 0; i < length; i++) {
        if (articleList[i] === undefined) {
            console.log("undefined: " + i);
            let repaired = false;
            for (let j = i; j < 100000; j++) {
                if (!repaired) {
                    if (articleList[j] !== undefined) {
                        articleList[i] = articleList[j];
                        delete articleList[j];
                        console.log("repaired " + i + " as " + j);
                        repaired = true;
                    }
                }
            }
        }
    }
    return articleList;
}
exports.fixNumbering = fixNumbering;
function resolveOrRejectIfError(resolve, reject) {
    return (err) => {
        if (err) {
            reject(err);
        }
        resolve();
    };
}
