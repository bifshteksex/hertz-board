<script lang="ts">
	/**
	 * ConnectorCreator - компонент для создания соединительных линий
	 */
	interface Props {
		startX: number;
		startY: number;
		currentX: number;
		currentY: number;
		connectorType?: 'straight' | 'curved' | 'elbow';
	}

	let { startX, startY, currentX, currentY, connectorType = 'straight' }: Props = $props();

	const previewStyle = {
		stroke: '#3b82f6',
		strokeWidth: 2,
		strokeDasharray: '5 5',
		fill: 'none'
	};

	// Arrow marker
	const arrowMarkerId = 'connector-preview-arrow';
</script>

<defs>
	<marker
		id={arrowMarkerId}
		markerWidth="10"
		markerHeight="10"
		refX="9"
		refY="3"
		orient="auto"
		markerUnits="strokeWidth"
	>
		<path d="M0,0 L0,6 L9,3 z" fill={previewStyle.stroke} />
	</marker>
</defs>

{#if connectorType === 'straight'}
	<line
		x1={startX}
		y1={startY}
		x2={currentX}
		y2={currentY}
		stroke={previewStyle.stroke}
		stroke-width={previewStyle.strokeWidth}
		stroke-dasharray={previewStyle.strokeDasharray}
		marker-end="url(#{arrowMarkerId})"
		pointer-events="none"
	/>
{:else if connectorType === 'curved'}
	{@const midX = (startX + currentX) / 2}
	{@const midY = (startY + currentY) / 2}
	{@const offsetX = (currentY - startY) * 0.2}
	{@const offsetY = (startX - currentX) * 0.2}
	<path
		d="M {startX} {startY} Q {midX + offsetX} {midY + offsetY} {currentX} {currentY}"
		fill={previewStyle.fill}
		stroke={previewStyle.stroke}
		stroke-width={previewStyle.strokeWidth}
		stroke-dasharray={previewStyle.strokeDasharray}
		marker-end="url(#{arrowMarkerId})"
		pointer-events="none"
	/>
{:else if connectorType === 'elbow'}
	{@const midX = (startX + currentX) / 2}
	<path
		d="M {startX} {startY} L {midX} {startY} L {midX} {currentY} L {currentX} {currentY}"
		fill={previewStyle.fill}
		stroke={previewStyle.stroke}
		stroke-width={previewStyle.strokeWidth}
		stroke-dasharray={previewStyle.strokeDasharray}
		marker-end="url(#{arrowMarkerId})"
		pointer-events="none"
	/>
{/if}

<style>
	/* Styles are inline for this component */
</style>
