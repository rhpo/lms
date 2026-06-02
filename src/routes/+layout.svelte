<script lang="ts">
    import "$lib/styles/app.css";
    import "aos/dist/aos.css";

    import AOS from "aos";
    import { onMount } from "svelte";
    import { toast, Toaster } from "svelte-sonner";
    import { navigating } from "$app/stores";
    import { Loader2 } from "lucide-svelte";

    onMount(() => {
        AOS.init({
            disable: "mobile",
            once: true,
        });
    });

    let { children } = $props();
</script>

<Toaster />

{#if $navigating}
    <div class="global-loading-overlay">
        <Loader2 class="spinner" size={40} />
    </div>
{/if}

{@render children()}

<style>
    :global([data-aos]) {
        transition-duration: var(--transition-duration) !important;
    }

    .global-loading-overlay {
        position: fixed;
        inset: 0;
        z-index: 9999;
        background: rgba(255, 255, 255, 0.4);
        backdrop-filter: blur(2px);
        display: flex;
        align-items: center;
        justify-content: center;
    }

    :global(.dark) .global-loading-overlay {
        background: rgba(15, 23, 42, 0.4);
    }

    .spinner {
        color: var(--color-accent);
        animation: spin 1s linear infinite;
    }

    @keyframes spin {
        from {
            transform: rotate(0deg);
        }
        to {
            transform: rotate(360deg);
        }
    }
</style>
