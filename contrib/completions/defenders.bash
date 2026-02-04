#!/usr/bin/env bash
# Bash completion for defenders CLI

_defenders_completions() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    # Main commands
    local commands="conf get-token cado prme release pr help"
    
    # Help flags
    local help_flags="-h --help"

    # If we're at the first argument (command level)
    if [ $COMP_CWORD -eq 1 ]; then
        COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
        return 0
    fi

    # Get the subcommand
    local subcommand="${COMP_WORDS[1]}"

    case "${subcommand}" in
        conf)
            local conf_opts="-s --show -p --path -r --reset -h --help"
            COMPREPLY=( $(compgen -W "${conf_opts}" -- ${cur}) )
            ;;
        get-token)
            COMPREPLY=( $(compgen -W "${help_flags}" -- ${cur}) )
            ;;
        cado)
            case "${prev}" in
                -u|--url|-j|--json)
                    # These require arguments, don't suggest anything
                    return 0
                    ;;
                *)
                    local cado_opts="-u --url -j --json -h --help"
                    COMPREPLY=( $(compgen -W "${cado_opts}" -- ${cur}) )
                    ;;
            esac
            ;;
        prme)
            case "${prev}" in
                -t|--title|-d|--description|-s|--source|-b|--target)
                    # These require arguments, don't suggest anything
                    return 0
                    ;;
                *)
                    local prme_opts="-t --title -d --description -s --source -b --target -h --help"
                    COMPREPLY=( $(compgen -W "${prme_opts}" -- ${cur}) )
                    ;;
            esac
            ;;
        release)
            case "${prev}" in
                -u|--url)
                    # Requires argument
                    return 0
                    ;;
                *)
                    local release_opts="-u --url -m --monitor -h --help"
                    COMPREPLY=( $(compgen -W "${release_opts}" -- ${cur}) )
                    ;;
            esac
            ;;
        pr)
            case "${prev}" in
                -u|--url)
                    # Requires argument
                    return 0
                    ;;
                *)
                    local pr_opts="-u --url -a --approve -c --complete -h --help"
                    COMPREPLY=( $(compgen -W "${pr_opts}" -- ${cur}) )
                    ;;
            esac
            ;;
        help)
            COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
            ;;
    esac
}

complete -F _defenders_completions defenders
