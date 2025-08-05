# Git Commit Scripts

Набор скриптов для автоматизации git коммитов с умной генерацией сообщений.

## 📁 Файлы

### Windows (PowerShell)
- `smart-commit.ps1` - умный коммит с анализом изменений
- `quick-commit.ps1` - быстрый коммит

### Linux/macOS (Bash)  
- `smart-commit.sh` - умный коммит для Unix-систем
- `quick-commit.sh` - быстрый коммит для Unix-систем

### Универсальные (Node.js)
- `smart-commit.js` - работает везде где есть Node.js
- `quick-commit.js` - универсальная версия

## 🚀 Использование

### Через VS Code Tasks
1. `Ctrl+Shift+P` → `Tasks: Run Task`
2. Выберите нужную задачу:
   - `Smart Commit (Windows)` - для Windows
   - `Smart Commit (Universal)` - Node.js версия
   - `Quick Commit (Windows)` - быстрый для Windows  
   - `Quick Commit (Universal)` - быстрый Node.js

### Напрямую

#### Windows
```powershell
# Умный коммит
.\scripts\smart-commit.ps1

# Быстрый коммит
.\scripts\quick-commit.ps1
```

#### Linux/macOS
```bash
# Сделать исполняемыми (один раз)
chmod +x scripts/*.sh

# Умный коммит
./scripts/smart-commit.sh

# Быстрый коммит  
./scripts/quick-commit.sh
```

#### Универсальный (Node.js)
```bash
# Умный коммит
node scripts/smart-commit.js

# Быстрый коммит
node scripts/quick-commit.js
```

## ✨ Возможности

### Smart Commit
- 🔍 Анализирует все изменения в проекте
- 🏷️ Автоматически определяет тип коммита (`feat`, `fix`, `refactor`)
- 📂 Категоризирует файлы (frontend, backend, config, docs)
- 💬 Генерирует осмысленное сообщение
- 👀 Показывает превью изменений с цветами
- ✏️ Позволяет редактировать сообщение перед коммитом
- 🚀 Предлагает запушить изменения

### Quick Commit
- ⚡ Быстрый workflow для мелких изменений
- 📋 Показывает краткий список изменений
- 🤖 Генерирует базовое сообщение на основе типов файлов
- 🎯 Минимум вопросов - максимум скорости

## 🎨 Цветовая схема

- 🟢 **Зеленый** - добавленные файлы
- 🟡 **Желтый** - измененные файлы  
- 🔴 **Красный** - удаленные файлы
- 🔵 **Синий** - процесс выполнения
- 🟦 **Голубой** - информационные сообщения

## 🔧 Требования

- **Windows**: PowerShell 5.1+
- **Linux/macOS**: Bash 4.0+
- **Универсальный**: Node.js 12+
- **Git**: установлен и настроен
