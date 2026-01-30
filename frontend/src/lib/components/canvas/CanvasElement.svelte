<script lang="ts">
	import type { CanvasElement as CanvasElementType } from '$lib/types/api';
	import { canvasStore } from '$lib/stores/canvas.svelte';
	import TextEditor from './TextEditor.svelte';

	interface Props {
		element: CanvasElementType;
		isSelected: boolean;
		onDoubleClick?: (_element: CanvasElementType) => void;
		onEditStart?: () => void;
		onEditEnd?: () => void;
	}

	let { element, isSelected, onDoubleClick, onEditStart, onEditEnd }: Props = $props();

	// Состояние редактирования
	let isEditing = $state(false);

	// Определяем цвета и стили
	const strokeColor = $derived(isSelected ? '#3b82f6' : element.style?.strokeColor || '#000000');
	const strokeWidth = $derived(isSelected ? 2 : element.style?.strokeWidth || 1);
	const fillColor = $derived(element.style?.backgroundColor || 'transparent');

	function handleDoubleClick(e: MouseEvent) {
		e.stopPropagation();
		if (element.type === 'text' || element.type === 'sticky' || element.type === 'list') {
			isEditing = true;
			if (onEditStart) onEditStart();
		}
		if (onDoubleClick) {
			onDoubleClick(element);
		}
	}

	function handleTextUpdate(content: string, html: string) {
		canvasStore.updateElement(element.id, {
			content: content,
			html_content: html
		});
	}

	function handleBlur() {
		// Use setTimeout to avoid state_unsafe_mutation error
		setTimeout(() => {
			isEditing = false;
			if (onEditEnd) onEditEnd();
		}, 0);
	}

	// Render list content
	function renderListContent(content: string, listType: string): string {
		if (!content) return '';
		const items = content.split('\n').filter((item) => item.trim());
		if (listType === 'checkbox') {
			return items
				.map((item) => {
					const checked = item.startsWith('[x]') || item.startsWith('[X]');
					const text = item.replace(/^\[(x|X| )\]\s*/, '');
					return `<div class="checkbox-item"><input type="checkbox" ${checked ? 'checked' : ''} disabled /><span>${text}</span></div>`;
				})
				.join('');
		}
		return content;
	}
</script>

<g class="canvas-element" data-element-id={element.id} class:selected={isSelected}>
	{#if element.type === 'text'}
		<!-- Text element with rich text editing -->
		<foreignObject
			x={element.pos_x}
			y={element.pos_y}
			width={element.width || 200}
			height={element.height || 100}
			ondblclick={handleDoubleClick}
			role="button"
			tabindex="0"
		>
			{#if isEditing}
				<TextEditor
					content={element.html_content || element.content || ''}
					onUpdate={handleTextUpdate}
					onBlur={handleBlur}
					fontSize={element.style?.fontSize || 16}
					fontFamily={element.style?.fontFamily || 'Inter'}
					color={element.style?.color || '#000000'}
					textAlign={element.style?.textAlign || 'left'}
				/>
			{:else}
				<div
					class="text-content"
					style:font-size="{element.style?.fontSize || 16}px"
					style:font-family={element.style?.fontFamily || 'Inter, sans-serif'}
					style:color={element.style?.color || '#000000'}
					style:text-align={element.style?.textAlign || 'left'}
					style:background-color={element.style?.backgroundColor || 'transparent'}
				>
					{@html element.html_content || element.content || 'Double-click to edit'}
				</div>
			{/if}
		</foreignObject>

		<!-- Selection outline для text -->
		{#if isSelected && !isEditing}
			<rect
				x={element.pos_x}
				y={element.pos_y}
				width={element.width || 200}
				height={element.height || 100}
				fill="none"
				stroke={strokeColor}
				stroke-width={strokeWidth}
				stroke-dasharray="4 4"
				pointer-events="none"
			/>
		{/if}
	{:else if element.type === 'rectangle'}
		<!-- Rectangle with rounded corners support -->
		<rect
			x={element.pos_x}
			y={element.pos_y}
			width={element.width || 100}
			height={element.height || 100}
			fill={fillColor}
			stroke={strokeColor}
			stroke-width={strokeWidth}
			rx={element.style?.borderRadius || 0}
			opacity={element.style?.opacity || 1}
		/>
	{:else if element.type === 'ellipse'}
		<!-- Ellipse -->
		<ellipse
			cx={element.pos_x + (element.width || 100) / 2}
			cy={element.pos_y + (element.height || 100) / 2}
			rx={(element.width || 100) / 2}
			ry={(element.height || 100) / 2}
			fill={fillColor}
			stroke={strokeColor}
			stroke-width={strokeWidth}
			opacity={element.style?.opacity || 1}
		/>
	{:else if element.type === 'triangle'}
		<!-- Triangle -->
		{#if element.width && element.height}
			<polygon
				points="{element.pos_x + element.width / 2},{element.pos_y} {element.pos_x +
					element.width},{element.pos_y + element.height} {element.pos_x},{element.pos_y +
					element.height}"
				fill={fillColor}
				stroke={strokeColor}
				stroke-width={strokeWidth}
				opacity={element.style?.opacity || 1}
			/>
		{/if}
	{:else if element.type === 'line'}
		<!-- Line -->
		<line
			x1={element.pos_x}
			y1={element.pos_y}
			x2={element.pos_x + (element.width || 100)}
			y2={element.pos_y + (element.height || 0)}
			stroke={strokeColor}
			stroke-width={element.style?.strokeWidth || 2}
			opacity={element.style?.opacity || 1}
		/>
	{:else if element.type === 'arrow'}
		<!-- Arrow (line with arrowhead) -->
		<defs>
			<marker
				id="arrowhead-{element.id}"
				markerWidth="10"
				markerHeight="10"
				refX="9"
				refY="3"
				orient="auto"
				markerUnits="strokeWidth"
			>
				<path d="M0,0 L0,6 L9,3 z" fill={strokeColor} />
			</marker>
		</defs>
		<line
			x1={element.pos_x}
			y1={element.pos_y}
			x2={element.pos_x + (element.width || 100)}
			y2={element.pos_y + (element.height || 0)}
			stroke={strokeColor}
			stroke-width={element.style?.strokeWidth || 2}
			marker-end="url(#arrowhead-{element.id})"
			opacity={element.style?.opacity || 1}
		/>
	{:else if element.type === 'sticky'}
		<!-- Sticky note with shadow effect -->
		<defs>
			<filter id="sticky-shadow-{element.id}">
				<feDropShadow dx="2" dy="2" stdDeviation="3" flood-opacity="0.3" />
			</filter>
		</defs>
		<rect
			x={element.pos_x}
			y={element.pos_y}
			width={element.width || 200}
			height={element.height || 200}
			fill={element.style?.backgroundColor || '#fef3c7'}
			stroke={strokeColor}
			stroke-width={strokeWidth}
			rx="4"
			filter="url(#sticky-shadow-{element.id})"
		/>
		<!-- Folded corner effect -->
		<path
			d="M {element.pos_x + (element.width || 200) - 20} {element.pos_y} L {element.pos_x +
				(element.width || 200)} {element.pos_y + 20} L {element.pos_x +
				(element.width || 200)} {element.pos_y} Z"
			fill="#00000010"
		/>
		<foreignObject
			x={element.pos_x + 10}
			y={element.pos_y + 10}
			width={(element.width || 200) - 20}
			height={(element.height || 200) - 20}
			ondblclick={handleDoubleClick}
			role="button"
			tabindex="0"
		>
			{#if isEditing}
				<TextEditor
					content={element.html_content || element.content || ''}
					onUpdate={handleTextUpdate}
					onBlur={handleBlur}
					fontSize={element.style?.fontSize || 14}
					fontFamily="Comic Sans MS"
					color={element.style?.color || '#000000'}
				/>
			{:else}
				<div
					class="sticky-content"
					style:font-size="{element.style?.fontSize || 14}px"
					style:color={element.style?.color || '#000000'}
				>
					{@html element.html_content || element.content || 'Double-click to edit'}
				</div>
			{/if}
		</foreignObject>
	{:else if element.type === 'image' && element.image_url}
		<!-- Image -->
		<image
			x={element.pos_x}
			y={element.pos_y}
			width={element.width || 200}
			height={element.height || 200}
			href={element.image_url}
			preserveAspectRatio="xMidYMid slice"
			opacity={element.style?.opacity || 1}
		/>
		{#if isSelected}
			<rect
				x={element.pos_x}
				y={element.pos_y}
				width={element.width || 200}
				height={element.height || 200}
				fill="none"
				stroke={strokeColor}
				stroke-width={strokeWidth}
				pointer-events="none"
			/>
		{/if}
	{:else if element.type === 'freehand'}
		<!-- Freehand drawing (smooth path) -->
		{#if element.path_data}
			<path
				d={element.path_data}
				fill="none"
				stroke={element.style?.strokeColor || '#000000'}
				stroke-width={element.style?.strokeWidth || 2}
				stroke-linecap="round"
				stroke-linejoin="round"
				opacity={element.style?.opacity || 1}
			/>
		{/if}
	{:else if element.type === 'list'}
		<!-- List (bullet, numbered, checkbox) -->
		<foreignObject
			x={element.pos_x}
			y={element.pos_y}
			width={element.width || 250}
			height={element.height || 200}
			ondblclick={handleDoubleClick}
			role="button"
			tabindex="0"
		>
			{#if isEditing}
				<TextEditor
					content={element.html_content || element.content || ''}
					onUpdate={handleTextUpdate}
					onBlur={handleBlur}
					fontSize={element.style?.fontSize || 14}
					fontFamily={element.style?.fontFamily || 'Inter'}
					color={element.style?.color || '#000000'}
				/>
			{:else}
				<div
					class="list-content"
					class:bullet-list={element.style?.listType === 'bullet'}
					class:numbered-list={element.style?.listType === 'numbered'}
					class:checkbox-list={element.style?.listType === 'checkbox'}
					style:font-size="{element.style?.fontSize || 14}px"
					style:color={element.style?.color || '#000000'}
				>
					{#if element.style?.listType === 'checkbox'}
						{@html renderListContent(element.content || '', 'checkbox')}
					{:else}
						{@html element.html_content || element.content || 'Double-click to edit'}
					{/if}
				</div>
			{/if}
		</foreignObject>

		{#if isSelected && !isEditing}
			<rect
				x={element.pos_x}
				y={element.pos_y}
				width={element.width || 250}
				height={element.height || 200}
				fill="none"
				stroke={strokeColor}
				stroke-width={strokeWidth}
				stroke-dasharray="4 4"
				pointer-events="none"
			/>
		{/if}
	{:else if element.type === 'connector'}
		<!-- Connector line between elements -->
		{#if element.connector_data}
			{@const data = element.connector_data}
			{@const curveType = element.style?.connectorType || 'straight'}

			<defs>
				{#if data.endArrow}
					<marker
						id="connector-end-{element.id}"
						markerWidth="10"
						markerHeight="10"
						refX="9"
						refY="3"
						orient="auto"
						markerUnits="strokeWidth"
					>
						<path d="M0,0 L0,6 L9,3 z" fill={strokeColor} />
					</marker>
				{/if}
				{#if data.startArrow}
					<marker
						id="connector-start-{element.id}"
						markerWidth="10"
						markerHeight="10"
						refX="1"
						refY="3"
						orient="auto"
						markerUnits="strokeWidth"
					>
						<path d="M9,0 L9,6 L0,3 z" fill={strokeColor} />
					</marker>
				{/if}
			</defs>

			{#if curveType === 'straight'}
				<line
					x1={data.startX}
					y1={data.startY}
					x2={data.endX}
					y2={data.endY}
					stroke={strokeColor}
					stroke-width={element.style?.strokeWidth || 2}
					marker-end={data.endArrow ? `url(#connector-end-${element.id})` : undefined}
					marker-start={data.startArrow ? `url(#connector-start-${element.id})` : undefined}
					opacity={element.style?.opacity || 1}
				/>
			{:else if curveType === 'curved'}
				{@const midX = (data.startX + data.endX) / 2}
				{@const midY = (data.startY + data.endY) / 2}
				{@const offsetX = (data.endY - data.startY) * 0.2}
				{@const offsetY = (data.startX - data.endX) * 0.2}
				<path
					d="M {data.startX} {data.startY} Q {midX + offsetX} {midY +
						offsetY} {data.endX} {data.endY}"
					fill="none"
					stroke={strokeColor}
					stroke-width={element.style?.strokeWidth || 2}
					marker-end={data.endArrow ? `url(#connector-end-${element.id})` : undefined}
					marker-start={data.startArrow ? `url(#connector-start-${element.id})` : undefined}
					opacity={element.style?.opacity || 1}
				/>
			{:else if curveType === 'elbow'}
				{@const midX = (data.startX + data.endX) / 2}
				<path
					d="M {data.startX} {data.startY} L {midX} {data.startY} L {midX} {data.endY} L {data.endX} {data.endY}"
					fill="none"
					stroke={strokeColor}
					stroke-width={element.style?.strokeWidth || 2}
					marker-end={data.endArrow ? `url(#connector-end-${element.id})` : undefined}
					marker-start={data.startArrow ? `url(#connector-start-${element.id})` : undefined}
					opacity={element.style?.opacity || 1}
				/>
			{/if}

			<!-- Label if exists -->
			{#if data.label}
				{@const labelX = (data.startX + data.endX) / 2}
				{@const labelY = (data.startY + data.endY) / 2}
				<text
					x={labelX}
					y={labelY}
					text-anchor="middle"
					dominant-baseline="middle"
					font-size="{element.style?.fontSize || 12}px"
					fill={element.style?.color || '#000000'}
					class="connector-label"
				>
					{data.label}
				</text>
			{/if}
		{/if}
	{:else}
		<!-- Fallback для неизвестных типов -->
		<rect
			x={element.pos_x}
			y={element.pos_y}
			width={element.width || 100}
			height={element.height || 100}
			fill="#f3f4f6"
			stroke="#9ca3af"
			stroke-width="1"
			stroke-dasharray="4 4"
		/>
		<text
			x={element.pos_x + (element.width || 100) / 2}
			y={element.pos_y + (element.height || 100) / 2}
			text-anchor="middle"
			dominant-baseline="middle"
			font-size="12"
			fill="#6b7280"
		>
			Unknown type: {element.type}
		</text>
	{/if}
</g>

<style>
	.canvas-element {
		cursor: pointer;
	}

	.canvas-element.selected {
		/* Additional selected styling if needed */
	}

	.text-content {
		width: 100%;
		height: 100%;
		padding: 8px;
		overflow: auto;
		word-wrap: break-word;
	}

	.text-content :global(p) {
		margin: 0 0 0.5em 0;
	}

	.text-content :global(p:last-child) {
		margin-bottom: 0;
	}

	.sticky-content {
		width: 100%;
		height: 100%;
		overflow: auto;
		word-wrap: break-word;
		white-space: pre-wrap;
		font-family: 'Comic Sans MS', cursive, sans-serif;
	}

	.list-content {
		width: 100%;
		height: 100%;
		padding: 8px;
		overflow: auto;
	}

	.list-content :global(.checkbox-item) {
		display: flex;
		align-items: center;
		gap: 8px;
		margin: 4px 0;
	}

	.list-content :global(.checkbox-item input) {
		margin: 0;
		cursor: default;
	}

	.bullet-list :global(ul) {
		list-style-type: disc;
		padding-left: 1.5em;
	}

	.numbered-list :global(ol) {
		list-style-type: decimal;
		padding-left: 1.5em;
	}

	.connector-label {
		background: white;
		padding: 2px 6px;
	}
</style>
