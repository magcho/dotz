#!/usr/bin/env bash

set -euo pipefail

CURRENT_PATH="$(pwd)"

# set args variables
USER_NAME="$1"
FORMULA_FILE_NAME="$2"

/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
brew tap $USER_NAME/$USER_NAME
brew update

brew install gawk

cp $FORMULA_FILE_NAME $(brew --repo $USER_NAME/$USER_NAME)

cd $(brew --repo $USER_NAME/$USER_NAME)

# variables
BIN_NAME="$(gawk 'match($0, /bin\.install\s+"(.*)"/, a) {print a[1]}' $FORMULA_FILE_NAME)"
BIN_VER="$(gawk 'match($0, /:tag\s+=>\s+"v([[:alnum:]+\.]+)"/, a) {print a[1]}' $FORMULA_FILE_NAME)"

brew audit $FORMULA_FILE_NAME --fix
brew style $FORMULA_FILE_NAME --fix

brew install -f $BIN_NAME
which $BIN_NAME
brew rm $BIN_NAME

cd $CURRENT_PATH
