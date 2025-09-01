<script lang="ts">
  import { onMount } from 'svelte';
  import { GetCurrentLocation, PerformPlayerAction } from '../../wailsjs/go/app/App.js';
  import SaveMenu from './SaveMenu.svelte';
  
  interface NPCInfo {
    id: string;
    name: string;
    race: string;
    location_id: string;
    description: string;
  }

  interface InteractionInfo {
    id: string;
    type: string;
    content: string;
    timestamp: string;
  }

  interface LocationInfo {
    id: string;
    name: string;
    description: string;
    current_state: string;
    type: string;
    exits: Record<string, string>;
    npcs: string[];
    npcs_detailed?: NPCInfo[];
    interactions?: InteractionInfo[];
  }
  
  let locationInfo: LocationInfo | null = null;
  let loading = true;
  let error = '';
  let saveMenu: SaveMenu;
  
  // Action input
  let actionText = '';
  let performingAction = false;

  async function loadCurrentLocation() {
    try {
      loading = true;
      error = '';
      
      const location = await GetCurrentLocation();
      locationInfo = location;
      
    } catch (err) {
      error = `Ошибка загрузки локации: ${err}`;
      console.error('Failed to load location:', err);
    } finally {
      loading = false;
    }
  }

  async function performAction() {
    if (!actionText.trim() || performingAction) {
      return;
    }

    try {
      performingAction = true;
      error = '';
      
      await PerformPlayerAction(actionText.trim());
      actionText = ''; // Очищаем поле ввода
      
      // Перезагружаем локацию для получения новых взаимодействий
      await loadCurrentLocation();
      
    } catch (err) {
      error = `Ошибка выполнения действия: ${err}`;
      console.error('Failed to perform action:', err);
    } finally {
      performingAction = false;
    }
  }

  function handleKeyPress(event: KeyboardEvent) {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      performAction();
    }
  }

  onMount(() => {
    loadCurrentLocation();
  });
</script>

<div class="location-screen">
  <div class="header">
    <h2>Текущая локация</h2>
  </div>
  
  {#if loading}
    <div class="loading">
      <p>Загрузка локации...</p>
    </div>
  {:else if error}
    <div class="error">
      <p>{error}</p>
      <button on:click={loadCurrentLocation} class="retry-btn">
        Попробовать снова
      </button>
    </div>
  {:else if locationInfo}
    <div class="main-layout">
      <!-- Левая часть: взаимодействия и поле ввода -->
      <div class="interactions-panel">
        <div class="interactions-content">
          <div class="location-header">
            <h3 class="location-name">{locationInfo.name}</h3>
            <span class="location-id">ID: {locationInfo.id}</span>
          </div>
          
          <div class="location-description">
            <p>{locationInfo.description}</p>
          </div>
          
          <!-- История взаимодействий -->
          <div class="interactions-history">
            <h4>История взаимодействий:</h4>
            <div class="interactions-scroll">
              {#if locationInfo.interactions && locationInfo.interactions.length > 0}
                {#each locationInfo.interactions as interaction}
                  <div class="interaction-item interaction-{interaction.type}">
                    <div class="interaction-header">
                      <span class="interaction-type">
                        {#if interaction.type === 'player_action'}
                          🎮 Ваше действие
                        {:else if interaction.type === 'location_response'}
                          🌍 Ответ локации
                        {:else if interaction.type === 'location_state'}
                          📍 Состояние локации
                        {/if}
                      </span>
                      <span class="interaction-time">{interaction.timestamp}</span>
                    </div>
                    <div class="interaction-content">{interaction.content}</div>
                  </div>
                {/each}
              {:else}
                <div class="no-interactions">
                  <p>Пока никаких взаимодействий не было...</p>
                </div>
              {/if}
            </div>
          </div>
          
          <!-- Поле ввода действий -->
          <div class="action-input-section">
            <h4>Введите ваше действие:</h4>
            <div class="action-input-container">
              <textarea 
                bind:value={actionText}
                on:keypress={handleKeyPress}
                placeholder="Опишите, что вы хотите сделать... (Enter для отправки, Shift+Enter для новой строки)"
                class="action-input"
                rows="3"
                disabled={performingAction}
              ></textarea>
              <button 
                on:click={performAction}
                class="action-submit-btn"
                disabled={!actionText.trim() || performingAction}
              >
                {#if performingAction}
                  Выполняется...
                {:else}
                  Выполнить действие
                {/if}
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Правая часть: NPC и кнопки -->
      <div class="sidebar">
        <div class="sidebar-section">
          <button class="save-menu-btn" on:click={() => saveMenu.openMenu()}>
            💾 Сохранения
          </button>
        </div>
        
        <div class="sidebar-section">
          <h4 class="npcs-title">
            {#if locationInfo.npcs_detailed && locationInfo.npcs_detailed.length > 0}
              НПЦ в локации ({locationInfo.npcs_detailed.length})
            {:else}
              Локация пуста
            {/if}
          </h4>
          
          {#if locationInfo.npcs_detailed && locationInfo.npcs_detailed.length > 0}
            <div class="npcs-list">
              {#each locationInfo.npcs_detailed as npc}
                <div class="npc-card">
                  <div class="npc-header">
                    <span class="npc-name">{npc.name}</span>
                  </div>
                  <div class="npc-race">Раса: {npc.race}</div>
                  <div class="npc-description">{npc.description}</div>
                </div>
              {/each}
            </div>
            {:else}
              <div class="no-npcs">
                <p>В этой локации никого нет...</p>
              </div>
            {/if}
        </div>
      </div>
    </div>
  {:else}
    <div class="no-location">
      <p>Информация о локации недоступна</p>
    </div>
  {/if}
</div>

<SaveMenu bind:this={saveMenu} />

<style>
  .location-screen {
    max-width: none;
    width: 100%;
    height: 100vh;
    margin: 0;
    padding: 0;
    font-family: 'Arial', sans-serif;
    background: #1a1a1a;
    color: #ecf0f1;
    display: flex;
    flex-direction: column;
    overflow: hidden; /* Убираем любые глобальные скроллы */
  }

  .header {
    text-align: center;
    margin-bottom: 20px;
    flex-shrink: 0;
    padding: 20px 20px 0 20px; /* Добавляем padding только сверху и по бокам */
  }

  h2 {
    color: #ecf0f1;
    margin: 0;
  }

  .main-layout {
    display: flex;
    gap: 20px;
    align-items: stretch; /* Растягиваем элементы по высоте */
    flex: 1; /* Занимаем всё оставшееся место */
    padding: 0 20px 20px 20px;
    min-height: 0; /* Позволяем flex элементам сжиматься */
    overflow: hidden; /* Убираем скролл на уровне layout */
  }

  /* Левая панель - взаимодействия (75% ширины) */
  .interactions-panel {
    flex: 75;
    background: #2c3e50;
    border-radius: 8px;
    padding: 25px;
    display: flex;
    flex-direction: column;
    min-height: 0; /* Позволяем сжиматься */
    overflow: hidden; /* Убираем скролл с самой панели */
  }

  .interactions-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-height: 0;
  }

  .location-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
    border-bottom: 2px solid #34495e;
    padding-bottom: 10px;
  }

  .location-name {
    color: #ecf0f1;
    margin: 0;
    font-size: 2em;
  }

  .location-id {
    color: #95a5a6;
    font-size: 1em;
  }

  .location-description, .location-state {
    margin: 25px 0;
  }

  .location-description {
    margin-bottom: 20px;
  }

  .location-description p {
    line-height: 1.8;
    color: #ecf0f1;
    font-size: 1.1em;
  }

  /* Правая панель - NPC и кнопки (25% ширины) */
  .sidebar {
    flex: 25;
    display: flex;
    flex-direction: column;
    gap: 20px;
    min-height: 0; /* Позволяем сжиматься */
    overflow: hidden; /* Убираем глобальный скролл */
  }

  .sidebar-section {
    background: #2c3e50;
    border-radius: 8px;
    padding: 20px;
  }

  .sidebar-section:last-child {
    flex: 1;
    overflow-y: auto; /* Скролл только для списка NPC */
    min-height: 0; /* Позволяем сжиматься */
  }

  .save-menu-btn {
    width: 100%;
    background: #3498db;
    color: white;
    border: none;
    padding: 12px 15px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 14px;
    font-weight: bold;
  }

  .save-menu-btn:hover {
    background: #2980b9;
  }

  .npcs-title {
    color: #3498db;
    margin: 0 0 20px 0;
    font-size: 1.3em;
  }

  .npcs-list {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }

  .npc-card {
    background: #34495e;
    border-radius: 8px;
    padding: 18px;
    border-left: 4px solid #e74c3c;
  }

  .npc-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }

  .npc-name {
    font-weight: bold;
    color: #ecf0f1;
    font-size: 1.1em;
  }

  .npc-race {
    color: #95a5a6;
    font-size: 1em;
    margin-bottom: 8px;
  }

  .npc-description {
    color: #ecf0f1;
    font-size: 1em;
    line-height: 1.5;
  }

  .no-npcs {
    text-align: center;
    color: #95a5a6;
    font-style: italic;
  }

  .loading, .error, .no-location {
    text-align: center;
    padding: 20px;
    border-radius: 8px;
    margin: 20px 0;
  }

  .loading {
    background: #2c3e50;
    color: #ecf0f1;
  }

  .error {
    background: #e74c3c;
    color: white;
  }

  .no-location {
    background: #95a5a6;
    color: white;
  }

  .retry-btn {
    background: #f39c12;
    color: white;
    border: none;
    padding: 8px 15px;
    border-radius: 4px;
    cursor: pointer;
    margin-top: 10px;
  }

  .retry-btn:hover {
    background: #e67e22;
  }

  /* Адаптивность для маленьких экранов */
  @media (max-width: 768px) {
    .main-layout {
      flex-direction: column;
    }
    
    .interactions-panel, .sidebar {
      flex: none;
    }
  }

  /* Стили для взаимодействий */
  .interactions-history {
    flex: 1;
    display: flex;
    flex-direction: column;
    margin-bottom: 20px;
    min-height: 0;
  }

  .interactions-history h4 {
    color: #3498db;
    margin-bottom: 15px;
    font-size: 1.3em;
  }

  .interactions-scroll {
    flex: 1;
    overflow-y: auto;
    border: 1px solid #34495e;
    border-radius: 6px;
    padding: 15px;
    background: #34495e;
    min-height: 0;
  }

  .interaction-item {
    margin-bottom: 15px;
    padding: 12px;
    border-radius: 6px;
    border-left: 4px solid;
  }

  .interaction-player_action {
    background: #1e3a8a20;
    border-left-color: #3b82f6;
  }

  .interaction-location_response {
    background: #15803d20;
    border-left-color: #22c55e;
  }

  .interaction-location_state {
    background: #92400e20;
    border-left-color: #f59e0b;
  }

  .interaction-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
    font-size: 0.9em;
  }

  .interaction-type {
    font-weight: bold;
    color: #ecf0f1;
  }

  .interaction-time {
    color: #95a5a6;
    font-size: 0.8em;
  }

  .interaction-content {
    color: #ecf0f1;
    line-height: 1.6;
  }

  .no-interactions {
    text-align: center;
    color: #95a5a6;
    font-style: italic;
    padding: 20px;
  }

  /* Стили для поля ввода действий */
  .action-input-section {
    flex-shrink: 0;
    margin-top: auto;
  }

  .action-input-section h4 {
    color: #3498db;
    margin-bottom: 10px;
    font-size: 1.2em;
  }

  .action-input-container {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .action-input {
    width: 100%;
    padding: 12px;
    border: 1px solid #34495e;
    border-radius: 6px;
    background: #34495e;
    color: #ecf0f1;
    font-family: inherit;
    font-size: 1em;
    resize: vertical;
    min-height: 60px;
  }

  .action-input:focus {
    outline: none;
    border-color: #3498db;
    box-shadow: 0 0 5px rgba(52, 152, 219, 0.3);
  }

  .action-input:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .action-submit-btn {
    background: #27ae60;
    color: white;
    border: none;
    padding: 12px 20px;
    border-radius: 6px;
    cursor: pointer;
    font-size: 1em;
    font-weight: bold;
    transition: background-color 0.3s;
  }

  .action-submit-btn:hover:not(:disabled) {
    background: #2ecc71;
  }

  .action-submit-btn:disabled {
    background: #95a5a6;
    cursor: not-allowed;
  }
</style>

