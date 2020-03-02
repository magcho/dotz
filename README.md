<p align="center">
  <img  src="https://github.com/magcho/dotz/raw/master/dotz-min.png">
</p>
<p align="center"><b>
Manage dotfiles for macOS  
</b></p>

![build for mac](https://github.com/magcho/dotz/workflows/build%20for%20mac/badge.svg)

## install

### Useing homebrew
```
brew tap magcho/magcho
brew install dotz
```

### Manually
Download dotz binary from [Github relase](https://github.com/magcho/dotz/releases),and move directory into $PATH.

### Self build
Cloning this repository. Exec `go build main.go -o dotz` and move directory into $PATH.


## Usage

### Init
  Before setting $DOTZ_ROOT or command parameter.
  ```
  export DOTZ_ROOT=~/.dotz
  dotz init
  ```
  
  1. Create dotz project folder into DOTZ_ROOT (default `~/.dotz`)
  1. Initialize git

### Track
  ```
  dotz track xx    // Tracking file
  dotz track -f xx // Tracking folder
  ```
  1. Move xx to dotz project folder.
  1. Create xx file or folder symbolic link.
  
### Backup
  Before setting dotz root git, dotz do backup.
  ```
  dotz backup [-p]
  ```
  1. Commit of dotzproject of git, optionally push.
  
---

### Restore
  Before cloning dotz project into $DOTZ_ROOT path and set env ＄DOTZ_ROOT.
  ```
  dotz restore
  ```
  1. Create symbolic links.
  
  
## Lisense

Apache License 2.0
  
