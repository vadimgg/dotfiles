#!/bin/bash

source "$CONFIG_DIR/colors.sh"

FOCUSED="$FOCUSED_WORKSPACE"

for sid in H J K L Y U I O; do
    if [ "$sid" = "$FOCUSED" ]; then
        sketchybar --set space.$sid \
            icon.color=$BASE \
            background.drawing=on \
            background.color=$SPACES_COLOR
    else
        sketchybar --set space.$sid \
            icon.color=$SPACES_COLOR \
            background.drawing=off
    fi
done
