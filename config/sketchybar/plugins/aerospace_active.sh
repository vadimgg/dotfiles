#!/bin/bash

FOCUSED=$(aerospace list-workspaces --focused 2>/dev/null)
[ -z "$FOCUSED" ] && FOCUSED="$AEROSPACE_FOCUSED_WORKSPACE"
[ -z "$FOCUSED" ] && FOCUSED="?"

APP=$(aerospace list-windows --workspace "$FOCUSED" 2>/dev/null \
  | awk -F'|' '{gsub(/ /, "", $2); print $2}' \
  | head -n 1)

[ -z "$APP" ] && APP="—"

sketchybar --set active_space \
    icon="$FOCUSED" \
    label="$APP"
