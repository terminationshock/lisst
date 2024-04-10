#!/usr/bin/env bash

lisst_completion() {
    COMPREPLY=($(lisst --completion "$COMP_LINE" "$2"))
    if [ "${COMPREPLY[0]}" == "-" ]; then
        COMPREPLY=($(compgen -A command "$2"))
    fi
}

complete -F "lisst_completion" lisst
