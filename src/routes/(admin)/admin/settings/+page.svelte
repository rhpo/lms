<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import {
    Settings,
    Calendar,
    BookOpen,
    Sun,
    Moon,
    Plus,
    Pencil,
    Trash2,
    CheckCircle,
    XCircle,
  } from "lucide-svelte";
  import { admin } from "$lib/api";

  import Button from "$lib/components/ui/Button.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";
  import DateInput from "$lib/components/ui/DateInput.svelte";
  import Badge from "$lib/components/ui/Badge.svelte";
  import Page from "$lib/components/ui/Page.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";

  let { data } = $props();

  const { settings, specialities } = $derived(data);

  // === Tab state ===============
  type Tab = "Général" | "deadlines" | "specialities";
  let activeTab = $state<Tab>("Général");

  const tabs: { id: Tab; label: string; icon: any }[] = [
    { id: "Général", label: "Général", icon: Settings },
    { id: "deadlines", label: "Deadlines", icon: Calendar },
    { id: "specialities", label: "Spécialités", icon: BookOpen },
  ];

  // === Général tab =============
  let theme = $state<"light" | "dark">("light");

  // === Toast / notification state ============
  let toastMessage = $state("");
  let toastType = $state<"success" | "error">("success");
  let toastVisible = $state(false);
  let toastTimer: ReturnType<typeof setTimeout> | null = null;

  function showToast(message: string, type: "success" | "error" = "success") {
    toastMessage = message;
    toastType = type;
    toastVisible = true;
    if (toastTimer) clearTimeout(toastTimer);
    toastTimer = setTimeout(() => {
      toastVisible = false;
    }, 3000);
  }

  // === Deadlines tab ===========
  // Capture initial values for form fields (intentionally not reactive)
  let submissionOpen = $state("");
  let submissionClose = $state("");
  let maxWishes = $state(5);

  $effect(() => {
    submissionOpen = settings.submission_open_at;
    submissionClose = settings.submission_close_at;
    maxWishes = settings.max_wishes;
  });
  let deadlinesSaving = $state(false);

  async function saveDeadlines() {
    deadlinesSaving = true;
    try {
      await admin.updateDeadlines({
        submission_open_at: submissionOpen,
        submission_close_at: submissionClose,
        max_wishes: maxWishes,
      });
      showToast("Deadlines enregistrees avec succes");
      await invalidateAll();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Erreur reseau", "error");
    } finally {
      deadlinesSaving = false;
    }
  }

  // === Specialities tab ========
  let showCreateSpecModal = $state(false);
  let showEditSpecModal = $state(false);
  let showDeleteSpecModal = $state(false);

  let newSpecName = $state("");
  let newSpecCode = $state("");
  let newSpecYearType = $state<"licence" | "master" | "ingenieur">("licence");

  let editingSpec = $state<{
    id: string;
    name: string;
    code: string;
    year_type: string;
  } | null>(null);
  let editSpecName = $state("");
  let editSpecCode = $state("");
  let editSpecYearType = $state<"licence" | "master" | "ingenieur">("licence");

  let deletingSpec = $state<{ id: string; name: string } | null>(null);

  function openEditSpec(spec: (typeof data.specialities)[0]) {
    editingSpec = spec;
    editSpecName = spec.name;
    editSpecCode = spec.code;
    editSpecYearType = spec.year_type as "licence" | "master" | "ingenieur";
    showEditSpecModal = true;
  }

  function openDeleteSpec(spec: (typeof data.specialities)[0]) {
    deletingSpec = spec;
    showDeleteSpecModal = true;
  }

  async function createSpeciality() {
    try {
      await admin.createSpeciality({
        name: newSpecName,
        code: newSpecCode,
        year_type: newSpecYearType,
      });
      showCreateSpecModal = false;
      newSpecName = "";
      newSpecCode = "";
      newSpecYearType = "licence";
      showToast("Specialite creee avec succes");
      await invalidateAll();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Erreur reseau", "error");
    }
  }

  async function updateSpeciality() {
    if (!editingSpec) return;
    try {
      await admin.updateSpeciality(editingSpec.id, {
        name: editSpecName,
        code: editSpecCode,
        year_type: editSpecYearType,
      });
      showEditSpecModal = false;
      editingSpec = null;
      showToast("Specialite modifiee avec succes");
      await invalidateAll();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Erreur reseau", "error");
    }
  }

  async function deleteSpeciality() {
    if (!deletingSpec) return;
    try {
      await admin.deleteSpeciality(deletingSpec.id);
      showDeleteSpecModal = false;
      deletingSpec = null;
      showToast("Specialite supprimee avec succes");
      await invalidateAll();
    } catch (err) {
      showToast(err instanceof Error ? err.message : "Erreur reseau", "error");
    }
  }
</script>

<Page title="Paramètres" subtitle="Configuration de l'application">
  {#snippet actions()}
    <div class="tab-bar">
      {#each tabs as tab}
        <button
          class="tab-btn"
          class:active={activeTab === tab.id}
          onclick={() => (activeTab = tab.id)}
        >
          <tab.icon size={16} />
          <span>{tab.label}</span>
        </button>
      {/each}
    </div>
  {/snippet}

  <!-- Toast notification -->
  {#if toastVisible}
    <div
      class="toast"
      class:toast-success={toastType === "success"}
      class:toast-error={toastType === "error"}
    >
      {#if toastType === "success"}
        <CheckCircle size={18} />
      {:else}
        <XCircle size={18} />
      {/if}
      <span>{toastMessage}</span>
    </div>
  {/if}

  <!-- ════════════════════════════════════════════════════════════════════════ -->
  <!-- Général TAB                                                           -->
  <!-- ════════════════════════════════════════════════════════════════════════ -->
  {#if activeTab === "Général"}
    <div class="panel">
      <h2 class="panel-title">Général</h2>
      <p class="panel-desc">Paramètres Généraux de l'application</p>

      <div class="settings-group">
        <div class="setting-row">
          <div class="setting-info">
            <span class="setting-label">Langue</span>
            <span class="setting-hint">Langue de l'interface (BETA)</span>
          </div>
          <select class="setting-input" disabled>
            <option value="fr" selected>Français</option>
            <option value="en">English</option>
            <option value="ar">العربية</option>
          </select>
        </div>

        <div class="setting-row">
          <div class="setting-info">
            <span class="setting-label">Theme</span>
            <span class="setting-hint">Apparence de l'interface</span>
          </div>
          <div class="theme-toggle">
            <button
              class="theme-btn"
              class:active={theme === "light"}
              onclick={() => (theme = "light")}
            >
              <Sun size={16} />
              <span>Clair</span>
            </button>
            <button
              class="theme-btn"
              class:active={theme === "dark"}
              onclick={() => (theme = "dark")}
            >
              <Moon size={16} />
              <span>Sombre</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- ════════════════════════════════════════════════════════════════════════ -->
    <!-- DEADLINES TAB                                                          -->
    <!-- ════════════════════════════════════════════════════════════════════════ -->
  {:else if activeTab === "deadlines"}
    <div class="panel">
      <h2 class="panel-title">Deadlines</h2>
      <p class="panel-desc">
        Paramètrer les dates d'ouverture et de clôture des dépôts de sujets
      </p>

      <div class="deadlines-form">
        <FormField label="Date d'ouverture des dépôts">
          <DateInput bind:value={submissionOpen} />
        </FormField>

        <FormField label="Date de clôture des dépôts">
          <DateInput bind:value={submissionClose} />
        </FormField>

        <FormField label="Nombre maximum de voeux par etudiant">
          <input
            type="number"
            bind:value={maxWishes}
            min="1"
            max="20"
            class="input"
          />
        </FormField>

        <div class="form-actions">
          <Button
            variant="primary"
            onclick={saveDeadlines}
            disabled={deadlinesSaving}
          >
            {deadlinesSaving ? "Enregistrement..." : "Enregistrer"}
          </Button>
        </div>
      </div>
    </div>

    <!-- ════════════════════════════════════════════════════════════════════════ -->
    <!-- SPECIALITIES TAB                                                       -->
    <!-- ════════════════════════════════════════════════════════════════════════ -->
  {:else if activeTab === "specialities"}
    <div class="panel">
      <div class="panel-header">
        <div>
          <h2 class="panel-title">Spécialités</h2>
          <p class="panel-desc">
            Gérer les spécialités et leurs niveaux associes
          </p>
        </div>
        <Button
          variant="primary"
          Icon={Plus}
          onclick={() => (showCreateSpecModal = true)}
        >
          Nouvelle specialite
        </Button>
      </div>

      <table>
        <thead>
          <tr>
            <th>Code</th>
            <th>Nom</th>
            <th>Niveau</th>
            <th>Cree le</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each data.specialities as spec}
            <tr>
              <td><code>{spec.code}</code></td>
              <td>{spec.name}</td>
              <td>
                <Badge
                  variant={spec.year_type === "licence"
                    ? "info"
                    : spec.year_type === "master"
                      ? "warning"
                      : "success"}
                  label={spec.year_type_label}
                />
              </td>
              <td>{spec.created_at_formatted}</td>
              <td class="actions-cell">
                <button
                  class="icon-btn"
                  title="Modifier"
                  onclick={() => openEditSpec(spec)}
                >
                  <Pencil size={16} />
                </button>
                <button
                  class="icon-btn danger"
                  title="Supprimer"
                  onclick={() => openDeleteSpec(spec)}
                >
                  <Trash2 size={16} />
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Page>

<!-- ══════════════════════════════════════════════════════════════════════════ -->
<!-- MODALS                                                                   -->
<!-- ══════════════════════════════════════════════════════════════════════════ -->

<!-- Create Speciality Modal -->
<Modal
  open={showCreateSpecModal}
  title="Nouvelle specialite"
  onClose={() => (showCreateSpecModal = false)}
>
  <div class="modal-form">
    <FormField label="Code" required>
      <input
        type="text"
        bind:value={newSpecCode}
        placeholder="ex. ISIL"
        required
        class="input"
      />
    </FormField>

    <FormField label="Nom" required>
      <input
        type="text"
        bind:value={newSpecName}
        placeholder="ex. Ingenierie des Systemes d Information"
        required
        class="input"
      />
    </FormField>

    <FormField label="Niveau" required>
      <select bind:value={newSpecYearType} required class="input">
        <option value="licence">Licence</option>
        <option value="master">Master</option>
        <option value="ingenieur">Ingenieur</option>
      </select>
    </FormField>
  </div>

  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showCreateSpecModal = false)}
      >Annuler</Button
    >
    <Button variant="primary" onclick={createSpeciality}>Creer</Button>
  </div>
</Modal>

<!-- Edit Speciality Modal -->
<Modal
  open={showEditSpecModal}
  title="Modifier la specialite"
  onClose={() => (showEditSpecModal = false)}
>
  <div class="modal-form">
    <FormField label="Code" required>
      <input type="text" bind:value={editSpecCode} required class="input" />
    </FormField>

    <FormField label="Nom" required>
      <input type="text" bind:value={editSpecName} required class="input" />
    </FormField>

    <FormField label="Niveau" required>
      <select bind:value={editSpecYearType} required class="input">
        <option value="licence">Licence</option>
        <option value="master">Master</option>
        <option value="ingenieur">Ingenieur</option>
      </select>
    </FormField>
  </div>

  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showEditSpecModal = false)}
      >Annuler</Button
    >
    <Button variant="primary" onclick={updateSpeciality}>Enregistrer</Button>
  </div>
</Modal>

<!-- Delete Speciality Confirmation Modal -->
<Modal
  open={showDeleteSpecModal}
  title="Confirmer la suppression"
  onClose={() => (showDeleteSpecModal = false)}
>
  <p class="delete-warning">
    Etes-vous sur de vouloir supprimer la specialite <strong
      >{deletingSpec?.name}</strong
    > ? Cette action est irreversible.
  </p>

  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showDeleteSpecModal = false)}
      >Annuler</Button
    >
    <Button variant="error" onclick={deleteSpeciality}>Supprimer</Button>
  </div>
</Modal>

<style>
  /* == Toast ==== */
  .toast {
    position: fixed;
    top: 1rem;
    right: 1rem;
    z-index: 9999;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1rem;
    border-radius: 10px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    font-weight: 500;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    animation: slideIn 0.2s ease-out;
  }

  .toast-success {
    background: color-mix(in srgb, var(--color-success) 15%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-success) 30%, transparent);
    color: var(--color-success);
  }

  .toast-error {
    background: color-mix(in srgb, var(--color-danger) 15%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-danger) 30%, transparent);
    color: var(--color-danger);
  }

  @keyframes slideIn {
    from {
      transform: translateX(100%);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }

  /* == Tab bar ================ */
  .tab-bar {
    display: flex;
    gap: 0.25rem;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 10px;
    padding: 0.25rem;
  }

  .tab-btn {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    border: none;
    border-radius: 8px;
    background: transparent;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .tab-btn:hover {
    background: var(--color-background-100);
    color: var(--color-text);
  }

  .tab-btn.active {
    background: var(--color-accent);
    color: #fff;
  }

  /* == Panel === */
  .panel {
    margin-top: var(--spacing-lg);
  }

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: var(--spacing-lg);
  }

  .panel-title {
    font-size: var(--text-lg);
    font-weight: 600;
    color: var(--color-text);
    margin: 0 0 0.25rem;
  }

  .panel-desc {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    margin: 0 0 var(--spacing-lg);
  }

  /* == Settings rows ========== */
  .settings-group {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    margin-top: var(--spacing-lg);
  }

  .setting-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: var(--spacing-md);
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 10px;
  }

  .setting-info {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }

  .setting-label {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text);
  }

  .setting-hint {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .setting-input {
    padding: 0.4rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    background: var(--color-background);
    color: var(--color-text);
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    min-width: 160px;
  }

  .theme-toggle {
    display: flex;
    gap: 0.25rem;
    background: var(--color-background);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 0.2rem;
  }

  .theme-btn {
    display: flex;
    align-items: center;
    gap: 0.35rem;
    padding: 0.35rem 0.75rem;
    border: none;
    border-radius: 6px;
    background: transparent;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .theme-btn.active {
    background: var(--color-accent);
    color: #fff;
  }

  /* == Deadlines form ========= */
  .deadlines-form {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 1.75rem;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
    width: 100%;
  }

  .deadlines-form :global(.input) {
    width: 100%;
    padding: 0.6rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    background: var(--color-background);
    color: var(--color-text);
    box-sizing: border-box;
  }

  .deadlines-form :global(.input:focus) {
    outline: none;
    border-color: var(--color-accent);
    box-shadow: 0 0 0 2px
      color-mix(in srgb, var(--color-accent) 20%, transparent);
  }

  /* == Table === */
  table {
    width: 100%;
    border-collapse: collapse;
    font-family: var(--font-sans);
    font-size: var(--text-sm);
  }

  thead th {
    text-align: left;
    padding: var(--spacing-sm) var(--spacing-md);
    font-weight: 600;
    color: var(--color-text-muted);
    border-bottom: 2px solid var(--color-border);
    white-space: nowrap;
  }

  tbody tr {
    border-bottom: 1px solid var(--color-border);
  }

  tbody tr:hover {
    background: var(--color-background-100);
  }

  tbody td {
    padding: var(--spacing-sm) var(--spacing-md);
    color: var(--color-text);
  }

  code {
    font-family: var(--font-mono);
    font-size: var(--text-xs);
    background: var(--color-background-100);
    padding: 0.15rem 0.4rem;
    border-radius: 4px;
  }

  .actions-cell {
    display: flex;
    gap: 0.35rem;
  }

  .icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    border-radius: 6px;
    background: transparent;
    color: var(--color-text-muted);
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .icon-btn:hover {
    background: var(--color-background-100);
    color: var(--color-accent);
  }

  .icon-btn.danger:hover {
    background: color-mix(in srgb, var(--color-danger) 10%, transparent);
    color: var(--color-danger);
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--spacing-sm);
    margin-top: var(--spacing-lg);
  }

  .delete-warning {
    font-size: var(--text-sm);
    color: var(--color-text);
    line-height: 1.5;
  }

  .modal-form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
</style>
