#!/bin/bash

FRONT_APP=(
    label.font="$FONT:Bold:12.0"
    label.color=$TEXT
    icon.font="sketchybar-app-font:Regular:14.0"
    icon.color=$TEXT
    icon.padding_left=6
    label.padding_left=4
    label.padding_right=8
    background.color=$TRANSPARENT
    script="$PLUGIN_DIR/front_app.sh"
)

sketchybar --add item front_app left \
           --set front_app "${FRONT_APP[@]}" \
           --subscribe front_app front_app_switched
