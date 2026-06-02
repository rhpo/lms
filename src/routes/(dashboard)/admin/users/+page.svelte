<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { admin } from "$lib/api";
  import { showToast } from "$lib/utils/toast";
  import { atomic } from "$lib/stores/atomic.svelte";
  import type { Profile, Company } from "$lib/types";
  import {
    TEACHER_GRADE_LABELS,
    TEACHER_GRADE_OPTIONS,
    YEAR_TYPE_OPTIONS,
    AVAILABILITY_LABELS,
    AVAILABILITY_VARIANTS,
  } from "$lib/constants/labels";
  import {
    Users,
    BookOpen,
    Building2,
    Plus,
    Upload,
    Check,
    X,
    Search,
    Pen,
    Ban,
  } from "lucide-svelte";

  import Button from "$lib/components/ui/Button.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import Badge from "$lib/components/ui/Badge.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";
  import DashboardPage from "$lib/components/admin/DashboardPage.svelte";
  import Avatar from "$lib/components/ui/Avatar.svelte";

  let { data } = $props();

  const profiles: Profile[] = $derived(data.profiles);
  const companies: Company[] = $derived(data.companies);
  const specialities = $derived(atomic.specialities);

  const teachers = $derived(profiles.filter((p) => p.role === "teacher"));
  const students = $derived(profiles.filter((p) => p.role === "student"));

  type TabId = "teachers" | "students" | "companies";
  let activeTab: TabId = $state("teachers");

  let showAddTeacher = $state(false);
  let showAddStudent = $state(false);
  let showImportCsv = $state(false);
  let searchQuery = $state("");

  let teacherFullName = $state("");
  let teacherEmail = $state("");
  let teacherGrade = $state("assistant");

  let studentFullName = $state("");
  let studentEmail = $state("");

  let csvFile: File | null = $state(null);
  let csvType = $state("teachers");
  let csvReplaceMode = $state(false);

  function filterBySearch<T>(items: T[], fields: (keyof T)[]): T[] {
    if (!searchQuery) return items;
    const q = searchQuery.toLowerCase();
    return items.filter((item) =>
      fields.some((f) => {
        const val = item[f];
        return val != null && String(val).toLowerCase().includes(q);
      }),
    );
  }

  let filteredTeachers = $derived.by(() => {
    let filtered = teachers;

    let removedDisabled = filtered.filter((e) => e.is_active);

    return removedDisabled;
  });

  let filteredStudents = $derived.by(() => {
    const q = searchQuery.toLowerCase();
    return students.filter((s) => {
      if (!q) return true;
      return (
        s.full_name.toLowerCase().includes(q) ||
        s.email.toLowerCase().includes(q) ||
        (s.student?.student_number ?? "").toLowerCase().includes(q) ||
        (s.student?.speciality?.name ?? "").toLowerCase().includes(q) ||
        (s.student?.speciality?.code ?? "").toLowerCase().includes(q) ||
        (s.student?.level ?? "").toLowerCase().includes(q) ||
        (s.student?.promotion?.label ?? "").toLowerCase().includes(q)
      );
    });
  });

  let filteredCompanies = $derived(
    filterBySearch(companies, [
      "company_name",
      "contact_phone",
      "contact_email",
    ]),
  );

  async function addTeacher() {
    if (!teacherEmail || !teacherFullName) return;
    try {
      await admin.createUser({
        role: "teacher",
        full_name: teacherFullName,
        email: teacherEmail,
      });
      showAddTeacher = false;
      teacherFullName = "";
      teacherEmail = "";
      teacherGrade = "assistant";
      showToast.success("Enseignant ajouté avec succès");
      await invalidateAll();
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }

  async function addStudent() {
    if (!studentEmail || !studentFullName) return;
    try {
      await admin.createUser({
        role: "student",
        full_name: studentFullName,
        email: studentEmail,
      });
      showAddStudent = false;
      studentFullName = "";
      studentEmail = "";
      showToast.success("Étudiant ajouté avec succès");
      await invalidateAll();
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }

  async function importCsv() {
    if (!csvFile) return;
    try {
      const text = await csvFile.text();
      await admin.importUsersCSV(text, csvType, csvReplaceMode);
      showImportCsv = false;
      csvFile = null;
      csvReplaceMode = false;
      showToast.success("Import effectué avec succès");
      await invalidateAll();
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }

  async function deactivateUser(id: number) {
    try {
      await admin.userAction(id, "deactivate");
      showToast.success("Utilisateur désactivé");
      await invalidateAll();
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }

  async function activateUser(id: number) {
    try {
      await admin.userAction(id, "activate");
      showToast.success("Utilisateur activé");
      await invalidateAll();
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }

  async function verifyCompany(companyId: number) {
    try {
      await admin.companyAction(companyId, "validate");
      showToast.success("Entreprise validée");
      await invalidateAll();
    } catch (err) {
      showToast.error(err instanceof Error ? err.message : "Erreur réseau");
    }
  }
</script>

<DashboardPage
  title="Gestion des utilisateurs"
  subtitle="Gérer les comptes enseignants, étudiants et entreprises."
>
  {#snippet actions()}
    <Button
      variant="ghost"
      Icon={Upload}
      onclick={() => (showImportCsv = true)}
    >
      Importer CSV
    </Button>
  {/snippet}

  <div class="tabs">
    <button
      class:active={activeTab === "teachers"}
      onclick={() => (activeTab = "teachers")}
    >
      <Users size={16} />
      <span>Enseignants</span>
      <span class="tab-count">{teachers.length}</span>
    </button>
    <button
      class:active={activeTab === "students"}
      onclick={() => (activeTab = "students")}
    >
      <BookOpen size={16} />
      <span>Étudiants</span>
      <span class="tab-count">{students.length}</span>
    </button>
    <button
      class:active={activeTab === "companies"}
      onclick={() => (activeTab = "companies")}
    >
      <Building2 size={16} />
      <span>Entreprises</span>
      <span class="tab-count">{companies.length}</span>
    </button>
  </div>

  <div class="toolbar">
    <div class="search-wrap">
      <Search size={16} />
      <input type="text" placeholder="Rechercher..." bind:value={searchQuery} />
    </div>

    {#if activeTab === "teachers"}
      <Button
        variant="primary"
        Icon={Plus}
        href="/admin/users/add?type=teacher"
      >
        Ajouter un enseignant
      </Button>
    {:else if activeTab === "students"}
      <Button
        variant="primary"
        Icon={Plus}
        href="/admin/users/add?type=student"
      >
        Ajouter un étudiant
      </Button>
    {/if}
  </div>

  {#if activeTab === "teachers"}
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Nom</th>
            <th>Email</th>
            <th>Grade</th>
            <th>Département</th>
            <th>Statut</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if filteredTeachers.length === 0}
            <tr>
              <td colspan="4" class="empty">Aucun enseignant trouvé.</td>
            </tr>
          {:else}
            {#each filteredTeachers as t}
              <tr>
                <td class="cell-name">
                  <div class="avatar-cell">
                    <Avatar size={36} user={t} />
                    {t.full_name}
                  </div>
                </td>
                <td class="cell-email">{t.email}</td>
                <td class="cell-email">{t.email}</td>
                <td class="cell-email">
                  {t.teacher?.department ? t.teacher?.department?.name : "N/A"}
                </td>
                <td>
                  <Badge
                    variant={t.is_active ? "success" : "danger"}
                    label={t.is_active ? "Actif" : "Inactif"}
                  />
                </td>
                <td class="cell-actions">
                  <Button
                    Icon={Pen}
                    variant="primary"
                    href="/admin/users/{t.id}">Modifier</Button
                  >
                  {#if t.is_active}
                    <Button
                      Icon={Ban}
                      variant="error"
                      onclick={() => deactivateUser(t.id)}
                    >
                      Supprimer
                    </Button>
                  {:else}
                    <Button
                      Icon={Check}
                      variant="success"
                      onclick={() => activateUser(t.id)}
                    >
                      Activer
                    </Button>
                  {/if}
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  {/if}

  {#if activeTab === "students"}
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Nom</th>
            <th>N° Étudiant</th>
            <th>Email</th>
            <th>Spécialité</th>
            <th>Niveau</th>
            <th>Promotion</th>
            <th>Statut</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if filteredStudents.length === 0}
            <tr>
              <td colspan="8" class="empty">Aucun étudiant trouvé.</td>
            </tr>
          {:else}
            {#each filteredStudents as s}
              <tr>
                <td class="cell-name">
                  <div class="avatar-cell">
                    <Avatar size={32} user={s} />
                    {s.full_name}
                  </div>
                </td>
                <td class="cell-mono">{s.student?.student_number ?? "-"}</td>
                <td class="cell-email">{s.email}</td>
                <td>
                  {#if s.student?.speciality}
                    <span class="spec-badge">{s.student.speciality.code}</span>
                    <span class="spec-name">{s.student.speciality.name}</span>
                  {:else}
                    <span class="cell-muted">-</span>
                  {/if}
                </td>
                <td class="cell-level">{s.student?.level ?? "-"}</td>
                <td class="cell-muted">{s.student?.promotion?.label ?? "-"}</td>
                <td>
                  <Badge
                    variant={s.is_active ? "success" : "danger"}
                    label={s.is_active ? "Actif" : "Inactif"}
                  />
                </td>
                <td class="cell-actions">
                  <Button
                    Icon={Pen}
                    variant="primary"
                    href="/admin/users/{s.id}">Modifier</Button
                  >
                  {#if s.is_active}
                    <Button
                      Icon={Ban}
                      variant="error"
                      onclick={() => deactivateUser(s.id)}
                    >
                      Désactiver
                    </Button>
                  {:else}
                    <Button
                      Icon={Check}
                      variant="success"
                      onclick={() => activateUser(s.id)}
                    >
                      Activer
                    </Button>
                  {/if}
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  {/if}

  {#if activeTab === "companies"}
    <div class="table-wrap">
      <table>
        <thead>
          <tr>
            <th>Raison sociale</th>
            <th>Secteur</th>
            <th>Téléphone</th>
            <th>Site web</th>
            <th>Statut</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if filteredCompanies.length === 0}
            <tr>
              <td colspan="6" class="empty">Aucune entreprise trouvée.</td>
            </tr>
          {:else}
            {#each filteredCompanies as c}
              <tr>
                <td class="cell-name">{c.company_name ?? "-"}</td>
                <td>{c.sector ?? "-"}</td>
                <td>{c.contact_phone ?? "-"}</td>
                <td>
                  {#if c.website}
                    <a
                      href={c.website}
                      target="_blank"
                      rel="noopener"
                      class="action-link"
                    >
                      {c.website}
                    </a>
                  {:else}
                    -
                  {/if}
                </td>
                <td>
                  <Badge
                    variant={c.is_verified ? "success" : "warning"}
                    label={c.is_verified ? "Vérifiée" : "En attente"}
                  />
                </td>
                <td class="cell-actions">
                  {#if !c.is_verified}
                    <button
                      class="action-link success"
                      onclick={() => verifyCompany(c.id)}
                    >
                      <Check size={14} /> Valider
                    </button>
                  {/if}
                  <button
                    class="action-link danger"
                    onclick={() => deactivateUser(c.profile_id)}
                  >
                    <X size={14} /> Désactiver
                  </button>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  {/if}

  <!-- Add Teacher Modal -->
  <Modal bind:open={showAddTeacher} title="Ajouter un enseignant" width="520px">
    <div class="form-grid">
      <FormField label="Nom complet" required>
        <input
          type="text"
          bind:value={teacherFullName}
          required
          class="input"
        />
      </FormField>
      <FormField label="Email" required>
        <input type="email" bind:value={teacherEmail} required class="input" />
      </FormField>
      <FormField label="Grade" required>
        <select bind:value={teacherGrade} required class="input">
          <option value="assistant">Assistant</option>
          <option value="maître_assistant">Maître Assistant</option>
          <option value="professeur">Professeur</option>
          <option value="professeur_assistant">Professeur Assistant</option>
        </select>
      </FormField>
    </div>
    <div class="form-actions">
      <Button variant="ghost" onclick={() => (showAddTeacher = false)}
        >Annuler</Button
      >
      <Button variant="primary" onclick={addTeacher}>Ajouter</Button>
    </div>
  </Modal>

  <!-- Add Student Modal -->
  <Modal bind:open={showAddStudent} title="Ajouter un étudiant" width="520px">
    <div class="form-grid">
      <FormField label="Nom complet" required>
        <input
          type="text"
          bind:value={studentFullName}
          required
          class="input"
        />
      </FormField>
      <FormField label="Email" required>
        <input type="email" bind:value={studentEmail} required class="input" />
      </FormField>
    </div>
    <div class="form-actions">
      <Button variant="ghost" onclick={() => (showAddStudent = false)}
        >Annuler</Button
      >
      <Button variant="primary" onclick={addStudent}>Ajouter</Button>
    </div>
  </Modal>

  <!-- Import CSV Modal -->
  <Modal
    bind:open={showImportCsv}
    title="Importer des comptes (CSV)"
    width="520px"
  >
    <div class="form-grid">
      <FormField label="Type de comptes" required>
        <select bind:value={csvType} required class="input">
          <option value="teachers">Enseignants</option>
          <option value="students">Étudiants</option>
        </select>
      </FormField>
      <FormField label="Fichier CSV" required>
        <input
          type="file"
          accept=".csv"
          required
          class="input"
          onchange={(e) => {
            const target = e.currentTarget as HTMLInputElement;
            csvFile = target.files?.[0] ?? null;
          }}
        />
      </FormField>
      <label class="checkbox-field">
        <input type="checkbox" bind:checked={csvReplaceMode} />
        <span
          >Remplacer les utilisateurs existants (sinon, ajout uniquement)</span
        >
      </label>
    </div>
    <p class="csv-hint">
      {#if csvType === "teachers"}
        Format : nom_complet, email, grade, departement, specialite (optionnel),
        domaines (optionnel, séparés par ;)
      {:else}
        Format : nom_complet, email, numero_etudiant, specialite, niveau,
        promotion
      {/if}
    </p>
    <p class="csv-hint">
      <a
        href="/examples/{csvType === 'teachers'
          ? 'enseignants'
          : 'etudiants'}.csv"
        download
        class="csv-download-link"
      >
        Télécharger le fichier exemple
      </a>
    </p>
    <div class="form-actions">
      <Button variant="ghost" onclick={() => (showImportCsv = false)}
        >Annuler</Button
      >
      <Button variant="primary" onclick={importCsv}>Importer</Button>
    </div>
  </Modal>
</DashboardPage>

<style>
  .tabs {
    display: flex;
    gap: 0.25rem;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 10px;
    padding: 0.25rem;
    overflow-x: auto;
  }

  button {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.6rem 1.25rem;
    border: none;
    border-radius: 8px;
    background: transparent;
    color: var(--color-text-muted);
    font-family: var(--font-sans);
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    white-space: nowrap;
    transition: all var(--transition-fast);
  }

  button.active {
    background: var(--color-accent);
    color: var(--color-background);
  }

  button:hover:not(.active) {
    background: var(--color-background-100);
    color: var(--color-text);
  }

  .tab-count {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 1.5rem;
    height: 1.5rem;
    padding: 0 0.35rem;
    border-radius: 999px;
    background: color-mix(in srgb, currentColor 15%, transparent);
    font-size: 0.75rem;
    font-weight: 600;
  }

  .toolbar {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    flex-wrap: wrap;
  }

  .search-wrap {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    flex: 1;
    min-width: 200px;
    padding: 0.5rem 0.75rem;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 8px;
  }

  .search-wrap input {
    flex: 1;
    border: none;
    background: transparent;
    color: var(--color-text);
    font-family: var(--font-sans);
    font-size: 0.9rem;
    outline: none;
  }

  .search-wrap input::placeholder {
    color: var(--color-text-muted);
  }

  .table-wrap {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    overflow: visible;
    box-shadow: var(--shadow-sm);
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  thead tr {
    background: var(--color-background);
  }

  th {
    font-size: var(--text-xs);
    font-weight: 600;
    font-family: var(--font-sans);
    color: var(--color-text-muted);
    padding: 0.75rem 1rem;
    white-space: nowrap;
    letter-spacing: 0.03em;
    text-transform: uppercase;
    text-align: left;
  }

  td {
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    color: var(--color-text);
    padding: 0.85rem 1rem;
    overflow: hidden;
    border-top: 1px solid var(--color-border);
    vertical-align: middle;
  }

  tbody tr {
    transition: background var(--transition-fast);
  }

  tbody tr:hover {
    background: var(--color-background);
  }

  .cell-name {
    font-weight: 600;
  }

  .avatar-cell {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .cell-email {
    color: var(--color-text-muted);
    font-size: var(--text-xs);
  }

  .cell-mono {
    font-family: monospace;
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    white-space: nowrap;
  }

  .cell-level {
    font-size: var(--text-xs);
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    color: var(--color-text-muted);
  }

  .cell-muted {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .spec-badge {
    display: inline-block;
    padding: 0.15rem 0.45rem;
    background: color-mix(in srgb, var(--color-accent) 12%, transparent);
    color: var(--color-accent);
    border-radius: 4px;
    font-size: 0.7rem;
    font-weight: 700;
    font-family: monospace;
    margin-right: 0.35rem;
  }

  .spec-name {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .cell-actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    white-space: nowrap;
  }

  th:last-child {
    white-space: nowrap;
  }

  .action-link {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    font-size: var(--text-xs);
    font-weight: 600;
    font-family: var(--font-sans);
    color: var(--color-accent);
    text-decoration: none;
    cursor: pointer;
    background: none;
    border: none;
    padding: 0;
    transition: opacity var(--transition-fast);
  }

  .action-link:hover {
    opacity: 0.75;
  }

  .action-link.danger {
    color: var(--color-danger);
  }

  .action-link.success {
    color: var(--color-success);
  }

  .empty {
    text-align: center;
    color: var(--color-text-muted);
    padding: 2.5rem;
    font-style: italic;
  }

  .form-grid {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
    margin-top: 1.25rem;
  }

  .csv-hint {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    font-family: var(--font-sans);
    margin: 0.75rem 0 0;
    line-height: 1.5;
  }

  .checkbox-field {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    color: var(--color-text);
    cursor: pointer;
  }

  .checkbox-field input[type="checkbox"] {
    width: 1rem;
    height: 1rem;
    accent-color: var(--color-accent);
    cursor: pointer;
  }

  :global(.csv-download-link) {
    color: var(--color-accent);
    font-weight: 600;
    text-decoration: none;
  }

  :global(.csv-download-link:hover) {
    text-decoration: underline;
  }

  @media (max-width: 768px) {
    .toolbar {
      flex-direction: column;
      align-items: stretch;
    }

    .search-wrap {
      min-width: 0;
    }

    table {
      font-size: var(--text-xs);
    }

    th,
    td {
      padding: 0.6rem 0.75rem;
    }
  }
</style>
