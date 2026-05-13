<script lang="ts">
  import type { Snippet } from "svelte";

  type Props = {
    checked?: boolean;
    id?: string;
    size?: number;
    label?: Snippet;
  };

  let {
    checked = $bindable(false),
    id = crypto.randomUUID(),
    size = 31,
    label,
  }: Props = $props();

  const trackWidth = $derived(size * (51 / 31));
  const thumbSize = $derived(size - 4);
  const thumbTravel = $derived(trackWidth - thumbSize - 4);
</script>

<label class="toggle-wrapper" for={id}>
  <div
    class="toggle"
    style="--track-width: {trackWidth}px; --track-height: {size}px; --thumb-size: {thumbSize}px; --thumb-travel: {thumbTravel}px;"
  >
    <input {id} type="checkbox" bind:checked />
    <span class="track">
      <span class="thumb"></span>
    </span>
  </div>

  {#if label}
    <span class="label">{@render label()}</span>
  {/if}
</label>

<style>
  .toggle-wrapper {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    cursor: pointer;
    user-select: none;
  }

  .label {
    font-size: 0.875rem;
    color: var(--color-text);
    line-height: 1.4;
  }

  .toggle {
    position: relative;
    flex-shrink: 0;

    input {
      position: absolute;
      opacity: 0;
      width: 0;
      height: 0;
    }
  }

  .track {
    display: block;
    width: var(--track-width);
    height: var(--track-height);
    background: var(--color-background-200);
    border-radius: 100px;
    transition: background 0.25s ease;
    position: relative;
  }

  .thumb {
    position: absolute;
    top: 2px;
    left: 2px;
    width: var(--thumb-size);
    height: var(--thumb-size);
    background: white;
    border-radius: 50%;
    box-shadow:
      0 2px 4px rgba(0, 0, 0, 0.2),
      0 1px 2px rgba(0, 0, 0, 0.1);
    transition: transform 0.25s ease;
  }

  .toggle input:checked ~ .track {
    background: var(--color-accent);
  }

  .toggle input:checked ~ .track .thumb {
    transform: translateX(var(--thumb-travel));
  }
</style>
