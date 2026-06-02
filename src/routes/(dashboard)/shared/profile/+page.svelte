<script lang="ts">
  import type { SessionUser } from "$lib/types";

  import { onMount } from "svelte";
  import { authStore } from "$lib/stores/auth";
  import { showToast } from "$lib/utils/toast";
  import { auth, upload } from "$lib/api";
  import { Camera, Mail, Shield, User } from "lucide-svelte";

  import Page from "$lib/components/ui/Page.svelte";

  let profile = $state<SessionUser | null>(null);
  let loading = $state(true);
  let uploading = $state(false);
  let fileInput: HTMLInputElement;

  const ROLE_LABELS: Record<string, string> = {
    admin: "Administrateur",
    teacher: "Enseignant",
    student: "Étudiant",
    company: "Entreprise",
  };

  onMount(async () => {
    try {
      profile = await auth.me();
    } catch {
      profile = authStore.profile;
    } finally {
      loading = false;
    }
  });

  function initials(name: string): string {
    return name
      .split(" ")
      .map((w) => w[0])
      .slice(0, 2)
      .join("")
      .toUpperCase();
  }

  function avatarColor(name: string): string {
    let hash = 0;
    for (let i = 0; i < name.length; i++) {
      hash = name.charCodeAt(i) + ((hash << 5) - hash);
    }
    const hue = Math.abs(hash) % 360;
    return `hsl(${hue}, 70%, 50%)`;
  }

  function openFilePicker() {
    fileInput?.click();
  }

  async function handleAvatarChange(e: Event) {
    const input = e.target as HTMLInputElement;
    const file = input.files?.[0];
    if (!file) return;


    if (!file.type.startsWith("image/")) {
      showToast.error("Fichier invalide", "Veuillez sélectionner une image.");
      return;
    }
    if (file.size > 5 * 1024 * 1024) {
      showToast.error("Fichier trop volumineux", "Taille max: 5 Mo.");
      return;
    }

    uploading = true;
    try {
      const formData = new FormData();
      formData.append("file", file);
      const result = await upload.avatar(formData);


      if (profile) {
        profile = { ...profile, avatar_url: result.url };
      }


      await authStore.refreshProfile();

      showToast.success("Photo de profil mise à jour !");
    } catch (err) {
      showToast.error("Échec du téléchargement", String(err));
    } finally {
      uploading = false;

      input.value = "";
    }
  }
</script>

<Page
  title="Mon profil"
  subtitle="Consultez et gérez vos informations personnelles."
>
  {#if loading}
    <p class="empty">Chargement...</p>
  {:else if profile}
    <div class="profile-page">
      <div class="card avatar-card">
        <button
          class="avatar-wrapper"
          onclick={openFilePicker}
          disabled={uploading}
          title="Cliquer pour changer la photo"
        >
          {#if profile.avatar_url}
            <img
              class="avatar-img"
              src={profile.avatar_url}
              alt={profile.full_name}
            />
          {:else}
            <div
              class="avatar-placeholder"
              style:background={avatarColor(profile.full_name)}
            >
              <span class="avatar-initials">{initials(profile.full_name)}</span>
            </div>
          {/if}

          <div class="avatar-overlay" class:uploading>
            {#if uploading}
              <span class="overlay-text">Envoi...</span>
            {:else}
              <Camera size={24} />
              <span class="overlay-text">Modifier</span>
            {/if}
          </div>
        </button>

        <input
          bind:this={fileInput}
          type="file"
          accept="image/*"
          class="hidden-input"
          onchange={handleAvatarChange}
        />

        <h2 class="profile-name">{profile.full_name}</h2>
        <span class="profile-role"
          >{ROLE_LABELS[profile.role] ?? profile.role}</span
        >
      </div>

      <div class="card info-card">
        <h3 class="card-title">Informations</h3>

        <div class="info-grid">
          <div class="info-item">
            <div class="info-icon"><User size={16} /></div>
            <div class="info-content">
              <span class="info-label">Nom complet</span>
              <span class="info-value">{profile.full_name}</span>
            </div>
          </div>

          <div class="info-item">
            <div class="info-icon"><Mail size={16} /></div>
            <div class="info-content">
              <span class="info-label">Email</span>
              <span class="info-value">{profile.email}</span>
            </div>
          </div>

          <div class="info-item">
            <div class="info-icon"><Shield size={16} /></div>
            <div class="info-content">
              <span class="info-label">Rôle</span>
              <span class="info-value"
                >{ROLE_LABELS[profile.role] ?? profile.role}</span
              >
            </div>
          </div>
        </div>
      </div>
    </div>
  {:else}
    <div class="empty">
      <p>Impossible de charger le profil.</p>
    </div>
  {/if}
</Page>

<style>
  .empty {
    text-align: center;
    padding: 3rem 1rem;
    color: var(--color-text-muted);
    font-size: var(--text-sm);
  }

  .profile-page {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-lg);
  }

  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 1.5rem;
  }

  /* ── Avatar card ── */
  .avatar-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
    padding: 2rem 1.5rem;
  }

  .avatar-wrapper {
    position: relative;
    width: 120px;
    height: 120px;
    border-radius: 50%;
    overflow: hidden;
    cursor: pointer;
    border: 3px solid var(--color-border);
    background: none;
    padding: 0;
    transition:
      border-color var(--transition-fast),
      transform var(--transition-fast);
  }

  .avatar-wrapper:hover {
    border-color: var(--color-accent);
    transform: scale(1.05);
  }

  .avatar-wrapper:active {
    transform: scale(0.97);
  }

  .avatar-wrapper:disabled {
    cursor: wait;
    opacity: 0.7;
  }

  .avatar-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .avatar-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .avatar-initials {
    font-size: 2.5rem;
    font-weight: 700;
    color: #fff;
    font-family: var(--font-sans);
  }

  .avatar-overlay {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.25rem;
    color: #fff;
    opacity: 0;
    transition: opacity var(--transition-fast);
  }

  .avatar-wrapper:hover .avatar-overlay,
  .avatar-overlay.uploading {
    opacity: 1;
  }

  .overlay-text {
    font-size: 0.75rem;
    font-weight: 600;
    font-family: var(--font-sans);
  }

  .hidden-input {
    display: none;
  }

  .profile-name {
    margin: 0;
    font-size: var(--text-xl);
    font-weight: 600;
    color: var(--color-text);
  }

  .profile-role {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    font-weight: 500;
  }

  /* ── Info card ── */
  .card-title {
    margin: 0 0 1rem;
    font-size: var(--text-lg);
    font-weight: 600;
    color: var(--color-text);
  }

  .info-grid {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .info-item {
    display: flex;
    align-items: center;
    gap: 0.85rem;
  }

  .info-icon {
    width: 36px;
    height: 36px;
    border-radius: 8px;
    background: var(--color-background-100);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--color-text-muted);
    flex-shrink: 0;
  }

  .info-content {
    display: flex;
    flex-direction: column;
    gap: 0.1rem;
  }

  .info-label {
    font-size: 0.7rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--color-text-muted);
  }

  .info-value {
    font-size: var(--text-sm);
    color: var(--color-text);
    font-weight: 500;
  }
</style>
