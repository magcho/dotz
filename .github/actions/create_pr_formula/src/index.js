const util = require("util");
const childProcess = require("child_process");
const exec = util.promisify(childProcess.exec);
const core = require("@actions/core");
const fs = require("fs");

function main() {
  core.debug("start main()");
  core.info(12345);
  const input = {
    formulaFilename: core.getInput("formula_filename"),
    githubUserName: core.getInput("github_username"),
    githubSecretsToken: core.getInput("github_secrets_token"),
    commitMail: core.getInput("commit_mail"),
    commitMessage: core.getInput("commit_message")
  };

  const userName = input.githubUserName;
  const pass = input.githubSecretsToken;
  const binName = input.formulaFilename.replace(/^(.*)\.rb/, "$1");

  let brewClonedPath;
  exec(`brew --repository ${userName}/${userName}`)
    .then(({ stdout, stderr }) => {
      core.info(stdout.replace(/\n/, ""));
      core.info(stderr.replace(/\n/, ""));
      brewClonedPath = stdout.replace(/\n/, "");
      return exec(`git -C ${brewClonedPath} config --get remote.origin.url`);
    })
    .then(({ stdout, stderr }) => {
      const gitConfigUrl = stdout
        .replace(/\n/, "")
        .replace("github.com", `${userName}:${pass}@github.com`);
      core.info(gitConfigUrl);
      return exec(
        `git -C ${brewClonedPath} config --local remote.origin.url ${gitConfigUrl}`
      );
    })
    .then(({ stdout, stderr }) => {
      core.info(stdout.replace(/\n/, ""));
      core.info(stderr.replace(/\n/, ""));
      return exec(
        `git -C ${brewClonedPath} config --global user.name '${input.githubUserName}'`
      );
    })
    .then(({ stdout, stderr }) => {
      core.info(stdout.replace(/\n/, ""));
      core.info(stderr.replace(/\n/, ""));
      return exec(
        `git -C ${brewClonedPath} config --global user.email '${input.commitMail}'`
      );
    })
    .then(({ stdout, stderr }) => {
      core.info(stdout.replace(/\n/, ""));
      core.info(stderr.replace(/\n/, ""));
      return exec(
        `git -C ${brewClonedPath} add ${brewClonedPath}/${input.formulaFilename}`
      );
    })
    .then(({ stdout, stderr }) => {
      core.info(stdout.replace(/\n/, ""));
      core.info(stderr.replace(/\n/, ""));
      return exec(`git -C ${brewClonedPath} commit -m "update ${binName}"`);
    })
    .then(({ stdout, stderr }) => {
      core.info(stdout.replace(/\n/, ""));
      core.info(stderr.replace(/\n/, ""));
      return exec(`git push`);
    })
    .then(({ stdout, stderr }) => {
      core.info(stdout.replace(/\n/, ""));
      core.info(stderr.replace(/\n/, ""));
      return;
    })
    .catch(err => {
      core.info(err.message);
      core.setFailed(err.message);
    });
}

main();
