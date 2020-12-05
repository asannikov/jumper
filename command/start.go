package command

import (
	"log"

	"github.com/urfave/cli/v2"
)

// CallStartProject runs docker project
func CallStartProject(cp func() (string, error)) *cli.Command {
	cmd := cli.Command{
		Name:    "start",
		Aliases: []string{"st"},
		Usage:   "Start project",
		Action: func(c *cli.Context) error {

			log.Println(cp())
			//projectPath := cp()

			return nil
			/*var binary = "docker"
			var initArgs = []string{"exec", "-it", "m24_phpfpm_1", "/usr/local/bin/php", "-d", "memory_limit=-1", "/usr/local/bin/composer", "update"}

			extraInitArgs := []string{}

			args := append(initArgs, extraInitArgs...)

			cmd := exec.Command(binary, args...)

			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			cmd.Run()

			return nil*/
		},
	}

	return &cmd
}

/*
#!/bin/bash
function parseYaml {
  local s='[[:space:]]*' w='[a-zA-Z0-9_]*' fs=$(echo @|tr @ '\034') # https://linuxhint.com/bash_tr_command/
  sed -ne "s|,$s\]$s\$|]|" \
      -e ":1;s|^\($s\)\($w\)$s:$s\[$s\(.*\)$s,$s\(.*\)$s\]|\1\2: [\3]\n\1  - \4|;t1" \
      -e "s|^\($s\)\($w\)$s:$s\[$s\(.*\)$s\]|\1\2:\n\1  - \3|;p" $1 | \
  sed -ne "s|,$s}$s\$|}|" \
      -e ":1;s|^\($s\)-$s{$s\(.*\)$s,$s\($w\)$s:$s\(.*\)$s}|\1- {\2}\n\1  \3: \4|;t1" \
      -e    "s|^\($s\)-$s{$s\(.*\)$s}|\1-\n\1  \2|;p" | \
  sed -ne "s|^\($s\):|\1|" \
      -e "s|^\($s\)-$s[\"']\(.*\)[\"']$s\$|\1$fs$fs\2|p" \
      -e "s|^\($s\)-$s\(.*\)$s\$|\1$fs$fs\2|p" \
      -e "s|^\($s\)\($w\)$s:$s[\"']\(.*\)[\"']$s\$|\1$fs\2$fs\3|p" \
      -e "s|^\($s\)\($w\)$s:$s\(.*\)$s\$|\1$fs\2$fs\3|p" | \
  awk -F$fs '{
    indent = length($1)/2;
    vname[indent] = $2;
    for (i in vname) {if (i > indent) {delete vname[i]; idx[i]=0}}
    if (length($2) == 0) {vname[indent] = ++idx[indent] };
    if (length($3) > 0) {
      vn=""; for (i=0; i<indent; i++) {vn = (vn)(vname[i])("_")}
      if ("'$2'_" == vn) {
         print substr($3 ,1 , match($3,":")-1)
      }
    }
  }'
}

# Check if volume files exist to avoid creating an empty folder
VOLUME_LIST=`parseYaml docker-compose.dev.yml services_app_volumes`
IGNORE_LIST="./src/app/code ./src/m2-hotfixes ./src/patches ./src/var/log ./src/var/report ./src"
IS_VALID=true
# Loop through all files missing from the docker-compose.dev.yml file
for file in $VOLUME_LIST; do
  if [ ! -e $file ] && [[ ! " $IGNORE_LIST " =~ " $file " ]]; then
    echo "$file: No such file or directory"
    IS_VALID=false
  fi
done
# Wait to exit until all missing files have been outputted to the user
[ $IS_VALID = false ] && echo "Failed to start docker for missing volume files" && exit

docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --force-recreate -d --remove-orphans "$@"

## Blackfire support
# ------------------
## First register the blackfire agent with:
#bin/root blackfire-agent --register --server-id={YOUR_SERVER_ID} --server-token={YOUR_SERVER_TOKEN}
## Then uncomment the below line (and leave uncommented) to start the agent automatically with bin/start:
#bin/root /etc/init.d/blackfire-agent start

*/
