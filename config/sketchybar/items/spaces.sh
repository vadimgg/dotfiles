#!/bin/bash

ACTIVE_SPACE=(
    icon.font="$FONT:Bold:12.0"
    icon.color=$TEXT
    label.font="$FONT:Bold:12.0"
    label.color=$TEXT
    icon.padding_left=10
    icon.padding_right=6
    label.padding_right=10

    background.color=$OVERLAY
    background.corner_radius=10
    background.height=24

    script="$PLUGIN_DIR/aerospace_active.sh"
    click_script="$PLUGIN_DIR/aerospace_popup.sh"

    # ✅ THIS enables popup (like your sound widget)
    popup.horizontal=off
    popup.align=left
    popup.height=28
    popup.y_offset=10
    # popup.background.color=$BASE
    popup.background.color=0xFF1E1E2E
    popup.background.corner_radius=10
    popup.background.border_width=1
    popup.background.border_color=$OVERLAY
)

# sketchybar --add item active_space left \
#            --set active_space "${ACTIVE_SPACE[@]}" \
#            --subscribe active_space aerospace_workspace_change front_app_switched
sketchybar --add item active_space left \
           --set active_space "${ACTIVE_SPACE[@]}" \
           --subscribe active_space aerospace_workspace_change front_app_switched mouse.exited.global
