const execSync = require("child_prosess").execSync;
const core = require("@actions/core");

execSync("sh ./entrypoint.sh").catch(err => core.setFailed(err.message));
