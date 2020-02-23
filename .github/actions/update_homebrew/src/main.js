const core = require("@actions/core");
const github = require("@actions/github");

const fetch = require("node-fetch");
const fs = require("fs");
const checksum = require("checksum");

function main() {
  const formulaUrl = core
    .getInput("formula_url")
    .replace(/\/$/, "")
    .replace(/.git$/, "");
  const formulaFileName = core.getInput("formula_file_name");
  const assetPath = core.getInput("asset_path");

  const version = core.getInput("tag_name").replace(/([\d.]+$)/, "$1");
  let sha256;
  check256(assetPath).then(sum => {
    sha256 = sum;
  });

  const url = `https://raw.githubusercontent.com/${formulaUrl.replace(
    "https://github.com/",
    ""
  )}/master/${formulaFileName}`;

  fetch(url)
    .then(res => res.text())
    .then(formula =>
      formula
        .split("\n")
        .map(line =>
          line
            .replace(/^(\s+)version\s+"([\d\.]+)"/, `$1version "${version}"`)
            .replace(/^(\s+)sha256\s+"([a-z\d\.]+)"/, `$1sha256 "${sha256}"`)
        )
    )
    .then(contents => {
      fs.writeFileSync(`${formulaFileName}`, contents.join("\n"));
    })
    .then(() => core.setOutput("formula_file_path", formulaFileName))
    .catch(err => {
      core.setFailed(err.message);
    });
}

/**
 * ファイルのチェックサムを計算
 * @arg filePath ファイルパス
 * @return string sha256
 */
function check256(filePath) {
  return new Promise((resolve, reject) => {
    checksum.file(filePath, (err, sum) => {
      resolve(sum);
    });
  });
}

try {
  main();
} catch (err) {
  core.setFailed(err.message);
}
