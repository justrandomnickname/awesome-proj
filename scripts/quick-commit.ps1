# Quick Commit Script
# Быстрый коммит с кратким анализом

$status = git status --porcelain
if (-not $status) {
    Write-Host "No changes to commit" -ForegroundColor Yellow
    exit 1
}

Write-Host "=== Changes ===" -ForegroundColor Cyan
git status --short

$message = Read-Host "`nEnter commit message (or press Enter for auto-generated)"

if (-not $message) {
    $fileCount = ($status | Measure-Object).Count
    $message = "update: $fileCount files changed"
    
    # Улучшаем сообщение на основе типов файлов
    if ($status -match "frontend/") {
        $message = "frontend: updates and improvements"
    } elseif ($status -match "\.go") {
        $message = "backend: code improvements"
    } elseif ($status -match "\.md|README") {
        $message = "docs: update documentation"
    } elseif ($status -match "\.json|\.config") {
        $message = "config: update project configuration"
    }
}

Write-Host "Committing: '$message'" -ForegroundColor Green

git add .
git commit -m $message

$push = Read-Host "Push? (Y/n)"
if ($push -ne "n") {
    git push
    Write-Host "Pushed! ✅" -ForegroundColor Green
}
