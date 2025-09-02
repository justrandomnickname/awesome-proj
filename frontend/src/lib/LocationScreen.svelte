<script lang="ts">
	import { onMount } from "svelte"
	import {
		GetCurrentLocation,
		PerformPlayerAction,
		GetCurrentPoint,
		GetAvailableConnections,
		MoveToPoint,
		GetNPCsForCurrentPoint,
	} from "../../wailsjs/go/app/App.js"
	import SaveMenu from "./SaveMenu.svelte"

	interface NPCInfo {
		id: string
		name: string
		race: string
		location_id: string
		description: string
	}

	interface InteractionInfo {
		id: string
		type: string
		content: string
		timestamp: string
		additional_content?: string
	}

	interface LocationInfo {
		id: string
		name: string
		description: string
		current_state: string
		type: string
		exits: Record<string, string>
		npcs: string[]
		npcs_detailed?: NPCInfo[]
		interactions?: InteractionInfo[]
	}

	interface Point {
		id: string
		name: string
		description: string
		sub_cluster_id: string
		type: string
		connections: string[]
		npcs: string[]
		is_entry_point: boolean
	}

	let locationInfo: LocationInfo | null = null
	let interactions: LocationInfo["interactions"] = []
	let currentPoint: Point | null = null
	let availableConnections: Point[] = []
	let currentPointNPCs: NPCInfo[] = []
	let loading = true
	let error = ""
	let saveMenu: SaveMenu
	let interactionsScrollElement: HTMLElement

	// Action input
	let actionText = ""
	let performingAction = false

	// Сохранение позиции скролла
	let savedScrollPosition = 0
	let isRestoring = false

	function saveScrollPosition() {
		if (interactionsScrollElement && !isRestoring) {
			savedScrollPosition = interactionsScrollElement.scrollTop
			console.log("Saved scroll position:", savedScrollPosition)
		}
	}

	function scrollToBottom() {
		if (interactionsScrollElement) {
			// Используем setTimeout чтобы дождаться обновления DOM
			setTimeout(() => {
				interactionsScrollElement.scrollTop = interactionsScrollElement.scrollHeight
				console.log("Scrolled to bottom")
			}, 50)
		}
	}

	async function loadCurrentLocation() {
		try {
			loading = true
			error = ""

			const location = await GetCurrentLocation()

			console.log("location", location)
			locationInfo = location

			interactions = location.interactions || []

			// Загружаем информацию о текущей точке и доступных переходах
			await loadNavigationInfo()
		} catch (err) {
			error = `Ошибка загрузки локации: ${err}`
			console.error("Failed to load location:", err)
		} finally {
			loading = false
		}
	}

	async function loadNavigationInfo() {
		try {
			const point = await GetCurrentPoint()
			currentPoint = point

			const connections = await GetAvailableConnections()
			availableConnections = connections

			// Загружаем NPCs для текущего Point
			const npcs = await GetNPCsForCurrentPoint()
			currentPointNPCs = npcs
		} catch (err) {
			console.error("Failed to load navigation info:", err)
			// Не показываем ошибку пользователю, просто логируем
			currentPoint = null
			availableConnections = []
			currentPointNPCs = []
		}
	}

	async function moveToPoint(pointId: string) {
		try {
			error = ""
			await MoveToPoint(pointId)

			// Перезагружаем информацию о локации и навигации
			await loadCurrentLocation()
		} catch (err) {
			error = `Ошибка перехода: ${err}`
			console.error("Failed to move to point:", err)
		}
	}

	// Функция для получения маршрутов из данной точки
	function getRoutePreview(pointId: string): string[] {
		if (!currentPoint) return []

		const hierarchy = locationInfo // Используем доступную информацию
		const routes: string[] = []

		// Находим точку в доступных соединениях
		const targetPoint = availableConnections.find(p => p.id === pointId)
		if (!targetPoint) return []

		// Добавляем прямые соединения этой точки
		targetPoint.connections.forEach(connectionId => {
			// Ищем название точки среди всех доступных соединений и текущей точки
			if (connectionId === currentPoint.id) {
				routes.push(`${currentPoint.name} (вернуться)`)
			} else {
				const connectedPoint = availableConnections.find(p => p.id === connectionId)
				if (connectedPoint) {
					routes.push(connectedPoint.name)
				}
			}
		})

		return routes.slice(0, 3) // Ограничиваем 3 маршрутами для компактности
	}

	async function performAction() {
		if (!actionText.trim() || performingAction) {
			return
		}

		try {
			performingAction = true
			error = ""

			await PerformPlayerAction(actionText.trim())
			actionText = "" // Очищаем поле ввода

			// Перезагружаем локацию для получения новых взаимодействий
			await loadCurrentLocation()

			// Устанавливаем скролл в самый низ после загрузки
			scrollToBottom()
		} catch (err) {
			error = `Ошибка выполнения действия: ${err}`
			console.error("Failed to perform action:", err)
		} finally {
			performingAction = false
		}
	}

	function handleKeyPress(event: KeyboardEvent) {
		if (event.key === "Enter" && !event.shiftKey) {
			event.preventDefault()
			performAction()
		}
	}

	onMount(() => {
		loadCurrentLocation()
	})
</script>

<div class="location-screen">
	<div class="header">
		<h2>Текущая локация</h2>
	</div>

	{#if locationInfo && currentPoint}
		<div class="main-layout">
			<!-- Левая часть: взаимодействия и поле ввода -->
			<div class="interactions-panel">
				<div class="interactions-content">
					<div class="location-info">
						<div class="location-header">
							<h3 class="location-name">{locationInfo.name}</h3>
							<div>
								<span class="location-id">ID: {currentPoint.id}</span>
								<span class="point-type point-type-{currentPoint.type}"
									>{currentPoint.type}</span>
							</div>
						</div>

						<div class="location-description">
							<p>{currentPoint.description}</p>
						</div>
					</div>

					<!-- История взаимодействий -->
					<div class="interactions-history">
						<h4>История взаимодействий:</h4>
						<div class="interactions-scroll" bind:this={interactionsScrollElement}>
							{#if interactions && interactions.length > 0}
								{#each interactions.slice(-10) as interaction (interaction.id)}
									<div class="interaction-item interaction-{interaction.type}">
										<div class="interaction-header">
											<span class="interaction-type">
												{#if interaction.type === "player_action"}
													🎮 Ваше действие
												{:else if interaction.type === "location_response"}
													<!-- 🌍 Ответ локации -->
												{:else if interaction.type === "location_state"}
													📍 Состояние локации
												{:else if interaction.type === "player_movement"}
													Переход
												{/if}
											</span>
											<span class="interaction-time"
												>{interaction.timestamp}</span>
										</div>
										<div class="interaction-content">
											<div class="interaction-main-content">
												{interaction.content}
											</div>
											{#if interaction.additional_content}
												<div class="interaction-additional-content">
													{interaction.additional_content}
												</div>
											{/if}
										</div>
									</div>
								{/each}
							{:else}
								<div class="no-interactions">
									<p>Пока никаких взаимодействий не было...</p>
								</div>
							{/if}
						</div>
					</div>

					<!-- Навигация по точкам -->
					{#if currentPoint && availableConnections.length > 0}
						<div class="navigation-section">
							<!-- <h4>Переходы:</h4> -->
							<div class="current-point">
								{currentPoint.name}
								<!-- <p class="point-description">{currentPoint.description}</p> -->
							</div>
							<div class="connections-grid">
								{#each availableConnections as connection}
									<div class="connection-container">
										<button
											class="connection-btn connection-{connection.type}"
											on:click={() => moveToPoint(connection.id)}
											title={connection.description}>
											<div class="connection-name">{connection.name}</div>
											<div class="connection-type">{connection.type}</div>
										</button>
										<div class="route-preview">
											{#each getRoutePreview(connection.id) as route, index}
												<div class="route-item">→ {route}</div>
											{/each}
										</div>
									</div>
								{/each}
							</div>
						</div>
					{/if}
				</div>

				<!-- Поле ввода действий - ВЫНЕСЕНО НА УРОВЕНЬ ПАНЕЛИ -->
				<div class="action-input-section">
					<h4>Введите ваше действие:</h4>
					<div class="action-input-container">
						<textarea
							bind:value={actionText}
							on:keypress={handleKeyPress}
							placeholder="Опишите, что вы хотите сделать... (Enter для отправки, Shift+Enter для новой строки)"
							class="action-input"
							rows="3"
							disabled={performingAction}></textarea>
						<button
							on:click={performAction}
							class="action-submit-btn"
							disabled={!actionText.trim() || performingAction}>
							{#if performingAction}
								Выполняется...
							{:else}
								Выполнить действие
							{/if}
						</button>
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
						{#if currentPointNPCs && currentPointNPCs.length > 0}
							NPC ({currentPointNPCs.length})
						{:else}
							Здесь никого нет
						{/if}
					</h4>

					{#if currentPointNPCs && currentPointNPCs.length > 0}
						<div class="npcs-list">
							{#each currentPointNPCs as npc}
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
							<p>В этой точке никого нет...</p>
						</div>
					{/if}
				</div>
			</div>
		</div>
	{/if}
</div>

<SaveMenu bind:this={saveMenu} />

<style>
	.location-screen {
		max-width: none;
		width: 100%;
		height: 85%;
		margin: 0;
		padding: 0;
		font-family: "Arial", sans-serif;
		background: #1a1a1a;
		color: #ecf0f1;
		display: flex;
		flex-direction: column;
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
	}

	/* Левая панель - взаимодействия (85% ширины) */
	.interactions-panel {
		flex: 85;
		background: #2c3e50;
		border-radius: 8px;
		padding: 25px;
		display: flex;
		overflow: auto;
		height: 100%;
		flex-direction: column;
	}

	.interactions-content {
		flex: 1;
		display: flex;
		flex-direction: column;
		margin-bottom: 20px; /* Отступ перед полем ввода */
	}

	.location-info {
		flex-shrink: 0; /* Информация о локации не сжимается */
		margin-bottom: 20px;
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

	.location-description,
	.location-state {
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

	/* Правая панель - NPC и кнопки (15% ширины) */
	.sidebar {
		flex: 15;
		display: flex;
		flex-direction: column;
		gap: 20px;
		min-height: 0; /* Позволяем сжиматься */
	}

	.sidebar-section {
		background: #2c3e50;
		border-radius: 8px;
		padding: 20px;
		overflow: auto;
	}

	.sidebar-section:last-child {
		flex: 1;
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

	.loading,
	.error,
	.no-location {
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

		.interactions-panel,
		.sidebar {
			flex: none;
		}
	}

	/* Стили для взаимодействий */
	.interactions-history {
		flex: 3; /* Увеличиваем размер - займет в 3 раза больше места чем навигация */
		display: flex;
		flex-direction: column;
	}

	.interactions-history h4 {
		color: #3498db;
		margin-bottom: 15px;
		font-size: 1.3em;
	}

	.interactions-scroll {
		flex: 1;
		border: 1px solid #34495e;
		border-radius: 6px;
		padding: 15px;
		background: #34495e;
		overflow-y: auto;
		overflow-x: hidden;
		scroll-behavior: auto; /* Отключаем плавный скролл для лучшего контроля */
	}

	/* Стилизация скроллбара */
	.interactions-scroll::-webkit-scrollbar {
		width: 8px;
	}

	.interactions-scroll::-webkit-scrollbar-track {
		background: #2c3e50;
		border-radius: 4px;
	}

	.interactions-scroll::-webkit-scrollbar-thumb {
		background: #3498db;
		border-radius: 4px;
	}

	.interactions-scroll::-webkit-scrollbar-thumb:hover {
		background: #2980b9;
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

	.interaction-main-content {
		margin-bottom: 0;
	}

	.interaction-additional-content {
		margin-top: 8px;
		font-size: 0.9em;
		color: #95a5a6;
		font-style: italic;
		padding-left: 12px;
		border-left: 2px solid #34495e;
	}

	.no-interactions {
		text-align: center;
		color: #95a5a6;
		font-style: italic;
		padding: 20px;
	}

	/* Стили для поля ввода действий */
	.action-input-section {
		flex-shrink: 0; /* Не даем сжиматься */
		background: #34495e; /* Добавляем фон для выделения */
		border-radius: 8px;
		padding: 15px;
		border-top: 2px solid #3498db; /* Добавляем верхнюю границу */
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
		min-height: 80px; /* Увеличиваем минимальную высоту */
		height: 80px; /* Добавляем фиксированную высоту */
		box-sizing: border-box; /* Учитываем padding в размере */
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
		min-height: 45px; /* Добавляем минимальную высоту */
		flex-shrink: 0; /* Не даем сжиматься */
	}

	.action-submit-btn:hover:not(:disabled) {
		background: #2ecc71;
	}

	.action-submit-btn:disabled {
		background: #95a5a6;
		cursor: not-allowed;
	}

	/* Навигационные стили */
	.navigation-section {
		background: #2c3e50;
		border-radius: 8px;
		padding: 15px;
		margin-bottom: 20px;
		border: 1px solid #34495e;
		flex-shrink: 0; /* Не даем сжиматься */
		flex: 0 0 auto; /* Не растягиваем, используем только необходимое место */
	}

	.navigation-section h4 {
		color: #ecf0f1;
		margin: 0 0 15px 0;
		font-size: 1.1em;
		border-bottom: 2px solid #3498db;
		padding-bottom: 5px;
	}

	.current-point {
		background: #34495e;
		border-radius: 6px;
		padding: 12px;
		margin-bottom: 15px;
		border-left: 4px solid #3498db;
	}

	.current-point strong {
		color: #3498db;
		display: block;
		margin-bottom: 5px;
	}

	.point-description {
		margin: 8px 0;
		color: #bdc3c7;
		font-style: italic;
	}

	.point-type {
		display: inline-block;
		padding: 3px 8px;
		border-radius: 12px;
		font-size: 0.8em;
		font-weight: bold;
		text-transform: uppercase;
	}

	.point-type-entry {
		background: #27ae60;
		color: white;
	}

	.point-type-regular {
		background: #3498db;
		color: white;
	}

	.point-type-special {
		background: #e74c3c;
		color: white;
	}

	.connections-grid {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
		gap: 10px;
	}

	.connection-btn {
		background: #34495e;
		border: 2px solid #3498db;
		border-radius: 8px;
		padding: 12px;
		color: #ecf0f1;
		cursor: pointer;
		transition: all 0.3s;
		text-align: left;
	}

	.connection-btn:hover {
		background: #3498db;
		transform: translateY(-2px);
		box-shadow: 0 4px 8px rgba(52, 152, 219, 0.3);
	}

	.connection-entry {
		border-color: #27ae60;
	}

	.connection-entry:hover {
		background: #27ae60;
	}

	.connection-special {
		border-color: #e74c3c;
	}

	.connection-special:hover {
		background: #e74c3c;
	}

	.connection-name {
		font-weight: bold;
		margin-bottom: 4px;
		font-size: 1em;
	}

	.connection-type {
		font-size: 0.8em;
		color: #bdc3c7;
		text-transform: uppercase;
	}

	/* Стили для контейнеров соединений и маршрутов */
	.connection-container {
		display: flex;
		flex-direction: column;
	}

	.route-preview {
		margin-top: 8px;
		padding: 6px 8px;
		background: rgba(44, 62, 80, 0.8);
		border-radius: 4px;
		border-left: 3px solid #3498db;
		text-align: left; /* Выравниваем по левому краю */
	}

	.route-item {
		font-size: 0.75em;
		color: #95a5a6;
		line-height: 1.3;
		margin-bottom: 2px;
		text-align: left; /* Выравниваем по левому краю */
	}

	.route-item:last-child {
		margin-bottom: 0;
	}

	/* Обновляем сетку для поддержки новой структуры */
	.connections-grid {
		display: grid;
		grid-template-columns: repeat(
			auto-fit,
			minmax(220px, 1fr)
		); /* Увеличиваем минимальную ширину */
		gap: 15px; /* Увеличиваем отступ */
	}
</style>
