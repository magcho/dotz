#!/usr/bin/env bash

set -euo pipefail

current_path="$(pwd)"

# set args variables
# - formula_repo_url   https://github.com/magcho/homebrew-magcho.git
# - formula_file_name  dotz.rb


# var
origin_formula_url=$()

# curl -o https://raw.githubusercontent.com/magcho/homebrew-magcho/master/dotz.rb
curl -O "${formula_repo_url}"
