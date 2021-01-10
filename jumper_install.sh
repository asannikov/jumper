#!/usr/bin/env sh

wget https://raw.githubusercontent.com/asannikov/jumper/master/homebrew/jumper.rb
rm -fr /usr/local/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/jumper.rb
brew create https://github.com/asannikov/jumper
yes | cp jumper.rb /usr/local/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/jumper.rb
rm jumper.rb
brew install --build-from-source jumper