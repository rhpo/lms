<script lang="ts">
    import type { Snippet } from "svelte";

    export interface Column {
        key: string;
        label: string;
        width?: string;
        align?: "left" | "center" | "right";
    }

    interface Props {
        columns: Column[];
        rows: Record<string, any>[];
        title?: string;
        viewAllHref?: string;
        emptyMessage?: string;
        cell?: Snippet<[{ row: Record<string, any>; col: Column }]>;
    }

    let {
        columns,
        rows,
        title,
        viewAllHref,
        emptyMessage = "No data available.",
        cell,
    }: Props = $props();
</script>

<div class="table-wrap">
    {#if title || viewAllHref}
        <div class="table-top">
            {#if title}<h2 class="table-title">{title}</h2>{/if}
            {#if viewAllHref}<a href={viewAllHref} class="view-all">View All</a
                >{/if}
        </div>
    {/if}

    <div class="table-scroll">
        <table class="tbl">
            <thead>
                <tr>
                    {#each columns as col}
                        <th
                            style:width={col.width}
                            style:text-align={col.align ?? "left"}
                            >{col.label}</th
                        >
                    {/each}
                </tr>
            </thead>
            <tbody>
                {#if rows.length === 0}
                    <tr>
                        <td colspan={columns.length} class="empty"
                            >{emptyMessage}</td
                        >
                    </tr>
                {:else}
                    {#each rows as row}
                        <tr>
                            {#each columns as col}
                                <td style:text-align={col.align ?? "left"}>
                                    {#if cell}
                                        {@render cell({ row, col })}
                                    {:else}
                                        {row[col.key] ?? "-"}
                                    {/if}
                                </td>
                            {/each}
                        </tr>
                    {/each}
                {/if}
            </tbody>
        </table>
    </div>
</div>

<style>
    .table-wrap {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        overflow: hidden;
        box-shadow: var(--shadow-sm);
        display: flex;
        flex-direction: column;
    }
    .table-top {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 1.25rem 1.5rem 1rem;
        border-bottom: 1px solid var(--color-border);
    }
    .table-title {
        font-size: var(--text-lg);
        font-weight: 600;
        font-family: var(--font-sans);
        color: var(--color-text);
        margin: 0;
    }
    .view-all {
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-accent);
        font-weight: 500;
        text-decoration: none;
        transition: opacity var(--transition-fast);
    }
    .view-all:hover {
        opacity: 0.75;
        text-decoration: none;
    }

    .table-scroll {
        overflow-x: auto;
        flex: 1;
    }

    .tbl {
        width: 100%;
        border-collapse: collapse;
    }

    thead tr {
        background: var(--color-background);
    }

    th {
        font-size: var(--text-xs);
        font-weight: 600;
        font-family: var(--font-sans);
        color: var(--color-text-muted);
        padding: 0.75rem 1.5rem;
        white-space: nowrap;
        letter-spacing: 0.03em;
        text-transform: uppercase;
    }
    td {
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-text);
        padding: 0.75rem 1.5rem;
        border-top: 1px solid var(--color-border);
        vertical-align: middle;
    }
    tbody tr {
        transition: background var(--transition-fast);
    }
    tbody tr:hover {
        background: var(--color-background);
    }

    .empty {
        text-align: center;
        color: var(--color-text-muted);
        padding: 2.5rem;
        font-style: italic;
        font-family: var(--font-sans);
    }
</style>
