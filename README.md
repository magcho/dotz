<p align="center">
  <img  src="https://github.com/magcho/dotz/raw/master/dotz-min.png">
</p>
<p align="center"><b>
Manage dotfiles for macOS  
</b></p>


## install

### Useing homebrew
```
brew tap magcho/magcho
brew install dotz
```

### Manually
Download dotz binary from [Github relase](https://github.com/magcho/dotz/releases),and move directory of $PATH.

### Self build
Cloning this repository. Exec `go build main.go -o dotz` and move directory of $PATH.


## Usage

### Init
  Before setting DOTZ_ROOT env OR command parameter.
  ```
  dotz init [--DOTZ_ROOT xx]
  ```
  
  1. Create dotz project folder into DOTZ_ROOT (default `~/.dotz`)
  1. Initialize git

### Track
  ```
  dotz track xx    // Tracking of file
  dotz track -f xx // Tracking of folder
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
  Before cloning dotz project into DOTZ_ROOT path and set env DOTZ_ROOT.
  ```
  dotz restore
  ```
  1. Create symbolic links.
  
  
## Lisense

Apache License 2.0
  
