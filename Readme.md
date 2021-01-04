# JUMPER

This tool has been created for helping developers out with daily docker routines. Once you configured the jumper file for a project, you can call easily run basic docker commands for any registered project and jump into main docker container using just `jumper b`.

I was inspired by [Mark's Shust](https://github.com/markshust/docker-magento) solutution for M2. Many thanks to Mark from my side for his ideas.

[![Build Status](https://travis-ci.com/asannikov/jumper.svg?branch=master)](https://travis-ci.com/asannikov/jumper)

It was not tested on Windows.

# WHAT THIS TOOL CAN DO?
common things:
* enter to main defined container
* run docker project from any path of your local machine.
* container management (start/stop/restart)
* stop all containers

PHP
* run composer install/update with/without memory constraint
* run composer commands

@todo
* xDebug on/off 
* mysqldump
* file sync
* redis and varnish management
* Magento2 command support
* Laravel artisan support
* Docker extended managment
* GIT routines
* Extendable jumper config options
* Kubernetes management

I'll highly appreciate for testing and any ideas where this tool can be useful. Please, write your suggestions or your experience in the issues.

# HOW IT WORKS?
Once you call `jumper`, it looks for the `jumper.json` file in the current path where you are at this moment. 

Case1:

if the file `jumper.json` exists, it checks if the main container defined. If not, you will be asked to select the related container. IMPORTANT: docker project has to be run to select the related container. It's necessary once until container name is changed.
If main container has been defined already, it runs the related command, ie start docker project.

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


#### 2. Using brew (for Mac OS only, does not work yet)
 - install brew
 ```
 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
 ```
 - run `brew install jumper` 

## Usage

`jumper [global options] command [command options] [arguments...]`

implemented commands:
```
cli, c                          Runs cli command in conatiner: {docker exec main_conatain} [command] [custom parameters]
   bash, b                         Runs cli bash command in conatiner: {docker exec main_conatain bash} [custom parameters]
   clinotty, cnt                   Runs command {docker exec -t main_container} [command] [custom parameters]
   cliroot, cr                     Runs command {docker exec -u root main_container} [command] [custom parameters]
   clirootnotty, crnt              Runs command {docker exec -u root -t main_container} [command] [custom parameters]
   composer, cmp                   Runs composer: {docker exec -it phpContainer composer} [custom parameters]
   composer:memory, cmpm           Runs composer with no memory constraint: {docker exec -i phpContainer /usr/bin/php -d memory_limit=-1 /usr/local/bin/composer} [custom parameters]
   composer:install, cmpi          Runs composer install: {docker exec -it phpContainer composer install} [custom parameters]
   composer:install:memory, cmpim  Runs composer install with no memory constraint: {docker exec -i phpContainer /usr/bin/php -d memory_limit=-1 /usr/local/bin/composer install} [custom parameters]
   composer:update, cmpu           Runs composer update: {docker exec -it phpContainer composer update} [custom parameters]
   composer:update:memory, cmpum   Runs composer update with no memory constraint: {docker exec -i phpContainer /usr/bin/php -d memory_limit=-1 /usr/local/bin/composer update} [custom parameters]
   start, st                       runs defined command: {docker-compose -f docker-compose.yml up} [custom parameters]
   start:force, s:f                runs defined command: {docker-compose -f docker-compose.yml up --force-recreat} [custom parameters]
   start:orphans, s:o              runs defined command: {docker-compose -f docker-compose.yml up --remove-orphans} [custom parameters]
   start:force-orphans, s:fo       runs defined command: {docker-compose -f docker-compose.yml up --force-recreate --remove-orphans} [custom parameters]
   start:maincontainer, startmc    runs defined command: {docker start main_container}
   start:containers, startc        runs defined command: {docker start} [container]
   restart:maincontainer, rmc      restarts main container
   restart:containers, rc          runs defined command: {docker start} [container]
   stopallcontainers, sac          Stops all docker containers
   stop:containers, scs            Stops docker containers
   stop:maincontainer, smc         Stops main docker container
   stop:container, stopc           Stops selected docker containers
   path                            gets project path
   copyright                       
   copyto, cpt                     Sync local -> docker container, set related path, ie `vendor/folder/` for syncing as a parameter, or use --all to sync all project
   copyfrom, cpf                   Sync docker container -> local, set related path, ie `vendor/folder/` for syncing as a parameter, or use --all to sync all project
   ```

# FAQ
## How to run bash commands inside container:
```jumper clirootnotty bash -c "cd src && ls"```

## How to hide long copyright text:
```$ jumper copyright disable```

You accept terms of the agreement the by doing this.

## How to add custom command:
Find method getCommandList in `commands.go` and write any custom method, ie:
`command.CustomMethod()`

getCommandList function has 2 parameters:
- c - config, contains general configuration and config for the current project.
- d - dialog, is used for communication with user, select container, remove file, etc.
- initf - function, that has to be called inside your command if you are going to interact with the config file or docker project.

## How to add extra config to project json file:
Add new project fields here `config/ProjectConfig.go` 
and add the related method to `config/config.go`, ie:
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
global config is visible only in `config/config.go` scope. Follows these steps to create custom options:

1. Add json field into `GlobalConfig` type in `globalConfig.go`.
2. Add new methods for new option handling in `globalConfig.go`.
3. Add duplicated methods in `config/config.go`.
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