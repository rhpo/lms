<script lang="ts">
    import { ChevronRight, ChevronDown, RefreshCw, Plus } from "lucide-svelte";
    import { authStore } from "$lib/stores/auth";

    type Role = "admin" | "teacher" | "student" | "company";

    type Persona = {
        id: string;
        email: string;
        role: Role;
        full_name: string;
        subtitle: string;
    };

    const PERSONAS: Persona[] = [
        {
            id: "seed-admin-001",
            email: "pfe@esst-sup.com",
            role: "admin",
            full_name: "Administrateur PFE",
            subtitle: "Acces complet au systeme",
        },
        {
            id: "seed-teacher-isil-001",
            email: "prof.isil1@esst-sup.com",
            role: "teacher",
            full_name: "Prof. ISIL 1",
            subtitle: "Professeur ISIL — disponible",
        },
        {
            id: "seed-teacher-isil-002",
            email: "prof.isil2@esst-sup.com",
            role: "teacher",
            full_name: "Prof. ISIL 2",
            subtitle: "MCA ISIL — indisponible",
        },
        {
            id: "seed-teacher-isil-003",
            email: "prof.isil3@esst-sup.com",
            role: "teacher",
            full_name: "Prof. ISIL 3",
            subtitle: "MCB ISIL — indisponible jusqu'au 01/08/2026",
        },
        {
            id: "seed-teacher-chim-001",
            email: "prof.chim1@esst-sup.com",
            role: "teacher",
            full_name: "Prof. CHIM 1",
            subtitle: "MCA CHIM — disponible",
        },
        {
            id: "seed-teacher-chim-002",
            email: "prof.chim2@esst-sup.com",
            role: "teacher",
            full_name: "Prof. CHIM 2",
            subtitle: "MAA CHIM — disponible",
        },
        {
            id: "seed-teacher-chim-003",
            email: "prof.chim3@esst-sup.com",
            role: "teacher",
            full_name: "Prof. CHIM 3",
            subtitle: "MAB CHIM — disponible",
        },
        {
            id: "seed-teacher-elec-001",
            email: "prof.elec1@esst-sup.com",
            role: "teacher",
            full_name: "Prof. ELEC 1",
            subtitle: "Professeur ELEC — disponible",
        },
        {
            id: "seed-teacher-elec-002",
            email: "prof.elec2@esst-sup.com",
            role: "teacher",
            full_name: "Prof. ELEC 2",
            subtitle: "Assistant ELEC — disponible",
        },
        {
            id: "seed-teacher-elec-003",
            email: "prof.elec3@esst-sup.com",
            role: "teacher",
            full_name: "Prof. ELEC 3",
            subtitle: "MCB ELEC — disponible",
        },
        {
            id: "seed-student-isil-001",
            email: "etudiant.isil1@esst-sup.com",
            role: "student",
            full_name: "Etudiant ISIL 1",
            subtitle: "ISIL — Ingenieur",
        },
        {
            id: "seed-student-isil-002",
            email: "etudiant.isil2@esst-sup.com",
            role: "student",
            full_name: "Etudiant ISIL 2",
            subtitle: "ISIL — Ingenieur",
        },
        {
            id: "seed-student-isil-003",
            email: "etudiant.isil3@esst-sup.com",
            role: "student",
            full_name: "Etudiant ISIL 3",
            subtitle: "ISIL — Ingenieur",
        },
        {
            id: "seed-student-chim-001",
            email: "etudiant.chim1@esst-sup.com",
            role: "student",
            full_name: "Etudiant CHIM 1",
            subtitle: "CHIM — Master",
        },
        {
            id: "seed-student-chim-002",
            email: "etudiant.chim2@esst-sup.com",
            role: "student",
            full_name: "Etudiant CHIM 2",
            subtitle: "CHIM — Master",
        },
        {
            id: "seed-company-techcorp",
            email: "contact@techcorp-dz.com",
            role: "company",
            full_name: "TechCorp Algeria",
            subtitle: "Entreprise validee — secteur Technology",
        },
    ];

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
            loginError =
                err instanceof Error ? err.message : "Erreur de connexion";
        }
    }

    async function loginCustom() {
        if (!cEmail) return;
        await loginAs(cEmail);
    }
</script>

<svelte:head>
    <title>Dev Auth — PFE Manager</title>
</svelte:head>

<div class="page">
    <div class="env-strip">
        Development environment — not accessible in production
    </div>

    <div class="container">
        <header class="header">
            <div>
                <p class="header-tag">DEV AUTH</p>
                <h1 class="header-title">Select an account</h1>
                <p class="header-sub">
                    Development personas for testing the UI. Click any persona
                    to view their role and profile details.
                </p>
            </div>
        </header>

        <section class="persona-list">
            {#each PERSONAS as persona (persona.id)}
                <button
                    type="button"
                    class="persona-row"
                    data-role={persona.role}
                    onclick={() => loginAs(persona.email)}
                    disabled={authStore.loading}
                >
                    <span class="role-bar"></span>
                    <span class="persona-main">
                        <span class="persona-name">{persona.full_name}</span>
                        <span class="persona-sub">{persona.subtitle}</span>
                    </span>
                    <span class="persona-email">{persona.email}</span>
                    <span class="persona-status">
                        <span class="status-arrow">
                            <ChevronRight size={15} strokeWidth={2} />
                        </span>
                    </span>
                </button>
            {/each}
            {#if loginError}
                <div class="login-error">{loginError}</div>
            {/if}
        </section>

        <section class="custom-section">
            <button
                type="button"
                class="custom-toggle"
                class:open={showCustom}
                onclick={() => (showCustom = !showCustom)}
            >
                <Plus size={13} strokeWidth={2.5} />
                Custom account
                <span class="toggle-chevron" class:rotated={showCustom}>
                    <ChevronDown size={13} strokeWidth={2} />
                </span>
            </button>

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
                        <div class="field">
                            <label for="c-email">Email</label>
                            <input
                                id="c-email"
                                type="email"
                                placeholder="user@dev.local"
                                bind:value={cEmail}
                            />
                        </div>
                        <div class="field">
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
                                >
                                    <RefreshCw size={12} strokeWidth={2} />
                                </button>
                            </span>
                        </div>
                    </div>
                    <div class="form-submit">
                        <button
                            type="button"
                            class="btn-login"
                            onclick={loginCustom}
                            disabled={!cEmail || authStore.loading}
                        >
                            Se connecter
                        </button>
                    </div>
                </div>
            {/if}
        </section>

        <footer class="footer">
            Development auth page — no server-side sign-in required
        </footer>
    </div>
</div>

<style>
    [data-role="admin"] {
        --role: #b91c1c;
    }
    [data-role="teacher"] {
        --role: #1d4ed8;
    }
    [data-role="student"] {
        --role: #15803d;
    }
    [data-role="company"] {
        --role: #6d28d9;
    }

    .page {
        min-height: 100vh;
        background: var(--color-background);
        display: flex;
        flex-direction: column;
    }

    .env-strip {
        background: #fefce8;
        border-bottom: 1px solid #fde047;
        color: #713f12;
        font-size: 0.7rem;
        font-family: var(--font-sans);
        font-weight: 500;
        letter-spacing: 0.04em;
        text-align: center;
        padding: 0.4rem 1rem;
        text-transform: uppercase;
    }

    .container {
        max-width: 640px;
        width: 100%;
        margin: 0 auto;
        padding: 3rem 1.5rem 5rem;
    }

    .header {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        gap: 2rem;
        margin-bottom: 2.5rem;
        flex-wrap: wrap;
    }

    .header-tag {
        font-family: monospace;
        font-size: 0.65rem;
        font-weight: 700;
        letter-spacing: 0.15em;
        text-transform: uppercase;
        color: var(--color-text-muted);
        border: 1px solid var(--color-border);
        border-radius: 3px;
        padding: 0.15rem 0.45rem;
        display: inline-block;
        margin: 0 0 0.65rem;
    }

    .header-title {
        font-family: var(--font-sans);
        font-size: 1.75rem;
        font-weight: 700;
        letter-spacing: -0.025em;
        color: var(--color-text);
        margin: 0 0 0.4rem;
        line-height: 1.15;
    }

    .header-sub {
        font-family: var(--font-sans);
        font-size: 0.85rem;
        color: var(--color-text-muted);
        margin: 0;
        line-height: 1.5;
        max-width: 38ch;
    }

    .persona-list {
        border: 1px solid var(--color-border);
        border-radius: 10px;
        overflow: hidden;
        margin-bottom: 0.875rem;
    }

    .persona-list button + button {
        border-top: 1px solid var(--color-border);
    }

    .persona-row {
        width: 100%;
        display: flex;
        align-items: center;
        gap: 1rem;
        padding: 0.85rem 1rem;
        background: var(--color-background);
        border: none;
        text-align: left;
        transition: background var(--transition-fast);
        cursor: pointer;
    }

    .persona-row:hover:not(:disabled) {
        background: var(--color-background-100);
    }

    .persona-row:disabled {
        opacity: 0.6;
        cursor: not-allowed;
    }

    .login-error {
        padding: 0.75rem 1rem;
        background: color-mix(in srgb, #ef4444 10%, transparent);
        color: #ef4444;
        font-size: 0.8rem;
        font-family: var(--font-sans);
        border-top: 1px solid color-mix(in srgb, #ef4444 20%, transparent);
    }

    .form-submit {
        margin-top: 0.75rem;
        display: flex;
        justify-content: flex-end;
    }

    .btn-login {
        padding: 0.5rem 1.25rem;
        background: var(--color-accent);
        color: #fff;
        border: none;
        border-radius: 7px;
        font-family: var(--font-sans);
        font-size: 0.85rem;
        font-weight: 600;
        cursor: pointer;
        transition: opacity var(--transition-fast);
    }

    .btn-login:disabled {
        opacity: 0.5;
        cursor: not-allowed;
    }

    .role-bar {
        display: block;
        width: 3px;
        height: 1.75rem;
        border-radius: 2px;
        background: var(--role);
        flex-shrink: 0;
    }

    .persona-main {
        flex: 1;
        min-width: 0;
        display: flex;
        flex-direction: column;
        gap: 0.1rem;
    }

    .persona-name {
        font-family: var(--font-sans);
        font-size: 0.875rem;
        font-weight: 600;
        color: var(--color-text);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .persona-sub {
        font-family: var(--font-sans);
        font-size: 0.72rem;
        color: var(--color-text-muted);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .persona-email {
        font-family: monospace;
        font-size: 0.72rem;
        color: var(--color-text-muted);
        white-space: nowrap;
    }

    .persona-status {
        display: flex;
        align-items: center;
        flex-shrink: 0;
        width: 5.5rem;
        justify-content: flex-end;
    }

    .status-arrow {
        display: flex;
        align-items: center;
        color: var(--color-text-muted);
        opacity: 0;
        transition: opacity var(--transition-fast);
    }

    .persona-row:hover .status-arrow {
        opacity: 1;
    }

    .custom-section {
        border: 1px solid var(--color-border);
        border-radius: 10px;
        overflow: hidden;
        margin-bottom: 2.5rem;
    }

    .custom-toggle {
        width: 100%;
        display: flex;
        align-items: center;
        gap: 0.45rem;
        padding: 0.8rem 1rem;
        background: none;
        border: none;
        cursor: pointer;
        font-family: var(--font-sans);
        font-size: 0.82rem;
        font-weight: 500;
        color: var(--color-text-muted);
        text-align: left;
        transition:
            background var(--transition-fast),
            color var(--transition-fast);
    }

    .custom-toggle:hover,
    .custom-toggle.open {
        background: var(--color-background-100);
        color: var(--color-text);
    }

    .toggle-chevron {
        margin-left: auto;
        display: flex;
        align-items: center;
        transition: transform var(--transition-fast);
    }

    .toggle-chevron.rotated {
        transform: rotate(180deg);
    }

    .custom-form {
        border-top: 1px solid var(--color-border);
        padding: 1.25rem 1rem;
        background: var(--color-background-100);
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }

    .custom-grid {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: 0.75rem;
    }

    .field {
        display: flex;
        flex-direction: column;
        gap: 0.3rem;
    }

    .field label {
        font-family: var(--font-sans);
        font-size: 0.72rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.05em;
        color: var(--color-text-muted);
        display: flex;
        align-items: center;
        gap: 0.5rem;
    }

    .label-hint {
        font-weight: 400;
        text-transform: none;
        letter-spacing: 0;
        color: var(--color-text-disabled);
    }

    .field input,
    .field select {
        height: 2rem;
        padding: 0 0.6rem;
        border: 1px solid var(--color-border);
        border-radius: 5px;
        background: var(--color-background);
        color: var(--color-text);
        font-family: var(--font-sans);
        font-size: 0.82rem;
        outline: none;
        width: 100%;
        transition: border-color var(--transition-fast);
    }

    .field input.mono {
        font-family: monospace;
        font-size: 0.75rem;
    }

    .field input:focus,
    .field select:focus {
        border-color: var(--color-accent);
    }

    .id-row {
        display: flex;
        gap: 0.4rem;
    }

    .id-row input {
        flex: 1;
        min-width: 0;
    }

    .btn-refresh {
        height: 2rem;
        width: 2rem;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
        border: 1px solid var(--color-border);
        border-radius: 5px;
        background: var(--color-background);
        color: var(--color-text-muted);
        cursor: pointer;
        transition:
            border-color var(--transition-fast),
            color var(--transition-fast);
    }

    .btn-refresh:hover {
        border-color: var(--color-border-strong);
        color: var(--color-text);
    }

    .footer {
        display: flex;
        align-items: center;
        gap: 0.6rem;
        flex-wrap: wrap;
        font-family: monospace;
        font-size: 0.7rem;
        color: var(--color-text-disabled);
    }

    @media (max-width: 520px) {
        .header {
            flex-direction: column;
        }

        .persona-email {
            display: none;
        }

        .persona-status {
            width: auto;
        }

        .custom-grid {
            grid-template-columns: 1fr;
        }
    }
</style>
