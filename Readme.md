# JUMPER

This tool has been created for helping developers out with daily docker routines. Once you configured the jumper file for a project, you can call easily run basic docker commands for any registered project and jump into main docker container using just `jumper sh`.

I was inspired by [Mark's Shust](https://github.com/markshust/docker-magento) solutution for M2. Many thanks to Mark from my side for his ideas.

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat)](https://github.com/asannikov/jumper/blob/master/LICENSE)
[![Build Status](https://travis-ci.com/asannikov/jumper.svg?branch=master)](https://travis-ci.com/asannikov/jumper)
[![Release](https://img.shields.io/badge/release-1.8.10-brightgreen.svg)](https://github.com/asannikov/jumper/releases/tag/v1.8.10)

It was not tested on Windows.

# WHAT THIS TOOL CAN DO?
common things:
* enter to main defined container
* run docker project from any path of your local machine.
* container management (start/stop/restart)
* stop all containers
* certain path or whole project synchronization
* get project path

PHP
* run composer install/update with/without memory constraint
* run composer commands
* xdebug enable/disable for cli or fpm

@todo
* mysqldump
* redis and varnish management
* Magento2 command support
* Laravel artisan support
* Docker extended managment
* GIT routines
* Kubernetes management

I'll highly appreciate for testing and any ideas where this tool can be useful. Please, write your suggestions or your experience in the issues.

# HOW IT WORKS?
Once you call `jumper`, it looks for the `jumper.json` file in the current path where you are at this moment. 

Case1:

if the file `jumper.json` exists, it checks if the main container defined. If not, you will be asked to select the related container. IMPORTANT: docker project has to be run to select the related container. It's necessary once until container name is changed.
If main container has been defined already, it Runs the related command, ie start docker project.

Case2:

if the file `jumper.json` does not exists, you will be asked to select the existing project or create a new one. You set the path to the project and select main container. After that `jumper.json` will be created in the defined project path. Besides, `.jumper.json` will be created in the user's path. It contains the list of created projects and is used for cli dialog when you want to select the project that you want to work with.

## Requirements 
- go1.13.x - go1.15.x

## Install
#### 1. Using Go
 - Please follow the link "[How to install go](https://golang.org/doc/install)".
 - Clone project to your custom directory `git clone https://github.com/asannikov/jumper.git`
 - Find the path to user's go folder `go env GOPATH` and copy it for the next command
 - Edit the file `~/.bash_profile` or `~/.zshrc` depending on enviroment and add/edit the the variables (This is an example path for Mac OS. Change the paths for Linux or Windows on your own):
 ```
export GOPATH="/Users/username/go" 
export GOROOT="/usr/local/bin/go"
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
 ```
 - run `. ~/.zshrc` or `source ~/.bash_profile`
 - Go to jumper project path and run `go install`
 

Now you can call `jumper` in any folder.

Useful links: [Command not found go — on Mac after installing Go](https://stackoverflow.com/questions/34708207/command-not-found-go-on-mac-after-installing-go)


#### 2. Using brew (does not work officially yet)

> I prepared branch for pull reqest https://github.com/asannikov/homebrew-core/tree/jumper, but it's going to be in brew when the next goals will be reached:
> 1) Magento, Mysql dumps and artisan management will be done
> 2) Repository reached 70 stars and 30 forks (requirement by Brew team)
> 3) Almost all code covered by tests and have no reported issues for a long time.


- install brew
 ```
 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
 ```
 - ~~run `brew install jumper`~~

but you can still use jumper in brew in test mode  (mac os and linux only). Download file:
```
wget https://raw.githubusercontent.com/asannikov/jumper/master/jumper_install.sh
```
script will install jumper in test mode on your machine. Run it:
```
sh ./jumper_install.sh
```
It can ask the "Formula name", leave it empty and press enter.

now you can call `jumper`.

#### 3. Download source directly from repository.

Every stable release has attached sources for "linux/amd64", "linux/386", "darwin/amd64", "windows/amd64" and "windows/386"

For example:
```
https://github.com/asannikov/jumper/releases/tag/v1.8.0
```
Find related source there and download it. Now you can place use source at any place you want on your machine or make it global in a standard way.

## Usage

`jumper [global options] command [command options] [arguments...]`

implemented commands:
```
   cli, c                          Runs cli command in conatiner: {docker exec main_conatain} [command] [custom parameters]
   sh, sh                          Runs cli sh command in conatiner: {docker exec main_conatain {shell_type}} [custom parameters]
   clinotty, cnt                   Runs command {docker exec -t main_container} [command] [custom parameters]
   cliroot, cr                     Runs command {docker exec -u root main_container} [command] [custom parameters]
   clirootnotty, crnt              Runs command {docker exec -u root -t main_container} [command] [custom parameters]
   composer, cmp                   Runs composer: {docker exec -it phpContainer composer} [custom parameters]
   composer:memory, cmpm           Runs composer with no memory constraint: {docker exec -i phpContainer /usr/bin/php -d memory_limit=-1 /usr/local/bin/composer} [custom parameters]
   composer:install, cmpi          Runs composer install: {docker exec -it phpContainer composer install} [custom parameters]
   composer:install:memory, cmpim  Runs composer install with no memory constraint: {docker exec -i phpContainer /usr/bin/php -d memory_limit=-1 /usr/local/bin/composer install} [custom parameters]
   composer:update, cmpu           Runs composer update: {docker exec -it phpContainer composer update} [custom parameters]
   composer:update:memory, cmpum   Runs composer update with no memory constraint: {docker exec -i phpContainer /usr/bin/php -d memory_limit=-1 /usr/local/bin/composer update} [custom parameters]
   start, st                       Runs defined command: {docker-compose -f docker-compose.yml up} [custom parameters]
   start:force, s:f                Runs defined command: {docker-compose -f docker-compose.yml up --force-recreat} [custom parameters]
   start:orphans, s:o              Runs defined command: {docker-compose -f docker-compose.yml up --remove-orphans} [custom parameters]
   start:force-orphans, s:fo       Runs defined command: {docker-compose -f docker-compose.yml up --force-recreate --remove-orphans} [custom parameters]
   start:maincontainer, startmc    Runs defined command: {docker start main_container}
   start:containers, startc        Runs defined command: {docker start} [container]
   restart:maincontainer, rmc      restarts main container
   restart:containers, rc          Runs defined command: {docker start} [container]
   stopallcontainers, sac          Stops all docker containers
   stop:containers, scs            Stops docker containers
   stop:maincontainer, smc         Stops main docker container
   stop:container, stopc           Stops selected docker containers
   path                            Gets project path
   copyright                       
   copyto, cpt                     Sync local -> docker container, set related path, ie `vendor/folder/` for syncing as a parameter, or use --all to sync all project
   copyfrom, cpf                   Sync docker container -> local, set related path, ie `vendor/folder/` for syncing as a parameter, or use --all to sync all project
   xdebug:fpm:enable, xe           Enable fpm xdebug
   xdebug:fpm:disable, xd          Disable fpm xdebug
   xdebug:fpm:toggle, x            Toggle fpm xdebug
   xdebug:cli:enable, xce          Enable cli xdebug
   xdebug:cli:disable, xcd         Disable cli xdebug
   xdebug:cli:toggle, xc           Toggle cli xdebug
   shell                           Change shell type for a project
   magento, m                      Call magento command bin/magento or magerun. This command has subcommands. Call jumper magento for more details.
```

# Project config example - jumper.json
```
{
 "name": "Project Name",                  
 "main_container": "php_container",       // container, where mostly you are working
 "start_command": "docker-compose up -d", // standard project run
 "path": "/var/www/src",                  // path to mounted project
 "xdebug_location": "local",              // Possible values: local and docker. Jumper will look for the xdebug file locally if was set "local"
 "xdebug_path_cli": "/path/to/project/config/xdebug_cli.ini", // absolute path to xdebug config locally or in docker
 "xdebug_path_fpm": "/path/to/project/config/xdebug_fpm.ini", // absolute path to xdebug config locally or in docker
 "shell": "bash"                          // shell by default
}
```

all these options are managed using jumper command. **It's not recommended to edit manually.**

# Global conifg file example - ~/.jumper.json
```
{
 "projects": [
  {
   "path": "/path/to/project1",
   "name": "project1"
  },
  {
   "path": "/path/to/project2",
   "name": "project2"
  }
 ],
 "copyright_text": false,
 "docker_service": "open --hide -a Docker",
 "inactive_command_types": [
  "compose",
  "xdebug"
 ]
}
```

# Available options for global config ~/.jumper.json

## Hide copyright text
Add `"copyright_text": false,` 

It can be done using jumper command though.

## Permanent hiding some types of commands
add `"inactive_command_types": ["type1","type2"],`

It can be done only directly in config file.

here is a list of avalable values:
* `composer` - hides all composer commands
* `xdebug`   - hides xdebug commands
* `magento`  - hides magento commands

# FAQ
## Where to start?
An example for Mac OS:
[![asciicast](https://asciinema.org/a/MUuwCpGh3Ty0zMChVL8K50BK9.svg)](https://asciinema.org/a/MUuwCpGh3Ty0zMChVL8K50BK9)

## How to force create path on sync:
`jumper copyto -f src/vendor/module/name`
assume that src/vendor does not exist in container. By this command it will be created recursively. 
The same you can do on host: `jumper copyfrom -f src/vendor/module/name`.

## How to change shell type:
Container might use only sh shell type and no bash. By default jumper uses `sh`. But you can use `bash` or even 	`csh`, `ksh` or `zsh`. Call `jumper shell` with no option and select shell type. It will be saved into project config `jumper.json` in `shell` node.

## How to configure xdebug:
Xdebug config ini files might be mounted to the container or able to be found directly in the container without mounting.
Depending on case, you have to select following option on xdebug:* running:

* mounted config file: 
     
   set path to config file on host. It has to be related to project root folder.

* config file in container:

   set absolute path to file in container

The config file must contain `;zend_extension=xdebug` option!

## How to run bash commands inside container:
```jumper clirootnotty bash -c "cd src && ls"```

## How to hide long copyright text:
```$ jumper copyright disable```

You accept terms of the agreement the by doing this.

## How to add custom command:
Find method getCommandList in `app/commands.go` and write any custom method, ie:
`command.CustomMethod()`

getCommandList function has 2 parameters:
- c - config, contains general configuration and config for the current project.
- d - dialog, is used for communication with user, select container, remove file, etc.
- initf - function, that has to be called inside your command if you are going to interact with the config file or docker project.

## How to add extra config to project json file:
Add new project fields here `app/config/ProjectConfig.go` 
and add the related method to `app/config/config.go`, ie:
```
// SaveContainerNameToProjectConfig saves container name into project file
func (c *Config) SaveContainerNameToProjectConfig(cn string) (err error) {
	c.projectConfig.MainContainer = cn
	return c.fileSystem.SaveConfigFile(c.projectConfig, c.getProjectFile())
}

// GetProjectMainContainer gets project main container
func (c *Config) GetProjectMainContainer() string {
	return c.projectConfig.GetMainContainer()
}
```

## How to add custom option to GlobalConfig:
global config is visible only in `app/config/config.go` scope. Follows these steps to create custom options:

1. Add json field into `GlobalConfig` type in `globalConfig.go`.
2. Add new methods for new option handling in `globalConfig.go`.
3. Add duplicated methods in `app/config/config.go`.
4. Inject your logic inside the code, see ie Copyright methods: `EnableCopyright/DisableCopyright/ShowCopyrightText`

## MIT License

Copyright (c) Anton Sannikov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.