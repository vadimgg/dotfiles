#!/bin/zsh

set -e

echo_header() {
  echo "\n\033[1;34m==> $1\033[0m"
}

is_installed() {
  command -v "$1" >/dev/null 2>&1
}

brew_installed() {
  brew list --formula "$1" >/dev/null 2>&1
}

cask_installed() {
  brew list --cask "$1" >/dev/null 2>&1
}

install_if_missing() {
  if ! brew_installed "$1"; then
    echo "Installing $1..."
    brew install "$1"
  else
    echo "$1 already installed."
  fi
}

install_cask_if_missing() {
  if ! cask_installed "$1"; then
    echo "Installing $1 (cask)..."
    brew install --cask "$1"
  else
    echo "$1 (cask) already installed."
  fi
}

# 1. Install Xcode CLI Tools if needed
echo_header "Checking Xcode Command Line Tools..."
if ! xcode-select -p >/dev/null 2>&1; then
  echo "Installing Xcode Command Line Tools..."
  xcode-select --install
else
  echo "Xcode CLI tools already installed."
fi

# 2. Install Homebrew
echo_header "Checking Homebrew..."
if ! is_installed brew; then
  echo "Installing Homebrew..."
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
else
  echo "Homebrew already installed."
fi

# Disable analytics
brew analytics off

# 3. Install Brew Formulae
echo_header "Installing brew formulae..."

FORMULAE=(
  helix
  ghostty
  veeso/termscp/termscp
  wget
  lsd
  bat
  yt-dlp
  lazygit
)

for pkg in "${FORMULAE[@]}"; do
  install_if_missing "$pkg"
done

# 4. Install Casks
echo_header "Installing brew casks..."

CASKS=(
  transmission
  whatsapp
  vlc
)

for cask in "${CASKS[@]}"; do
  install_cask_if_missing "$cask"
done

# 5. Cleanup
echo_header "Cleaning up..."
brew cleanup

