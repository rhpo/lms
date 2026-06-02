<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import {
    Settings,
    Calendar,
    CalendarDays,
    BookOpen,
    Globe,
    Building2,
    Sun,
    Moon,
    Plus,
    Pencil,
    Trash2,
    Lock,
  } from "lucide-svelte";
  import { admin } from "$lib/api";
  import { atomic } from "$lib/stores/atomic.svelte";
  import { showToast } from "$lib/utils/toast";
  import type { Speciality, Domain, Department } from "$lib/types";
  import {
    YEAR_TYPE_LABELS,
    YEAR_TYPE_OPTIONS,
    YEAR_TYPE_VARIANTS,
    ACADEMIC_YEAR_STATUS_LABELS,
    formatDate,
  } from "$lib/constants/labels";

  import Button from "$lib/components/ui/Button.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";
  import DateInput from "$lib/components/ui/DateInput.svelte";
  import Badge from "$lib/components/ui/Badge.svelte";
  import Page from "$lib/components/ui/Page.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import { setTheme, theme } from "$lib/stores/theme.js";

  let { data } = $props();

  const { settings, specialities, domains, departments, academicYears } =
    $derived(data);

  type Tab =
    | "general"
    | "deadlines"
    | "academic-year"
    | "departments"
    | "specialities"
    | "domains";
  let activeTab = $state<Tab>("general");

  const tabs: { id: Tab; label: string; icon: any }[] = [
    { id: "general", label: "Général", icon: Settings },
    { id: "deadlines", label: "Deadlines", icon: Calendar },
    { id: "academic-year", label: "Année académique", icon: CalendarDays },
    { id: "departments", label: "Départements", icon: Building2 },
    { id: "specialities", label: "Spécialités", icon: BookOpen },
    { id: "domains", label: "Domaines", icon: Globe },
  ];

  let submissionOpen = $state("");
  let submissionClose = $state("");
  let maxWishes = $state(5);
  let deadlinesSaving = $state(false);

  $effect(() => {
    if (settings) {
      submissionOpen = settings.submission_open_at ?? "";
      submissionClose = settings.submission_close_at ?? "";
      maxWishes = settings.max_wishes ?? 5;
    }
  });

  async function saveDeadlines() {
    if (
      submissionOpen &&
      submissionClose &&
      new Date(submissionClose) <= new Date(submissionOpen)
    ) {
      showToast.error(
        "La date de clôture doit être après la date d'ouverture.",
      );
      return;
    }
    if (maxWishes < 1 || maxWishes > 20) {
      showToast.error("Le nombre maximum de vœux doit être entre 1 et 20.");
      return;
    }
    deadlinesSaving = true;
    try {
      await admin.updateDeadlines({
        submission_open_at: submissionOpen,
        submission_close_at: submissionClose,
        max_wishes: maxWishes,
      });
      showToast.success("Deadlines enregistrées avec succès");
      await invalidateAll();
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    } finally {
      deadlinesSaving = false;
    }
  }

  let showCreateDeptModal = $state(false);
  let showDeleteDeptModal = $state(false);
  let newDeptName = $state("");
  let deletingDept = $state<Department | null>(null);

  async function createDepartment() {
    if (newDeptName.trim().length < 2) {
      showToast.error(
        "Le nom du département doit contenir au moins 2 caractères.",
      );
      return;
    }
    try {
      await admin.createDepartment({ name: newDeptName.trim() });
      showCreateDeptModal = false;
      newDeptName = "";
      showToast.success("Département créé avec succès");
      await Promise.all([invalidateAll(), atomic.reload()]);
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }

  async function deleteDepartment() {
    if (!deletingDept) return;
    try {
      await admin.deleteDepartment(deletingDept.id);
      showDeleteDeptModal = false;
      deletingDept = null;
      showToast.success("Département supprimé avec succès");
      await Promise.all([invalidateAll(), atomic.reload()]);
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }

  let showCreateSpecModal = $state(false);
  let showDeleteSpecModal = $state(false);
  let newSpecName = $state("");
  let newSpecCode = $state("");
  let newSpecYearType = $state<"licence" | "master">("licence");
  let newSpecDeptId = $state(0);
  let deletingSpec = $state<Speciality | null>(null);

  async function createSpeciality() {
    if (newSpecName.trim().length < 2) {
      showToast.error(
        "Le nom de la spécialité doit contenir au moins 2 caractères.",
      );
      return;
    }
    if (newSpecCode.trim().length < 2) {
      showToast.error(
        "Le code de la spécialité doit contenir au moins 2 caractères.",
      );
      return;
    }
    try {
      await admin.createSpeciality({
        name: newSpecName.trim(),
        code: newSpecCode.trim(),
        year_type: newSpecYearType,
        department_id: newSpecDeptId || undefined,
      });
      showCreateSpecModal = false;
      newSpecName = "";
      newSpecCode = "";
      newSpecYearType = "licence";
      newSpecDeptId = 0;
      showToast.success("Spécialité créée avec succès");
      await Promise.all([invalidateAll(), atomic.reload()]);
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }

  async function deleteSpeciality() {
    if (!deletingSpec) return;
    try {
      await admin.deleteSpeciality(deletingSpec.id);
      showDeleteSpecModal = false;
      deletingSpec = null;
      showToast.success("Spécialité supprimée avec succès");
      await Promise.all([invalidateAll(), atomic.reload()]);
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }

  let showCreateDomainModal = $state(false);
  let showDeleteDomainModal = $state(false);
  let newDomainName = $state("");
  let deletingDomain = $state<Domain | null>(null);

  async function createDomain() {
    if (newDomainName.trim().length < 2) {
      showToast.error("Le nom du domaine doit contenir au moins 2 caractères.");
      return;
    }
    try {
      await admin.createDomain({ name: newDomainName.trim() });
      showCreateDomainModal = false;
      newDomainName = "";
      showToast.success("Domaine créé avec succès");
      await Promise.all([invalidateAll(), atomic.reload()]);
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }

  async function deleteDomain() {
    if (!deletingDomain) return;
    try {
      await admin.deleteDomain(deletingDomain.id);
      showDeleteDomainModal = false;
      deletingDomain = null;
      showToast.success("Domaine supprimé avec succès");
      await Promise.all([invalidateAll(), atomic.reload()]);
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }

  let showCreateYearModal = $state(false);
  let showCloseYearModal = $state(false);
  let newYearLabel = $state("");
  let newYearMaxWishes = $state(5);
  let closingYearId = $state<number | null>(null);
  let closingYearLabel = $state("");
  let yearSaving = $state(false);

  async function createAcademicYear() {
    if (newYearLabel.trim().length < 4) {
      showToast.error(
        "Le label de l'année doit contenir au moins 4 caractères (ex: 2024-2025).",
      );
      return;
    }
    if (newYearMaxWishes < 1 || newYearMaxWishes > 20) {
      showToast.error("Le nombre maximum de vœux doit être entre 1 et 20.");
      return;
    }
    yearSaving = true;
    try {
      await admin.createAcademicYear({
        label: newYearLabel.trim(),
        status: "active",
        max_wishes: newYearMaxWishes,
      });
      showCreateYearModal = false;
      newYearLabel = "";
      newYearMaxWishes = 5;
      showToast.success("Année académique créée avec succès");
      await Promise.all([invalidateAll(), atomic.reload()]);
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    } finally {
      yearSaving = false;
    }
  }

  async function closeAcademicYear() {
    if (!closingYearId) return;
    yearSaving = true;
    try {
      await admin.closeAcademicYear(closingYearId);
      showCloseYearModal = false;
      closingYearId = null;
      closingYearLabel = "";
      showToast.success("Année académique clôturée avec succès");
      await Promise.all([invalidateAll(), atomic.reload()]);
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    } finally {
      yearSaving = false;
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

  <!-- General -->
  {#if activeTab === "general"}
    <div class="panel">
      <h2 class="panel-title">Général</h2>
      <p class="panel-desc">Paramètres généraux de l'application</p>

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
            <span class="setting-label">Thème</span>
            <span class="setting-hint">Apparence de l'interface</span>
          </div>
          <div class="theme-toggle">
            <button
              class="theme-btn"
              class:active={$theme === "light"}
              onclick={() => setTheme("light")}
            >
              <Sun size={16} />
              <span>Clair</span>
            </button>
            <button
              class="theme-btn"
              class:active={$theme === "dark"}
              onclick={() => setTheme("dark")}
            >
              <Moon size={16} />
              <span>Sombre</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Deadlines -->
  {:else if activeTab === "deadlines"}
    <div class="panel">
      <h2 class="panel-title">Deadlines</h2>
      <p class="panel-desc">
        Paramétrer les dates d'ouverture et de clôture des dépôts de sujets
      </p>

      <div class="deadlines-form">
        <FormField label="Date d'ouverture des dépôts">
          <DateInput bind:value={submissionOpen} />
        </FormField>

        <FormField label="Date de clôture des dépôts">
          <DateInput bind:value={submissionClose} />
        </FormField>

        <FormField label="Nombre maximum de vœux par étudiant">
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

    <!-- Academic Year -->
  {:else if activeTab === "academic-year"}
    <div class="panel">
      <div class="panel-header">
        <div>
          <h2 class="panel-title">Année académique</h2>
          <p class="panel-desc">
            Gérer les années académiques. Une seule année peut être active à la
            fois.
          </p>
        </div>
        <Button
          variant="primary"
          Icon={Plus}
          onclick={() => (showCreateYearModal = true)}
        >
          Nouvelle année
        </Button>
      </div>

      {#if academicYears.length === 0}
        <div class="year-empty">
          <CalendarDays size={40} />
          <p>Aucune année académique configurée.</p>
          <p class="year-empty-hint">
            Créez une année académique pour commencer à recevoir des vœux des
            étudiants.
          </p>
        </div>
      {:else}
        <table>
          <thead>
            <tr>
              <th>Label</th>
              <th>Statut</th>
              <th>Max vœux</th>
              <th>Ouverture dépôts</th>
              <th>Clôture dépôts</th>
              <th>Créé le</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {#each academicYears as year}
              <tr>
                <td><strong>{year.label}</strong></td>
                <td>
                  <Badge
                    variant={year.status === "active" ? "success" : "default"}
                    label={ACADEMIC_YEAR_STATUS_LABELS[year.status] ??
                      year.status}
                  />
                </td>
                <td>{year.max_wishes}</td>
                <td
                  >{year.submission_open_at
                    ? formatDate(year.submission_open_at)
                    : "-"}</td
                >
                <td
                  >{year.submission_close_at
                    ? formatDate(year.submission_close_at)
                    : "-"}</td
                >
                <td>{formatDate(year.created_at)}</td>
                <td class="actions-cell">
                  {#if year.status === "active"}
                    <button
                      class="icon-btn danger"
                      title="Clôturer"
                      onclick={() => {
                        closingYearId = year.id;
                        closingYearLabel = year.label;
                        showCloseYearModal = true;
                      }}
                    >
                      <Lock size={16} />
                    </button>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    </div>

    <!-- Departments -->
  {:else if activeTab === "departments"}
    <div class="panel">
      <div class="panel-header">
        <div>
          <h2 class="panel-title">Départements</h2>
          <p class="panel-desc">Gérer les départements de l'établissement</p>
        </div>
        <Button
          variant="primary"
          Icon={Plus}
          onclick={() => (showCreateDeptModal = true)}
        >
          Nouveau département
        </Button>
      </div>

      <table>
        <thead>
          <tr>
            <th>Nom</th>
            <th>Spécialités</th>
            <th>Créé le</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each departments as dept}
            {@const deptSpecs = specialities.filter(
              (s) => s.department_id === dept.id,
            )}
            <tr>
              <td><strong>{dept.name}</strong></td>
              <td>
                {#if deptSpecs.length > 0}
                  <div class="spec-badges">
                    {#each deptSpecs as spec}
                      <Badge
                        variant={YEAR_TYPE_VARIANTS[spec.year_type] ?? "info"}
                        label="{spec.code} ({YEAR_TYPE_LABELS[spec.year_type] ??
                          spec.year_type})"
                      />
                    {/each}
                  </div>
                {:else}
                  <span class="text-muted">Aucune</span>
                {/if}
              </td>
              <td>{formatDate(dept.created_at)}</td>
              <td class="actions-cell">
                <button
                  class="icon-btn danger"
                  title="Supprimer"
                  onclick={() => {
                    deletingDept = dept;
                    showDeleteDeptModal = true;
                  }}
                >
                  <Trash2 size={16} />
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <!-- Specialities -->
  {:else if activeTab === "specialities"}
    <div class="panel">
      <div class="panel-header">
        <div>
          <h2 class="panel-title">Spécialités</h2>
          <p class="panel-desc">
            Gérer les spécialités et leurs niveaux associés
          </p>
        </div>
        <Button
          variant="primary"
          Icon={Plus}
          onclick={() => (showCreateSpecModal = true)}
        >
          Nouvelle spécialité
        </Button>
      </div>

      <table>
        <thead>
          <tr>
            <th>Code</th>
            <th>Nom</th>
            <th>Département</th>
            <th>Niveau</th>
            <th>Créé le</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each specialities as spec}
            <tr>
              <td><code>{spec.code}</code></td>
              <td>{spec.name}</td>
              <td
                >{departments.find((d) => d.id === spec.department_id)?.name ??
                  "-"}</td
              >
              <td>
                <Badge
                  variant={YEAR_TYPE_VARIANTS[spec.year_type] ?? "info"}
                  label={YEAR_TYPE_LABELS[spec.year_type] ?? spec.year_type}
                />
              </td>
              <td>{formatDate(spec.created_at)}</td>
              <td class="actions-cell">
                <button
                  class="icon-btn danger"
                  title="Supprimer"
                  onclick={() => {
                    deletingSpec = spec;
                    showDeleteSpecModal = true;
                  }}
                >
                  <Trash2 size={16} />
                </button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>

    <!-- Domains -->
  {:else if activeTab === "domains"}
    <div class="panel">
      <div class="panel-header">
        <div>
          <h2 class="panel-title">Domaines</h2>
          <p class="panel-desc">Gérer les domaines de recherche</p>
        </div>
        <Button
          variant="primary"
          Icon={Plus}
          onclick={() => (showCreateDomainModal = true)}
        >
          Nouveau domaine
        </Button>
      </div>

      <table>
        <thead>
          <tr>
            <th>Nom</th>
            <th>Créé le</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each domains as domain}
            <tr>
              <td>{domain.name}</td>
              <td>{formatDate(domain.created_at)}</td>
              <td class="actions-cell">
                <button
                  class="icon-btn danger"
                  title="Supprimer"
                  onclick={() => {
                    deletingDomain = domain;
                    showDeleteDomainModal = true;
                  }}
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

<!-- Create Speciality Modal -->
<Modal
  open={showCreateSpecModal}
  title="Nouvelle spécialité"
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
        placeholder="ex. Ingénierie des Systèmes d'Information"
        required
        class="input"
      />
    </FormField>
    <FormField label="Département" required>
      <select
        required
        class="input"
        value={newSpecDeptId}
        onchange={(e) => (newSpecDeptId = Number(e.currentTarget.value))}
      >
        <option value={0}>- Choisir -</option>
        {#each departments as dept}
          <option value={dept.id}>{dept.name}</option>
        {/each}
      </select>
    </FormField>
    <FormField label="Niveau" required>
      <select bind:value={newSpecYearType} required class="input">
        {#each YEAR_TYPE_OPTIONS as opt}
          <option value={opt.value}>{opt.label}</option>
        {/each}
      </select>
    </FormField>
  </div>
  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showCreateSpecModal = false)}
      >Annuler</Button
    >
    <Button variant="primary" onclick={createSpeciality}>Créer</Button>
  </div>
</Modal>

<!-- Delete Speciality Confirmation -->
<Modal
  open={showDeleteSpecModal}
  title="Confirmer la suppression"
  onClose={() => (showDeleteSpecModal = false)}
>
  <p class="delete-warning">
    Êtes-vous sûr de vouloir supprimer la spécialité <strong
      >{deletingSpec?.name}</strong
    > ? Cette action est irréversible.
  </p>
  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showDeleteSpecModal = false)}
      >Annuler</Button
    >
    <Button variant="error" onclick={deleteSpeciality}>Supprimer</Button>
  </div>
</Modal>

<!-- Create Domain Modal -->
<Modal
  open={showCreateDomainModal}
  title="Nouveau domaine"
  onClose={() => (showCreateDomainModal = false)}
>
  <div class="modal-form">
    <FormField label="Nom" required>
      <input
        type="text"
        bind:value={newDomainName}
        placeholder="ex. Intelligence Artificielle"
        required
        class="input"
      />
    </FormField>
  </div>
  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showCreateDomainModal = false)}
      >Annuler</Button
    >
    <Button variant="primary" onclick={createDomain}>Créer</Button>
  </div>
</Modal>

<!-- Create Department Modal -->
<Modal
  open={showCreateDeptModal}
  title="Nouveau département"
  onClose={() => (showCreateDeptModal = false)}
>
  <div class="modal-form">
    <FormField label="Nom" required>
      <input
        type="text"
        bind:value={newDeptName}
        placeholder="ex. Informatique"
        required
        class="input"
      />
    </FormField>
  </div>
  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showCreateDeptModal = false)}
      >Annuler</Button
    >
    <Button variant="primary" onclick={createDepartment}>Créer</Button>
  </div>
</Modal>

<!-- Delete Department Confirmation -->
<Modal
  open={showDeleteDeptModal}
  title="Confirmer la suppression"
  onClose={() => (showDeleteDeptModal = false)}
>
  <p class="delete-warning">
    Êtes-vous sûr de vouloir supprimer le département <strong
      >{deletingDept?.name}</strong
    > ? Cette action est irréversible.
  </p>
  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showDeleteDeptModal = false)}
      >Annuler</Button
    >
    <Button variant="error" onclick={deleteDepartment}>Supprimer</Button>
  </div>
</Modal>

<!-- Delete Domain Confirmation -->
<Modal
  open={showDeleteDomainModal}
  title="Confirmer la suppression"
  onClose={() => (showDeleteDomainModal = false)}
>
  <p class="delete-warning">
    Êtes-vous sûr de vouloir supprimer le domaine <strong
      >{deletingDomain?.name}</strong
    > ? Cette action est irréversible.
  </p>
  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showDeleteDomainModal = false)}
      >Annuler</Button
    >
    <Button variant="error" onclick={deleteDomain}>Supprimer</Button>
  </div>
</Modal>

<!-- Create Academic Year Modal -->
<Modal
  open={showCreateYearModal}
  title="Nouvelle année académique"
  onClose={() => (showCreateYearModal = false)}
>
  <div class="modal-form">
    <FormField label="Label" required>
      <input
        type="text"
        bind:value={newYearLabel}
        placeholder="ex. 2024-2025"
        required
        class="input"
      />
    </FormField>
    <FormField label="Nombre maximum de vœux par étudiant">
      <input
        type="number"
        bind:value={newYearMaxWishes}
        min="1"
        max="20"
        class="input"
      />
    </FormField>
    <p class="year-modal-note">
      L'année sera créée avec le statut <strong>Active</strong>. Toute autre
      année déjà active devra être clôturée d'abord.
    </p>
  </div>
  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showCreateYearModal = false)}
      >Annuler</Button
    >
    <Button
      variant="primary"
      onclick={createAcademicYear}
      disabled={yearSaving || !newYearLabel.trim()}
    >
      {yearSaving ? "Création..." : "Créer"}
    </Button>
  </div>
</Modal>

<!-- Close Academic Year Confirmation -->
<Modal
  open={showCloseYearModal}
  title="Clôturer l'année académique"
  onClose={() => (showCloseYearModal = false)}
>
  <p class="delete-warning">
    Êtes-vous sûr de vouloir clôturer l'année <strong>{closingYearLabel}</strong
    > ? Les étudiants ne pourront plus soumettre de vœux et aucune nouvelle affectation
    ne sera possible.
  </p>
  <div class="form-actions">
    <Button variant="ghost" onclick={() => (showCloseYearModal = false)}
      >Annuler</Button
    >
    <Button variant="error" onclick={closeAcademicYear} disabled={yearSaving}>
      {yearSaving ? "Clôture..." : "Clôturer"}
    </Button>
  </div>
</Modal>

<style>
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

  .spec-badges {
    display: flex;
    flex-wrap: wrap;
    gap: 0.25rem;
  }

  .text-muted {
    color: var(--color-text-muted);
    font-size: var(--text-xs);
  }

  .year-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    color: var(--color-text-muted);
    text-align: center;
    gap: 0.75rem;
    border: 1px dashed var(--color-border);
    border-radius: 12px;

    & p {
      margin: 0;
      font-size: var(--text-sm);
      font-weight: 500;
    }
  }

  .year-empty-hint {
    font-weight: 400 !important;
    font-size: var(--text-xs) !important;
    max-width: 360px;
  }

  .year-modal-note {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    background: var(--color-background-100);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    padding: 0.75rem 1rem;
    margin: 0;
    line-height: 1.5;
  }
</style>
