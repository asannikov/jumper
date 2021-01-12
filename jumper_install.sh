#!/usr/bin/env sh

wget https://raw.githubusercontent.com/asannikov/jumper/master/homebrew/jumper.rb
rm -fr /usr/local/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/jumper.rb
brew create https://github.com/asannikov/jumper

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    yes | cp jumper.rb /home/linuxbrew/.linuxbrew/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/jumper.rb
elif [[ "$OSTYPE" == "darwin"* ]]; then
    yes | cp jumper.rb /usr/local/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/jumper.rb
elif [[ "$OSTYPE" == "cygwin" ]]; then
    yes | cp jumper.rb /home/linuxbrew/.linuxbrew/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/jumper.rb
elif [[ "$OSTYPE" == "freebsd"* ]]; then
    yes | cp jumper.rb /home/linuxbrew/.linuxbrew/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/jumper.rb
else
    echo "unknown OS $OSTYPE"
    exit 1
fi

rm jumper.rb
brew install --build-from-source jumper