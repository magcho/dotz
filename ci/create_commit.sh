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

REPO="$(git config --get remote.origin.url)"
REPO=${REPO#https://}
REMOTE_REPO="https://${GITHUB_USERNAME}:${GITHUB_SECRETS_TOKEN}@${REPO}"

git add .
echo $(git status)
echo $(git commit -m "update ${BIN_NAME}")
git push "${REMOTE_REPO}"

cd $CURRENT_PATH
