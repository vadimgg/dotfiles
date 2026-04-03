#!/bin/bash
source "$CONFIG_DIR/icons.sh"

CITY="Delhi"

MAX_RETRIES=3
DELAY=5 # 秒
DEFAULT_TEMPERATURE="N/A"
DEFAULT_ICON="$WEATHER_DEFAULT" 

TEMPERATURE=$DEFAULT_TEMPERATURE
CUSTOM_ICON=$DEFAULT_ICON

for ((i = 1; i <= MAX_RETRIES; i++)); do
    WEATHER_JSON=$(curl -s "wttr.in/${CITY}?format=j1")
    if [[ -n "$WEATHER_JSON" && "$(echo $WEATHER_JSON | jq -r '.current_condition[0].temp_C')" != "null" ]]; then
        TEMPERATURE=$(echo $WEATHER_JSON | jq -r '.current_condition[0].temp_C')
        DESCRIPTION=$(echo $WEATHER_JSON | jq -r '.current_condition[0].weatherDesc[0].value')
        # echo $TEMPERATURE
        # echo $DESCRIPTION
        break
    else
        echo "Retry #$i: Failed to fetch weather data." >> "/tmp/weather.log"
        sleep $DELAY
    fi
done

if [[ -z "$WEATHER_JSON" || "$TEMPERATURE" == "null" ]]; then
    TEMPERATURE=$DEFAULT_TEMPERATURE
    CUSTOM_ICON=$DEFAULT_ICON
else
    case $DESCRIPTION in
    *[Ss]un*|*[Cc]lear*) CUSTOM_ICON="$WEATHER_SUNNY" ;;                           
    *[Cc]loud*|*[Oo]vercast*|*[Pp]artly*) CUSTOM_ICON="$WEATHER_CLOUDY" ;;         
    *[Rr]ain*|*[Ss]hower*) CUSTOM_ICON="$WEATHER_RAINY" ;;                         
    *[Tt]hunder*|*[Ss]torm*) CUSTOM_ICON="$WEATHER_STORMY" ;;                      
    *[Ss]now*|*[Bb]lizzard*) CUSTOM_ICON="$WEATHER_SNOWY" ;;                       
    *[Mm]ist*) CUSTOM_ICON="$WEATHER_MISTY" ;;                                     
    *[Hh]aze*) CUSTOM_ICON="$WEATHER_HAZE" ;;                                      
    *) CUSTOM_ICON="$WEATHER_DEBUG" ;;                                             
    esac
fi

sketchybar --set $NAME icon="$CUSTOM_ICON" label="${TEMPERATURE}°C"
