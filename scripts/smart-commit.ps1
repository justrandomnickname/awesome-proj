# Smart Git Commit Script
# Анализирует изменения и генерирует осмысленное сообщение коммита

param(
    [string]$message = ""
)

function Get-GitChanges {
    $status = git status --porcelain
    if (-not $status) {
        Write-Host "No changes to commit" -ForegroundColor Yellow
        return $false
    }
    
    $changes = @{
        Added = @()
        Modified = @()
        Deleted = @()
        Renamed = @()
    }
    
    foreach ($line in $status) {
        $statusCode = $line.Substring(0, 2)
        $file = $line.Substring(3)
        
        switch ($statusCode.Trim()) {
            "A" { $changes.Added += $file }
            "M" { $changes.Modified += $file }
            "D" { $changes.Deleted += $file }
            "R" { $changes.Renamed += $file }
            "??" { $changes.Added += $file }
        }
    }
    
    return $changes
}

function Generate-CommitMessage {
    param($changes)
    
    $messages = @()
    $type = "feat"
    
    # Определяем тип коммита
    if ($changes.Added.Count -gt $changes.Modified.Count) {
        $type = "feat"
    } elseif ($changes.Modified.Count -gt 0) {
        $type = "fix"
    } elseif ($changes.Deleted.Count -gt 0) {
        $type = "refactor"
    }
    
    # Анализируем файлы
    $categories = @{
        Frontend = @()
        Backend = @()
        Config = @()
        Docs = @()
    }
    
    foreach ($file in ($changes.Added + $changes.Modified + $changes.Deleted)) {
        if ($file -match "frontend/|\.svelte|\.ts|\.js|\.css|\.html") {
            $categories.Frontend += $file
        } elseif ($file -match "\.go|go\.mod|go\.sum") {
            $categories.Backend += $file
        } elseif ($file -match "\.json|\.config|\.yml|\.yaml|\.prettierrc|tasks\.json|settings\.json") {
            $categories.Config += $file
        } elseif ($file -match "README|\.md") {
            $categories.Docs += $file
        }
    }
    
    # Генерируем описание
    $description = ""
    
    if ($categories.Frontend.Count -gt 0) {
        $description += "frontend updates"
    }
    if ($categories.Backend.Count -gt 0) {
        if ($description) { $description += ", " }
        $description += "backend changes"
    }
    if ($categories.Config.Count -gt 0) {
        if ($description) { $description += ", " }
        $description += "configuration"
    }
    if ($categories.Docs.Count -gt 0) {
        if ($description) { $description += ", " }
        $description += "documentation"
    }
    
    if (-not $description) {
        $description = "project updates"
    }
    
    return "${type}: ${description}"
}

function Show-ChangesPreview {
    param($changes)
    
    Write-Host "=== Git Changes Preview ===" -ForegroundColor Cyan
    
    if ($changes.Added.Count -gt 0) {
        Write-Host "Added files:" -ForegroundColor Green
        foreach ($file in $changes.Added) {
            Write-Host "  + $file" -ForegroundColor Green
        }
    }
    
    if ($changes.Modified.Count -gt 0) {
        Write-Host "Modified files:" -ForegroundColor Yellow
        foreach ($file in $changes.Modified) {
            Write-Host "  ~ $file" -ForegroundColor Yellow
        }
    }
    
    if ($changes.Deleted.Count -gt 0) {
        Write-Host "Deleted files:" -ForegroundColor Red
        foreach ($file in $changes.Deleted) {
            Write-Host "  - $file" -ForegroundColor Red
        }
    }
    
    Write-Host ""
}

# Основная логика
$changes = Get-GitChanges
if (-not $changes) {
    exit 1
}

Show-ChangesPreview $changes

if (-not $message) {
    $message = Generate-CommitMessage $changes
    Write-Host "Generated commit message: " -NoNewline -ForegroundColor Cyan
    Write-Host "'$message'" -ForegroundColor White
    
    $confirm = Read-Host "Use this message? (Y/n/edit/append)"
    
    if ($confirm -eq "n") {
        $message = Read-Host "Enter custom commit message"
    } elseif ($confirm -eq "edit" -or $confirm -eq "e") {
        $newMessage = Read-Host "Edit message [$message]"
        if ($newMessage.Trim()) {
            $message = $newMessage
        }
    } elseif ($confirm -eq "append" -or $confirm -eq "a") {
        $addition = Read-Host "Add to message '$message'"
        if ($addition.Trim()) {
            $message = "$message`n`n$addition"
        }
    }
}

Write-Host "Staging all changes..." -ForegroundColor Blue
git add .

Write-Host "Committing with message: '$message'" -ForegroundColor Blue
git commit -m $message

$pushConfirm = Read-Host "Push to remote? (Y/n)"
if ($pushConfirm -ne "n") {
    Write-Host "Pushing to remote..." -ForegroundColor Blue
    git push
}

Write-Host "Done! ✅" -ForegroundColor Green
