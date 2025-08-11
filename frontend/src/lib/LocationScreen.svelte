<script lang="ts">
  import { onMount } from 'svelte';
  import { GetCurrentLocation } from '../../wailsjs/go/app/App.js';
  
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

  function refreshLocation() {
    loadCurrentLocation();
  }

  onMount(() => {
    loadCurrentLocation();
  });
</script>

<div class="location-screen">
  <h2>Текущая локация</h2>
  
  {#if loading}
    <div class="loading">
      <p>Загрузка локации...</p>
    </div>
  {:else if error}
    <div class="error">
      <p>{error}</p>
      <button on:click={refreshLocation} class="retry-btn">
        Попробовать снова
      </button>
    </div>
  {:else if locationInfo}
    <div class="location-info">
      <div class="location-header">
        <h3 class="location-name">{locationInfo.name}</h3>
        <span class="location-id">ID: {locationInfo.id}</span>
      </div>
      
      <div class="location-description">
        <p>{locationInfo.description}</p>
      </div>
      
      <div class="location-npcs">
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
      
      <div class="location-actions">
        <button on:click={refreshLocation} class="refresh-btn">
          🔄 Обновить
        </button>
      </div>
    </div>
  {:else}
    <div class="no-location">
      <p>Информация о локации недоступна</p>
    </div>
  {/if}
</div>

<style>
  .location-screen {
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
    font-family: 'Arial', sans-serif;
  }

  h2 {
    text-align: center;
    color: #2c3e50;
    margin-bottom: 20px;
  }

  .loading, .error, .no-location {
    text-align: center;
    padding: 20px;
    border-radius: 8px;
    margin: 20px 0;
  }

  .loading {
    background: #e3f2fd;
    color: #1976d2;
  }

  .error {
    background: #ffebee;
    color: #c62828;
  }

  .no-location {
    background: #f5f5f5;
    color: #666;
  }

  .retry-btn, .refresh-btn {
    background: #3498db;
    color: white;
    border: none;
    padding: 8px 16px;
    border-radius: 4px;
    cursor: pointer;
    margin: 10px 5px 0;
    font-size: 0.9em;
  }

  .retry-btn:hover, .refresh-btn:hover {
    background: #2980b9;
  }

  .location-info {
    border: 1px solid #ddd;
    border-radius: 10px;
    padding: 20px;
    background: #f9f9f9;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  }

  .location-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
    border-bottom: 2px solid #e74c3c;
    padding-bottom: 10px;
  }

  .location-name {
    color: #2c3e50;
    margin: 0;
    font-size: 1.5em;
  }

  .location-id {
    color: #7f8c8d;
    font-size: 0.9em;
    background: #ecf0f1;
    padding: 4px 8px;
    border-radius: 4px;
  }

  .location-description {
    margin: 15px 0;
    padding: 15px;
    background: #e8f5e8;
    border-radius: 8px;
  }

  .location-description p {
    margin: 0;
    color: #2c3e50;
    line-height: 1.6;
    font-style: italic;
  }

  .location-npcs {
    margin: 20px 0;
  }

  .npcs-title {
    color: #2c3e50;
    margin-bottom: 15px;
    font-size: 1.2em;
    border-bottom: 2px solid #e74c3c;
    padding-bottom: 5px;
  }

  .npcs-list {
    display: flex;
    flex-direction: column;
    gap: 15px;
  }

  .npc-card {
    border: 1px solid #bdc3c7;
    border-radius: 8px;
    padding: 15px;
    background: white;
    box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
    transition: transform 0.2s ease, box-shadow 0.2s ease;
  }

  .npc-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.15);
  }

  .npc-header {
    display: flex;
    align-items: center;
    gap: 10px;
    margin-bottom: 8px;
  }

  .npc-name {
    font-weight: bold;
    color: #2c3e50;
    flex-grow: 1;
  }

  .npc-race {
    font-size: 0.9em;
    color: #7f8c8d;
    margin-bottom: 6px;
    font-weight: 500;
  }

  .npc-description {
    color: #34495e;
    font-style: italic;
    line-height: 1.4;
  }

  .no-npcs {
    text-align: center;
    padding: 20px;
    color: #7f8c8d;
    background: #f8f9fa;
    border-radius: 8px;
    border: 1px dashed #bdc3c7;
  }

  .location-actions {
    margin-top: 20px;
    text-align: center;
  }
</style>
