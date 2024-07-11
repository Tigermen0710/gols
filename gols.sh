#!/bin/bash
if [ ! -z "$1" ]; then
    cd $1
fi
exec ~/gols/ls
