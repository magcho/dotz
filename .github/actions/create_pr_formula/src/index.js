const util = require("util");
const childProcess = require("child_process");
const exec = util.promisify(childProcess.exec);
const core = require("@actions/core");
const fs = require("fs");

function setAuth(userName, pass) {
  const text = `machine github.com
login ${userName}
password ${pass}
`;
  fs.writeFile("~/.netrc", text, err => core.setFailed(err.message));
}

async function main() {
  const input = {
    formulaPath: core.getInput("formula_path"),
    githubSecretsToken: core.getInput("github_secrets_token"),
    githubUserName: core.getInput("github_username"),
    formulaFilePath: core.getInput("formula_file_path"),
    authorName: core.getInput("author_name"),
    authorEmail: core.getInput("author_email"),
    commitMessage: core.getInput("commit_message")
  };

  setAuth(input.githubUserName, input.githubSecretsToken);
  await exec(`git config --global user.name '${authorName}'`);
  await exec(`git config --global user.email '${authorEmail}'`);
  await exec(
    `git -C ${input.formulaPath} add ${input.formulaPath}/${input.formulaFilePath}`
  );
  await exec(`git -C ${input.formulaPath} commit -m '${commitMessage}'`);
  await exec(`git -C ${input.formulaPath} push`);
  return;
}

main().catch(err => core.setFailed(err.message));
