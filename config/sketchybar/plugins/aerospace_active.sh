#!/bin/bash

# ✅ Get focused workspace reliably
FOCUSED=$(aerospace list-workspaces --focused)

# Fallback (just in case)
[ -z "$FOCUSED" ] && FOCUSED="?"

# Get first app in that workspace
APP=$(aerospace list-windows --workspace "$FOCUSED" 2>/dev/null \
  | awk -F'|' '{gsub(/ /, "", $2); print $2}' \
  | head -n 1)

[ -z "$APP" ] && APP="—"

# Update UI
sketchybar --set active_space \
    icon="$FOCUSED" \
    label="$APP"
