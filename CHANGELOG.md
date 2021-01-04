# CHANGELOG

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
