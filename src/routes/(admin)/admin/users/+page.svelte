<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import {
    Users,
    BookOpen,
    Building2,
    Plus,
    Upload,
    Check,
    X,
    Search,
  } from "lucide-svelte";

  import Button from "$lib/components/ui/Button.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import Badge from "$lib/components/ui/Badge.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";
  import DashboardPage from "$lib/components/admin/DashboardPage.svelte";

  let { data } = $props();

  type Speciality = {
    id: string;
    name: string;
    code: string;
    year_type: "licence" | "master" | "ingenieur";
  };

  const specialities: Speciality[] = $derived(data.specialities);

  const niveaux = [
    { value: "licence", label: "Licence" },
    { value: "master", label: "Master" },
    { value: "ingenieur", label: "Ingenieur" },
  ] as const;

  const { formattedTeachers, formattedStudents, formattedCompanies } =
    $derived(data);

  type TabId = "teachers" | "students" | "companies";
  let activeTab: TabId = $state("teachers");

  let showAddTeacher = $state(false);
  let showAddStudent = $state(false);
  let showImportCsv = $state(false);
  let searchQuery = $state("");

  // Teacher form
  let teacherEmail = $state("");
  let teacherGrade = $state("assistant");
  let teacherDepartment = $state("");
  let teacherNiveau = $state<string>("");
  let teacherSpecialtyId = $state<string>("");

  const teacherFilteredSpecialities = $derived(
    teacherNiveau
      ? specialities.filter((s) => s.year_type === teacherNiveau)
      : [],
  );

  // Student form
  let studentEmail = $state("");
  let studentNumber = $state("");
  let studentNiveau = $state<string>("");
  let studentSpecialtyId = $state<string>("");
  let studentPromotionYear = $state(new Date().getFullYear());

  const studentFilteredSpecialities = $derived(
    studentNiveau
      ? specialities.filter((s) => s.year_type === studentNiveau)
      : [],
  );

  let csvType = $state("teachers");
  let csvFile: File | null = $state(null);

  function filterBySearch<T extends Record<string, unknown>>(
    items: T[],
    fields: (keyof T)[],
  ): T[] {
    if (!searchQuery) return items;
    const q = searchQuery.toLowerCase();
    return items.filter((item) =>
      fields.some((f) => {
        const val = item[f];
        return val != null && String(val).toLowerCase().includes(q);
      }),
    );
  }

  let filteredTeachers = $derived(
    filterBySearch(formattedTeachers, ["email", "department", "profile_name"]),
  );

  let filteredStudents = $derived(
    filterBySearch(formattedStudents, [
      "email",
      "student_number",
      "specialty",
      "profile_name",
    ]),
  );

  let filteredCompanies = $derived(
    filterBySearch(formattedCompanies, [
      "company_name",
      "contact_phone",
      "profile_email",
    ]),
  );

  async function addTeacher() {
    try {
      const res = await fetch("/api/users/add-teacher", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: teacherEmail,
          grade: teacherGrade,
          department: teacherDepartment,
          specialty_ids: teacherSpecialtyId ? [teacherSpecialtyId] : [],
        }),
      });
      if (res.ok) {
        showAddTeacher = false;
        teacherEmail = "";
        teacherGrade = "assistant";
        teacherDepartment = "";
        teacherNiveau = "";
        teacherSpecialtyId = "";
        await invalidateAll();
      }
    } catch (err) {
      console.error("Erreur reseau", err);
    }
  }

  async function addStudent() {
    try {
      const res = await fetch("/api/users/add-student", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          email: studentEmail,
          student_number: studentNumber,
          specialty: studentSpecialtyId,
          level: studentNiveau,
          promotion_year: studentPromotionYear,
        }),
      });
      if (res.ok) {
        showAddStudent = false;
        studentEmail = "";
        studentNumber = "";
        studentNiveau = "";
        studentSpecialtyId = "";
        studentPromotionYear = new Date().getFullYear();
        await invalidateAll();
      }
    } catch (err) {
      console.error("Erreur reseau", err);
    }
  }

  async function importCsv() {
    if (!csvFile) return;
    try {
      const formData = new FormData();
      formData.append("type", csvType);
      formData.append("file", csvFile);
      const res = await fetch("/api/users/import-csv", {
        method: "POST",
        body: formData,
      });
      if (res.ok) {
        showImportCsv = false;
        csvFile = null;
        await invalidateAll();
      }
    } catch (err) {
      console.error("Erreur reseau", err);
    }
  }

  async function deactivateUser(id: string) {
    try {
      const res = await fetch("/api/users/deactivate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id }),
      });
      if (!res.ok) {
        const err = await res.json();
        console.error(err.message);
        return;
      }
      await invalidateAll();
    } catch (err) {
      console.error("Erreur reseau", err);
    }
  }

  async function verifyCompany(profileId: string) {
    try {
      const res = await fetch("/api/users/verify-company", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ profileId }),
      });
      if (!res.ok) {
        const err = await res.json();
        console.error(err.message);
        return;
      }
      await invalidateAll();
    } catch (err) {
      console.error("Erreur reseau", err);
    }
  }
</script>

<DashboardPage
  title="Gestion des utilisateurs"
  subtitle="Gérer les comptes enseignants, etudiants et entreprises."
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
      <span class="tab-count">{formattedTeachers.length}</span>
    </button>
    <button
      class:active={activeTab === "students"}
      onclick={() => (activeTab = "students")}
    >
      <BookOpen size={16} />
      <span>Etudiants</span>
      <span class="tab-count">{formattedStudents.length}</span>
    </button>
    <button
      class:active={activeTab === "companies"}
      onclick={() => (activeTab = "companies")}
    >
      <Building2 size={16} />
      <span>Entreprises</span>
      <span class="tab-count">{formattedCompanies.length}</span>
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
        onclick={() => (showAddTeacher = true)}
      >
        Ajouter un enseignant
      </Button>
    {:else if activeTab === "students"}
      <Button
        variant="primary"
        Icon={Plus}
        onclick={() => (showAddStudent = true)}
      >
        Ajouter un etudiant
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
            <th>Departement</th>
            <th>Spécialités</th>
            <th>Disponibilite</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if filteredTeachers.length === 0}
            <tr>
              <td colspan="7" class="empty">Aucun enseignant trouve.</td>
            </tr>
          {:else}
            {#each filteredTeachers as teacher}
              <tr>
                <td class="cell-name">{teacher.profile_name}</td>
                <td class="cell-email">{teacher.email ?? "—"}</td>
                <td><Badge variant="info" label={teacher.grade_label} /></td>
                <td>{teacher.department}</td>
                <td>{teacher.specialties?.join(", ") ?? "—"}</td>
                <td>
                  <Badge
                    variant={teacher.availability_variant}
                    label={teacher.availability_label}
                  />
                </td>
                <td class="cell-actions">
                  <a href="/admin/users/{teacher.id}" class="action-link"
                    >Modifier</a
                  >
                  <button
                    class="action-link danger"
                    onclick={() =>
                      deactivateUser(teacher.profile_id ?? teacher.id)}
                  >
                    Desactiver
                  </button>
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
            <th>Email</th>
            <th>Numero etudiant</th>
            <th>Specialite</th>
            <th>Niveau</th>
            <th>Promotion</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if filteredStudents.length === 0}
            <tr>
              <td colspan="7" class="empty">Aucun etudiant trouve.</td>
            </tr>
          {:else}
            {#each filteredStudents as student}
              <tr>
                <td class="cell-name">{student.profile_name}</td>
                <td class="cell-email">{student.email ?? "—"}</td>
                <td>{student.student_number}</td>
                <td>{student.specialty}</td>
                <td><Badge variant="info" label={student.level_label} /></td>
                <td>{student.promotion_label}</td>
                <td class="cell-actions">
                  <a href="/admin/users/{student.id}" class="action-link"
                    >Modifier</a
                  >
                  <button
                    class="action-link danger"
                    onclick={() =>
                      deactivateUser(student.profile_id ?? student.id)}
                  >
                    Desactiver
                  </button>
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
            <th>Email</th>
            <th>Secteur</th>
            <th>Telephone</th>
            <th>Site web</th>
            <th>Statut</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#if filteredCompanies.length === 0}
            <tr>
              <td colspan="7" class="empty">Aucune entreprise trouvée.</td>
            </tr>
          {:else}
            {#each filteredCompanies as company}
              <tr>
                <td class="cell-name">{company.company_name}</td>
                <td class="cell-email">{company.profile_email}</td>
                <td>{company.sector ?? "—"}</td>
                <td>{company.contact_phone ?? "—"}</td>
                <td>
                  {#if company.website}
                    <a
                      href={company.website}
                      target="_blank"
                      rel="noopener"
                      class="action-link"
                    >
                      {company.website}
                    </a>
                  {:else}
                    —
                  {/if}
                </td>
                <td>
                  <Badge
                    variant={company.status_variant}
                    label={company.status_label}
                  />
                </td>
                <td class="cell-actions">
                  {#if !company.is_verified}
                    <button
                      class="action-link success"
                      onclick={() => verifyCompany(company.profile_id)}
                    >
                      <Check size={14} /> Valider
                    </button>
                  {/if}
                  <button
                    class="action-link danger"
                    onclick={() => deactivateUser(company.profile_id)}
                  >
                    <X size={14} /> Desactiver
                  </button>
                </td>
              </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
  {/if}

  <Modal bind:open={showAddTeacher} title="Ajouter un enseignant" width="520px">
    <div class="form-grid">
      <FormField label="Email" required>
        <input type="email" bind:value={teacherEmail} required class="input" />
      </FormField>
      <FormField label="Grade" required>
        <select bind:value={teacherGrade} required class="input">
          <option value="assistant">Assistant</option>
          <option value="mab">MAB</option>
          <option value="maa">MAA</option>
          <option value="mcb">MCB</option>
          <option value="mca">MCA</option>
          <option value="professeur">Professeur</option>
        </select>
      </FormField>
      <FormField label="Departement" required>
        <input
          type="text"
          bind:value={teacherDepartment}
          required
          class="input"
        />
      </FormField>
      <FormField label="Niveau" required>
        <select
          bind:value={teacherNiveau}
          required
          class="input"
          onchange={() => (teacherSpecialtyId = "")}
        >
          <option value="">Selectionner un niveau</option>
          {#each niveaux as n}
            <option value={n.value}>{n.label}</option>
          {/each}
        </select>
      </FormField>
      <FormField label="Specialite" required>
        <select
          bind:value={teacherSpecialtyId}
          required
          class="input"
          disabled={!teacherNiveau}
        >
          <option value="">
            {teacherNiveau
              ? "Selectionner une specialite"
              : "Choisissez d'abord un niveau"}
          </option>
          {#each teacherFilteredSpecialities as spec}
            <option value={spec.id}>{spec.code} — {spec.name}</option>
          {/each}
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

  <Modal bind:open={showAddStudent} title="Ajouter un etudiant" width="520px">
    <div class="form-grid">
      <FormField label="Email" required>
        <input type="email" bind:value={studentEmail} required class="input" />
      </FormField>
      <FormField label="Numero etudiant" required>
        <input type="text" bind:value={studentNumber} required class="input" />
      </FormField>
      <FormField label="Niveau" required>
        <select
          bind:value={studentNiveau}
          required
          class="input"
          onchange={() => (studentSpecialtyId = "")}
        >
          <option value="">Selectionner un niveau</option>
          {#each niveaux as n}
            <option value={n.value}>{n.label}</option>
          {/each}
        </select>
      </FormField>
      <FormField label="Specialite" required>
        <select
          bind:value={studentSpecialtyId}
          required
          class="input"
          disabled={!studentNiveau}
        >
          <option value="">
            {studentNiveau
              ? "Selectionner une specialite"
              : "Choisissez d'abord un niveau"}
          </option>
          {#each studentFilteredSpecialities as spec}
            <option value={spec.id}>{spec.code} — {spec.name}</option>
          {/each}
        </select>
      </FormField>
      <FormField label="Annee promotion" required>
        <input
          type="number"
          bind:value={studentPromotionYear}
          min="2000"
          max="2100"
          required
          class="input"
        />
      </FormField>
    </div>
    <div class="form-actions">
      <Button variant="ghost" onclick={() => (showAddStudent = false)}
        >Annuler</Button
      >
      <Button variant="primary" onclick={addStudent}>Ajouter</Button>
    </div>
  </Modal>

  <Modal
    bind:open={showImportCsv}
    title="Importer des comptes (CSV)"
    width="520px"
  >
    <div class="form-grid">
      <FormField label="Type de comptes" required>
        <select bind:value={csvType} required class="input">
          <option value="teachers">Enseignants</option>
          <option value="students">Etudiants</option>
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
    </div>
    <p class="csv-hint">
      Format attendu : email, prenom, nom, specialite, grade (pour enseignants)
      ou email, prenom, nom, specialite, promotion_id (pour etudiants).
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
    overflow: hidden;
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

  .cell-email {
    color: var(--color-text-muted);
    font-size: var(--text-xs);
  }

  .cell-actions {
    display: flex;
    align-items: center;
    gap: 0.5rem;
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
