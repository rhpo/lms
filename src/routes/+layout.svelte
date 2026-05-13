<script lang="ts">
  import "$lib/styles/app.css";

  import AOS from "aos";
  import "aos/dist/aos.css";

  import { onMount } from "svelte";

  import favicon from "$lib/assets/favicon.svg";
  import { Toaster } from "svelte-sonner";

  let { children } = $props();

  onMount(() => {
    AOS.init();
  });
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<!-- Global toast container — placed once here so every route can trigger toasts -->
<Toaster
  theme="dark"
  richColors
  position="bottom-right"
  closeButton
  toastOptions={{
    style: "font-family: var(--font-sans, sans-serif); font-size: 0.9rem;",
  }}
/>

{@render children()}
