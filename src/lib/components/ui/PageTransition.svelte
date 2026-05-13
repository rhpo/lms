<script lang="ts">
  import { fade as animation } from "svelte/transition";
  import AOS from "aos";

  let {
    key = "",
    duration = 300,
    children,
  }: {
    key?: string;
    duration?: number;
    children?: any;
  } = $props();

  $effect(() => {
    AOS.init({
      duration: 800,
      once: true,
    });
  });
</script>

{#key key}
  <main
    in:animation={{ duration, delay: duration }}
    out:animation={{ duration }}
  >
    {@render children?.()}
  </main>
{/key}

<style>
  main {
    inset: 0;
  }
</style>
