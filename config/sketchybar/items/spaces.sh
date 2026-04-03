#!/bin/bash

SPACE_ICONS=("H" "J" "K" "L" "Y" "U" "I" "O")

SPACE=(
    icon.padding_left=8
    icon.padding_right=8
    icon.font="$FONT:Bold:12.0"
    icon.color=$SPACES_COLOR
    label.drawing=off
    background.height=24
    background.color=$OVERLAY
    background.corner_radius=8
    background.drawing=off
    script="$PLUGIN_DIR/aerospace.sh"
    click_script="$PLUGIN_DIR/aerospace_popup.sh"
)

for i in "${!SPACE_ICONS[@]}"
do
    sid=${SPACE_ICONS[i]}
    sketchybar --add item space.$sid left \
               --set space.$sid "${SPACE[@]}" \
                                icon=$sid \
               --subscribe space.$sid aerospace_workspace_change
done

sketchybar --add bracket spaces '/space\..*/' \
           --set spaces background.color=$TRANSPARENT \
                        background.corner_radius=8 \
                        background.height=24
