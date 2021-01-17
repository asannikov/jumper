# CHANGELOG

## 1.7.2 (17.01.2021)
- bugfix
    - refactoring for commands, added general type that contains all func parameters
    - added xdebug description

## 1.7.0 (14.01.2021)
- features
    - add magento command (bin/magento)
    - add magerun2 command

## 1.6.3 (13.01.2021)
- bugfix 
    - remove jumper before install it

## 1.6.2 (13.01.2021)
- bugfix 
    - fix jumper.rb template, improved jumper_install.sh

## 1.6.1 (13.01.2021)
- bugfix 
    - add readme info und version

## 1.6.0 (13.01.2021)
- features
    - 13 Add scope command config. User can hide unnecessary commands in global config file (composer, xdebug). 

## 1.5.7 (13.01.2021)
- bugfix
    - 27 Fix shell type detection outside project

## 1.5.6 (13.01.2021)
- bugfix
    - refactoring: changed files hierarchy in project
    - fix unittest for getSyncArgs

## 1.5.5 (12.01.2021)
- bugfix
    - remove legacy interfaces

## 1.5.4 (12.01.2021)
- bugfix
    - 27 shell command now detect project path. it could not detect project and had no project dialog. fix

## 1.5.3 (11.01.2021)
- features
    - add jumper_install script

## 1.5.2 (10.01.2021)
- bugfix
    - add homebrew macos support (linux was not tested)

## 1.5.1 (09.01.2021)
- features
    - add go build bash script
    - prepared for built releases

## 1.5.0 (06.01.2021)
- features
    - 25 Add shell command, which set the default shell type for usage in main contaner
- bugfix
    - 25 Split logic for config package
    - 25 Add missed xdebug unittests

## 1.4.1 (06.01.2021)
- bugfix
    - 25 refactored unittests 
    - 25 add xdebug unittests
    - 25 refactored interfaces for commands, main interface splitted into different interfaces

## 1.4.0 (05.01.2021)
- features
    - 4 Add xdebug enable/disable
- bugfix
    - 4 Added sync test
    - 4 some containers do not use bash. sh was set by default

## 1.3.2 (05.01.2021)
- bugfix
    - 22 User is always asked for a project select even he is in the folder with jumper.json.
      Now jumper checks if the file exists and if the configuration appropriate.

## 1.3.1 (04.01.2021)
- bugfix
    - 17 Was not possible to stop containers, cause docker object was not initialized. Fixed.

## 1.3.0 (30.12.2020)
- features
    - 3 Add sync container -> host and host -> container
- bugfix
    - 3 Fix typo in dialog titles for project path naming 
    - 3 Add docker status check for stop commands
    - 3 Initf function now has project path return 

## 1.2.1 (29.12.2020)
- bugfix
    - 17 Add unittests for docker wrapper object
    - 17 Add unittests for docker start dialog object

## 1.2.0 (28.12.2020)
- features
    - 16 Add docker start option
    - 16 Add docker service dialog for docker starting using jumper
    - 16 Add unittests for Docker object 25%
- bugfix
    - 16 refacored docker instance
    - update docker API

## 1.1.0 (27.12.2020)
- features
    - 14 Added option, which hides long copyright text
    - 14 Added unittest for this option
- bugfix
    - 14 Global config is loaded on `jumper` call now, even without command call. It helps to get an access to conifg options at app instance initialization. 

## 1.0.2
    - 12 Add project path output
    - 12 Add changelog file
