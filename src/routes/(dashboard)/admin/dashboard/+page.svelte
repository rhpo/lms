<script lang="ts">
  import {
    Users,
    BookOpen,
    Download,
    Briefcase,
    CheckCircle,
    AlertTriangle,
    FileSpreadsheet,
    FileText,
  } from "lucide-svelte";

  import Button from "$lib/components/ui/Button.svelte";
  import Modal from "$lib/components/ui/Modal.svelte";
  import StatCard from "$lib/components/admin/StatCard.svelte";
  import QuickActions from "$lib/components/admin/QuickActions.svelte";
  import DashboardPage from "$lib/components/admin/DashboardPage.svelte";
  import CycleProgress from "$lib/components/admin/CycleProgress.svelte";
  import DashboardCharts from "$lib/components/admin/DashboardCharts.svelte";

  import type { ProgressStep } from "$lib/components/admin/CycleProgress.svelte";
  import type { IconComponent } from "$lib/types/domain";
  import { theme } from "$lib/stores/theme.js";

  let { data } = $props();

  const { stats, activeYear, timeline } = $derived(data);

  interface StatCardData {
    title: string;
    value: number;
    badge?: string;
    badgeVariant: "green" | "orange" | "purple" | "blue";
    icon: IconComponent;
    barColor: string;
  }

  const statCards: StatCardData[] = $derived([
    {
      title: "Etudiants",
      value: stats.totalStudents,
      badge:
        stats.totalAssignments > 0
          ? `${stats.totalAssignments} affectes`
          : undefined,
      badgeVariant: "blue",
      icon: Users,
      barColor: "#3b82f6",
    },
    {
      title: "Sujets",
      value: stats.totalSubjects,
      badge: `${stats.pendingSubjects} en attente`,
      badgeVariant: "orange",
      icon: BookOpen,
      barColor: "#f59e0b",
    },
    {
      title: "Sujets valides",
      value: stats.validatedSubjects,
      badge:
        stats.totalSubjects > 0
          ? `${Math.round((stats.validatedSubjects / stats.totalSubjects) * 100)}% taux`
          : undefined,
      badgeVariant: "green",
      icon: CheckCircle,
      barColor: "#10b981",
    },
    {
      title: "Enseignants",
      value: stats.totalTeachers,
      badge: `${stats.totalDefenses} soutenances`,
      badgeVariant: "purple",
      icon: Briefcase,
      barColor: "#8b5cf6",
    },
  ]);


  const cycleSteps: ProgressStep[] = $derived([
    {
      label: "Depot des sujets",
      subtitle: activeYear?.submission_open_at
        ? `Ouvert le ${activeYear.submission_open_at_formatted}`
        : "Non configure",
      status:
        activeYear?.submission_open_at &&
        new Date(activeYear.submission_open_at) <= new Date()
          ? "completed"
          : "upcoming",
    },
    {
      label: "Validation",
      subtitle: activeYear?.submission_close_at
        ? `Clôture le ${activeYear.submission_close_at_formatted}`
        : "En cours",
      status: "active",
    },
    {
      label: "Soutenances",
      subtitle: `${stats.totalDefenses} planifiees`,
      status: stats.totalDefenses > 0 ? "active" : "upcoming",
    },
    {
      label: "Clôture",
      subtitle: activeYear?.submission_close_at
        ? `Ferme le ${activeYear.submission_close_at_formatted}`
        : "Non definie",
      status: "upcoming",
    },
  ]);

  let exportModalOpen = $state(false);

  function exportToCSV() {
    const csvContent =
      "Indicateur,Valeur\n" +
      `Etudiants Total,${stats.totalStudents || 0}\n` +
      `Enseignants Total,${stats.totalTeachers || 0}\n` +
      `Entreprises Total,${stats.totalCompanies || 0}\n` +
      `Entreprises en Attente,${stats.pendingCompanies || 0}\n` +
      `Sujets PFE Totaux,${stats.totalSubjects || 0}\n` +
      `Sujets Valides,${stats.validatedSubjects || 0}\n` +
      `Sujets Rejetes,${stats.rejectedSubjects || 0}\n` +
      `Sujets en Attente,${stats.pendingSubjects || 0}\n` +
      `Affectations PFE,${stats.totalAssignments || 0}\n` +
      `Soutenances Planifiees,${stats.totalDefenses || 0}`;

    const blob = new Blob([csvContent], { type: "text/csv;charset=utf-8;" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.setAttribute("href", url);
    link.setAttribute(
      "download",
      `dashboard_PFE_${new Date().toISOString().split("T")[0]}.csv`,
    );
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    exportModalOpen = false;
  }

  function exportToExcel() {
    const htmlContent = `
      <html xmlns:x="urn:schemas-microsoft-com:office:excel">
      <head><meta charset="UTF-8"></head>
      <body>
        <table border="1">
          <tr><th>Indicateur</th><th>Valeur</th></tr>
          <tr><td>Etudiants Total</td><td>${stats.totalStudents || 0}</td></tr>
          <tr><td>Enseignants Total</td><td>${stats.totalTeachers || 0}</td></tr>
          <tr><td>Entreprises Total</td><td>${stats.totalCompanies || 0}</td></tr>
          <tr><td>Entreprises en Attente</td><td>${stats.pendingCompanies || 0}</td></tr>
          <tr><td>Sujets PFE Totaux</td><td>${stats.totalSubjects || 0}</td></tr>
          <tr><td>Sujets Valides</td><td>${stats.validatedSubjects || 0}</td></tr>
          <tr><td>Sujets Rejetes</td><td>${stats.rejectedSubjects || 0}</td></tr>
          <tr><td>Sujets en Attente</td><td>${stats.pendingSubjects || 0}</td></tr>
          <tr><td>Affectations PFE</td><td>${stats.totalAssignments || 0}</td></tr>
          <tr><td>Soutenances Planifiees</td><td>${stats.totalDefenses || 0}</td></tr>
        </table>
      </body>
      </html>
    `;

    const blob = new Blob([htmlContent], { type: "application/vnd.ms-excel" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.setAttribute("href", url);
    link.setAttribute(
      "download",
      `dashboard_PFE_${new Date().toISOString().split("T")[0]}.xls`,
    );
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    exportModalOpen = false;
  }
</script>

<DashboardPage
  title="Tableau de bord"
  subtitle="Bienvenue sur la plateforme PFE. Voici un apercu de l'activite en cours."
>
  {#snippet actions()}
    <Button
      variant="ghost"
      Icon={Briefcase}
      href="/teacher/dashboard"
    >
      Dashboard Enseignant
    </Button>
    <Button
      variant="ghost"
      Icon={Download}
      id="admin-export-btn"
      onclick={() => (exportModalOpen = true)}
    >
      Exporter
    </Button>
  {/snippet}

  <div class="stat-grid">
    {#each statCards as card}
      <StatCard {...card} />
    {/each}
  </div>

  {#if stats.pendingCompanies > 0}
    <div class="alert-banner">
      <AlertTriangle size={18} />
      <span
        >{stats.pendingCompanies} entreprise(s) en attente de validation.</span
      >
      <a href="/admin/users">Gérer les entreprises</a>
    </div>
  {/if}

  <section data-aos="fade-up">
    <h2>Progression du cycle PFE</h2>
    <CycleProgress steps={cycleSteps} />
  </section>

  <DashboardCharts {stats} {timeline} theme={$theme} />

  <div class="dashboard-middle">
    <!-- Journal d'activité masqué -->

    <aside>
      <QuickActions />
    </aside>
  </div>
</DashboardPage>

<Modal
  bind:open={exportModalOpen}
  title="Exporter le tableau de bord"
  description="Choisissez le format du fichier pour l'export. Les donnees contiendront les metriques globales."
  width="450px"
>
  <div class="export-options">
    <button class="export-btn excel" onclick={exportToExcel}>
      <!-- <FileSpreadsheet size={36} /> -->
      <img src="/media/icons/excel.png" alt="Excel" width={36} />

      <div class="export-text">
        <strong>Export Excel</strong>
        <span>Fichier .xls compatible Microsoft</span>
      </div>
    </button>
    <button class="export-btn csv" onclick={exportToCSV}>
      <!-- <FileText size={36} /> -->
      <img src="/media/icons/csv.png" alt="CSV" width={36} />

      <div class="export-text">
        <strong>Export CSV</strong>
        <span>Fichier texte standard separateur virgule</span>
      </div>
    </button>
  </div>
</Modal>

<style>
  .stat-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1rem;
  }

  .alert-banner {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    background: color-mix(in srgb, #f59e0b 10%, var(--color-surface));
    border: 1px solid color-mix(in srgb, #f59e0b 20%, transparent);
    border-radius: 10px;
    color: #92400e;
    font-family: var(--font-sans);
    font-size: 0.9rem;

    a {
      margin-left: auto;
      color: var(--color-accent);
      font-weight: 600;
      text-decoration: none;

      &:hover {
        text-decoration: underline;
      }
    }
  }

  :global(html.dark) .alert-banner {
    color: #fcd34d;
  }

  .dashboard-middle {
    display: grid;
    grid-template-columns: 1fr 340px;
    gap: 1rem;
    align-items: start;
  }

  section {
    min-width: 0;
  }

  h2 {
    font-size: var(--text-xl);
    font-weight: 600;
    font-family: var(--font-sans);
    color: var(--color-text);
    margin: 0 0 1rem;
  }

  span {
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    color: var(--color-text);
  }

  @media (max-width: 1100px) {
    .dashboard-middle {
      grid-template-columns: 1fr;
    }
  }

  @media (max-width: 768px) {
    .stat-grid {
      grid-template-columns: repeat(2, 1fr);
    }
  }

  @media (max-width: 480px) {
    .stat-grid {
      grid-template-columns: 1fr;
    }
  }

  .export-options {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-top: 10px;
  }

  .export-btn {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
    background: var(--color-background-100);
    border: 2px solid var(--color-border);
    border-radius: 12px;
    cursor: pointer;
    text-align: left;
    transition: all 0.2s ease;
    color: var(--color-text);
  }

  .export-btn:hover {
    background: var(--color-background-200);
    transform: translateY(-2px);
  }

  .export-btn.excel:hover {
    border-color: #10b981;
    color: #10b981;
  }

  .export-btn.csv:hover {
    border-color: #3b82f6;
    color: #3b82f6;
  }

  .export-text {
    display: flex;
    flex-direction: column;
  }

  .export-text strong {
    font-size: 1.1rem;
    font-family: var(--font-sans, system-ui);
  }

  .export-text span {
    font-size: 0.85rem;
    color: var(--color-text-muted);
    font-family: var(--font-sans, system-ui);
    margin-top: 2px;
  }
</style>
