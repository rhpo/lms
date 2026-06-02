<script lang="ts">
    import { tippy } from "$lib/scripts/tippy";

    let {
        href,
        variant = "primary",
        type,
        label = "",
        Icon = null,
        iconPosition = "left",
        autoWidth = false,
        onClick = () => {},
        children = null,
        color = "var(--color-accent)",
        isSubmit = false,
        disabled = false,
        disabledLabel = "",

        large = false,
        uppercase = false,
        round = false,
        blink = false,
        square = false,

        ...rest
    }: {
        href?: string;
        variant?:
            | "primary"
            | "secondary"
            | "third"
            | "cta"
            | "error"
            | "ghost"
            | "success";
        type?: string;
        label?: string;
        Icon?: any;
        iconPosition?: "left" | "right";
        autoWidth?: boolean;
        onClick?: () => void;
        children?: any;
        color?: string;
        disabled?: boolean;
        isSubmit?: boolean;

        round?: boolean;
        large?: boolean;
        uppercase?: boolean;
        blink?: boolean;
        square?: boolean;
    } & Record<string, any> = $props();

    let buttonType = $derived(
        (type ?? (isSubmit ? "submit" : "button")) as
            | "button"
            | "submit"
            | "reset",
    );
</script>

{#if href}
    <a
        class="button {variant}"
        class:large
        class:uppercase
        class:round
        class:blink
        class:square
        class:disabled
        use:tippy={{
            content: disabled ? disabledLabel : null,
            allowHTML: true,
            followCursor: "horizontal",
        }}
        class:auto-width={autoWidth}
        {href}
        {...rest}
        onclick={disabled ? undefined : onClick}
    >
        {#if Icon && iconPosition === "left"}
            <Icon />
        {/if}

        {#if label}
            <span class="label">{label}</span>
        {/if}

        {@render children?.()}

        {#if Icon && iconPosition === "right"}
            <Icon />
        {/if}
    </a>
{:else}
    <button
        class="button {variant}"
        class:large
        class:uppercase
        class:round
        class:blink
        class:square
        class:disabled
        type={buttonType}
        use:tippy={{
            content: disabled ? disabledLabel : null,

            allowHTML: true,
            followCursor: "horizontal",
        }}
        class:auto-width={autoWidth}
        onclick={onClick}
        {disabled}
        {...rest}
    >
        {#if Icon && iconPosition === "left"}
            <Icon />
        {/if}

        {#if label}
            <span class="label">{label}</span>
        {/if}

        {@render children?.()}

        {#if Icon && iconPosition === "right"}
            <Icon />
        {/if}
    </button>
{/if}

<style>
    .button {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        padding: 0.75rem 1.5rem;
        border: none;
        border-radius: 0.7em;
        font-family: var(--font-sans);
        font-weight: 500;
        font-size: 0.95rem;
        text-decoration: none;
        white-space: nowrap;
        cursor: pointer;
        transition: all var(--transition-normal);
        position: relative;
        direction: ltr;
    }

    .button.round {
        border-radius: 50px;
    }

    .button.blink .label {
        animation: margin-right 1s infinite;
    }

    .button.blink {
        animation: blink 2s infinite;
    }

    @keyframes blink {
        0% {
            opacity: 1;
            transform: scale(1);
        }
        50% {
            opacity: 0.5;
            transform: scale(1.02);
        }
        100% {
            opacity: 1;
            transform: scale(1);
        }
    }

    @keyframes margin-right {
        0% {
            margin-right: 0;
        }
        100% {
            margin-right: 0.2rem;
        }
    }

    .button.large {
        font-size: 1.1rem;
        padding: 1rem 2rem;
    }

    .button.uppercase {
        text-transform: uppercase;
    }

    .button[disabled],
    .button.disabled {
        opacity: 0.4;
        /* pointer-events: none; */
        transform: none !important;
        box-shadow: none !important;
        filter: none !important;

        cursor: not-allowed;
    }

    .button:hover {
        transform: translateY(-2px);
    }

    /* primary */
    .primary {
        background: var(--color-accent);
        color: var(--color-background);
    }

    .primary:hover {
        filter: brightness(1.1);
        box-shadow: 0 6px 16px
            color-mix(in srgb, var(--color-accent) 40%, transparent);
    }

    /* secondary */
    .secondary {
        background: transparent;
        color: var(--color-accent);
        border: 2px solid var(--color-accent);
    }

    .secondary:hover {
        background: var(--color-accent);
        color: var(--color-background);
    }

    /* third */
    .button.third {
        background: var(--color-background-100);
        color: var(--color-text);
        box-shadow: 0 4px 12px
            color-mix(in srgb, var(--color-text) 10%, transparent);
    }

    .button.third:hover {
        background: var(--color-background-200);
        box-shadow: 0 6px 16px
            color-mix(in srgb, var(--color-text) 15%, transparent);
    }

    /* cta */
    .cta {
        font-size: 1.1rem;
        padding: 1rem 2rem;
        background: var(--color-accent);
        color: var(--color-background);
    }

    /* ghost - surface background with border, neutral hover */
    .ghost {
        background: var(--color-surface);
        color: var(--color-text);
        border: 1px solid var(--color-border);
        box-shadow: var(--shadow-sm);
    }
    .ghost:hover {
        background: var(--color-background-100);
        color: var(--color-text);
        filter: none;
        box-shadow: var(--shadow-sm);
    }

    /* success */
    .success {
        background: var(--color-success);
        color: var(--color-background);
    }
    .success:hover {
        filter: brightness(1.1);
        box-shadow: 0 6px 16px
            color-mix(in srgb, var(--color-success) 40%, transparent);
    }

    /* square - equal dimensions for icon-only buttons */
    .button.square {
        padding: 0.5rem;
        width: 2.125rem;
        height: 2.125rem;
        border-radius: 8px;
        flex-shrink: 0;
    }
    .button.square.large {
        width: 2.75rem;
        height: 2.75rem;
        padding: 0.65rem;
    }

    /* error */
    .error {
        background: var(--color-error, #e74c3c);
        color: var(--color-background);
        box-shadow: 0 4px 12px
            color-mix(in srgb, var(--color-error, #e74c3c) 30%, transparent);
    }

    @media (max-width: 768px) {
        .button.auto-width {
            width: 100%;
        }

        .button.large {
            font-size: 1rem;
            padding: 0.75rem 1.5rem;
        }
    }
</style>
