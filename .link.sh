#!/bin/bash

# Colors
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m" # No Color

# Set paths
DOTFILES_HOME="$HOME/dotfiles/home"

# Ensure the dotfiles/home directory exists
if [ ! -d "$DOTFILES_HOME" ]; then
  echo -e "${RED}Error: $DOTFILES_HOME does not exist.${NC}"
  exit 1
fi

# Loop through all files in dotfiles/home
for SOURCE in "$DOTFILES_HOME"/.*; do
  BASENAME=$(basename "$SOURCE")

  # Skip special entries
  [[ "$BASENAME" == "." || "$BASENAME" == ".." ]] && continue

  TARGET="$HOME/$BASENAME"

  # Backup existing non-symlink files
  if [ -e "$TARGET" ] && [ ! -L "$TARGET" ]; then
    BACKUP="$TARGET.backup"
    echo -e "${YELLOW}⚠️  Backing up existing $TARGET to $BACKUP${NC}"
    mv "$TARGET" "$BACKUP"
  fi

  # Remove existing symlinks
  if [ -L "$TARGET" ]; then
    echo -e "${YELLOW}🔁 Removing existing symlink $TARGET${NC}"
    rm "$TARGET"
  fi

  # Create the symlink
  ln -s "$SOURCE" "$TARGET"
  echo -e "${GREEN}✔ Linked $SOURCE → $TARGET${NC}"
done
