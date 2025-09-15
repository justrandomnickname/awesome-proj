<script lang="ts">
	import { ExitGame } from "../../wailsjs/go/app/App"
	import { Quit } from "../../wailsjs/runtime/runtime"
	import { createEventDispatcher } from "svelte"
	import SaveMenu from "./SaveMenu.svelte"

	const dispatch = createEventDispatcher()

	let isOpen = false
	let saveMenuVisible = false
	let loading = false

	export function openMenu() {
		isOpen = true
	}

	function closeMenu() {
		isOpen = false
	}

	function showSaveMenu() {
		saveMenuVisible = true
	}

	function returnToMainMenu() {
		closeMenu()
		dispatch("returnToMainMenu")
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

	// Обработчик загрузки игры
	function onGameLoaded() {
		saveMenuVisible = false
		closeMenu()
		dispatch("gameLoaded")
	}

	function onSaveMenuClosed() {
		saveMenuVisible = false
	}
</script>

{#if isOpen}
	<div class="overlay" on:click={closeMenu}>
		<div class="menu" on:click|stopPropagation>
			<div class="menu-header">
				<h2>🎮 Игровое меню</h2>
				<button class="close-btn" on:click={closeMenu}>×</button>
			</div>

			<div class="menu-content">
				<div class="menu-buttons">
					<button class="menu-btn resume-btn" on:click={closeMenu} disabled={loading}>
						▶️ Продолжить игру
					</button>

					<button class="menu-btn save-btn" on:click={showSaveMenu} disabled={loading}>
						💾 Сохранения
					</button>

					<button
						class="menu-btn main-menu-btn"
						on:click={returnToMainMenu}
						disabled={loading}>
						🏠 Главное меню
					</button>

					<button class="menu-btn exit-btn" on:click={exitGame} disabled={loading}>
						🚪 Выйти из игры
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}

<!-- SaveMenu для сохранений и загрузки -->
<SaveMenu
	mode="full"
	visible={saveMenuVisible}
	on:gameLoaded={onGameLoaded}
	on:close={onSaveMenuClosed} />

<style>
	.overlay {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background: rgba(0, 0, 0, 0.8);
		display: flex;
		justify-content: center;
		align-items: center;
		z-index: 1000;
		backdrop-filter: blur(5px);
	}

	.menu {
		background: linear-gradient(145deg, #2c3e50, #34495e);
		border-radius: 15px;
		width: 90%;
		max-width: 400px;
		box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
		color: #ecf0f1;
		border: 1px solid rgba(255, 255, 255, 0.1);
		overflow: hidden;
	}

	.menu-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 20px;
		background: rgba(52, 152, 219, 0.1);
		border-bottom: 1px solid rgba(255, 255, 255, 0.1);
	}

	.menu-header h2 {
		margin: 0;
		color: #ecf0f1;
		font-size: 1.3em;
	}

	.close-btn {
		background: none;
		border: none;
		font-size: 24px;
		color: #ecf0f1;
		cursor: pointer;
		padding: 5px;
		line-height: 1;
		border-radius: 50%;
		width: 35px;
		height: 35px;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.3s ease;
	}

	.close-btn:hover {
		background: rgba(231, 76, 60, 0.2);
		color: #e74c3c;
	}

	.menu-content {
		padding: 30px;
	}

	.menu-buttons {
		display: flex;
		flex-direction: column;
		gap: 15px;
	}

	.menu-btn {
		background: linear-gradient(145deg, #3498db, #2980b9);
		color: white;
		border: none;
		padding: 15px 20px;
		border-radius: 10px;
		cursor: pointer;
		font-size: 1.1em;
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

	.resume-btn {
		background: linear-gradient(145deg, #27ae60, #229954);
	}

	.resume-btn:hover:not(:disabled) {
		background: linear-gradient(145deg, #2ecc71, #27ae60);
	}

	.save-btn {
		background: linear-gradient(145deg, #f39c12, #e67e22);
	}

	.save-btn:hover:not(:disabled) {
		background: linear-gradient(145deg, #f1c40f, #f39c12);
	}

	.main-menu-btn {
		background: linear-gradient(145deg, #9b59b6, #8e44ad);
	}

	.main-menu-btn:hover:not(:disabled) {
		background: linear-gradient(145deg, #a569bd, #9b59b6);
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

	/* Адаптивность */
	@media (max-width: 600px) {
		.menu {
			width: 95%;
			margin: 10px;
		}

		.menu-content {
			padding: 20px;
		}

		.menu-btn {
			padding: 12px 15px;
			font-size: 1em;
		}
	}
</style>
