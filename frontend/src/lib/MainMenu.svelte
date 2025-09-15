<script lang="ts">
	import { NewGame, ExitGame } from "../../wailsjs/go/app/App"
	import { Quit } from "../../wailsjs/runtime/runtime"
	import { createEventDispatcher } from "svelte"
	import SaveMenu from "./SaveMenu.svelte"

	const dispatch = createEventDispatcher()

	let saveMenuVisible = false
	let loading = false
	let error = ""
	let gameSeed = ""

	// Функция для генерации случайного seed
	function generateRandomSeed(): string {
		return Date.now().toString()
	}

	// Инициализируем seed при загрузке компонента
	gameSeed = generateRandomSeed()

	function showNewGameMenu() {
		gameSeed = generateRandomSeed()
		// Можно расширить в будущем, пока просто запускаем новую игру
		startNewGame()
	}

	async function startNewGame() {
		try {
			loading = true
			// Парсим seed как число
			const seedNumber = parseInt(gameSeed)
			if (isNaN(seedNumber)) {
				error = "Некорректный seed. Введите числовое значение."
				return
			}

			await NewGame(seedNumber)
			// Оповещаем родительский компонент о начале игры
			dispatch("gameStarted")
		} catch (err) {
			error = `Ошибка создания новой игры: ${err}`
		} finally {
			loading = false
		}
	}

	function showLoadMenu() {
		saveMenuVisible = true
	}

	async function exitGame() {
		try {
			// Пробуем сначала через Wails runtime
			await Quit()
		} catch (err) {
			try {
				// Затем через backend метод
				await ExitGame()
			} catch (err2) {
				// И в конце fallback
				window.close()
			}
		}
	}

	// Обработчик загрузки игры из SaveMenu
	function onGameLoaded() {
		saveMenuVisible = false
		dispatch("gameStarted")
	}
</script>

<div class="main-menu">
	<div class="menu-container">
		<div class="game-title">
			<h1>⚔️ RPG Приключение</h1>
			<p class="version">Альфа 0.0.0.1</p>
		</div>

		{#if error}
			<div class="error">{error}</div>
		{/if}

		<div class="menu-buttons">
			<button class="menu-btn new-game-btn" on:click={showNewGameMenu} disabled={loading}>
				🎮 Начать новую игру
			</button>

			<button class="menu-btn load-game-btn" on:click={showLoadMenu} disabled={loading}>
				📂 Загрузить игру
			</button>

			<button class="menu-btn exit-btn" on:click={exitGame} disabled={loading}>
				🚪 Выйти
			</button>
		</div>

		<!-- Дополнительные настройки для новой игры -->
		<div class="game-settings">
			<h3>Настройки новой игры</h3>
			<div class="seed-container">
				<label for="seed-input">Seed игры:</label>
				<div class="seed-input-group">
					<input
						id="seed-input"
						type="text"
						bind:value={gameSeed}
						placeholder="Введите seed или оставьте случайный"
						disabled={loading} />
					<button
						type="button"
						class="regenerate-btn"
						on:click={() => (gameSeed = generateRandomSeed())}
						disabled={loading}
						title="Сгенерировать новый случайный seed">
						🎲
					</button>
				</div>
			</div>
		</div>

		{#if loading}
			<div class="loading">Загрузка...</div>
		{/if}
	</div>
</div>

<!-- SaveMenu для загрузки игр -->
<SaveMenu
	mode="load-only"
	visible={saveMenuVisible}
	on:gameLoaded={onGameLoaded}
	on:close={() => (saveMenuVisible = false)} />

<style>
	.main-menu {
		display: flex;
		justify-content: center;
		align-items: center;
		min-height: 100vh;
		background: linear-gradient(135deg, #1a1a1a 0%, #2c3e50 100%);
		padding: 20px;
		background-image:
			radial-gradient(circle at 25% 25%, rgba(52, 152, 219, 0.1) 0%, transparent 50%),
			radial-gradient(circle at 75% 75%, rgba(231, 76, 60, 0.1) 0%, transparent 50%);
	}

	.menu-container {
		background: rgba(255, 255, 255, 0.1);
		backdrop-filter: blur(10px);
		border-radius: 20px;
		padding: 40px;
		width: 100%;
		max-width: 500px;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
		border: 1px solid rgba(255, 255, 255, 0.2);
	}

	.game-title {
		text-align: center;
		margin-bottom: 40px;
	}

	.game-title h1 {
		color: #ecf0f1;
		margin: 0 0 10px 0;
		font-size: 2.5em;
		text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
	}

	.version {
		color: #95a5a6;
		margin: 0;
		font-size: 1.1em;
		font-style: italic;
	}

	.menu-buttons {
		display: flex;
		flex-direction: column;
		gap: 20px;
		margin-bottom: 40px;
	}

	.menu-btn {
		background: linear-gradient(145deg, #3498db, #2980b9);
		color: white;
		border: none;
		padding: 18px 25px;
		border-radius: 12px;
		cursor: pointer;
		font-size: 1.2em;
		font-weight: bold;
		transition: all 0.3s ease;
		box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.menu-btn:hover:not(:disabled) {
		transform: translateY(-2px);
		box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
	}

	.menu-btn:active:not(:disabled) {
		transform: translateY(0);
	}

	.new-game-btn {
		background: linear-gradient(145deg, #27ae60, #229954);
	}

	.new-game-btn:hover:not(:disabled) {
		background: linear-gradient(145deg, #2ecc71, #27ae60);
	}

	.load-game-btn {
		background: linear-gradient(145deg, #f39c12, #e67e22);
	}

	.load-game-btn:hover:not(:disabled) {
		background: linear-gradient(145deg, #f1c40f, #f39c12);
	}

	.exit-btn {
		background: linear-gradient(145deg, #e74c3c, #c0392b);
	}

	.exit-btn:hover:not(:disabled) {
		background: linear-gradient(145deg, #e67e22, #e74c3c);
	}

	.menu-btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
		transform: none;
	}

	.game-settings {
		background: rgba(255, 255, 255, 0.05);
		padding: 20px;
		border-radius: 12px;
		border: 1px solid rgba(255, 255, 255, 0.1);
	}

	.game-settings h3 {
		color: #bdc3c7;
		margin: 0 0 15px 0;
		font-size: 1.1em;
	}

	.seed-container {
		display: flex;
		flex-direction: column;
		gap: 10px;
	}

	.seed-container label {
		color: #ecf0f1;
		font-weight: 500;
	}

	.seed-input-group {
		display: flex;
		gap: 10px;
	}

	.seed-input-group input {
		flex: 1;
		padding: 12px;
		border: 1px solid rgba(255, 255, 255, 0.2);
		border-radius: 8px;
		background: rgba(255, 255, 255, 0.1);
		color: #ecf0f1;
		font-size: 1em;
	}

	.seed-input-group input::placeholder {
		color: #95a5a6;
	}

	.seed-input-group input:focus {
		outline: none;
		border-color: #3498db;
		background: rgba(255, 255, 255, 0.15);
	}

	.regenerate-btn {
		background: linear-gradient(145deg, #9b59b6, #8e44ad);
		color: white;
		border: none;
		padding: 12px 15px;
		border-radius: 8px;
		cursor: pointer;
		font-size: 1.1em;
		transition: all 0.3s ease;
	}

	.regenerate-btn:hover:not(:disabled) {
		background: linear-gradient(145deg, #a569bd, #9b59b6);
		transform: scale(1.05);
	}

	.regenerate-btn:disabled {
		opacity: 0.6;
		cursor: not-allowed;
	}

	.error {
		background: rgba(231, 76, 60, 0.9);
		color: white;
		padding: 15px;
		border-radius: 8px;
		margin-bottom: 20px;
		text-align: center;
		font-weight: 500;
	}

	.loading {
		text-align: center;
		color: #3498db;
		font-size: 1.1em;
		margin-top: 20px;
		padding: 15px;
		background: rgba(52, 152, 219, 0.1);
		border-radius: 8px;
	}

	/* Адаптивность */
	@media (max-width: 600px) {
		.menu-container {
			padding: 30px 20px;
			margin: 10px;
		}

		.game-title h1 {
			font-size: 2em;
		}

		.menu-btn {
			padding: 15px 20px;
			font-size: 1.1em;
		}
	}
</style>
