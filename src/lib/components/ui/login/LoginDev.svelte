<script lang="ts">
    let { autoWidth = false }: { autoWidth?: boolean } = $props();

    import { Terminal } from "lucide-svelte";
</script>

<a class="button" href="/dev/login" class:auto-width={autoWidth}>
    <div class="icon">
        <Terminal size="22px" strokeWidth={2} />
    </div>

    <div class="content">
        Continue with <span class="dev">Dev</span>
    </div>

    {#if autoWidth}
        <span></span>
    {/if}
</a>

<style>
    :root {
        --padding: 1rem;
    }

    .button {
        border: 2px dashed
            color-mix(in srgb, var(--color-text) 25%, transparent);
        background-color: color-mix(
            in srgb,
            var(--color-accent) 6%,
            var(--color-background-200)
        );
        color: var(--color-text);
        display: flex;
        align-items: center;
        gap: calc(var(--padding) / 1.5);

        width: fit-content;
        border-radius: 999rem;
        padding: calc(var(--padding) / 1.7) var(--padding);
        font-size: 1.2rem;
        font-family: var(--font-sans);
        cursor: pointer;
        text-decoration: none;
        transition: all var(--transition-normal);

        &:hover {
            text-decoration: none;
            transform: scale(1.05) translateY(-2px);
            border-style: solid;
        }

        &:hover .dev {
            text-decoration: underline;
        }

        &:active {
            transform: none;
        }

        &.auto-width {
            width: 100%;

            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .icon {
            opacity: 0.65;
        }

        .dev {
            font-weight: 600;
            opacity: 0.8;
        }
    }
</style>
