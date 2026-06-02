<script lang="ts">
    import { authStore } from "$lib/stores/auth";
    import { onMount } from "svelte";

    onMount(() => {
        authStore.logout();
    });
</script>
