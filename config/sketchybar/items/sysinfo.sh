#!/bin/bash

# CPU Component
CPU=(
    icon=$SYSINFO_CPU
    icon.font="$FONT:Bold:16.0"
    icon.color=$ROSE
    icon.padding_left=4
    icon.padding_right=4
    label.font="$FONT:Bold:12.0"
    label.color=$ROSE
    label.padding_left=4
    label.padding_right=12
    update_freq=2
    script="$PLUGIN_DIR/cpu.sh"
)

# Memory Component
MEMORY=(
    icon=$SYSINFO_MEMORY
    icon.font="$FONT:Bold:16.0"
    icon.color=$GOLD
    icon.padding_left=4
    icon.padding_right=4
    label.font="$FONT:Bold:12.0"
    label.color=$GOLD
    label.padding_left=4
    label.padding_right=4
    update_freq=2
    script="$PLUGIN_DIR/memory.sh"
)

# Disk Component
DISK=(
    icon=$SYSINFO_DISK
    icon.font="$FONT:Bold:16.0"
    icon.color=$FOAM
    icon.padding_left=4
    icon.padding_right=4
    label.font="$FONT:Bold:12.0"
    label.color=$FOAM
    label.padding_left=4
    label.padding_right=4
    update_freq=60
    script="$PLUGIN_DIR/disk.sh"
)

# Speed Down Component
SPEED_DOWN_ITEM=(
    icon=$SPEED_DOWN
    icon.font="$FONT:Bold:14.0"      # reduced from 16
    icon.color=$IRIS
    icon.padding_left=4
    icon.padding_right=2             # tightened
    label.font="$FONT:Bold:11.0"     # reduced from 12
    label.color=$IRIS
    label.padding_left=2             # tightened
    label.padding_right=4
    label.width=60                   # reduced from 72
)

# Speed Up Component
SPEED_UP_ITEM=(
    icon=$SPEED_UP
    icon.font="$FONT:Bold:14.0"
    icon.color=$ROSE
    icon.padding_left=4
    icon.padding_right=2
    label.font="$FONT:Bold:11.0"
    label.color=$ROSE
    label.padding_left=2
    label.padding_right=4
    label.width=60
    update_freq=2
    script="$PLUGIN_DIR/speed.sh"
)

sketchybar --add item spacer_left left \
           --set spacer_left width=4 \
                             drawing=on

sketchybar --add item cpu left \
           --set cpu "${CPU[@]}" \
           \
           --add item memory left \
           --set memory "${MEMORY[@]}" \
           --subscribe memory mouse.clicked \
           \
           --add item disk left \
           --set disk "${DISK[@]}" \
           \
           --add item speed_down left \
           --set speed_down "${SPEED_DOWN_ITEM[@]}" \
           \
           --add item speed_up left \
           --set speed_up "${SPEED_UP_ITEM[@]}" \
           \
           --add bracket sysinfo cpu memory disk speed_down speed_up surge \
           --set sysinfo background.color=$OVERLAY \
                        background.height=26 \
                        background.corner_radius=10 \
                        background.padding_left=4 \
                        background.padding_right=4

sketchybar --add item spacer_right left \
           --set spacer_right width=4 \
                              drawing=on
