#!/bin/zsh

# Install xCode cli tools
echo "Installing commandline tools..."
xcode-select --install

# Homebrew
## Install
echo "Installing Brew..."
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
brew analytics off

# Install more recent versions of some macOS tools.
brew helix
brew 'ghostty'
brew 'helix'
brew 'veeso/termscp/termscp'
brew 'wget'
brew 'yt-dlp'

cask 'transmission'
cask 'whatsapp'
cask 'vlc'

# Remove outdated versions from the cellar.
brew cleanup
