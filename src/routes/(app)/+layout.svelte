<script lang="ts">
  import { BRAND } from "$lib/constants/branding";

  import NavigationBar from "$lib/components/NavigationBar.svelte";

  let { children } = $props();
</script>

<svelte:head>
  <title>
    {BRAND.name} - Main page
  </title>
</svelte:head>

<NavigationBar />

{@render children()}
