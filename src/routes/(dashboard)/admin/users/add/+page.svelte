<script lang="ts">
  import { admin } from "$lib/api";
  import { goto } from "$app/navigation";
  import Page from "$lib/components/ui/Page.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";
  import { toast } from "svelte-sonner";
  import { UserPlus, GraduationCap, Users } from "lucide-svelte";

  let { data } = $props();


  let type = $state(data.type as "teacher" | "student");


  let fullName = $state("");
  let email = $state("");
  let saving = $state(false);


  let grade = $state("assistant");
  let departmentId = $state<number | null>(null);
  let selectedDomains = $state<number[]>([]);


  let studentNumber = $state("");
  let level = $state("");
  let specialityId = $state<number | null>(null);
  let promotionId = $state<number | null>(null);

  function toggleDomain(id: number) {
    if (selectedDomains.includes(id)) {
      selectedDomains = selectedDomains.filter((d) => d !== id);
    } else {
      selectedDomains = [...selectedDomains, id];
    }
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();

    if (fullName.trim().length < 3) {
      toast.error("Le nom complet doit contenir au moins 3 caractères.");
      return;
    }
    if (fullName.trim().length > 100) {
      toast.error("Le nom complet ne doit pas dépasser 100 caractères.");
      return;
    }
    if (!email.trim() || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.trim())) {
      toast.error("Veuillez entrer une adresse email valide.");
      return;
    }
    if (type === "student" && studentNumber.trim() && !/^\d{4,20}$/.test(studentNumber.trim())) {
      toast.error("Le numéro d'étudiant doit contenir entre 4 et 20 chiffres.");
      return;
    }

    saving = true;
    try {
      if (type === "teacher") {
        await admin.createTeacher({
          full_name: fullName,
          email,
          grade,
          department_id: departmentId ? Number(departmentId) : 0,
        });


      } else {
        await admin.createStudent({
          full_name: fullName,
          email,
          student_number: studentNumber,
          speciality_id: specialityId ? Number(specialityId) : undefined,
          level: level || undefined,
          promotion_id: promotionId ? Number(promotionId) : undefined,
        });
      }

      toast.success(
        type === "teacher"
          ? "Enseignant créé avec succès !"
          : "Étudiant créé avec succès !",
      );
      goto("/admin/users");
    } catch (err: any) {
      toast.error(err.message || "Erreur lors de la création");
    } finally {
      saving = false;
    }
  }

  const typeLabel = $derived(type === "teacher" ? "Enseignant" : "Étudiant");
  const TypeIcon = $derived(type === "teacher" ? Users : GraduationCap);
</script>

<Page
  title="Ajouter un {typeLabel}"
  subtitle="Créer un nouveau compte {typeLabel.toLowerCase()} dans le système."
>
  <div class="add-container">
    <!-- Left sidebar: type switcher + info -->
    <div class="sidebar-section card">
      <div class="type-icon-wrap">
        <div class="type-icon">
          <TypeIcon size={40} />
        </div>
        <h3>Type de compte</h3>
      </div>

      <div class="type-switcher">
        <button
          type="button"
          class="type-btn"
          class:active={type === "teacher"}
          onclick={() => (type = "teacher")}
        >
          <Users size={16} />
          Enseignant
        </button>
        <button
          type="button"
          class="type-btn"
          class:active={type === "student"}
          onclick={() => (type = "student")}
        >
          <GraduationCap size={16} />
          Étudiant
        </button>
      </div>

      <div class="info-box">
        <UserPlus size={20} />
        <p>
          {#if type === "teacher"}
            Un email de bienvenue sera envoyé à l'enseignant avec ses
            identifiants de connexion.
          {:else}
            Un email de bienvenue sera envoyé à l'étudiant avec ses identifiants
            de connexion.
          {/if}
        </p>
      </div>
    </div>

    <!-- Right: form -->
    <form class="details-section card" onsubmit={handleSubmit}>
      <h3>Informations de Base</h3>
      <div class="form-grid">
        <FormField label="Nom Complet" required>
          <input
            type="text"
            class="input"
            bind:value={fullName}
            required
            minlength={3}
            maxlength={100}
            placeholder="Ex: Ahmed Benali"
          />
        </FormField>
        <FormField label="Email" required>
          <input
            type="email"
            class="input"
            bind:value={email}
            required
            placeholder="exemple@université.dz"
          />
        </FormField>
      </div>

      {#if type === "teacher"}
        <h3 class="section-title">Informations Enseignant</h3>
        <div class="form-grid">
          <FormField label="Grade" required>
            <select class="input" bind:value={grade} required>
              <option value="">Sélectionner</option>
              <option value="assistant">Assistant</option>
              <option value="mab">MAB</option>
              <option value="maa">MAA</option>
              <option value="mcb">MCB</option>
              <option value="mca">MCA</option>
              <option value="professeur">Professeur</option>
            </select>
          </FormField>

          <FormField label="Département">
            <select class="input" bind:value={departmentId}>
              <option value={null}>Sélectionner un département</option>
              {#each data.departments as dept}
                <option value={dept.id}>{dept.name}</option>
              {/each}
            </select>
          </FormField>
        </div>

        <FormField
          label="Domaines de Spécialité ({selectedDomains.length} sélectionnés)"
        >
          <div class="domains-grid">
            {#each data.domains as dom}
              <label class="checkbox-label">
                <input
                  type="checkbox"
                  checked={selectedDomains.includes(dom.id)}
                  onchange={() => toggleDomain(dom.id)}
                />
                <span class="checkbox-text">{dom.name}</span>
              </label>
            {/each}
          </div>
        </FormField>
      {:else if type === "student"}
        <h3 class="section-title">Informations Étudiant</h3>
        <div class="form-grid">
          <FormField label="Numéro d'Étudiant">
            <input
              type="text"
              class="input"
              bind:value={studentNumber}
              placeholder="Ex: 20210001"
            />
          </FormField>

          <FormField label="Niveau">
            <select class="input" bind:value={level}>
              <option value="">Sélectionner</option>
              <option value="L3">L3</option>
              <option value="M1">M1</option>
              <option value="M2">M2</option>
            </select>
          </FormField>

          <FormField label="Spécialité">
            <select class="input" bind:value={specialityId}>
              <option value={null}>Sélectionner une spécialité</option>
              {#each data.specialities as spec}
                <option value={spec.id}>{spec.code} - {spec.name}</option>
              {/each}
            </select>
          </FormField>

          <FormField label="Promotion">
            <select class="input" bind:value={promotionId}>
              <option value={null}>Sélectionner une promotion</option>
              {#each data.promotions as promo}
                <option value={promo.id}>{promo.label}</option>
              {/each}
            </select>
          </FormField>
        </div>
      {/if}

      <div class="form-actions">
        <Button
          variant="secondary"
          outline
          onclick={() => goto("/admin/users")}
          disabled={saving}
        >
          Annuler
        </Button>
        <Button type="submit" disabled={saving} Icon={UserPlus}>
          {saving ? "Création..." : `Créer le compte`}
        </Button>
      </div>
    </form>
  </div>
</Page>

<style>
  .add-container {
    display: grid;
    grid-template-columns: 280px 1fr;
    gap: var(--spacing-lg);
    align-items: start;
  }

  @media (max-width: 900px) {
    .add-container {
      grid-template-columns: 1fr;
    }
  }

  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: var(--spacing-lg);
  }

  h3 {
    margin: 0 0 var(--spacing-md) 0;
    font-size: var(--text-lg);
    font-weight: 600;
    border-bottom: 1px solid var(--color-border);
    padding-bottom: var(--spacing-sm);
    font-family: var(--font-sans);
    color: var(--color-text);
  }

  .section-title {
    margin-top: var(--spacing-xl);
  }

  /* ── Sidebar ── */
  .sidebar-section {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-md);
  }

  .type-icon-wrap {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    text-align: center;
    padding-bottom: var(--spacing-md);
    border-bottom: 1px solid var(--color-border);
  }

  .type-icon-wrap h3 {
    margin: 0;
    border: none;
    padding: 0;
    font-size: var(--text-base);
  }

  .type-icon {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: color-mix(in srgb, var(--color-accent) 12%, transparent);
    border: 2px solid color-mix(in srgb, var(--color-accent) 30%, transparent);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--color-accent);
    margin-bottom: 0.25rem;
  }

  .type-switcher {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .type-btn {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    padding: 0.65rem 1rem;
    border-radius: 8px;
    border: 1px solid var(--color-border);
    background: var(--color-background);
    color: var(--color-text-muted);
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-fast);
    text-align: left;
    width: 100%;
  }

  .type-btn:hover:not(.active) {
    border-color: var(--color-accent);
    color: var(--color-text);
    background: color-mix(in srgb, var(--color-accent) 5%, transparent);
  }

  .type-btn.active {
    background: var(--color-accent);
    border-color: var(--color-accent);
    color: var(--color-background);
  }

  .info-box {
    display: flex;
    gap: 0.75rem;
    align-items: flex-start;
    padding: 0.75rem 1rem;
    background: color-mix(in srgb, var(--color-accent) 8%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-accent) 20%, transparent);
    border-radius: 8px;
    color: var(--color-accent);
  }

  .info-box p {
    margin: 0;
    font-size: var(--text-xs);
    font-family: var(--font-sans);
    line-height: 1.5;
    color: var(--color-text-muted);
  }

  /* ── Form ── */
  .form-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-md);
    margin-bottom: var(--spacing-md);
  }

  @media (max-width: 600px) {
    .form-grid {
      grid-template-columns: 1fr;
    }
  }

  .input {
    width: 100%;
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    background: var(--color-background);
    color: var(--color-text);
    box-sizing: border-box;
    transition: border-color var(--transition-fast);
  }

  .input:focus {
    outline: none;
    border-color: var(--color-accent);
  }

  .domains-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 0.5rem;
    padding: 1rem;
    background: color-mix(
      in srgb,
      var(--color-surface) 50%,
      var(--color-background)
    );
    border: 1px solid var(--color-border);
    border-radius: 8px;
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    width: fit-content;
  }

  .checkbox-label input[type="checkbox"] {
    accent-color: var(--color-accent);
  }

  .checkbox-text {
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    color: var(--color-text);
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--spacing-sm);
    margin-top: var(--spacing-xl);
    padding-top: var(--spacing-md);
    border-top: 1px solid var(--color-border);
  }
</style>
