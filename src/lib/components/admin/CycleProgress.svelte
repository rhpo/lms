<script lang="ts">
    import { Check } from "lucide-svelte";

    export type StepStatus = "completed" | "active" | "upcoming";

    export interface ProgressStep {
        label: string;
        subtitle: string;
        status: StepStatus;
    }

    interface Props {
        steps: ProgressStep[];
    }

    let { steps }: Props = $props();
</script>

<div class="cp-wrap">
    <div class="cp-track">
        {#each steps as step, i}
            <!-- connector line BEFORE the node (skip for first) -->
            {#if i > 0}
                <div
                    class="cp-connector"
                    class:filled={steps[i - 1].status === "completed"}
                ></div>
            {/if}

            <div class="cp-step">
                <div class="cp-node" data-status={step.status}>
                    {#if step.status === "completed"}
                        <Check size={13} color="#fff" strokeWidth={3} />
                    {:else if step.status === "active"}
                        <div class="cp-dot"></div>
                    {/if}
                </div>
                <div class="cp-info">
                    <p class="cp-step-label">{step.label}</p>
                    <p
                        class="cp-step-sub"
                        class:active-sub={step.status === "active"}
                    >
                        {step.subtitle}
                    </p>
                </div>
            </div>
        {/each}
    </div>
</div>

<style>
    .cp-wrap {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: 1.75rem 2rem;
        box-shadow: var(--shadow-sm);
    }
    .cp-track {
        display: flex;
        align-items: flex-start;
        gap: 0;
    }
    .cp-connector {
        flex: 1;
        height: 2px;
        background: var(--color-border);
        margin-top: 17px; /* align with node center */
        transition: background var(--transition-normal);
    }
    .cp-connector.filled {
        background: var(--color-accent);
    }

    .cp-step {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.6rem;
        flex-shrink: 0;
    }
    .cp-node {
        width: 34px;
        height: 34px;
        border-radius: 50%;
        border: 2px solid var(--color-border);
        background: var(--color-surface);
        display: flex;
        align-items: center;
        justify-content: center;
        transition: all var(--transition-normal);
    }
    .cp-node[data-status="completed"] {
        background: var(--color-accent);
        border-color: var(--color-accent);
    }
    .cp-node[data-status="active"] {
        border-color: var(--color-accent);
        border-width: 2px;
    }
    .cp-dot {
        width: 10px;
        height: 10px;
        border-radius: 50%;
        background: var(--color-accent);
    }
    .cp-info {
        display: flex;
        flex-direction: column;
        align-items: center;
        text-align: center;
        gap: 0.2rem;
        min-width: 100px;
    }
    .cp-step-label {
        font-size: var(--text-sm);
        font-weight: 600;
        font-family: var(--font-sans);
        color: var(--color-text);
        margin: 0;
        white-space: nowrap;
    }
    .cp-step-sub {
        font-size: var(--text-xs);
        font-family: var(--font-sans);
        color: var(--color-text-muted);
        margin: 0;
        white-space: nowrap;
    }
    .active-sub {
        color: var(--color-accent);
        font-weight: 500;
    }

    @media (max-width: 600px) {
        .cp-track {
            flex-direction: column;
            align-items: flex-start;
        }
        .cp-connector {
            width: 2px;
            height: 32px;
            margin-top: 0;
            margin-left: 16px;
            flex: unset;
        }
        .cp-step {
            flex-direction: row;
            align-items: center;
        }
        .cp-info {
            text-align: left;
            align-items: flex-start;
        }
    }
</style>
