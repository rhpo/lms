<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { upload } from "$lib/api";
  import {
    Building2,
    ImagePlus,
    X,
    CheckCircle,
    AlertCircle,
  } from "lucide-svelte";

  let { data } = $props();
  let profile = $derived(data.profile);
  let company = $derived(profile?.company);

  let logoFile = $state<File | null>(null);
  let logoPreview = $state<string>("");
  let uploading = $state(false);
  let uploadSuccess = $state(false);
  let uploadError = $state("");

  function handleLogoChange(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const file = target.files?.[0];
    if (!file) return;
    logoFile = file;
    uploadSuccess = false;
    uploadError = "";
    const reader = new FileReader();
    reader.onload = (ev) => {
      logoPreview = ev.target?.result as string;
    };
    reader.readAsDataURL(file);
  }

  function clearLogo() {
    logoFile = null;
    logoPreview = "";
    uploadError = "";
    uploadSuccess = false;
  }

  async function uploadLogo() {
    if (!logoFile) return;
    uploading = true;
    uploadError = "";
    uploadSuccess = false;
    try {
      const fd = new FormData();
      fd.append("file", logoFile);
      await upload.companyLogo(fd);
      uploadSuccess = true;
      logoFile = null;
      logoPreview = "";
      await invalidateAll();
    } catch (err: any) {
      uploadError = err.message || "Erreur lors de l'upload";
    } finally {
      uploading = false;
    }
  }
</script>

<div class="profile-page">
  <div class="page-header">
    <Building2 size={24} />
    <h1>Profil entreprise</h1>
  </div>

  <!-- Company info card -->
  <div class="card">
    <div class="card-header">
      <div class="company-avatar">
        {#if company?.logo_url}
          <img src={company.logo_url} alt={company.company_name ?? ""} />
        {:else}
          <Building2 size={32} />
        {/if}
      </div>
      <div class="company-meta">
        <h2>{company?.company_name ?? "-"}</h2>
        {#if company?.sector}
          <span class="sector">{company.sector}</span>
        {/if}
        {#if company?.is_verified}
          <span class="badge verified">Vérifiée</span>
        {:else}
          <span class="badge pending">En attente de validation</span>
        {/if}
      </div>
    </div>

    {#if company?.description}
      <p class="description">{company.description}</p>
    {/if}

    <div class="info-grid">
      {#if company?.contact_email}
        <div class="info-item">
          <span class="info-label">Email</span>
          <span class="info-value">{company.contact_email}</span>
        </div>
      {/if}
      {#if company?.contact_phone}
        <div class="info-item">
          <span class="info-label">Téléphone</span>
          <span class="info-value">{company.contact_phone}</span>
        </div>
      {/if}
    </div>
  </div>

  <!-- Logo upload card -->
  <div class="card">
    <h3 class="section-title">
      <ImagePlus size={18} />
      Logo de l'entreprise
    </h3>
    <p class="section-desc">
      Uploadez ou mettez à jour le logo de votre entreprise. Il sera visible par
      les étudiants lors de la sélection d'entreprise.
    </p>

    {#if uploadSuccess}
      <div class="alert success">
        <CheckCircle size={16} />
        Logo mis à jour avec succès.
      </div>
    {/if}

    {#if uploadError}
      <div class="alert error">
        <AlertCircle size={16} />
        {uploadError}
      </div>
    {/if}

    {#if logoPreview}
      <div class="logo-preview-wrap">
        <img src={logoPreview} alt="Aperçu" class="logo-preview" />
        <div class="logo-preview-info">
          <span class="logo-filename">{logoFile?.name}</span>
          <button class="logo-clear" type="button" onclick={clearLogo}>
            <X size={14} /> Supprimer
          </button>
        </div>
      </div>
      <button class="btn-upload" onclick={uploadLogo} disabled={uploading}>
        {uploading ? "Upload en cours..." : "Enregistrer le logo"}
      </button>
    {:else}
      <label class="logo-upload-zone" for="logo-input">
        {#if company?.logo_url}
          <img src={company.logo_url} alt="Logo actuel" class="current-logo" />
          <span>Cliquez pour changer le logo</span>
        {:else}
          <ImagePlus size={32} />
          <span>Cliquez pour choisir un logo</span>
        {/if}
        <span class="logo-hint">PNG, JPG ou WEBP · max 5 Mo</span>
      </label>
      <input
        id="logo-input"
        type="file"
        accept="image/png,image/jpeg,image/webp"
        class="logo-input-hidden"
        onchange={handleLogoChange}
      />
    {/if}
  </div>
</div>

<style>
  .profile-page {
    max-width: 720px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    padding: 1.5rem;
  }

  .page-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    color: var(--color-text);
  }

  .page-header h1 {
    font-size: 1.5rem;
    font-weight: 700;
    margin: 0;
  }

  /* ─── Card ─── */
  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 1rem;
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .card-header {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .company-avatar {
    width: 72px;
    height: 72px;
    border-radius: 1rem;
    background: var(--color-background-100);
    border: 1px solid var(--color-border);
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    flex-shrink: 0;
    color: var(--color-text-muted);
  }

  .company-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .company-meta {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .company-meta h2 {
    font-size: 1.2rem;
    font-weight: 700;
    margin: 0;
    color: var(--color-text);
  }

  .sector {
    font-size: 0.85rem;
    color: var(--color-text-muted);
  }

  .badge {
    display: inline-block;
    padding: 0.2rem 0.6rem;
    border-radius: 2rem;
    font-size: 0.75rem;
    font-weight: 600;
    width: fit-content;
  }

  .badge.verified {
    background: color-mix(in srgb, #22c55e 15%, transparent);
    color: #16a34a;
  }

  .badge.pending {
    background: color-mix(in srgb, #f59e0b 15%, transparent);
    color: #b45309;
  }

  .description {
    font-size: 0.9rem;
    color: var(--color-text-muted);
    line-height: 1.6;
    margin: 0;
  }

  .info-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.75rem;
  }

  .info-item {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }

  .info-label {
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--color-text-muted);
    font-weight: 600;
  }

  .info-value {
    font-size: 0.9rem;
    color: var(--color-text);
  }

  /* ─── Logo Upload ─── */
  .section-title {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1rem;
    font-weight: 600;
    color: var(--color-text);
    margin: 0;
  }

  .section-desc {
    font-size: 0.875rem;
    color: var(--color-text-muted);
    margin: 0;
    line-height: 1.5;
  }

  .alert {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1rem;
    border-radius: 0.5rem;
    font-size: 0.875rem;
    font-weight: 500;
  }

  .alert.success {
    background: color-mix(in srgb, #22c55e 12%, transparent);
    color: #16a34a;
    border: 1px solid color-mix(in srgb, #22c55e 30%, transparent);
  }

  .alert.error {
    background: color-mix(in srgb, #ef4444 10%, transparent);
    color: #dc2626;
    border: 1px solid color-mix(in srgb, #ef4444 25%, transparent);
  }

  .logo-upload-zone {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.6rem;
    padding: 2rem;
    border: 2px dashed var(--color-border);
    border-radius: 0.75rem;
    cursor: pointer;
    color: var(--color-text-muted);
    transition: all 0.2s ease;
    background: var(--color-background-100);
    text-align: center;
  }

  .logo-upload-zone:hover {
    border-color: var(--color-accent);
    color: var(--color-accent);
    background: color-mix(
      in srgb,
      var(--color-accent) 4%,
      var(--color-background-100)
    );
  }

  .logo-upload-zone span {
    font-size: 0.9rem;
    font-weight: 500;
  }

  .logo-hint {
    font-size: 0.78rem !important;
    font-weight: 400 !important;
    color: var(--color-text-muted) !important;
  }

  .current-logo {
    width: 64px;
    height: 64px;
    border-radius: 0.5rem;
    object-fit: cover;
    border: 1px solid var(--color-border);
    margin-bottom: 0.25rem;
  }

  .logo-input-hidden {
    display: none;
  }

  .logo-preview-wrap {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.75rem 1rem;
    border: 1px solid var(--color-border);
    border-radius: 0.75rem;
    background: var(--color-background-100);
  }

  .logo-preview {
    width: 56px;
    height: 56px;
    border-radius: 0.5rem;
    object-fit: cover;
    border: 1px solid var(--color-border);
    flex-shrink: 0;
  }

  .logo-preview-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    min-width: 0;
  }

  .logo-filename {
    font-size: 0.875rem;
    color: var(--color-text);
    font-weight: 500;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .logo-clear {
    display: inline-flex;
    align-items: center;
    gap: 0.3rem;
    font-size: 0.8rem;
    color: #ef4444;
    background: none;
    border: none;
    cursor: pointer;
    padding: 0;
    font-family: var(--font-sans);
    font-weight: 500;
  }

  .logo-clear:hover {
    text-decoration: underline;
  }

  .btn-upload {
    align-self: flex-start;
    padding: 0.65rem 1.25rem;
    background: var(--color-accent);
    color: white;
    border: none;
    border-radius: 0.5rem;
    font-family: var(--font-sans);
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: opacity 0.2s;
  }

  .btn-upload:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .btn-upload:hover:not(:disabled) {
    opacity: 0.9;
  }
</style>
