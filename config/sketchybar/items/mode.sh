#!/bin/bash

sketchybar --add event aerospace_mode_changed

sketchybar --add item aerospace_mode left \
           --set aerospace_mode \
               drawing=off \
               icon.drawing=off \
               label="" \
               label.color=$LOVE \
               background.color=$OVERLAY \
               background.corner_radius=10 \
               background.height=24 \
               background.drawing=off \
               label.padding_left=8 \
               label.padding_right=8 \
               padding_left=4 \
               padding_right=4 \
               script="$PLUGIN_DIR/aerospace_mode.sh"

sketchybar --subscribe aerospace_mode aerospace_mode_changed

sketchybar --set aerospace_mode drawing=off background.drawing=off label=""
