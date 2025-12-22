#!/usr/bin/env bash

# -------- CONFIG --------
HUGO_ROOT="$HOME/personal_blog"
NOTES_DIR="$HUGO_ROOT/content/daily"
EDITOR_CMD="hx"
DATE_FORMAT="%Y-%m-%d"
TIME_FORMAT="%H:%M"
# ------------------------

TODAY=$(date +"$DATE_FORMAT")
NOW=$(date +"$TIME_FORMAT")
FILE_PATH="$NOTES_DIR/$TODAY.md"

# Ensure notes directory exists
mkdir -p "$NOTES_DIR"

# Create file if it does not exist
if [[ ! -f "$FILE_PATH" ]]; then
  cat <<EOF > "$FILE_PATH"
---
title: "$TODAY"
date: $(date -Iseconds)
draft: false
---
EOF
fi

# Append time heading
echo -e "\n## $NOW\n" >> "$FILE_PATH"

# Open in Helix
exec "$EDITOR_CMD" "$FILE_PATH"

