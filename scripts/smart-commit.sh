#!/bin/bash

# Smart Git Commit Script for Linux/macOS
# Анализирует изменения и генерирует осмысленное сообщение коммита

show_help() {
    echo "Usage: $0 [message]"
    echo "Smart git commit with auto-generated messages"
    echo ""
    echo "Options:"
    echo "  -h, --help    Show this help"
    echo "  message       Custom commit message"
}

get_git_changes() {
    local status=$(git status --porcelain 2>/dev/null)
    if [ -z "$status" ]; then
        echo "No changes to commit" >&2
        return 1
    fi
    echo "$status"
}

analyze_changes() {
    local status="$1"
    local added=0
    local modified=0
    local deleted=0
    local frontend=0
    local backend=0
    local config=0
    local docs=0
    
    while IFS= read -r line; do
        if [ -z "$line" ]; then continue; fi
        
        local file_status="${line:0:2}"
        local file_path="${line:3}"
        
        case "$file_status" in
            "A "|"??") ((added++)) ;;
            "M "|" M") ((modified++)) ;;
            "D "|" D") ((deleted++)) ;;
        esac
        
        case "$file_path" in
            frontend/*|*.svelte|*.ts|*.js|*.css|*.html) ((frontend++)) ;;
            *.go|go.mod|go.sum) ((backend++)) ;;
            *.json|*.config|*.yml|*.yaml|*prettierrc*|*tasks.json|*settings.json) ((config++)) ;;
            README*|*.md) ((docs++)) ;;
        esac
    done <<< "$status"
    
    echo "$added $modified $deleted $frontend $backend $config $docs"
}

generate_commit_message() {
    local stats=($1)
    local added=${stats[0]}
    local modified=${stats[1]}
    local deleted=${stats[2]}
    local frontend=${stats[3]}
    local backend=${stats[4]}
    local config=${stats[5]}
    local docs=${stats[6]}
    
    local type="feat"
    if [ $modified -gt $added ]; then
        type="fix"
    elif [ $deleted -gt 0 ]; then
        type="refactor"
    fi
    
    local description=""
    if [ $frontend -gt 0 ]; then
        description="frontend updates"
    fi
    if [ $backend -gt 0 ]; then
        if [ -n "$description" ]; then
            description="$description, backend changes"
        else
            description="backend changes"
        fi
    fi
    if [ $config -gt 0 ]; then
        if [ -n "$description" ]; then
            description="$description, configuration"
        else
            description="configuration"
        fi
    fi
    if [ $docs -gt 0 ]; then
        if [ -n "$description" ]; then
            description="$description, documentation"
        else
            description="documentation"
        fi
    fi
    
    if [ -z "$description" ]; then
        description="project updates"
    fi
    
    echo "$type: $description"
}

show_changes_preview() {
    local status="$1"
    echo "=== Git Changes Preview ==="
    
    while IFS= read -r line; do
        if [ -z "$line" ]; then continue; fi
        
        local file_status="${line:0:2}"
        local file_path="${line:3}"
        
        case "$file_status" in
            "A "|"??") echo -e "\033[32m  + $file_path\033[0m" ;;
            "M "|" M") echo -e "\033[33m  ~ $file_path\033[0m" ;;
            "D "|" D") echo -e "\033[31m  - $file_path\033[0m" ;;
        esac
    done <<< "$status"
    
    echo ""
}

# Main logic
if [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
    show_help
    exit 0
fi

status=$(get_git_changes)
if [ $? -ne 0 ]; then
    exit 1
fi

show_changes_preview "$status"

message="$1"
if [ -z "$message" ]; then
    stats=$(analyze_changes "$status")
    message=$(generate_commit_message "$stats")
    
    echo -e "\033[36mGenerated commit message: \033[0m'$message'"
    echo -n "Use this message? (Y/n/edit): "
    read -r confirm
    
    case "$confirm" in
        n|N)
            echo -n "Enter custom commit message: "
            read -r message
            ;;
        e|edit)
            echo -n "Edit message [$message]: "
            read -r custom_message
            if [ -n "$custom_message" ]; then
                message="$custom_message"
            fi
            ;;
    esac
fi

echo -e "\033[34mStaging all changes...\033[0m"
git add .

echo -e "\033[34mCommitting with message: '$message'\033[0m"
git commit -m "$message"

echo -n "Push to remote? (Y/n): "
read -r push_confirm

if [ "$push_confirm" != "n" ] && [ "$push_confirm" != "N" ]; then
    echo -e "\033[34mPushing to remote...\033[0m"
    git push
fi

echo -e "\033[32mDone! ✅\033[0m"
