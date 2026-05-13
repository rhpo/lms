import type { Writable } from "svelte/store";

import { get, writable } from "svelte/store";

type Theme = 'dark' | 'light';

export const theme: Writable<Theme> = writable('dark');

export function setTheme(value: Theme) {
    theme.set(value);

    localStorage.setItem("theme", value);

    if (get(theme) === 'dark') {
        document.documentElement.classList.add("dark")
    } else document.documentElement.classList.remove("dark")

}

export function toggleTheme() {
    setTheme(get(theme) === 'light' ? 'dark' : 'light')
}

export function setupTheme() {
    let rawValue = localStorage.getItem("theme");

    if (rawValue === null)
        setTheme("light");
    else {
        const value = rawValue as Theme;
        setTheme(value);
    }

}
