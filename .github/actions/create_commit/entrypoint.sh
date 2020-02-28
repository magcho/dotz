#!/usr/bin/env bash

set -euo pipefail

# Ensure all variables are present
formulaFilename="$1"
githubUsername="$2"
githubSecretsToken="$3"
commitMail="$4"
commitMessage="$5"


# Create Temporary Directory
TEMP=$(mktemp -d)

# Setup git
git config --global user.email $commitMail
git config --global user.name $githubUsername
git clone $REPO $TEMP
cd $TEMP


# Sync $TARGET folder to $REPO state repository with excludes
echo "running 'rsync -avh --delete "${EXCLUDES[@]}" $GITHUB_WORKSPACE/$SOURCE/ $TEMP/$TARGET'"
rsync -avh --delete "${EXCLUDES[@]}" $GITHUB_WORKSPACE/$SOURCE/ $TEMP/$TARGET

# Success finish early if there are no changes
if [ -z "$(git status --porcelain)" ]; then
  echo "no changes to sync"
  exit 0
fi

# Add changes and push commit
git add .
SHORT_SHA=$(echo $GITHUB_SHA | head -c 6)
git commit -F- <<EOF
Automatic CI SYNC Commit $SHORT_SHA
Syncing with $GITHUB_REPOSITORY commit $GITHUB_SHA
EOF
git push
