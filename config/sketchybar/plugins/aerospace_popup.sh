#!/bin/bash

source "$CONFIG_DIR/colors.sh"

# Check if popup is already open
DRAWING=$(sketchybar --query active_space | jq -r '.popup.drawing')

if [ "$DRAWING" = "on" ]; then
    sketchybar --set active_space popup.drawing=off
    exit 0
fi

# Hide first (prevents flicker)
sketchybar --set active_space popup.drawing=off

# Remove old popup items
sketchybar --remove '/popup.space\..*/'

WORKSPACES=$(aerospace list-workspaces --all)

while IFS= read -r ws; do
    APPS=$(aerospace list-windows --workspace "$ws" 2>/dev/null \
        | awk -F'|' '{gsub(/ /, "", $2); print $2}' \
        | sort -u)

    APP_ICONS=""
    while IFS= read -r app; do
        [ -z "$app" ] && continue
        ICON=$("$CONFIG_DIR/plugins/icon_map.sh" "$app" 2>/dev/null)
        APP_ICONS="$APP_ICONS $ICON"
    done <<< "$APPS"

    [ -z "$APP_ICONS" ] && APP_ICONS=" —"

    # Highlight current workspace
    if [ "$ws" = "$AEROSPACE_FOCUSED_WORKSPACE" ]; then
        BG_COLOR=$SPACES_COLOR   # strong accent (your beige)
        ICON_COLOR=$BASE         # dark text on light bg
        LABEL_COLOR=$BASE
    else
        BG_COLOR=0xFF1E1E2E      # solid background
        ICON_COLOR=$TEXT
        LABEL_COLOR=$SUBTLE
    fi

    sketchybar --add item popup.space.$ws popup.active_space \
               --set popup.space.$ws \
                   click_script="aerospace workspace $ws; sketchybar --set active_space popup.drawing=off" \
                   icon="$ws" \
                   icon.color=$ICON_COLOR \
                   label="$APP_ICONS" \
                   label.color=$LABEL_COLOR \
                   background.color=$BG_COLOR \
                   icon.font="$FONT:Bold:12.0" \
                   label.font="sketchybar-app-font:Regular:12.0" \
                   padding_left=10 \
                   padding_right=10
done <<< "$WORKSPACES"

# Show popup AFTER building (fixes glitch)
sketchybar --set active_space popup.drawing=on
