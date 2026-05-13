<script lang="ts">
    import { ChevronDown } from "lucide-svelte";

    interface Option {
        value: string;
        label: string;
    }

    let {
        value,
        options,
        placeholder = "Tous",
        label,
        onChange,
    }: {
        value: string;
        options: Option[];
        placeholder?: string;
        label?: string;
        onChange?: (value: string) => void;
    } = $props();

    let open = $state(false);
    let safeOptions = $derived(options ?? []);
    let selectedLabel = $derived(
        value
            ? (safeOptions.find((o) => o.value === value)?.label ?? placeholder)
            : placeholder,
    );
    let selectEl: HTMLDivElement | undefined = $state();

    function toggle() {
        open = !open;
    }

    function select(val: string) {
        value = val;
        open = false;
        onChange?.(val);
    }

    $effect(() => {
        if (!open) return;
        function handleClick(e: MouseEvent) {
            if (selectEl && !selectEl.contains(e.target as Node)) {
                open = false;
            }
        }
        function handleKey(e: KeyboardEvent) {
            if (e.key === "Escape") open = false;
        }
        requestAnimationFrame(() => {
            document.addEventListener("click", handleClick);
            document.addEventListener("keydown", handleKey);
        });
        return () => {
            document.removeEventListener("click", handleClick);
            document.removeEventListener("keydown", handleKey);
        };
    });
</script>

<div class="select-wrapper" bind:this={selectEl}>
    {#if label}
        <span class="select-label">{label}</span>
    {/if}
    <!-- svelte-ignore a11y_no_noninteractive_tabindex -->
    <button
        class="select-trigger"
        onclick={toggle}
        aria-haspopup="listbox"
        aria-expanded={open}
        type="button"
    >
        <span class="select-value" class:placeholder={!value}
            >{selectedLabel}</span
        >
        <span class="select-chevron" class:open>
            <ChevronDown size={14} />
        </span>
    </button>

    {#if open}
        <div class="select-dropdown" role="listbox">
            <button
                class="select-option"
                class:selected={!value}
                role="option"
                aria-selected="false"
                onclick={() => select("")}
                type="button"
            >
                {placeholder}
            </button>
            {#each safeOptions as option}
                <button
                    class="select-option"
                    class:selected={value === option.value}
                    role="option"
                    aria-selected="false"
                    onclick={() => select(option.value)}
                    type="button"
                >
                    {option.label}
                </button>
            {/each}
        </div>
    {/if}
</div>

<style>
    .select-wrapper {
        position: relative;
        min-width: 160px;
    }

    .select-label {
        display: block;
        font-size: 11px;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.06em;
        color: var(--color-text-muted);
        margin-bottom: 0.3rem;
    }

    .select-trigger {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 0.5rem;
        width: 100%;
        padding: 0.55rem 0.85rem;
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 10px;
        color: var(--color-text);
        font-family: var(--font-sans);
        font-size: var(--text-sm);
        cursor: pointer;
        transition:
            border-color 0.15s,
            box-shadow 0.15s;

        &:hover {
            border-color: var(--color-accent);
        }

        &:focus-visible {
            outline: none;
            border-color: var(--color-accent);
            box-shadow: 0 0 0 3px
                color-mix(in srgb, var(--color-accent) 20%, transparent);
        }
    }

    .select-value {
        flex: 1;
        text-align: left;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;

        &.placeholder {
            color: var(--color-text-muted);
        }
    }

    .select-chevron {
        flex-shrink: 0;
        transition: transform 0.2s;
        color: var(--color-text-muted);

        &.open {
            transform: rotate(180deg);
        }
    }

    .select-dropdown {
        position: absolute;
        top: calc(100% + 4px);
        left: 0;
        right: 0;
        z-index: 50;
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 10px;
        box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
        overflow: hidden;
        max-height: 240px;
        overflow-y: auto;
    }

    .select-option {
        display: block;
        width: 100%;
        padding: 0.55rem 0.85rem;
        background: none;
        border: none;
        color: var(--color-text);
        font-family: var(--font-sans);
        font-size: var(--text-sm);
        text-align: left;
        cursor: pointer;
        transition: background 0.1s;

        &:hover {
            background: var(--color-background-100);
        }

        &.selected {
            background: color-mix(
                in srgb,
                var(--color-accent) 10%,
                var(--color-surface)
            );
            color: var(--color-accent);
            font-weight: 600;
        }

        &:not(:last-child) {
            border-bottom: 1px solid var(--color-border);
        }
    }
</style>
