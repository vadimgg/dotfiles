#!/bin/bash

CLOCK=(
    icon.padding_left=12
    label.padding_right=12
    label.font="$FONT:Bold:12.0"
    label.color=$CLOCK_COLOR
    icon.font="$FONT:Bold:16.0"
    icon.color=$CLOCK_COLOR
    background.height=26
    background.color=$OVERLAY
    background.corner_radius=10
    script="$PLUGIN_DIR/clock.sh"
)

sketchybar --add item clock right --set clock "${CLOCK[@]}" update_freq=15
