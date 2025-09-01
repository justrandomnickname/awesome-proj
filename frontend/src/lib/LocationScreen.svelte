<script lang="ts">
  import { onMount } from 'svelte';
  import { GetCurrentLocation } from '../../wailsjs/go/app/App.js';
  import SaveMenu from './SaveMenu.svelte';
  
  interface NPCInfo {
    id: string;
    name: string;
    race: string;
    description: string;
  }

  interface LocationInfo {
    id: string;
    name: string;
    description: string;
    npcs: NPCInfo[];
  }
  
  let locationInfo: LocationInfo | null = null;
  let loading = true;
  let error = '';
  let saveMenu: SaveMenu;

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
      <!-- Левая часть: информация о локации -->
      <div class="location-panel">
        <div class="location-content">
          <div class="location-header">
            <h3 class="location-name">{locationInfo.name}</h3>
            <span class="location-id">ID: {locationInfo.id}</span>
          </div>
          
          <div class="location-description">
            <h4>Описание:</h4>
            <p>{locationInfo.description}</p>
          </div>
          
          <!-- Здесь будем добавлять информацию о location state -->
          <div class="location-state">
            <h4>Состояние локации:</h4>
            <p class="placeholder">Информация о состоянии будет добавлена позже...</p>
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
            {#if locationInfo.npcs.length > 0}
              НПЦ в локации ({locationInfo.npcs.length})
            {:else}
              Локация пуста
            {/if}
          </h4>
          
          {#if locationInfo.npcs.length > 0}
            <div class="npcs-list">
              {#each locationInfo.npcs as npc}
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

  /* Левая панель - информация о локации (75% ширины) */
  .location-panel {
    flex: 75;
    background: #2c3e50;
    border-radius: 8px;
    padding: 25px;
    display: flex;
    flex-direction: column;
    min-height: 0; /* Позволяем сжиматься */
    overflow: hidden; /* Убираем скролл с самой панели */
  }

  .location-content {
    flex: 1;
    overflow-y: auto; /* Скролл только для контента */
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

  .location-description h4, .location-state h4 {
    color: #3498db;
    margin-bottom: 15px;
    font-size: 1.3em;
  }

  .location-description p, .location-state p {
    line-height: 1.8;
    color: #ecf0f1;
    font-size: 1.1em;
  }

  .placeholder {
    color: #95a5a6;
    font-style: italic;
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
    
    .location-panel, .sidebar {
      flex: none;
    }
  }
</style>
