const util = require("util");
const childProcess = require("child_process");
const exec = util.promisify(childProcess.exec);
const core = require("@actions/core");

function main() {
  input = [
    core.getInput("formula_filename"),
    core.getInput("github_username"),
    core.getInput("github_secrets_token"),
    core.getInput("commit_mail")
  ];
  const args = input.join(" ");
  exec(`sh ./.github/actions/create_commit_sh/main.sh ${args}`)
    .then(({ stdout, stderr }) => {
      core.info(stdout);
      if (stderr != "") {
        core.warning(stderr);
      }
    })
    .catch(err => {
      core.setFailed(err.message);
    });
}
main();
