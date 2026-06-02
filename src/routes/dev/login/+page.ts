import { admin, shared } from '$lib/api';

export async function load() {
    try {
        const accounts = await shared.accounts();
        return { accounts: (accounts as any[]) ?? [] };
    } catch {
        return { accounts: [] };
    }
}
