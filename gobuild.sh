#!/usr/bin/env sh

# https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04

version=$1
token=$2
if [[ -z "$token" || -z "$version" ]]; then
  echo "usage: $0 <version> <github-token>"
  exit 1
fi

package=github.com/asannikov/jumper
package_name=jumper

platforms=("linux/amd64" "linux/386" "darwin/amd64" "windows/amd64" "windows/386")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    if [ $GOOS = "darwin" ]; then
        output_name+='.zip'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done

mkdir ./pkg/$version/
cp LICENSE ./pkg/$version/

echo "mv jumper-windows-amd64.exe ./pkg/$version/"
mv jumper-windows-amd64.exe ./pkg/$version/
echo "mv jumper-windows-386.exe ./pkg/$version/"
mv jumper-windows-386.exe ./pkg/$version/

echo "mv jumper-linux-amd64 ./pkg/$version/jumper"
mv jumper-linux-amd64 ./pkg/$version/jumper
cd ./pkg/$version
tar -czvf jumper-linux-amd64.tar.gz jumper LICENSE
rm jumper

cd ../../

echo "mv jumper-linux-386 ./pkg/$version/jumper"
mv jumper-linux-386 ./pkg/$version/jumper
cd ./pkg/$version
tar -czvf jumper-linux-386.tar.gz jumper LICENSE
rm jumper

cd ../../

echo "mv jumper-darwin-amd64.zip ./pkg/$version/jumper"
mv jumper-darwin-amd64.zip ./pkg/$version/jumper
cd ./pkg/$version/
zip jumper-darwin-amd64.zip jumper LICENSE
rm jumper

# darwin 386 is end of life but we need it to use fgo tool, which requires darwin-386
cp jumper-darwin-amd64.zip jumper-darwin-386.zip

rm LICENSE

cd ../../

brew unlink go && brew link go@1.14
echo 'export PATH="/usr/local/opt/go@1.14/bin:$PATH"' >> ~/.zshrc
. ~/.zshrc

fgo --pkg ./pkg build --token $token $version

brew unlink go && brew link go@1.15
echo 'export PATH="/usr/local/opt/go@1.15/bin:$PATH"' >> ~/.zshrc
. ~/.zshrc
