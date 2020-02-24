const util = require("util");
const childProcess = require("child_process");
const exec = util.promisify(childProcess.exec);
const core = require("@actions/core");

async function main() {
  // await exec("sh .github/actions/homebrew_test/main.sh");
  // https://github.com/magcho/homebrew-magcho.git
  const formulaName = core
    .getInput("formula_url")
    .replace("https://github.com/", "")
    .replace(/\/homebrew-/, "/")
    .replace(/\/$/, "")
    .replace(/.git$/, "");
  const formulaFilePath = core.getInput("formula_file_path");
  const commandName = core.getInput("command_name");

  await exec(
    '/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"'
  );
  await exec(`brew tap ${formulaName}`);
  await exec("brew update");
  const formulaPath = await exec(`brew --repository ${formulaName}`).stdout;
  await exec(`cp ${formulaFilePath} ${formulaPath}`);

  await exec(`brew audit ${commandName} --fix`);
  await exec(`brew style ${formulaFilePath} --fix`);

  await exec(`brew install ${commandName}`);

  await exec(`which ${commandName}`);
  await exec(`brew rm ${commandName}`);
  return;
}

main().catch(err => core.setFailed(err.message));
