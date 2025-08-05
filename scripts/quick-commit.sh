#!/bin/bash

# Quick Commit Script for Linux/macOS
# Быстрый коммит с кратким анализом

status=$(git status --porcelain 2>/dev/null)
if [ -z "$status" ]; then
    echo "No changes to commit"
    exit 1
fi

echo "=== Changes ==="
git status --short

echo ""
echo -n "Enter commit message (or press Enter for auto-generated): "
read -r message

if [ -z "$message" ]; then
    file_count=$(echo "$status" | wc -l)
    message="update: $file_count files changed"
    
    # Улучшаем сообщение на основе типов файлов
    if echo "$status" | grep -q "frontend/"; then
        message="frontend: updates and improvements"
    elif echo "$status" | grep -q "\.go"; then
        message="backend: code improvements"
    elif echo "$status" | grep -q "\.md\|README"; then
        message="docs: update documentation"
    elif echo "$status" | grep -q "\.json\|\.config"; then
        message="config: update project configuration"
    fi
fi

echo -e "\033[32mCommitting: '$message'\033[0m"

git add .
git commit -m "$message"

echo -n "Push? (Y/n): "
read -r push_confirm

if [ "$push_confirm" != "n" ] && [ "$push_confirm" != "N" ]; then
    git push
    echo -e "\033[32mPushed! ✅\033[0m"
fi
