<script lang="ts">
	let {
		headers = [],
		rows = [],
	}: {
		headers: string[];
		rows: Array<Array<string | number | boolean | null | undefined>>;
	} = $props();
</script>

<div class="table-wrapper">
	<table>
		<thead>
			<tr>
				{#each headers as header}
					<th>{header}</th>
				{/each}
			</tr>
		</thead>
		<tbody>
			{#each rows as row}
				<tr>
					{#each row as cell}
						<td>{cell ?? '-'}</td>
					{/each}
				</tr>
			{/each}
		</tbody>
	</table>
</div>

<style>
	.table-wrapper {
		width: 100%;
		overflow-x: auto;
		border-radius: 12px;
		border: 1px solid var(--color-background-200);
	}

	table {
		width: 100%;
		border-collapse: collapse;
		font-family: var(--font-sans);
		font-size: 0.9rem;
	}

	thead {
		background: var(--color-background-100);
	}

	th {
		text-align: left;
		padding: 0.75rem 1rem;
		font-weight: 600;
		color: var(--color-text);
		font-size: 0.8rem;
		text-transform: uppercase;
		letter-spacing: 0.04em;
		border-bottom: 1px solid var(--color-background-200);
	}

	td {
		padding: 0.75rem 1rem;
		color: var(--color-text);
		border-bottom: 1px solid var(--color-background-200);
	}

	tr:last-child td {
		border-bottom: none;
	}

	tbody tr:hover {
		background: var(--color-background-100);
	}
</style>
