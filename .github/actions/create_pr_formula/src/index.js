const util = require("util");
const childProcess = require("child_process");
const exec = util.promisify(childProcess.exec);
const core = require("@actions/core");
const fs = require("fs");

function main() {
  const input = {
    formulaFilename: core.getInput("formula_filename"),
    githubUserName: core.getInput("github_username"),
    githubSecretsToken: core.getInput("github_secrets_token"),
    formulaUrl: core.getInput("formula_url"),
    authorName: core.getInput("author_name"),
    authorEmail: core.getInput("author_email"),
    commitMessage: core.getInput("commit_message")
  };

  const userName = input.githubUserName;
  const pass = input.githubSecretsToken;
  const binName = input.formulaFilename.replace(/^(.*)\.rb/, $1);

  let brewClonedPath;
  exec(`brew --repository ${userName}/${userName}`)
    .then(({ stdout, stderr }) => {
      brewClonedPath = stdout.replace(/\n/, "");
      return exec(`git -C ${brewClonedPath} config --get remote.origin.url`);
    })
    .then(({ stdout, stderr }) => {
      const gitConfigUrl = stdout
        .replace(/\n/, "")
        .replace("github.com", `${userName}:${pass}@github.com`);
      return exec(
        `git -C ${brewClonedPath} config --local remote.origin.url ${gitConfigUrl}`
      );
    })
    .then(() =>
      exec(
        `git -C ${brewClonedPath} config --global user.name '${input.authorName}'`
      )
    )
    .then(() =>
      exec(
        `git -C ${brewClonedPath} config --global user.email '${input.authorEmail}'`
      )
    )
    .then(() =>
      exec(
        `git -C ${brewClonedPath} add ${brewClonedPath}/${input.formulaFilename}`
      )
    )
    .then(() => exec(`git -C ${brewClonedPath} commit -m 'update ${binName}'`))
    .then(() => exec(`git push`))
    .catch(err => core.setFailed(err.message));
}

main();

// setAuth("asdf", "asdf");
