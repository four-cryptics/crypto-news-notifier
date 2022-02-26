import fs from 'fs'
import Site from './Site'
import Article from './Article'

export async function readIfFileExists()  {
  if (fs.existsSync('./data/news_sites.json')) { 
    console.log("A file for sites exists");
    return JSON.parse(fs.readFileSync('./data/news_sites.json', 'utf-8'));
  } else { 
    console.log("Sites file doesn't exist yet.");
    return {};
  }
}

export async function writeFile(filePath: string, value: string): Promise<void> {
  return new Promise((resolve, reject) => {
    fs.writeFile(filePath, value, 'utf-8', resolveOrRejectIfError(resolve,reject))
  })
}

export function fixNumbering(articleList: Article[]) {
  let length = Object.keys(articleList).length
  console.log(length)
  for(let i = 0; i < length; i++) {
    if(articleList[i] === undefined) {
      console.log("undefined: "+i)
      let repaired = false;
      for(let j = i; j < 100000; j++) {
        if(!repaired) {
          if(articleList[j] !== undefined) {
            articleList[i] = articleList[j]
            delete articleList[j]
            console.log("repaired " + i + " as " + j)
            repaired = true;
          }
        }
      }
    }
  }
  return articleList
}

function resolveOrRejectIfError(
  resolve: (value: void | PromiseLike<void>) => void,
  reject: (reason?: any) => void
): fs.NoParamCallback {
  return(err) => {
    if(err) {
      reject(err);
    }
    resolve()
  }
}