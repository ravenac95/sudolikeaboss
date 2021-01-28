#!/bin/bash
TITLE=$(osascript -e 'tell application "iTerm" to get the name of the front window')
if [[ ${TITLE} =~ "@" ]] ; then
    FIRST="${TITLE#*@}"
    if [[ ${FIRST} =~ ":" ]] ; then
        HOST="${FIRST%%:*}"
        export SUDOLIKEABOSS_DEFAULT_HOST="sudolikeaboss://${HOST}"
    fi;
fi;

/usr/local/bin/sudolikeaboss
