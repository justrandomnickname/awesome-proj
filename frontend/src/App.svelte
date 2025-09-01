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
    height: 100vh; /* Фиксированная высота */
    background-color: #1a1a1a;
    padding: 0;
    display: flex;
    flex-direction: column;
    overflow: hidden; /* Полностью убираем скролл на уровне main */
    width: 100%;
    max-width: 100vw;
  }

  .game-interface {
    flex: 1;
    display: flex;
    flex-direction: column;
    width: 100%;
    max-width: 100vw;
    overflow: hidden; /* Убираем скролл на уровне интерфейса */
    min-height: 0; /* Позволяем flex элементам сжиматься */
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
