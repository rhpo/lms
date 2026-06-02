<script lang="ts">
  import { goto } from "$app/navigation";
  import { BRAND } from "$lib/constants/branding";
  import { auth, upload, setToken } from "$lib/api";
  import type { Company } from "$lib/types";

  import Logo from "$lib/components/Logo.svelte";
  import View from "$lib/components/ui/View.svelte";
  import Input from "$lib/components/ui/Input/Input.svelte";
  import InputPassword from "$lib/components/ui/Input/InputPassword.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import {
    Building2,
    UserPlus,
    Search,
    ChevronRight,
    ChevronLeft,
    ArrowRight,
    CheckCircle,
    Plus,
    Briefcase,
    Mail,
    Phone,
    User,
    GraduationCap,
    ImagePlus,
    X,
  } from "lucide-svelte";


  let step = $state(1);
  let loading = $state(false);
  let error = $state("");
  let success = $state(false);


  let fullName = $state("");
  let email = $state("");
  let password = $state("");
  let confirmPassword = $state("");
  let phone = $state("");
  let position = $state("");


  let companySearch = $state("");
  let verifiedCompanies = $state<Company[]>([]);
  let selectedCompany = $state<Company | null>(null);
  let companiesLoaded = $state(false);


  let companyName = $state("");
  let sector = $state("");
  let description = $state("");
  let contactEmail = $state("");
  let contactPhone = $state("");


  let logoFile = $state<File | null>(null);
  let logoPreview = $state<string>("");

  function handleLogoChange(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const file = target.files?.[0];
    if (!file) return;
    logoFile = file;
    const reader = new FileReader();
    reader.onload = (ev) => {
      logoPreview = ev.target?.result as string;
    };
    reader.readAsDataURL(file);
  }

  function clearLogo() {
    logoFile = null;
    logoPreview = "";
  }


  const positionSuggestions = [
    "Directeur Général",
    "Directeur Technique",
    "Chef de Projet",
    "Ingénieur Logiciel",
    "Responsable RH",
    "Chargé de Recrutement",
    "Responsable Formation",
    "Ingénieur R&D",
    "Consultant",
    "Chef de Département",
    "Superviseur de Stage",
    "Tuteur Entreprise",
    "Responsable IT",
    "Architecte Logiciel",
    "DevOps Engineer",
    "Data Scientist",
    "Product Manager",
    "Scrum Master",
  ];
  let showPositionDropdown = $state(false);
  let filteredPositions = $derived(
    position.length > 0
      ? positionSuggestions.filter((p) =>
          p.toLowerCase().includes(position.toLowerCase()),
        )
      : positionSuggestions,
  );


  let filteredCompanies = $derived(
    companySearch.length > 0
      ? verifiedCompanies.filter(
          (c) =>
            c.company_name
              ?.toLowerCase()
              .includes(companySearch.toLowerCase()) ||
            c.sector?.toLowerCase().includes(companySearch.toLowerCase()),
        )
      : verifiedCompanies,
  );


  let step1Valid = $derived(
    fullName.trim().length >= 3 &&
      /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.trim()) &&
      password.length >= 6 &&
      password === confirmPassword,
  );

  let step3Valid = $derived(companyName.trim().length >= 2);


  async function loadCompanies() {
    if (companiesLoaded) return;
    try {
      verifiedCompanies = await auth.listVerifiedCompanies();
      companiesLoaded = true;
    } catch (err: any) {
      console.error("Failed to load companies:", err);
      verifiedCompanies = [];
      companiesLoaded = true;
    }
  }


  function goToStep2() {
    if (!step1Valid) return;
    step = 2;
    loadCompanies();
  }

  function selectCompany(c: Company) {
    selectedCompany = c;
  }

  function goToCreateCompany() {
    selectedCompany = null;
    step = 3;
  }


  async function handleSubmit() {
    loading = true;
    error = "";


    if (fullName.trim().length < 3) {
      error = "Le nom complet doit contenir au moins 3 caractères.";
      loading = false;
      return;
    }
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.trim())) {
      error = "Veuillez entrer une adresse email valide.";
      loading = false;
      return;
    }
    if (step === 3 && companyName.trim().length < 2) {
      error = "Le nom de l'entreprise doit contenir au moins 2 caractères.";
      loading = false;
      return;
    }
    if (
      step === 3 &&
      contactEmail &&
      !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(contactEmail.trim())
    ) {
      error = "L'email de l'entreprise n'est pas valide.";
      loading = false;
      return;
    }

    try {
      let result;
      if (step === 2 && selectedCompany) {

        result = await auth.registerCompany({
          full_name: fullName,
          email,
          position,
          phone,
          company_id: selectedCompany.id,
        });
      } else if (step === 3) {

        result = await auth.registerCompany({
          full_name: fullName,
          email,
          position,
          phone,
          company_name: companyName,
          sector,
          description,
          contact_email: contactEmail || email,
          contact_phone: contactPhone || phone,
        });
      }

      if (result?.token) {
        setToken(result.token);
      }

      if (step === 2 && selectedCompany) {

        goto("/company/dashboard");
        return;
      }


      if (logoFile) {
        try {
          const fd = new FormData();
          fd.append("file", logoFile);
          await upload.companyLogo(fd);
        } catch {

        }
      }


      success = true;
    } catch (err: any) {
      error = err.message || "Échec de l'inscription";
    } finally {
      loading = false;
    }
  }
</script>

<View fullScreen center>
  <div class="registration-container">
    {#if success}
      <!-- ─── Success Screen ─── -->
      <div class="success-card" data-aos="zoom-in">
        <div class="success-icon">
          <CheckCircle size={64} />
        </div>
        <h1>Compte créé avec succès !</h1>
        <p class="success-msg">
          {#if step === 3}
            Votre compte a été créé. Votre entreprise est <strong
              >en attente de validation</strong
            > par l'administration. Vous recevrez un email une fois votre entreprise
            approuvée.
          {:else}
            Votre compte a été créé et rattaché à l'entreprise <strong
              >{selectedCompany?.company_name}</strong
            >.
          {/if}
        </p>
        <div class="success-actions">
          <Button onclick={() => goto("/company/dashboard")} Icon={ArrowRight}>
            Accéder au tableau de bord
          </Button>
        </div>
      </div>
    {:else}
      <!-- ─── Registration Form ─── -->
      <div class="form-card" data-aos="fade-up">
        <!-- Header -->
        <div class="card-header">
          <Logo width="140px" />
          <div class="title-row">
            <Building2 size={28} />
            <h1>Créer un compte entreprise</h1>
          </div>
          <p class="subtitle">
            Inscrivez-vous pour gérer vos PFE externes sur la plateforme {BRAND
              .university.name}
          </p>
        </div>

        <!-- Progress Stepper -->
        <div class="stepper">
          <button
            class="step-item"
            class:active={step >= 1}
            class:current={step === 1}
            onclick={() => {
              if (step > 1) step = 1;
            }}
          >
            <span class="step-number">1</span>
            <span class="step-label">Vos informations</span>
          </button>
          <div class="step-line" class:active={step >= 2}></div>
          <button
            class="step-item"
            class:active={step >= 2}
            class:current={step === 2}
            onclick={() => {
              if (step > 2) step = 2;
            }}
          >
            <span class="step-number">2</span>
            <span class="step-label">Votre entreprise</span>
          </button>
          {#if step === 3}
            <div class="step-line" class:active={step >= 3}></div>
            <button
              class="step-item"
              class:active={step >= 3}
              class:current={step === 3}
            >
              <span class="step-number">3</span>
              <span class="step-label">Nouvelle entreprise</span>
            </button>
          {/if}
        </div>

        {#if error}
          <div class="error-banner">
            <p>{error}</p>
          </div>
        {/if}

        <!-- Step 1: Employee Info -->
        {#if step === 1}
          <div class="form-step" data-aos="fade-left">
            <div class="section-header">
              <User size={20} />
              <h2>Informations personnelles</h2>
            </div>

            <div class="form-grid">
              <div class="form-group full">
                <Input
                  name="fullName"
                  label="Nom complet"
                  placeholder="Prénom et Nom"
                  required
                  bind:value={fullName}
                />
              </div>

              <div class="form-group full">
                <Input
                  name="email"
                  type="email"
                  label="Email professionnel"
                  placeholder="vous@entreprise.com"
                  required
                  bind:value={email}
                />
              </div>

              <div class="form-group">
                <InputPassword
                  name="password"
                  label="Mot de passe"
                  placeholder="••••••••"
                  required
                  bind:value={password}
                />
              </div>

              <div class="form-group">
                <InputPassword
                  name="confirmPassword"
                  label="Confirmer le mot de passe"
                  placeholder="••••••••"
                  required
                  bind:value={confirmPassword}
                />
              </div>

              <div class="form-group">
                <Input
                  name="phone"
                  type="tel"
                  label="Téléphone"
                  placeholder="+213 5XX XX XX XX"
                  bind:value={phone}
                />
              </div>

              <div class="form-group full position-group">
                <Input
                  name="position"
                  label="Poste / Fonction"
                  placeholder="Ex: Chef de Projet, Ingénieur..."
                  bind:value={position}
                  onfocus={() => (showPositionDropdown = true)}
                  onblur={() =>
                    setTimeout(() => (showPositionDropdown = false), 200)}
                />
                {#if showPositionDropdown && filteredPositions.length > 0}
                  <div class="autocomplete-dropdown">
                    {#each filteredPositions.slice(0, 6) as suggestion}
                      <button
                        class="autocomplete-item"
                        type="button"
                        onmousedown={() => {
                          position = suggestion;
                          showPositionDropdown = false;
                        }}
                      >
                        <Briefcase size={14} />
                        {suggestion}
                      </button>
                    {/each}
                  </div>
                {/if}
              </div>
            </div>

            <div class="step-actions">
              <div></div>
              <Button
                onclick={goToStep2}
                disabled={!step1Valid}
                Icon={ChevronRight}
              >
                Suivant
              </Button>
            </div>
          </div>
        {/if}

        <!-- Step 2: Company Search / Select -->
        {#if step === 2}
          <div class="form-step" data-aos="fade-left">
            <div class="section-header">
              <Search size={20} />
              <h2>Rechercher votre entreprise</h2>
            </div>

            <div class="search-box">
              <Search size={18} />
              <input
                type="text"
                placeholder="Tapez le nom de votre entreprise..."
                bind:value={companySearch}
                class="search-input"
              />
            </div>

            {#if !companiesLoaded}
              <div class="loading-state">
                <div class="spinner"></div>
                <p>Chargement des entreprises...</p>
              </div>
            {:else if filteredCompanies.length > 0}
              <div class="company-list">
                {#each filteredCompanies as company}
                  <button
                    class="company-card"
                    class:selected={selectedCompany?.id === company.id}
                    type="button"
                    onclick={() => selectCompany(company)}
                  >
                    <div class="company-avatar">
                      {#if company.logo_url}
                        <img
                          src={company.logo_url}
                          alt={company.company_name ?? ""}
                        />
                      {:else}
                        <Building2 size={24} />
                      {/if}
                    </div>
                    <div class="company-info">
                      <span class="company-name">{company.company_name}</span>
                      {#if company.sector}
                        <span class="company-sector">{company.sector}</span>
                      {/if}
                    </div>
                    {#if selectedCompany?.id === company.id}
                      <div class="check-icon">
                        <CheckCircle size={20} />
                      </div>
                    {/if}
                  </button>
                {/each}
              </div>
            {:else}
              <div class="empty-state">
                <Building2 size={40} />
                <p>
                  Aucune entreprise trouvée{companySearch
                    ? ` pour "${companySearch}"`
                    : ""}
                </p>
              </div>
            {/if}

            <div class="divider-row">
              <div class="divider"></div>
              <span>ou</span>
              <div class="divider"></div>
            </div>

            <button
              class="create-company-btn"
              type="button"
              onclick={goToCreateCompany}
            >
              <Plus size={18} />
              Créer mon entreprise
            </button>

            <div class="step-actions">
              <Button onclick={() => (step = 1)} Icon={ChevronLeft}>
                Retour
              </Button>
              <Button
                onclick={handleSubmit}
                disabled={!selectedCompany || loading}
                Icon={UserPlus}
              >
                {loading ? "Inscription..." : "S'inscrire"}
              </Button>
            </div>
          </div>
        {/if}

        <!-- Step 3: Create New Company -->
        {#if step === 3}
          <div class="form-step" data-aos="fade-left">
            <div class="section-header">
              <Building2 size={20} />
              <h2>Créer votre entreprise</h2>
            </div>

            <div class="info-banner">
              <p>
                Votre entreprise sera soumise à une <strong
                  >validation par l'administration</strong
                > avant d'être active sur la plateforme.
              </p>
            </div>

            <div class="form-grid">
              <div class="form-group full">
                <Input
                  name="companyName"
                  label="Nom de l'entreprise"
                  placeholder="Ex: Sonatrach, Cevital..."
                  required
                  bind:value={companyName}
                />
              </div>

              <!-- Logo Upload -->
              <div class="form-group full">
                <label class="logo-label"
                  >Logo de l'entreprise <span class="optional">(optionnel)</span
                  ></label
                >
                {#if logoPreview}
                  <div class="logo-preview-wrap">
                    <img
                      src={logoPreview}
                      alt="Aperçu logo"
                      class="logo-preview"
                    />
                    <div class="logo-preview-info">
                      <span class="logo-filename">{logoFile?.name}</span>
                      <button
                        class="logo-clear"
                        type="button"
                        onclick={clearLogo}
                      >
                        <X size={14} /> Supprimer
                      </button>
                    </div>
                  </div>
                {:else}
                  <label class="logo-upload-zone" for="logo-input">
                    <ImagePlus size={28} />
                    <span>Cliquez pour choisir un logo</span>
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

              <div class="form-group">
                <Input
                  name="sector"
                  label="Domaine d'activité"
                  placeholder="Ex: Pétrolier, IT, Agroalimentaire..."
                  bind:value={sector}
                />
              </div>

              <div class="form-group">
                <Input
                  name="contactEmail"
                  type="email"
                  label="Email de l'entreprise"
                  placeholder="contact@entreprise.com"
                  bind:value={contactEmail}
                />
              </div>

              <div class="form-group">
                <Input
                  name="contactPhone"
                  type="tel"
                  label="Téléphone de l'entreprise"
                  placeholder="+213 XXX XX XX XX"
                  bind:value={contactPhone}
                />
              </div>

              <div class="form-group full">
                <Input
                  name="description"
                  category="textarea"
                  label="Description"
                  placeholder="Décrivez brièvement l'activité de votre entreprise..."
                  bind:value={description}
                />
              </div>
            </div>

            <div class="step-actions">
              <Button onclick={() => (step = 2)} Icon={ChevronLeft}>
                Retour
              </Button>
              <Button
                onclick={handleSubmit}
                disabled={!step3Valid || loading}
                Icon={UserPlus}
              >
                {loading ? "Inscription..." : "Créer et s'inscrire"}
              </Button>
            </div>
          </div>
        {/if}

        <!-- Footer -->
        <div class="form-footer">
          <p>
            Déjà un compte ?
            <a href="./login">Se connecter</a>
          </p>
          <p>
            <a href="../login" class="intern-link">
              <GraduationCap size={16} />
              Je suis un interne {BRAND.university.name}
            </a>
          </p>
        </div>
      </div>
    {/if}
  </div>
</View>

<style>
  .registration-container {
    width: 100%;
    max-width: 640px;
    padding: 2rem 0;
  }

  /* ─── Success Card ─── */
  .success-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 1.5rem;
    padding: 3rem 2rem;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 1.5rem;
  }

  .success-icon {
    color: #22c55e;
    animation: scaleIn 0.5s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  }

  @keyframes scaleIn {
    from {
      transform: scale(0);
      opacity: 0;
    }
    to {
      transform: scale(1);
      opacity: 1;
    }
  }

  .success-card h1 {
    font-size: 1.6rem;
    font-weight: 700;
    color: var(--color-text);
  }

  .success-msg {
    color: var(--color-text-muted);
    font-size: 1rem;
    line-height: 1.6;
    max-width: 400px;
  }

  .success-actions {
    margin-top: 0.5rem;
  }

  /* ─── Form Card ─── */
  .form-card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 1.5rem;
    padding: 2.5rem;
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .card-header {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 0.75rem;
  }

  .title-row {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    color: var(--color-accent);
  }

  .title-row h1 {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--color-text);
    margin: 0;
  }

  .subtitle {
    color: var(--color-text-muted);
    font-size: 0.95rem;
    margin: 0;
  }

  /* ─── Stepper ─── */
  .stepper {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0;
  }

  .step-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 0.75rem;
    border-radius: 2rem;
    border: none;
    background: none;
    cursor: pointer;
    transition: all 0.3s ease;
    font-family: var(--font-sans);
    font-size: 0.85rem;
    color: var(--color-text-muted);
  }

  .step-item.active {
    color: var(--color-accent);
  }

  .step-item.current {
    background: color-mix(in srgb, var(--color-accent) 10%, transparent);
  }

  .step-number {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    border-radius: 50%;
    font-weight: 700;
    font-size: 0.8rem;
    background: var(--color-background-100);
    border: 2px solid var(--color-border);
    transition: all 0.3s ease;
  }

  .step-item.active .step-number {
    background: var(--color-accent);
    border-color: var(--color-accent);
    color: white;
  }

  .step-label {
    font-weight: 500;
  }

  .step-line {
    width: 40px;
    height: 2px;
    background: var(--color-border);
    transition: all 0.3s ease;
    flex-shrink: 0;
  }

  .step-line.active {
    background: var(--color-accent);
  }

  /* ─── Section Header ─── */
  .section-header {
    margin-bottom: 1.25rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--color-text);
    padding-bottom: 0.75rem;
    border-bottom: 1px solid var(--color-border);
  }

  .section-header h2 {
    font-size: 1.1rem;
    font-weight: 600;
    margin: 0;
  }

  /* ─── Form Grid ─── */
  .form-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  .form-group.full {
    grid-column: 1 / -1;
  }

  .position-group {
    position: relative;
  }

  /* ─── Autocomplete ─── */
  .autocomplete-dropdown {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    z-index: 50;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 0.75rem;
    overflow: hidden;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    margin-top: 4px;
  }

  .autocomplete-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    width: 100%;
    padding: 0.65rem 1rem;
    border: none;
    background: none;
    color: var(--color-text);
    font-family: var(--font-sans);
    font-size: 0.9rem;
    cursor: pointer;
    text-align: left;
    transition: background 0.15s;
  }

  .autocomplete-item:hover {
    background: var(--color-background-100);
  }

  .autocomplete-item:not(:last-child) {
    border-bottom: 1px solid var(--color-border);
  }

  /* ─── Search Box ─── */
  .search-box {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0 1rem;
    background: var(--color-background-100);
    border: 1px solid var(--color-border);
    border-radius: 1rem;
    transition: border-color 0.2s;
    color: var(--color-text-muted);
  }

  .search-box:focus-within {
    border-color: var(--color-accent);
  }

  .search-input {
    flex: 1;
    border: none;
    background: none;
    padding: 0.85rem 0;
    font-family: var(--font-sans);
    font-size: 1rem;
    color: var(--color-text);
    outline: none;
  }

  .search-input::placeholder {
    color: var(--color-text-muted);
  }

  /* ─── Company List ─── */
  .company-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    max-height: 280px;
    overflow-y: auto;
    padding-right: 0.25rem;
  }

  .company-card {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    border: 1px solid var(--color-border);
    border-radius: 1rem;
    cursor: pointer;
    transition: all 0.2s ease;
    background: var(--color-surface);
    text-align: left;
    font-family: var(--font-sans);
    width: 100%;
  }

  .company-card:hover {
    border-color: var(--color-accent);
    background: color-mix(
      in srgb,
      var(--color-accent) 3%,
      var(--color-surface)
    );
  }

  .company-card.selected {
    border-color: var(--color-accent);
    background: color-mix(
      in srgb,
      var(--color-accent) 8%,
      var(--color-surface)
    );
    box-shadow: 0 0 0 3px
      color-mix(in srgb, var(--color-accent) 15%, transparent);
  }

  .company-avatar {
    width: 48px;
    height: 48px;
    border-radius: 0.75rem;
    background: var(--color-background-100);
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    color: var(--color-text-muted);
    overflow: hidden;
  }

  .company-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .company-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }

  .company-name {
    font-weight: 600;
    color: var(--color-text);
    font-size: 0.95rem;
  }

  .company-sector {
    font-size: 0.8rem;
    color: var(--color-text-muted);
  }

  .check-icon {
    color: var(--color-accent);
    flex-shrink: 0;
    animation: scaleIn 0.3s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  }

  /* ─── Empty/Loading States ─── */
  .empty-state,
  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.75rem;
    padding: 2rem;
    color: var(--color-text-muted);
    text-align: center;
  }

  .spinner {
    width: 32px;
    height: 32px;
    border: 3px solid var(--color-border);
    border-top-color: var(--color-accent);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }

  /* ─── Divider Row ─── */
  .divider-row {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin: 0.5rem 0;
  }

  .divider {
    flex: 1;
    height: 1px;
    background: var(--color-border);
  }

  .divider-row span {
    color: var(--color-text-muted);
    font-size: 0.85rem;
    font-weight: 500;
  }

  /* ─── Create Company Button ─── */
  .create-company-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    width: 100%;
    padding: 0.9rem;
    border: 2px dashed var(--color-border);
    border-radius: 1rem;
    background: none;
    color: var(--color-accent);
    font-family: var(--font-sans);
    font-size: 0.95rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .create-company-btn:hover {
    border-color: var(--color-accent);
    background: color-mix(in srgb, var(--color-accent) 5%, transparent);
  }

  /* ─── Info Banner ─── */
  .info-banner {
    background: color-mix(in srgb, var(--color-accent) 8%, transparent);
    border: 1px solid color-mix(in srgb, var(--color-accent) 25%, transparent);
    border-radius: 0.75rem;
    padding: 0.85rem 1rem;
    margin-bottom: 1.25rem;
  }

  .info-banner p {
    color: var(--color-text);
    font-size: 0.9rem;
    margin: 0;
    line-height: 1.5;
  }

  /* ─── Error Banner ─── */
  .error-banner {
    background: color-mix(in srgb, #ef4444 10%, transparent);
    border: 1px solid color-mix(in srgb, #ef4444 30%, transparent);
    border-radius: 0.75rem;
    padding: 0.85rem 1rem;
    text-align: center;
  }

  .error-banner p {
    color: #ef4444;
    font-size: 0.9rem;
    margin: 0;
  }

  /* ─── Step Actions ─── */
  .step-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 1.25rem;
  }

  /* ─── Footer ─── */
  .form-footer {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    padding-top: 0.5rem;
    border-top: 1px solid var(--color-border);
  }

  .form-footer p {
    font-size: 0.9rem;
    color: var(--color-text-muted);
    margin: 0;
  }

  .form-footer a {
    color: var(--color-accent);
    font-weight: 600;
    text-decoration: none;
  }

  .form-footer a:hover {
    text-decoration: underline;
  }

  .intern-link {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
  }

  /* ─── Logo Upload ─── */
  .logo-label {
    display: block;
    font-size: 0.875rem;
    font-weight: 500;
    color: var(--color-text);
    margin-bottom: 0.5rem;
  }

  .optional {
    font-weight: 400;
    color: var(--color-text-muted);
    font-size: 0.8rem;
  }

  .logo-upload-zone {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 1.5rem;
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

  /* ─── Responsive ─── */
  @media (max-width: 600px) {
    .form-card {
      padding: 1.5rem;
      border-radius: 1rem;
    }

    .form-grid {
      grid-template-columns: 1fr;
    }

    .stepper {
      flex-wrap: wrap;
      gap: 0.25rem;
    }

    .step-label {
      display: none;
    }

    .title-row h1 {
      font-size: 1.2rem;
    }
  }
</style>
