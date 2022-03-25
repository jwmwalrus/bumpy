ChangeLog
=========

All noticeable changes in the project  are documented in this file.

Format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).

This project uses [semantic versions](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.56.2] 2022-03-24

Detach BuildInfo from the version struct

## [0.56.1] 2022-03-24

Dummy release

## [0.56.0] 2022-03-24

Dummy release

## [0.55.0] 2022-03-24

Misc. Changes

### Added

* Exposing raw vcs-related values from runtime
    * Done through `version.LastCommit()` method.
* Flag to tag with a message

### Modified

* Made the git package public
* Updated Go version
* Reviewed tasks' help
* Updated README
    * Fixed typos and removed TODO list

## [0.54.1] 2021-12-04

Cleanup

### Fixed

* Requirements in the README

### Added

* config-related entry to the README

## [0.54.0] 2021-12-03

Add config task

### Fixed

* Displaying version, again
* The `init` task not affecting `.bumpy-ride`

### Added

* The `config` task

### Modified

* The `--do-init-commit` changed to `--persist` for the init task

## [0.53.1] 2021-12-03

Dummy release

## [0.53.0] 2021-12-03

Rename repository

### Fixed

* Version displayed by the `--version` flag

### Modified

* Renamed repository in order to shorten binary.
* Bumped go.mod version to 1.16

## [0.52.0] 2021-12-02

Add config file

## Added

* A configuration file, `.bumpy-ride`, at the root directory. Its format is JSON and has the following options:
    * `noFetch` (bool): if true, inhibits the feh operations
    * `noCommit` (bool): if true, inhibits the commit operations
    * `versionPrefix` (string): path to the `version.json` file (default: ".")
    * `npmPrefixes` (array of strings): if non-empty, the `package(-lock)?.json` file(s) at the given relative directories will be affected by the `bump` task

## Modified

* All tasks now get executed from git's root directory

## Removed

* Flag `--no-commit`, replaced by the `noCommit` config entry
* Flag `--no-fetch`, replaced by the `noFetch` config entry
* Flag `--npm-prefix`, replaced by the `npmPrefixes` config entry

## [0.51.2] 2021-09-20

Update dependencies

## [0.51.1] 2021-09-20

Cleanup

### Modified

* Printlin usage in favor of Printf

### Removed

* Dead code

## [0.51.0] 2021-06-20

Version task and commands help

### Modified

* Commands help

### Added

* The version task

## [0.50.1] 2021-06-20

Fixes and improvements

### Fixed

* Bump error after init

### Added

* `--no-commit` flag to the bump command

## [0.50.0] 2021-06-19

Implement sync, bump and tag

### Added

* The sync task
* The bump task
* The tag task

## [0.10.0] 2021-06-14

Implement init

### Added

* The init task

