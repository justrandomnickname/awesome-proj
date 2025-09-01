<script lang="ts">
  import { GetSavesList, SaveGame, LoadGame, DeleteSave, NewGame } from '../../wailsjs/go/app/App';
  import { app } from '../../wailsjs/go/models';

  type SaveInfo = app.SaveInfo;

  let isOpen = false;
  let saves: SaveInfo[] = [];
  let newSaveName = '';
  let loading = false;
  let error = '';

  export function openMenu() {
    isOpen = true;
    loadSaves();
  }

  function closeMenu() {
    console.log('closeMenu called, isOpen before:', isOpen);
    isOpen = false;
    error = '';
    console.log('closeMenu done, isOpen after:', isOpen);
  }

  async function loadSaves() {
    try {
      loading = true;
      const result = await GetSavesList();
      saves = result as SaveInfo[];
      // Сортируем по дате создания - самые новые вверху
      saves.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
      error = '';
    } catch (err) {
      console.error('Error loading saves:', err);
      error = `Ошибка загрузки сохранений: ${err}`;
      saves = [];
    } finally {
      loading = false;
    }
  }

  async function saveGame() {
    if (!newSaveName.trim()) {
      error = 'Введите имя сохранения';
      return;
    }

    try {
      loading = true;
      await SaveGame(newSaveName.trim());
      newSaveName = '';
      await loadSaves();
      error = '';
    } catch (err) {
      error = `Ошибка сохранения: ${err}`;
    } finally {
      loading = false;
    }
  }

  async function loadGame(filename: string) {
    try {
      loading = true;
      await LoadGame(filename);
      closeMenu();
      // Перезагрузка страницы для обновления состояния
      window.location.reload();
    } catch (err) {
      error = `Ошибка загрузки: ${err}`;
    } finally {
      loading = false;
    }
  }

  async function deleteSave(filename: string, saveName: string) {
    try {
      loading = true;
      await DeleteSave(filename);
      await loadSaves();
    } catch (err) {
      error = `Ошибка удаления: ${err}`;
    } finally {
      loading = false;
    }
  }

  async function newGame() {
    try {
      loading = true;
      await NewGame();
      closeMenu();
      // Перезагружаем страницу чтобы обновить интерфейс с новой игрой
      window.location.reload();
    } catch (err) {
      error = `Ошибка создания новой игры: ${err}`;
    } finally {
      loading = false;
    }
  }

  function formatDate(dateValue: any): string {
    try {
      if (!dateValue) return 'Неизвестно';
      const date = new Date(dateValue);
      return date.toLocaleString();
    } catch {
      return 'Неизвестно';
    }
  }
</script>

{#if isOpen}
  <div class="overlay" on:click={closeMenu}>
    <div class="menu" on:click|stopPropagation>
      <div class="menu-header">
        <h2>Сохранения</h2>
        <button class="close-btn" on:click={closeMenu}>×</button>
      </div>

      {#if error}
        <div class="error">{error}</div>
      {/if}

      <div class="menu-content">
        <!-- Новая игра -->
        <div class="section">
          <button class="new-game-btn" on:click={newGame} disabled={loading}>
            🎮 Новая игра
          </button>
        </div>

        <!-- Сохранить текущую игру -->
        <div class="section">
          <h3>Сохранить текущую игру</h3>
          <div class="save-form">
            <input
              type="text"
              bind:value={newSaveName}
              placeholder="Имя сохранения"
              disabled={loading}
            />
            <button on:click={saveGame} disabled={loading || newSaveName.trim().length === 0}>
              💾 Сохранить
            </button>
          </div>
        </div>

        <!-- Список сохранений -->
        <div class="section">
          <h3>Загрузить сохранение</h3>
          
          {#if loading}
            <div class="loading">Загрузка...</div>
          {:else if saves.length === 0}
            <div class="no-saves">Нет сохранений</div>
          {:else}
            <div class="saves-list">
              {#each saves as save}
                <div class="save-item">
                  <div class="save-info">
                    <div class="save-name">{save.name}</div>
                    <div class="save-date">
                      {formatDate(save.created_at)}
                    </div>
                  </div>
                  <div class="save-actions">
                    <button class="load-btn" on:click={() => loadGame(save.filename)} disabled={loading}>
                      📂 Загрузить
                    </button>
                    <button class="delete-btn" on:click={() => deleteSave(save.filename, save.name)} disabled={loading}>
                      🗑️ Удалить
                    </button>
                  </div>
                </div>
              {/each}
            </div>
          {/if}
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  .overlay {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 1000;
  }

  .menu {
    background: #2c3e50;
    border-radius: 12px;
    width: 90%;
    max-width: 600px;
    max-height: 80vh;
    overflow-y: auto;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
    color: #ecf0f1;
  }

  .menu-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    border-bottom: 1px solid #34495e;
  }

  .menu-header h2 {
    margin: 0;
    color: #ecf0f1;
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 24px;
    color: #ecf0f1;
    cursor: pointer;
    padding: 5px;
    line-height: 1;
  }

  .close-btn:hover {
    color: #e74c3c;
  }

  .menu-content {
    padding: 20px;
  }

  .section {
    margin-bottom: 30px;
  }

  .section h3 {
    margin: 0 0 15px 0;
    color: #bdc3c7;
    font-size: 16px;
  }

  .error {
    background: #e74c3c;
    color: white;
    padding: 10px;
    border-radius: 6px;
    margin: 0 20px 20px 20px;
  }

  .new-game-btn {
    background: #27ae60;
    color: white;
    border: none;
    padding: 12px 20px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 16px;
    width: 100%;
  }

  .new-game-btn:hover:not(:disabled) {
    background: #2ecc71;
  }

  .save-form {
    display: flex;
    gap: 10px;
  }

  .save-form input {
    flex: 1;
    padding: 10px;
    border: 1px solid #34495e;
    border-radius: 6px;
    background: #34495e;
    color: #ecf0f1;
  }

  .save-form input::placeholder {
    color: #95a5a6;
  }

  .save-form button {
    background: #3498db;
    color: white;
    border: none;
    padding: 10px 15px;
    border-radius: 6px;
    cursor: pointer;
  }

  .save-form button:hover:not(:disabled) {
    background: #2980b9;
  }

  .loading, .no-saves {
    text-align: center;
    color: #95a5a6;
    padding: 20px;
  }

  .saves-list {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .save-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 15px;
    background: #34495e;
    border-radius: 8px;
  }

  .save-info {
    flex: 1;
  }

  .save-name {
    font-weight: bold;
    color: #ecf0f1;
    margin-bottom: 5px;
  }

  .save-date {
    font-size: 12px;
    color: #95a5a6;
  }

  .save-actions {
    display: flex;
    gap: 10px;
  }

  .load-btn {
    background: #27ae60;
    color: white;
    border: none;
    padding: 8px 12px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 12px;
  }

  .load-btn:hover:not(:disabled) {
    background: #2ecc71;
  }

  .delete-btn {
    background: #e74c3c;
    color: white;
    border: none;
    padding: 8px 12px;
    border-radius: 4px;
    cursor: pointer;
    font-size: 12px;
  }

  .delete-btn:hover:not(:disabled) {
    background: #c0392b;
  }

  button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }
</style>
