<script lang="ts">
  import {
    Users,
    BookOpen,
    Download,
    Briefcase,
    CheckCircle,
    MoreVertical,
    AlertTriangle,
  } from "lucide-svelte";

  import Button from "$lib/components/ui/Button.svelte";
  import StatCard from "$lib/components/admin/StatCard.svelte";
  import DataTable from "$lib/components/admin/DataTable.svelte";
  import QuickActions from "$lib/components/admin/QuickActions.svelte";
  import DashboardPage from "$lib/components/admin/DashboardPage.svelte";
  import CycleProgress from "$lib/components/admin/CycleProgress.svelte";

  import type { ProgressStep } from "$lib/components/admin/CycleProgress.svelte";
  import type { IconComponent } from "$lib/types/domain";

  let { data } = $props();

  const { stats, activeYear, recentAuditLogs } = $derived(data);

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

  const tableColumns = [
    { key: "description", label: "Action" },
    { key: "actor", label: "Acteur" },
    { key: "created_at", label: "Date" },
    {
      key: "action",
      label: "",
      width: "60px",
      align: "center" as const,
    },
  ];

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
</script>

<DashboardPage
  title="Tableau de bord"
  subtitle="Bienvenue sur la plateforme PFE. Voici un apercu de l'activite en cours."
>
  {#snippet actions()}
    <Button variant="ghost" Icon={Download} id="admin-export-btn">
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

  <div class="dashboard-middle">
    <section>
      <DataTable
        columns={tableColumns}
        rows={recentAuditLogs.slice(0, 5)}
        title="Activite recente"
        viewAllHref="/admin/audit-log"
      >
        {#snippet cell({ row, col })}
          {#if col.key === "description"}
            <span>{row.description ?? "-"}</span>
          {:else if col.key === "actor"}
            <span>{row.actor ?? "-"}</span>
          {:else if col.key === "created_at"}
            <span>{row.created_at_formatted ?? "-"}</span>
          {:else if col.key === "action"}
            <Button
              variant="ghost"
              square
              Icon={MoreVertical}
              aria-label="Details"
            />
          {:else}
            {row[col.key] ?? "-"}
          {/if}
        {/snippet}
      </DataTable>
    </section>

    <aside>
      <QuickActions />
    </aside>
  </div>

  <section>
    <h2>Progression du cycle PFE</h2>
    <CycleProgress steps={cycleSteps} />
  </section>
</DashboardPage>

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
</style>
