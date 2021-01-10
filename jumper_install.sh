#!/usr/bin/env sh

brew create https://github.com/asannikov/jumper
cp ./homebrew/jumper.rb /usr/local/Homebrew/Library/Taps/homebrew/homebrew-core/Formula/jumper.rb
brew install --build-from-source jumper