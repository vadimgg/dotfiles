#!/bin/bash
echo "MODE=$MODE" >> /tmp/mode_debug.log

case "$MODE" in
    main)
        sketchybar --set aerospace_mode drawing=off background.drawing=off label=""
        ;;
    service)
        sketchybar --set aerospace_mode drawing=on background.drawing=on label="SVC"
        ;;
    open)
        sketchybar --set aerospace_mode drawing=on background.drawing=on label="OPEN"
        ;;
    *)
        sketchybar --set aerospace_mode drawing=on background.drawing=on label="  $MODE"
        ;;
esac
