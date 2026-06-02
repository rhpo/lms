<script lang="ts">
    import { RefreshCw, UserCircle } from "lucide-svelte";
    import { authStore } from "$lib/stores/auth";

    let { data } = $props();

    type Role = "admin" | "teacher" | "student" | "company";

    type Persona = {
        id: string;
        email: string;
        role: Role;
        full_name: string;
        subtitle: string;
    };

    const PERSONAS: Persona[] = $derived([...data.accounts]);

    let showCustom = $state(false);
    let cRole = $state<Role>("student");
    let cName = $state("");
    let cEmail = $state("");
    let cId = $state(`dev-custom-${Math.random().toString(36).slice(2, 8)}`);

    function newId() {
        cId = `dev-custom-${Math.random().toString(36).slice(2, 8)}`;
    }

    let loginError = $state("");

    async function loginAs(email: string) {
        loginError = "";
        try {
            await authStore.devLogin(email);
        } catch (err) {
            loginError = err instanceof Error ? err.message : "Login error";
        }
    }

    async function loginCustom() {
        if (!cEmail) return;
        await loginAs(cEmail);
    }

    function initials(name: string) {
        return name
            .split(/\s+/)
            .filter(Boolean)
            .slice(0, 2)
            .map((s) => s[0]?.toUpperCase() ?? "")
            .join("");
    }

    const ROLE_COLORS: Record<Role, string> = {
        admin: "#b91c1c",
        teacher: "#1d4ed8",
        student: "#15803d",
        company: "#6d28d9",
    };
</script>

<svelte:head>
    <title>Sign in with Google</title>
</svelte:head>

<div class="page">
    <div class="card">
        <!-- Card header bar -->
        <div class="card-header">
            <img src="/media/icons/google.png" alt="Google" class="g-icon" />
            <span class="header-text">Sign in with Google</span>
        </div>

        <!-- Card body -->
        <div class="card-body">
            <div class="left">
                <h1 class="title">Choose an account</h1>
                <p class="sub">
                    to continue to <a href="/" class="app-link">PFE Manager</a>
                </p>
            </div>

            <div class="right">
                <ul class="account-list" role="list">
                    {#each PERSONAS as persona (persona.id)}
                        <li class="account-item">
                            <button
                                type="button"
                                class="account-row"
                                onclick={() => loginAs(persona.email)}
                                disabled={authStore.loading}
                            >
                                <span
                                    class="avatar"
                                    style="background:{ROLE_COLORS[
                                        persona.role
                                    ]}"
                                    aria-hidden="true"
                                    >{initials(persona.full_name)}</span
                                >
                                <span class="account-info">
                                    <span class="account-name"
                                        >{persona.full_name}</span
                                    >
                                    <span class="account-email"
                                        >{persona.email}</span
                                    >
                                </span>
                            </button>
                        </li>
                    {/each}

                    <li class="account-item">
                        <button
                            type="button"
                            class="account-row use-another"
                            onclick={() => (showCustom = !showCustom)}
                        >
                            <span class="avatar avatar-icon" aria-hidden="true">
                                <UserCircle size="1.4rem" color="#444" />
                            </span>
                            <span class="account-info">
                                <span class="account-name"
                                    >Use Another account</span
                                >
                            </span>
                        </button>
                    </li>
                </ul>

                {#if loginError}
                    <p class="login-error" role="alert">{loginError}</p>
                {/if}

                {#if showCustom}
                    <div class="custom-form">
                        <div class="custom-grid">
                            <div class="field">
                                <label for="c-name">Full name</label>
                                <input
                                    id="c-name"
                                    type="text"
                                    placeholder="Name Surname"
                                    bind:value={cName}
                                />
                            </div>
                            <div class="field">
                                <label for="c-role">Role</label>
                                <select id="c-role" bind:value={cRole}>
                                    <option value="admin">Admin</option>
                                    <option value="teacher">Teacher</option>
                                    <option value="student">Student</option>
                                    <option value="company">Company</option>
                                </select>
                            </div>
                            <div class="field field-full">
                                <label for="c-email">Email</label>
                                <input
                                    id="c-email"
                                    type="email"
                                    placeholder="user@dev.local"
                                    bind:value={cEmail}
                                />
                            </div>
                            <div class="field field-full">
                                <label for="c-id">
                                    Identifier
                                    <span class="label-hint"
                                        >stable across logins</span
                                    >
                                </label>
                                <span class="id-row">
                                    <input
                                        id="c-id"
                                        type="text"
                                        class="mono"
                                        bind:value={cId}
                                    />
                                    <button
                                        type="button"
                                        class="btn-refresh"
                                        onclick={newId}
                                        title="Regenerate"
                                        aria-label="Regenerate identifier"
                                    >
                                        <RefreshCw size={13} strokeWidth={2} />
                                    </button>
                                </span>
                            </div>
                        </div>
                        <div class="form-actions">
                            <button
                                type="button"
                                class="btn-next"
                                onclick={loginCustom}
                                disabled={!cEmail || authStore.loading}
                                >Next</button
                            >
                        </div>
                    </div>
                {/if}

                <div class="divider"></div>

                <p class="privacy-note">
                    To continue, Google will share your name, email address,
                    language preference, and profile picture with PFE Manager.
                    Before using this app, you can review PFE Manager's
                    <a href="/privacy" class="app-link">privacy policy</a> and
                    <a href="/terms" class="app-link">terms of service</a>.
                </p>
            </div>
        </div>
    </div>

    <footer class="page-footer">
        <div class="lang">
            <button type="button" class="lang-btn">
                English (United States)
                <svg viewBox="0 0 10 6" width="10" height="6" fill="none">
                    <path
                        d="M1 1l4 4 4-4"
                        stroke="#5f6368"
                        stroke-width="1.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    />
                </svg>
            </button>
        </div>
        <ul class="footer-links">
            <li><a href="#help">Help</a></li>
            <li><a href="#privacy">Privacy</a></li>
            <li><a href="#terms">Terms</a></li>
        </ul>
    </footer>
</div>

<style>
    @import url("https://fonts.googleapis.com/css2?family=Google+Sans:ital,opsz,wght@0,17..18,400..700;1,17..18,400..700&display=swap");

    :root {
        --card-height: 600px;
        --google-header-height: 34px;
    }

    .page {
        min-height: 100vh;
        background: #eff1f5;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        font-family: "Google Sans", Roboto, "Helvetica Neue", Arial, sans-serif;
        color: #202124;
        padding: 24px 16px;
    }

    /* Card */
    .card {
        width: 100%;
        max-width: 1200px;
        height: var(--card-height);
        overflow: hidden;
        background: #fff;
        border-radius: 28px;
        box-shadow:
            0 1px 3px rgba(0, 0, 0, 0.08),
            0 4px 16px rgba(0, 0, 0, 0.06);
    }

    .card-header {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 0 20px;
        height: var(--google-header-height);
        border-bottom: 1px solid #e8eaed;
    }

    .g-icon {
        width: 18px;
        height: 18px;
        object-fit: contain;
    }

    .header-text {
        font-size: 14px;
        color: #444746;
        font-weight: 400;
    }

    .card-body {
        display: grid;
        grid-template-columns: 1fr 1fr;
        padding: 48px 56px 48px 48px;
        gap: 0;

        height: calc(100% - var(--google-header-height));

        > * {
            height: calc(
                var(--card-height) - var(--google-header-height) - 48px * 2
            );
        }
    }

    /* Left */
    .left {
        display: flex;
        flex-direction: column;
        gap: 12px;
        padding-top: 4px;
        padding-right: 32px;

        .title {
            font-size: 3rem;
            font-weight: 400;
            line-height: 1.3;
            margin: 0;
            color: #202124;
            letter-spacing: -0.01em;
        }

        .sub {
            font-size: 15px;
            color: #444746;
            margin: 0;
        }

        .app-link {
            color: #1a73e8;
            text-decoration: none;
        }

        .app-link:hover {
            text-decoration: underline;
        }
    }

    /* Right */
    .right {
        display: flex;
        flex-direction: column;
    }

    .account-list {
        list-style: none;
        padding: 0;
        margin: 0;

        height: 100%;

        overflow-y: scroll;
    }

    .account-item {
        border-bottom: 1px solid #e8eaed;
    }

    .account-item:last-child {
        border-bottom: none;
    }

    .account-row {
        width: 100%;
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 12px 0;
        background: transparent;
        border: none;
        cursor: pointer;
        text-align: left;
        font-family: inherit;
        color: inherit;
        transition: opacity 0.15s;
    }

    .account-row:hover:not(:disabled) {
        opacity: 0.75;
    }

    .account-row:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .avatar {
        width: 32px;
        height: 32px;
        border-radius: 50%;
        color: #fff;
        font-size: 13px;
        font-weight: 500;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
    }

    .avatar-icon {
        background: transparent;
    }

    .account-info {
        display: flex;
        flex-direction: column;
        gap: 1px;
        min-width: 0;
    }

    .account-name {
        font-size: 14px;
        font-weight: 400;
        color: #202124;
        line-height: 1.4;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .account-email {
        font-size: 13px;
        color: #5f6368;
        line-height: 1.4;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .login-error {
        font-size: 13px;
        color: #c5221f;
        margin: 8px 0 0;
        padding: 8px 10px;
        background: #fce8e6;
        border-radius: 4px;
    }

    /* Custom form */
    .custom-form {
        margin-top: 12px;
        padding: 14px;
        background: #f8f9fa;
        border-radius: 8px;
        border: 1px solid #e8eaed;
        display: flex;
        flex-direction: column;
        gap: 10px;
    }

    .custom-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 10px;
    }

    .field-full {
        grid-column: 1 / -1;
    }

    .field {
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .field label {
        font-size: 11px;
        font-weight: 500;
        color: #5f6368;
        display: flex;
        align-items: center;
        gap: 6px;
    }

    .label-hint {
        font-weight: 400;
        color: #80868b;
    }

    .field input,
    .field select {
        height: 34px;
        padding: 0 10px;
        border: 1px solid #dadce0;
        border-radius: 4px;
        background: #fff;
        color: #202124;
        font-family: inherit;
        font-size: 13px;
        outline: none;
        width: 100%;
        transition:
            border-color 0.15s,
            box-shadow 0.15s;
    }

    .field input.mono {
        font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
        font-size: 12px;
    }

    .field input:focus,
    .field select:focus {
        border-color: #1a73e8;
        box-shadow: inset 0 0 0 1px #1a73e8;
    }

    .id-row {
        display: flex;
        gap: 6px;
    }

    .id-row input {
        flex: 1;
        min-width: 0;
    }

    .btn-refresh {
        height: 34px;
        width: 34px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
        border: 1px solid #dadce0;
        border-radius: 4px;
        background: #fff;
        color: #5f6368;
        cursor: pointer;
        transition: background 0.15s;
    }

    .btn-refresh:hover {
        background: #f1f3f4;
    }

    .form-actions {
        display: flex;
        justify-content: flex-end;
    }

    .btn-next {
        height: 34px;
        padding: 0 20px;
        background: #1a73e8;
        color: #fff;
        border: none;
        border-radius: 4px;
        font-family: inherit;
        font-size: 13px;
        font-weight: 500;
        cursor: pointer;
        transition: background 0.15s;
    }

    .btn-next:hover:not(:disabled) {
        background: #1765cc;
    }

    .btn-next:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .divider {
        height: 1px;
        background: #e8eaed;
        margin: 16px 0;
    }

    .privacy-note {
        font-size: 12px;
        color: #5f6368;
        line-height: 1.65;
        margin: 0;
    }

    /* Footer */
    .page-footer {
        width: 100%;
        max-width: 900px;
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 16px 4px 0;
        flex-wrap: wrap;
        gap: 12px;
    }

    .lang-btn {
        display: flex;
        align-items: center;
        gap: 8px;
        background: none;
        border: none;
        font-family: inherit;
        font-size: 13px;
        color: #5f6368;
        cursor: pointer;
        padding: 6px 8px;
        border-radius: 4px;
    }

    .lang-btn:hover {
        background: rgba(0, 0, 0, 0.05);
    }

    .footer-links {
        display: flex;
        gap: 24px;
        list-style: none;
        padding: 0;
        margin: 0;
    }

    .footer-links a {
        font-size: 13px;
        color: #5f6368;
        text-decoration: none;
    }

    .footer-links a:hover {
        text-decoration: underline;
    }

    @media (max-width: 680px) {
        .card {
            border-radius: 16px;
        }

        .card-body {
            grid-template-columns: 1fr;
            padding: 32px 24px;
            gap: 24px;

            > * {
                height: auto;
            }
        }

        .left {
            padding-right: 0;
        }

        .account-list {
            overflow: none !important;
        }

        .right {
            overflow-y: scroll;
        }

        .custom-grid {
            grid-template-columns: 1fr;
        }
    }
</style>
