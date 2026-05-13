<script lang="ts">
    export type StatusType = "pending" | "validated" | "rejected" | "in-review";

    interface Props {
        status: StatusType;
    }

    let { status }: Props = $props();

    const map: Record<StatusType, string> = {
        pending: "PENDING",
        validated: "VALIDATED",
        rejected: "REJECTED",
        "in-review": "IN REVIEW",
    };
</script>

<span class="badge {status}">{map[status] ?? status.toUpperCase()}</span>

<style>
    .badge {
        display: inline-block;
        font-size: 0.68rem;
        font-weight: 700;
        font-family: var(--font-sans);
        letter-spacing: 0.05em;
        padding: 0.22rem 0.65rem;
        border-radius: 5px;
        white-space: nowrap;
    }
    .pending {
        background: color-mix(in srgb, #f59e0b 15%, var(--color-surface));
        color: #92400e;
    }
    .validated {
        background: color-mix(in srgb, #10b981 15%, var(--color-surface));
        color: #065f46;
    }
    .rejected {
        background: color-mix(in srgb, #ef4444 15%, var(--color-surface));
        color: #991b1b;
    }
    .in-review {
        background: color-mix(in srgb, #8b5cf6 15%, var(--color-surface));
        color: #5b21b6;
    }

    :global(html.dark) .pending {
        color: #fcd34d;
    }
    :global(html.dark) .validated {
        color: #6ee7b7;
    }
    :global(html.dark) .rejected {
        color: #fca5a5;
    }
    :global(html.dark) .in-review {
        color: #c4b5fd;
    }
</style>
