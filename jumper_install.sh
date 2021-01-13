#!/usr/bin/env sh

brew uninstall jumper

wget https://raw.githubusercontent.com/asannikov/jumper/master/homebrew/jumper.rb

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    PATH_TO_FORMULA="/home/linuxbrew/.linuxbrew/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/jumper.rb"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    PATH_TO_FORMULA="/usr/local/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/jumper.rb"
elif [[ "$OSTYPE" == "cygwin" ]]; then
    PATH_TO_FORMULA="/home/linuxbrew/.linuxbrew/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/jumper.rb"
elif [[ "$OSTYPE" == "freebsd"* ]]; then
    PATH_TO_FORMULA="/home/linuxbrew/.linuxbrew/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/jumper.rb"
else
    echo "unknown OS $OSTYPE"
    exit 1
fi

rm -fr $PATH_TO_FORMULA
brew create https://github.com/asannikov/jumper
yes | cp jumper.rb $PATH_TO_FORMULA
rm jumper.rb
brew install --build-from-source jumper