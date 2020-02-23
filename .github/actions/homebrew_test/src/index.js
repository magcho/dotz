const util = require("util");
const childProcess = require("child_process");
const exec = util.promisify(childProcess.exec);
const core = require("@actions/core");

async function main() {
  await console.log("pwd");
  await exec("sh .github/actions/homebrew_test/main.sh");
}

main().catch(err => core.setFailed(err.message));
