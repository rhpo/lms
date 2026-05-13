<script lang="ts">
    import type { Component, Snippet } from "svelte";

    import { page } from "$app/stores";
    import { BRAND } from "$lib/constants/branding";
    import { LogOut, Bell, User, Settings } from "lucide-svelte";

    import Logo from "../Logo.svelte";
    import Profile from "./Profile.svelte";
    import Button from "./Button.svelte";
    import ThemeToggle from "./ThemeToggle.svelte";

    interface NavLink {
        href: string;
        label: string;
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        icon: any;
        count?: number;
    }

    interface Props {
        links: NavLink[];
        user: { full_name: string; role: string };
        notificationsHref: string;
        children: Snippet;
    }

    let { links, user, notificationsHref, children }: Props = $props();

    // Compute total count from sidebar nav links for the profile trigger
    let sidebarTotalCount = $derived(
        links.reduce((sum, link) => sum + (link.count ?? 0), 0) + 1,
    );

    const profileLinks = $derived([
        { href: "/status", label: "Statut", icon: User },
        {
            href: "/notifications",
            label: "Notifications",
            icon: Bell,
            count: sidebarTotalCount,
        },
        { href: "/admin/settings", label: "Paramètres", icon: Settings },
    ]);

    function isActive(href: string) {
        return (
            $page.url.pathname === href ||
            $page.url.pathname.startsWith(href + "/")
        );
    }
</script>

<div class="app-shell">
    <!-- == Topnav ============= -->
    <header class="topnav">
        <div class="topnav-left">
            <a class="brand" href={links[0]?.href ?? "/"}>
                <div class="logo">
                    <Logo
                        height="65%"
                        textColor="transparent"
                        logoColor="white"
                    />
                </div>

                <h1>{BRAND.name}</h1>
            </a>
        </div>

        <div class="topnav-right">
            <ThemeToggle />

            <Profile {user} links={profileLinks} />
        </div>
    </header>

    <!-- == Body: sidebar + main ============= -->
    <div class="app-body">
        <!-- == Sidebar ======== -->
        <aside class="sidebar">
            <main>
                <div class="sidebar-brand">
                    <span class="brand-title">PLATEFORME {BRAND.name}</span>
                    <span class="brand-sub">
                        PANNEAU {user.role.toUpperCase()}
                    </span>
                </div>

                <nav class="sidebar-nav" aria-label="Navigation principale">
                    {#each links as item}
                        <a
                            href={item.href}
                            class="nav-item"
                            class:active={isActive(item.href)}
                            aria-current={isActive(item.href)
                                ? "page"
                                : undefined}
                        >
                            <item.icon size={18} />
                            <span class="nav-label">{item.label}</span>
                            {#if item.count && item.count > 0}
                                <span class="nav-badge"
                                    >{item.count > 99
                                        ? "99+"
                                        : item.count}</span
                                >
                            {/if}
                        </a>
                    {/each}
                </nav>

                <div class="sidebar-footer">
                    <a href="/auth/logout" class="nav-item logout-item">
                        <LogOut size={18} />
                        <span class="nav-label">Deconnexion</span>
                    </a>
                </div>
            </main>
        </aside>

        <!-- == Main content ================= -->
        <main class="app-main">
            {@render children()}
        </main>
    </div>
</div>

<style>
    :global(body) {
        font-family: var(--font-sans);
        background: var(--color-background);
    }

    .app-shell {
        display: flex;
        flex-direction: column;
        height: 100vh;
        overflow: hidden;
        background: var(--color-background);
    }

    /* == Topnav == */
    .topnav {
        height: var(--nav-height);
        background: var(--color-accent-fixed);
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0 1.5rem;
        position: sticky;
        top: 0;
        z-index: 50;
    }

    .topnav-left {
        display: flex;
        align-items: center;
        gap: 2rem;
        height: 100%;
    }

    .brand {
        display: flex;
        gap: 1rem;
        align-items: center;
        text-decoration: none;
        height: 100%;

        .logo {
            height: 100%;
            display: flex;
            align-items: center;
            justify-content: center;
            transition: all var(--transition-fast);
        }

        h1 {
            font-weight: 500;
            font-family: var(--font-serif);
        }

        &:hover .logo {
            transform: scale(1.1);
        }

        &:active {
            transform: scale(0.9);
        }

        &:active .logo {
            transform: none;
        }
    }

    .brand * {
        color: var(--color-gray-50);
    }

    .topnav-right {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        color: var(--color-gray-50);

        :global(button) {
            color: var(--color-gray-50);
        }
    }

    /* == Body ==== */
    .app-body {
        display: flex;
        flex: 1;
        min-height: 0;
        overflow: hidden;
    }

    /* == Sidebar ================ */
    .sidebar {
        width: var(--sidebar-width, 240px);
        height: 100%;
        flex-shrink: 0;
        display: flex;
        flex-direction: column;
        background: var(--color-surface);
        border-right: 1px solid var(--color-border);
        overflow-y: auto;

        main {
            display: flex;
            flex-direction: column;
            flex: 1;
            min-height: 0;
        }
    }

    .sidebar-brand {
        padding: 1.25rem 1.25rem 0.75rem;
        border-bottom: 1px solid var(--color-border);
    }
    .brand-title {
        display: block;
        font-size: 0.65rem;
        font-weight: 700;
        font-family: var(--font-sans);
        color: var(--color-accent);
        letter-spacing: 0.1em;
    }
    .brand-sub {
        display: block;
        font-size: 0.6rem;
        font-family: var(--font-sans);
        color: var(--color-text-muted);
        letter-spacing: 0.08em;
        margin-top: 0.15rem;
    }

    .sidebar-nav {
        display: flex;
        flex-direction: column;
        padding: 0.75rem 0.75rem;
        gap: 0.15rem;
        flex: 1;
    }

    .nav-item {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.65rem 0.85rem;
        border-radius: 8px;
        text-decoration: none;
        color: var(--color-text-muted);
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        font-weight: 500;
        transition:
            background var(--transition-fast),
            color var(--transition-fast);
        cursor: pointer;
    }
    .nav-item:hover {
        background: var(--color-background-100);
        color: var(--color-text);
        text-decoration: none;
    }
    .nav-item.active {
        background: var(--color-accent);
        color: #fff;
    }

    .nav-label {
        flex: 1;
    }

    .nav-badge {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        min-width: 18px;
        height: 18px;
        padding: 0 5px;
        border-radius: 999px;
        background: #ef4444;
        color: #fff;
        font-size: 10px;
        font-weight: 700;
        line-height: 1;
    }

    .nav-item.active .nav-badge {
        background: rgba(255, 255, 255, 0.85);
        color: #ef4444;
    }

    .sidebar-footer {
        padding: 0.75rem;
        border-top: 1px solid var(--color-border);
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .logout-item {
        color: var(--color-text-muted);
    }

    /* == Main content =========== */
    .app-main {
        flex: 1;
        min-width: 0;
        overflow-y: auto;
        padding: var(--spacing-lg);
    }

    /* == Responsive ============= */
    @media (max-width: 768px) {
        .sidebar {
            display: none;
        }

        .app-main {
            padding: var(--spacing-md);
        }
    }
</style>
