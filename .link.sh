#!/bin/bash

# Colors
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m" # No Color

DOTFILES_DIR="$HOME/dotfiles"
HOME_SOURCE="$DOTFILES_DIR/home"
CONFIG_SOURCE="$DOTFILES_DIR/config"
CONFIG_TARGET="$HOME/.config"

# Ensure directories exist
[[ ! -d "$HOME_SOURCE" ]] && echo -e "${RED}Missing: $HOME_SOURCE${NC}" && exit 1
[[ ! -d "$CONFIG_SOURCE" ]] && echo -e "${RED}Missing: $CONFIG_SOURCE${NC}" && exit 1
mkdir -p "$CONFIG_TARGET"

echo -e "${YELLOW}Linking dotfiles into \$HOME...${NC}"

# Link dotfiles from ~/dotfiles/home to ~/
for SOURCE in "$HOME_SOURCE"/.*; do
  BASENAME=$(basename "$SOURCE")
  [[ "$BASENAME" == "." || "$BASENAME" == ".." ]] && continue

  TARGET="$HOME/$BASENAME"

  if [ -e "$TARGET" ] && [ ! -L "$TARGET" ]; then
    echo -e "${YELLOW}⚠️  Backing up $TARGET to $TARGET.backup${NC}"
    mv "$TARGET" "$TARGET.backup"
  fi

  [ -L "$TARGET" ] && rm "$TARGET"

  ln -s "$SOURCE" "$TARGET"
  echo -e "${GREEN}✔ Linked $SOURCE → $TARGET${NC}"
done

echo -e "${YELLOW}Linking config files into ~/.config...${NC}"

# Link configs from ~/dotfiles/config to ~/.config/
for SOURCE in "$CONFIG_SOURCE"/*; do
  BASENAME=$(basename "$SOURCE")
  TARGET="$CONFIG_TARGET/$BASENAME"

  if [ -e "$TARGET" ] && [ ! -L "$TARGET" ]; then
    echo -e "${YELLOW}⚠️  Backing up $TARGET to $TARGET.backup${NC}"
    mv "$TARGET" "$TARGET.backup"
  fi

  [ -L "$TARGET" ] && rm "$TARGET"

  ln -s "$SOURCE" "$TARGET"
  echo -e "${GREEN}✔ Linked $SOURCE → $TARGET${NC}"
done



echo -e "${YELLOW}Linking scripts from dotfiles/bin into ~/.local/bin...${NC}"

BIN_SOURCE="$DOTFILES_DIR/bin"
BIN_TARGET="$HOME/.local/bin"

mkdir -p "$BIN_TARGET"

for SCRIPT in "$BIN_SOURCE"/*; do
  BASENAME=$(basename "$SCRIPT")
  TARGET="$BIN_TARGET/$BASENAME"

  if [ -e "$TARGET" ] && [ ! -L "$TARGET" ]; then
    echo -e "${YELLOW}⚠️  Backing up $TARGET to $TARGET.backup${NC}"
    mv "$TARGET" "$TARGET.backup"
  fi

  [ -L "$TARGET" ] && rm "$TARGET"

  ln -s "$SCRIPT" "$TARGET"
  chmod +x "$SCRIPT" # Ensure the original script is executable
  echo -e "${GREEN}✔ Linked $SCRIPT → $TARGET${NC}"
done

# Check if ~/.local/bin is in PATH
if [[ ":$PATH:" != *":$BIN_TARGET:"* ]]; then
  echo -e "${RED}⚠️  ~/.local/bin is not in your PATH.${NC}"
  echo -e "${YELLOW}👉 Add this line to your shell config (~/.zshrc or ~/.bashrc):${NC}"
  echo -e "   ${GREEN}export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
fi
