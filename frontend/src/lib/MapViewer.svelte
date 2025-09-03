<script>
	import { onMount, createEventDispatcher } from "svelte"
	import { GetLocationHierarchy } from "../../wailsjs/go/app/App.js"

	const dispatch = createEventDispatcher()

	export let visible = false

	let hierarchy = null
	let mapContainer
	let currentPointId = ""

	onMount(() => {
		if (visible) {
			loadHierarchy()
		}
	})

	$: if (visible) {
		loadHierarchy()
	}

	async function loadHierarchy() {
		try {
			console.log("Loading hierarchy...")
			hierarchy = await GetLocationHierarchy()
			console.log("Hierarchy loaded:", hierarchy)
			console.log("Clusters:", hierarchy?.clusters)
		} catch (error) {
			console.error("Failed to load hierarchy:", error)
		}
	}

	function closeMap() {
		dispatch("close")
	}

	function getPointTypeIcon(pointType) {
		switch (pointType) {
			case "entry":
				return "🚪"
			case "regular":
				return "📍"
			case "special":
				return "⭐"
			case "exit":
				return "🚪"
			default:
				return "📍"
		}
	}

	function getClusterTypeIcon(clusterType) {
		switch (clusterType) {
			case "mountain":
				return "🏔️"
			case "forest":
				return "🌲"
			case "cave":
				return "🕳️"
			case "village":
				return "🏘️"
			case "ruins":
				return "🏛️"
			case "swamp":
				return "🦎"
			default:
				return "🗺️"
		}
	}

	function getSubClusterColor(subClusterId) {
		// Генерируем цвет на основе ID субкластера
		const colors = ["#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4", "#FECA57", "#FF9FF3", "#54A0FF"]
		let hash = 0
		for (let i = 0; i < subClusterId.length; i++) {
			hash = subClusterId.charCodeAt(i) + ((hash << 5) - hash)
		}
		return colors[Math.abs(hash) % colors.length]
	}
</script>

{#if visible}
	<div class="map-overlay" on:click={closeMap}>
		<div class="map-container" on:click|stopPropagation bind:this={mapContainer}>
			<div class="map-header">
				<h2>🗺️ Карта мира</h2>
				<button class="close-btn" on:click={closeMap}>✕</button>
			</div>

			{#if hierarchy}
				<div class="map-content">
					{#each Object.values(hierarchy.clusters) as cluster}
						<div class="cluster" data-cluster-id={cluster.id}>
							<div class="cluster-header">
								<span class="cluster-icon">{getClusterTypeIcon(cluster.type)}</span>
								<h3>{cluster.name}</h3>
							</div>
							<div class="cluster-description">{cluster.description}</div>

							<div class="subclusters">
								{#each Object.values(cluster.sub_clusters) as subCluster}
									<div
										class="subcluster"
										style="border-color: {getSubClusterColor(subCluster.id)}">
										<div
											class="subcluster-header"
											style="background-color: {getSubClusterColor(
												subCluster.id
											)}20">
											<h4>{subCluster.name}</h4>
										</div>
										<div class="subcluster-description">
											{subCluster.description}
										</div>

										<div class="points-container">
											{#each Object.values(subCluster.points) as point}
												<div
													class="point"
													class:entry-point={point.is_entry_point}
													data-point-id={point.id}>
													<div class="point-icon">
														{getPointTypeIcon(point.type)}
													</div>
													<div class="point-info">
														<div class="point-name">{point.name}</div>
														<div class="point-description">
															{point.description}
														</div>
														{#if point.connections && point.connections.length > 0}
															<div class="connections">
																Соединения: {point.connections
																	.length}
															</div>
														{/if}
													</div>
												</div>
											{/each}
										</div>
									</div>
								{/each}
							</div>
						</div>
					{/each}
				</div>
			{:else}
				<div class="loading">🔄 Загрузка карты...</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	.map-overlay {
		position: fixed;
		top: 0;
		left: 0;
		width: 100vw;
		height: 100vh;
		background: rgba(0, 0, 0, 0.8);
		display: flex;
		justify-content: center;
		align-items: center;
		z-index: 1000;
	}

	.map-container {
		background: white;
		border-radius: 12px;
		padding: 24px;
		max-width: 90vw;
		max-height: 90vh;
		overflow-y: auto;
		box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
		position: relative;
	}

	.map-header {
		display: flex;
		justify-content: space-between;
		align-items: center;
		margin-bottom: 24px;
		border-bottom: 2px solid #eee;
		padding-bottom: 16px;
	}

	.map-header h2 {
		margin: 0;
		color: #333;
		font-size: 1.5em;
	}

	.close-btn {
		background: #ff4757;
		color: white;
		border: none;
		border-radius: 50%;
		width: 32px;
		height: 32px;
		cursor: pointer;
		font-size: 16px;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	.close-btn:hover {
		background: #ff3838;
	}

	.map-content {
		min-width: 800px;
		color: #2d3436; /* Добавляем тёмный цвет текста */
	}

	.loading {
		text-align: center;
		font-size: 1.2em;
		color: #636e72;
		padding: 40px;
	}

	.cluster {
		border: 3px solid #ddd;
		border-radius: 16px;
		margin-bottom: 32px;
		padding: 20px;
		background: #f8f9fa;
	}

	.cluster-header {
		display: flex;
		align-items: center;
		gap: 12px;
		margin-bottom: 12px;
	}

	.cluster-icon {
		font-size: 2em;
	}

	.cluster-header h3 {
		margin: 0;
		color: #2d3436;
		font-size: 1.4em;
	}

	.cluster-description {
		color: #636e72;
		margin-bottom: 20px;
		font-style: italic;
	}

	.subclusters {
		display: flex;
		flex-wrap: wrap;
		gap: 16px;
	}

	.subcluster {
		border: 2px solid;
		border-radius: 12px;
		padding: 16px;
		background: white;
		min-width: 250px;
		flex: 1;
	}

	.subcluster-header {
		padding: 8px 12px;
		border-radius: 8px;
		margin: -16px -16px 16px -16px;
	}

	.subcluster-header h4 {
		margin: 0;
		color: #2d3436;
		font-size: 1.1em;
	}

	.subcluster-description {
		color: #636e72;
		font-size: 0.9em;
		margin-bottom: 16px;
		font-style: italic;
	}

	.points-container {
		display: flex;
		flex-direction: column;
		gap: 12px;
	}

	.point {
		display: flex;
		align-items: flex-start;
		gap: 12px;
		padding: 12px;
		border-radius: 8px;
		background: #f8f9fa;
		border: 1px solid #dee2e6;
		transition: all 0.2s ease;
	}

	.point:hover {
		background: #e9ecef;
		border-color: #adb5bd;
	}

	.point.entry-point {
		background: #d4edda;
		border-color: #28a745;
	}

	.point-icon {
		font-size: 1.5em;
		min-width: 32px;
		text-align: center;
	}

	.point-info {
		flex: 1;
	}

	.point-name {
		font-weight: bold;
		color: #2d3436;
		margin-bottom: 4px;
	}

	.point-description {
		color: #636e72;
		font-size: 0.9em;
		margin-bottom: 6px;
	}

	.connections {
		font-size: 0.8em;
		color: #74b9ff;
		font-weight: 500;
	}

	.loading {
		text-align: center;
		padding: 40px;
		font-size: 1.2em;
		color: #636e72;
	}
</style>
