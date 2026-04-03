#!/bin/bash

source "$CONFIG_DIR/colors.sh"

POPUP_OFF="sketchybar --set space.popup popup.drawing=off"

if [ "$SENDER" = "mouse.clicked" ]; then
    # Toggle popup
    DRAWING=$(sketchybar --query space.popup | jq -r '.popup.drawing')
    if [ "$DRAWING" = "on" ]; then
        $POPUP_OFF
        exit 0
    fi

    # Clear existing popup items
    sketchybar --remove '/popup.space\..*/'

    # Get all workspaces and their apps
    WORKSPACES=$(aerospace list-workspaces --all)

    while IFS= read -r ws; do
        APPS=$(aerospace list-windows --workspace "$ws" 2>/dev/null | awk -F'|' '{gsub(/ /, "", $2); print $2}' | sort -u)
        
        APP_ICONS=""
        while IFS= read -r app; do
            [ -z "$app" ] && continue
            ICON=$("$CONFIG_DIR/plugins/icon_map.sh" "$app" 2>/dev/null)
            APP_ICONS="$APP_ICONS $ICON"
        done <<< "$APPS"

        [ -z "$APP_ICONS" ] && APP_ICONS=" —"

        sketchybar --add item popup.space.$ws popup.space.H \
                   --set popup.space.$ws \
                       icon="$ws" \
                       icon.font="$FONT:Bold:12.0" \
                       icon.color=$SPACES_COLOR \
                       icon.padding_left=12 \
                       icon.padding_right=8 \
                       label="$APP_ICONS" \
                       label.font="sketchybar-app-font:Regular:12.0" \
                       label.color=$TEXT \
                       label.padding_right=12 \
                       background.color=$TRANSPARENT
    done <<< "$WORKSPACES"

    sketchybar --set space.H popup.drawing=on \
                             popup.align=left \
                             popup.horizontal=false \
                             popup.height=28 \
                             popup.background.color=$OVERLAY \
                             popup.background.corner_radius=10 \
                             popup.background.border_width=1 \
                             popup.background.border_color=$MUTED
fi
