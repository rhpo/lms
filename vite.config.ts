import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
    plugins: [sveltekit()],
    ssr: {
        noExternal: ['svelte-fa', 'svelte-toggles', 'svelte-sonner']
    }
});
