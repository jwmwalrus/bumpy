Bumpy Ride
==========

A versioning tool.

## Table of Contents
* [Requirements](#requirements)
* [Installation](#installation)
* [Usage](#usage)
* [TODO](#todo)

## Requirements

* Go version 1.16 or higher. 
You can probably install it through your system's package manager (`apt`, `brew`, etc.). 
For general instructions, go [here](https://golang.org/doc/install) --no pun intended.

## Installation

To install, open a terminal and execute the following
```bash
go install github.com/jwmwalrus/bumpy@latest
```

The same command can be used for subsequent updates.

## Usage

**Bumpy Ride** stores a repository's current version in the `version.json` file at any specified location in said repository (the default location is the root directory). The idea is to allow for the version to be embedded in the executable, easing further manipulation, with no regards for link-time flags and such.

An overview of the available commands and options can be obtained by executing
```bash
bumpy --help
```

The available commands can be categorized into three groups: **control**, **git-affecting** and **informational**. 

### Control Commands

These commands allow for manipulation of the `version.json` file that stores the repository's current version.

#### init

The `init` command initializes the repository's version --i.e., it creates the configuration and the version files.

Detailed information aobut the `init` command can be otained with:
```bash
bumpy help init
```

#### config

The `config` command updates the repository's version configuration file, according to the provided options, and displays its resulting contents.

Detailed information aobut the `config` command can be otained with:
```bash
bumpy help config
```

#### sync

The `sync` command makes sure that the version stored in `version.json` matches the latest git tag available.

Detailed information aobut the `sync` command can be otained with:
```bash
bumpy help sync
```

### Git-affecting Commants

These commands may perform operations on the `version.json` file, and cause at least one commit and/or other git-related operations.

#### bump

The `bump` command updates the `version.json` file, according to the given set of options.

Detailed information aobut the `bump` command can be otained with:
```bash
bumpy help bump
```

#### tag

The `tag` command commits the ChangeLog file, and tags its commit with the latest version from `version.json`.

Detailed information aobut the `tag` command can be otained with:
```bash
bumpy help tag
```

### Informational Commands

These commands simply display information.

#### version

The `version` command shows the current version stored in the `version.json` file.

Detailed information aobut the `version` command can be otained with:
```bash
bumpy help version
```


## TODO

- [x] ~Implement the init task~
- [x] ~Implement `bump [--major|--minor|--patch]` task~
- [x] ~Implement `bump [--pre PRE] [--build BUILD]` task~
- [x] ~Implement tag task~
- [x] ~Implement sync task~
- [x] ~Implement version task~
- [x] ~Implement version configuration~
- [ ] Implement incremental patterns for `--pre` and `--build`
- [ ] Implement version operators
- [ ] Allow tagging with an arbitrary list of files

