<script lang="ts">
	let {
		text = '',
		position = 'top',
		children,
	}: {
		text?: string;
		position?: 'top' | 'bottom' | 'left' | 'right';
		children: any;
	} = $props();
</script>

<div class="tooltip-wrapper">
	<div class="tooltip-trigger">
		{@render children?.()}
	</div>
	<div class="tooltip-content" data-position={position}>
		{text}
	</div>
</div>

<style>
	.tooltip-wrapper {
		position: relative;
		display: inline-flex;
	}

	.tooltip-content {
		position: absolute;
		z-index: 9999;
		padding: 0.4rem 0.75rem;
		border-radius: 6px;
		background: var(--color-text);
		color: var(--color-background);
		font-size: 0.8rem;
		font-family: var(--font-sans);
		white-space: nowrap;
		pointer-events: none;
		opacity: 0;
		transition: opacity 0.15s ease;
	}

	.tooltip-wrapper:hover .tooltip-content {
		opacity: 1;
	}

	.tooltip-content[data-position="top"] {
		bottom: calc(100% + 6px);
		left: 50%;
		transform: translateX(-50%);
	}

	.tooltip-content[data-position="bottom"] {
		top: calc(100% + 6px);
		left: 50%;
		transform: translateX(-50%);
	}

	.tooltip-content[data-position="left"] {
		right: calc(100% + 6px);
		top: 50%;
		transform: translateY(-50%);
	}

	.tooltip-content[data-position="right"] {
		left: calc(100% + 6px);
		top: 50%;
		transform: translateY(-50%);
	}
</style>
