<script lang="ts">
  import { onMount } from "svelte";
  import { Moon, Sun } from "lucide-svelte";
  import { setupTheme, theme, toggleTheme } from "$lib/stores/theme";
  import { scale } from "svelte/transition";

  onMount(() => {
    setupTheme();
  });
</script>

<button class="theme-toggle" onclick={toggleTheme} aria-label="Toggle theme">
  {#if $theme === "dark"}
    <div class="icon" transition:scale>
      <Sun size="24" />
    </div>
  {:else}
    <div class="icon" transition:scale>
      <Moon size="24" />
    </div>
  {/if}

  <div class="icon hidden" transition:scale>
    <Moon size="24" />
  </div>
</button>

<style>
  .hidden {
    opacity: 0;
  }

  .theme-toggle {
    color: inherit;

    z-index: 1000;
    border-radius: 50px;
    width: 50px;
    height: 50px;
    cursor: pointer;
    transition: all var(--transition-normal);

    display: flex;
    justify-content: center;
    align-items: center;

    &:hover {
      transform: scale(1.1);
    }

    &:active {
      transform: scale(0.8);
    }

    .icon {
      position: absolute;

      &.hidden {
        opacity: 0;
      }
    }
  }

  @media (max-width: 768px) {
    .theme-toggle {
      right: 1rem;
      padding: 0.5rem;
    }
  }
</style>
