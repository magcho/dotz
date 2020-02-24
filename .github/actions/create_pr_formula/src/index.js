const util = require("util");
const childProcess = require("child_process");
const exec = util.promisify(childProcess.exec);
const core = require("@actions/core");
const fs = require("fs");

async function setAuth(formulaPath, repoUrl, userName, pass) {
  const repo =
    repoUrl
      .replace(/\/$/, "")
      .replace(/\.git$/, "")
      .replace(/^https:\/\//) + ".git";
  const gitConfFile = fs.readFile(`${formulaPath}/.git/config`, (err, data) => {
    if (err) {
      core.setFailed(err.message);
    }
    data.replace(repo, `${userName}:${pass}${repoUrl}`);
  });
}

async function main() {
  const input = {
    formulaPath: core.getInput("formula_path"),
    githubSecretsToken: core.getInput("github_secrets_token"),
    githubUserName: core.getInput("github_username"),
    formulaFilePath: core.getInput("formula_file_path"),
    formulaUrl: core.getInput("formula_url"),
    authorName: core.getInput("author_name"),
    authorEmail: core.getInput("author_email"),
    commitMessage: core.getInput("commit_message")
  };

  setAuth(
    input.formulaFilePath,
    input.formulaUrl,
    input.githubUserName,
    input.githubSecretsToken
  );
  await exec(`git config --global user.name '${input.authorName}'`);
  await exec(`git config --global user.email '${input.authorEmail}'`);
  await exec(
    `git -C ${input.formulaPath} add ${input.formulaPath}/${input.formulaFilePath}`
  );
  await exec(`git -C ${input.formulaPath} commit -m '${commitMessage}'`);
  await exec(`git -C ${input.formulaPath} push`);
  return;
}

main().catch(err => core.setFailed(err.message));

// setAuth("asdf", "asdf");
