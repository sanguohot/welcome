#!/bin/bash
function SetDockerImg()
{
    IMG="error: WELCOME_IMG env or param required, e.g. sanguohot/welcome:1.0,sanguohot/welcome:latest"
    if [ -n "$1" ]; then
      IMG=$1
    elif [ -n "WELCOME_IMG" ]; then
      IMG=$WELCOME_IMG
    fi
    echo "${IMG}"
}
