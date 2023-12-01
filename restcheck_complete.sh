#!/bin/bash

_restcheck()
{
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="-init -msg -initapi -check -checkall"

    if [[ ${cur} == -* ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi

    if [[ ${prev} == -check || ${prev} == -checkall ]] ; then
        COMPREPLY=( $(compgen -d -- ${cur}) )
        return 0
    fi

    if [[ ${COMP_WORDS[COMP_CWORD-2]} == -check || ${COMP_WORDS[COMP_CWORD-2]} == -checkall ]] ; then
        COMPREPLY=( "-save" )
        return 0
    fi
}

complete -F _restcheck restcheck
