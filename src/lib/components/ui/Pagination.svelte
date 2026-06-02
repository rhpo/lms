<script lang="ts">
	interface Props {
		page: number;
		max: number;
	}

	const MAX_ALLOWED = 10;

	let { page = $bindable(1), max = $bindable(1) }: Props = $props();

	$effect(() => {
		if (max > MAX_ALLOWED) {
			max = MAX_ALLOWED;
		}
	});


	function getVisiblePages(currentPage: number, totalPages: number): (number | string)[] {
		const delta = 1;
		const range: number[] = [];
		const rangeWithDots: (number | string)[] = [];


		range.push(1);


		const start = Math.max(2, currentPage - delta);
		const end = Math.min(totalPages - 1, currentPage + delta);


		for (let i = start; i <= end; i++) {
			range.push(i);
		}


		if (totalPages > 1) {
			range.push(totalPages);
		}


		const uniqueRange = [...new Set(range)].sort((a, b) => a - b);


		let prev = 0;
		for (const pageNum of uniqueRange) {
			if (pageNum - prev > 1) {
				rangeWithDots.push('...');
			}
			rangeWithDots.push(pageNum);
			prev = pageNum;
		}

		return rangeWithDots;
	}

	function goToPage(newPage: number): void {
		if (newPage >= 1 && newPage <= max && newPage !== page) {
			page = newPage;
		}
	}

	function nextPage(): void {
		if (page < max) {
			page = page + 1;
		}
	}

	function prevPage(): void {
		if (page > 1) {
			page = page - 1;
		}
	}

	const visiblePages = $derived(getVisiblePages(page, max));
</script>

<nav class="pagination-container" aria-label="Pagination Navigation">
	<button
		class="pagination-btn prev-btn"
		class:disabled={page <= 1}
		onclick={prevPage}
		disabled={page <= 1}
		aria-label="Previous page"
		type="button"
	>
		<svg
			width="16"
			height="16"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
		>
			<polyline points="15,18 9,12 15,6"></polyline>
		</svg>
	</button>

	<div class="pages-container">
		{#each visiblePages as pageNum}
			{#if pageNum === '...'}
				<span class="ellipsis" aria-hidden="true">...</span>
			{:else if typeof pageNum === 'number'}
				<button
					class="pagination-btn page-btn"
					class:active={pageNum === page}
					onclick={() => goToPage(pageNum)}
					aria-label="Go to page {pageNum}"
					aria-current={pageNum === page ? 'page' : undefined}
					type="button"
				>
					{pageNum}
				</button>
			{/if}
		{/each}
	</div>

	<button
		class="pagination-btn next-btn"
		class:disabled={page >= max}
		onclick={nextPage}
		disabled={page >= max}
		aria-label="Next page"
		type="button"
	>
		<svg
			width="16"
			height="16"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
		>
			<polyline points="9,18 15,12 9,6"></polyline>
		</svg>
	</button>
</nav>

<style>
	.pagination-container {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 4px;
		padding: 8px;
		background: var(--color-surface);
		border-radius: 50px;
		box-shadow: var(--shadow-sm);
		border: 1px solid var(--color-border);
		width: 100%;
	}

	.pagination-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 40px;
		height: 32px;
		border: none;
		background: transparent;
		border-radius: 50px;
		font-size: 14px;
		font-weight: 500;
		color: var(--color-text);
		cursor: pointer;
		transition: all 0.2s ease;
		user-select: none;
		padding: 0 12px;
	}

	.pagination-btn:hover:not(.disabled):not(.active) {
		background: var(--color-background-100);
		color: var(--color-text);
	}

	.pagination-btn:active:not(.disabled) {
		transform: scale(0.95);
	}

	.pagination-btn.active {
		background: var(--color-accent);
		color: #fff;
		font-weight: 600;
	}

	.pagination-btn.disabled {
		color: var(--color-text-disabled);
		cursor: not-allowed;
	}

	.pages-container {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 100%;
		max-width: 200px;
		gap: 0.5rem;
	}

	.page-btn {
		min-width: 32px;
	}

	.ellipsis {
		display: flex;
		align-items: center;
		justify-content: center;
		min-width: 32px;
		height: 32px;
		color: var(--color-text-muted);
		font-weight: 500;
		user-select: none;
	}

	/* Responsive Design */
	@media (max-width: 640px) {
		.pagination-container {
			padding: 3px;
			gap: 1px;
		}

		.pagination-btn {
			min-width: 28px;
			height: 28px;
			font-size: 12px;
			padding: 0 6px;
		}

		.page-btn {
			min-width: 28px;
		}

		.ellipsis {
			min-width: 28px;
			height: 28px;
		}

		/* Hide some page numbers on very small screens */
		.pages-container {
			overflow: hidden;
			max-width: 180px;
		}
	}

	@media (max-width: 480px) {
		.pagination-container {
			padding: 2px;
			gap: 1px;
		}

		.pagination-btn {
			min-width: 26px;
			height: 26px;
			font-size: 11px;
			padding: 0 4px;
		}

		.pages-container {
			width: 100%;
		}

		.page-btn {
			min-width: 26px;
		}

		.ellipsis {
			min-width: 26px;
			height: 26px;
		}
	}

	/* Focus styles for accessibility */
	.pagination-btn:focus-visible {
		outline: 2px solid var(--color-accent);
		outline-offset: 2px;
	}

	/* Animation for smooth transitions */
	.pagination-btn {
		position: relative;
		overflow: hidden;
	}

	.pagination-btn::before {
		content: '';
		position: absolute;
		top: 50%;
		left: 50%;
		width: 0;
		height: 0;
		background: rgba(59, 130, 246, 0.1);
		border-radius: 50%;
		transform: translate(-50%, -50%);
		transition:
			width 0.3s ease,
			height 0.3s ease;
	}

	.pagination-btn:hover:not(.disabled):not(.active)::before {
		width: 100%;
		height: 100%;
	}
</style>
