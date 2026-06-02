<script lang="ts">
  import { untrack } from "svelte";
  import { admin } from "$lib/api";
  import { goto, invalidateAll } from "$app/navigation";
  import Page from "$lib/components/ui/Page.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import FormField from "$lib/components/ui/FormField.svelte";
  import Avatar from "$lib/components/ui/Avatar.svelte";
  import { toast } from "svelte-sonner";

  let { data } = $props();


  let profile = $state(JSON.parse(JSON.stringify(data.profile)));
  let saving = $state(false);
  let uploadingAvatar = $state(false);


  let fullName = $state(data.profile.full_name || "");
  let email = $state(data.profile.email || "");


  let grade = $state(data.profile.teacher?.grade ?? "");
  let departmentId = $state<number | null>(data.profile.teacher?.department_id ?? null);
  let selectedDomains = $state<number[]>(
    data.profile.teacher?.domaines?.map((d: any) => d.id) ?? [],
  );


  let studentNumber = $state(data.profile.student?.student_number ?? "");
  let level = $state(data.profile.student?.level ?? "");
  let specialityId = $state<number | null>(data.profile.student?.speciality_id ?? null);
  let promotionId = $state<number | null>(data.profile.student?.promotion_id ?? null);


  $effect(() => {
    const p = data.profile;
    untrack(() => {
      profile = JSON.parse(JSON.stringify(p));
      fullName = p.full_name || "";
      email = p.email || "";
      grade = p.teacher?.grade ?? "";
      departmentId = p.teacher?.department_id ?? null;
      selectedDomains = p.teacher?.domaines?.map((d: any) => d.id) ?? [];
      studentNumber = p.student?.student_number ?? "";
      level = p.student?.level ?? "";
      specialityId = p.student?.speciality_id ?? null;
      promotionId = p.student?.promotion_id ?? null;
    });
  });

  async function handleAvatarPick(e: Event) {
    const el = e.target as HTMLInputElement;
    if (!el.files || el.files.length === 0) return;

    uploadingAvatar = true;
    const formData = new FormData();
    formData.append("file", el.files[0]);

    try {
      const res = await admin.updateUserAvatar(profile.id, formData);
      profile.avatar_url = res.url;
      toast.success("Avatar mis à jour");
      await invalidateAll();
    } catch (err: any) {
      toast.error(err.message || "Erreur lors de la mise à jour de l'avatar");
    } finally {
      uploadingAvatar = false;
      el.value = "";
    }
  }

  function toggleDomain(id: number) {
    if (selectedDomains.includes(id)) {
      selectedDomains = selectedDomains.filter((d: number) => d !== id);
    } else {
      selectedDomains = [...selectedDomains, id];
    }
  }

  async function handleSubmit(e: Event) {
    e.preventDefault();
    saving = true;
    try {
      const payload: any = {
        full_name: fullName,
        email: email,
      };

      if (profile.role === "teacher" || profile.role === "admin") {
        payload.grade = grade;
        payload.department_id = departmentId ? Number(departmentId) : null;
        payload.domain_ids = selectedDomains;
      } else if (profile.role === "student") {
        payload.student_number = studentNumber;
        payload.level = level;
        payload.speciality_id = specialityId ? Number(specialityId) : null;
        payload.promotion_id = promotionId ? Number(promotionId) : null;
      }

      await admin.updateUser(profile.id, payload);
      toast.success("Utilisateur mis à jour avec succès");
      await invalidateAll();
      goto("/admin/users");
    } catch (err: any) {
      toast.error(err.message || "Erreur lors de la mise à jour");
    } finally {
      saving = false;
    }
  }
</script>

<Page
  title="Modifier {profile.full_name}"
  subtitle="Gérer les informations de ce profil ({profile.role})"
>
  <div class="edit-container">
    <div class="avatar-section card">
      <h3>Avatar</h3>
      <div class="avatar-wrapper">
        <Avatar user={profile} size={120} />
      </div>
      <div class="avatar-actions">
        <label class="btn-upload" class:disabled={uploadingAvatar}>
          {uploadingAvatar ? "Téléchargement..." : "Changer l'avatar"}
          <input
            type="file"
            accept="image/png, image/jpeg, image/webp"
            onchange={handleAvatarPick}
            disabled={uploadingAvatar}
          />
        </label>
      </div>
    </div>

    <form class="details-section card" onsubmit={handleSubmit}>
      <h3>Informations de Base</h3>
      <div class="form-grid">
        <FormField label="Nom Complet" required>
          <input type="text" class="input" bind:value={fullName} required />
        </FormField>
        <FormField label="Email" required>
          <input type="email" class="input" bind:value={email} required />
        </FormField>
      </div>

      {#if profile.role === "teacher" || profile.role === "admin"}
        <h3 class="section-title">Informations Enseignant</h3>
        <div class="form-grid">
          <FormField label="Grade">
            <select class="input" bind:value={grade}>
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
      {:else if profile.role === "student"}
        <h3 class="section-title">Informations Étudiant</h3>
        <div class="form-grid">
          <FormField label="Numéro d'Étudiant">
            <input type="text" class="input" bind:value={studentNumber} />
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
        <Button type="submit" disabled={saving}>
          {saving ? "Enregistrement..." : "Enregistrer les modifications"}
        </Button>
      </div>
    </form>
  </div>
</Page>

<style>
  .edit-container {
    display: grid;
    grid-template-columns: 300px 1fr;
    gap: var(--spacing-lg);
    align-items: start;
  }

  @media (max-width: 900px) {
    .edit-container {
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
  }

  .section-title {
    margin-top: var(--spacing-xl);
  }

  .avatar-section {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
  }

  .avatar-wrapper {
    margin-bottom: var(--spacing-md);
  }

  .avatar-actions {
    display: flex;
    justify-content: center;
  }

  .btn-upload {
    display: inline-block;
    padding: 0.5rem 1rem;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    font-size: var(--text-sm);
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
  }

  .btn-upload:hover:not(.disabled) {
    border-color: var(--color-accent);
    color: var(--color-accent);
  }

  .btn-upload.disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-upload input[type="file"] {
    display: none;
  }

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
  }

  .checkbox-text {
    font-size: var(--text-sm);
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
