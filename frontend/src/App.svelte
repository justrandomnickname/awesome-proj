<script lang="ts">
  import LocationScreen from './lib/LocationScreen.svelte'

  const playerName = "developer" // Имя для разработки
</script>

<main>
  <!-- Игровой экран сразу -->
  <div class="game-interface">
    <div class="game-header">
      <h1>⚔️ RPG Приключение</h1>
      <div class="header-controls">
        <span class="player-name">Игрок: {playerName}</span>
      </div>
    </div>
    
    <LocationScreen />
  </div>
</main>

<style>
  main {
    min-height: 100vh;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    padding: 20px;
    display: flex;
    flex-direction: column;
  }

  .game-interface {
    flex: 1;
    display: flex;
    flex-direction: column;
  }

  .game-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: rgba(255, 255, 255, 0.1);
    padding: 15px 25px;
    border-radius: 10px;
    margin-bottom: 20px;
    backdrop-filter: blur(10px);
  }

  .game-header h1 {
    color: #fff;
    margin: 0;
    font-size: 1.8em;
    text-shadow: 1px 1px 2px rgba(0,0,0,0.5);
  }

  .header-controls {
    display: flex;
    align-items: center;
    gap: 20px;
  }

  .player-name {
    color: #fff;
    font-weight: bold;
    background: rgba(255, 255, 255, 0.2);
    padding: 8px 15px;
    border-radius: 20px;
  }
</style>
