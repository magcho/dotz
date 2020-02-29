#!/usr/bin/env bash

set -euo pipefail

CURRENT_PATH="$(pwd)"

# set args variables
FORMULA_FILENAME="$1"
GITHUB_USERNAME="$2"
GITHUB_SECRETS_TOKEN="$3"
COMMIT_MAIL="$4"

# set env variables
BIN_NAME="${FORMULA_FILENAME%%.rb}"


cd "$(brew --repository ${GITHUB_USERNAME}/${GITHUB_USERNAME})"

git config --global user.name $GITHUB_USERNAME
git config --global user.email $COMMIT_MAIL

git config remote.origin.url $(git config --get remote.origin.url | sed -e "s/github.com/${GITHUB_USERNAME}:${GITHUB_SECRETS_TOKEN}@gitub.com/")

git add .
echo $(git status)
echo $(git commit -m "update ${BIN_NAME}")
git push

cd $CURRENT_PATH