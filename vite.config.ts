import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

const API_BASE = 'http://localhost:5000'

export default defineConfig({
    plugins: [sveltekit()],
    ssr: {
        noExternal: ['svelte-fa', 'svelte-toggles', 'svelte-sonner']
    },
    server: {
        allowedHosts: true,

        proxy: {
            '/api': {
                target: API_BASE,
                changeOrigin: true,
            },

            '/uploads': {
                target: API_BASE,
                changeOrigin: true,
            },
        }
    }
});
