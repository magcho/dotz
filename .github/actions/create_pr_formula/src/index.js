const util = require("util");
const childProcess = require("child_process");
const exec = util.promisify(childProcess.exec);
const core = require("@actions/core");
const fs = require("fs");

async function setAuth(userName, pass) {
  const repoGithubUrl = await exec(`brew --repository ${userName}/${userName}`)
    .stdout;
  const replaceGitUrl = repoGithubUrl.replace(
    "github.com",
    `${userName}:${pass}@github.com`
  );
  const formulaInstalledRepoPath = await exec(
    `${userName}/homebrew-${userName}`
  ).stdout;
  await exec(`git conifg remote.origin.url ${replaceGitUrl}`);
  return;
}

async function main() {
  const input = {
    githubSecretsToken: core.getInput("github_secrets_token"),
    githubUserName: core.getInput("github_username"),
    formulaFilePath: core.getInput("formula_file_path"),
    formulaUrl: core.getInput("formula_url"),
    authorName: core.getInput("author_name"),
    authorEmail: core.getInput("author_email"),
    commitMessage: core.getInput("commit_message")
  };

  await setAuth(input.githubUserName, input.githubSecretsToken);

  await exec(`git config --global user.name '${input.authorName}'`);
  await exec(`git config --global user.email '${input.authorEmail}'`);
  await exec(
    `git -C ${input.formulaFilePath} add ${input.formulaFilePath}/${input.formulaFilePath}`
  );
  await exec(`git -C ${repoGithubUrl} commit -m '${commitMessage}'`);
  await exec(`git -C ${repoGithubUrl} push`);
  return;
}

main().catch(err => core.setFailed(err.message));

// setAuth("asdf", "asdf");
