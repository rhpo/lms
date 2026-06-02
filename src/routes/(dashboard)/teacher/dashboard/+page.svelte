<script lang="ts">
  import { goto } from "$app/navigation";
  import { Users, BookOpen, CheckCircle, Calendar, Bell } from "lucide-svelte";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const {
    supervisedCount,
    pendingValidationCount,
    proposedCount,
    upcomingJuryDuties,
    availabilityStatus,
    unavailableUntil,
    unreadCount,
  } = $derived(data);

  const availabilityLabel = $derived.by(() => {
    switch (availabilityStatus) {
      case "disponible":
        return "Disponible";
      case "indisponible":
        return "Indisponible";
      case "indisponible_jusqu_au": {
        if (unavailableUntil) {
          const d = new Date(unavailableUntil).toLocaleDateString("fr-FR");
          return `Indisponible jusqu'au ${d}`;
        }
        return "Indisponible";
      }
      default:
        return "Disponible";
    }
  });

  const availabilityVariant = $derived.by<
    "success" | "danger" | "warning" | "info" | "neutral"
  >(() => {
    switch (availabilityStatus) {
      case "disponible":
        return "success";
      case "indisponible":
        return "danger";
      case "indisponible_jusqu_au":
        return "warning";
      default:
        return "neutral";
    }
  });
</script>

<Page title="Tableau de bord" subtitle="Bienvenue sur votre espace enseignant.">
  <div class="stat-grid">
    <div
      class="stat-card"
      onclick={() => goto("/teacher/supervised-pfes")}
      role="button"
      tabindex="0"
    >
      <div class="stat-icon">
        <Users size={20} />
      </div>
      <p class="stat-label">Etudiants encadres</p>
      <p class="stat-value">{supervisedCount}</p>
    </div>

    <div
      class="stat-card"
      onclick={() => goto("/teacher/subjects-to-validate")}
      role="button"
      tabindex="0"
    >
      <div class="stat-icon">
        <CheckCircle size={20} />
      </div>
      <p class="stat-label">Sujets a valider</p>
      <p class="stat-value">{pendingValidationCount}</p>
    </div>

    <div
      class="stat-card"
      onclick={() => goto("/teacher/proposed-subjects")}
      role="button"
      tabindex="0"
    >
      <div class="stat-icon">
        <BookOpen size={20} />
      </div>
      <p class="stat-label">Sujets proposes</p>
      <p class="stat-value">{proposedCount}</p>
    </div>

    <div
      class="stat-card"
      onclick={() => goto("/teacher/availability")}
      role="button"
      tabindex="0"
    >
      <div class="stat-icon">
        <Calendar size={20} />
      </div>
      <p class="stat-label">Disponibilité</p>
      <Badge variant={availabilityVariant} label={availabilityLabel} />
    </div>
  </div>

  <div class="dashboard-grid">
    <section>
      <h2>Prochaines soutenances</h2>
      {#if upcomingJuryDuties.length > 0}
        <div class="list">
          {#each upcomingJuryDuties as duty}
            <div
              class="list-item"
              onclick={() => goto(`/teacher/jury-duties`)}
              role="button"
              tabindex="0"
            >
              <div class="list-item-icon">
                <Calendar size={16} />
              </div>
              <div>
                <p class="list-item-title">{duty.scheduled_at}</p>
                <p class="list-item-subtitle">Salle {duty.room}</p>
              </div>
            </div>
          {/each}
        </div>
        <div class="section-actions">
          <Button variant="ghost" onclick={() => goto("/teacher/jury-duties")}>
            Voir toutes mes soutenances
          </Button>
        </div>
      {:else}
        <p class="empty">Aucune soutenance planifiee</p>
      {/if}
    </section>

    <section>
      <h2>Notifications</h2>
      {#if unreadCount > 0}
        <div class="unread-banner">
          <Bell size={16} />
          <span>{unreadCount} notification(s) non lue(s)</span>
        </div>
      {:else}
        <p class="empty">Aucune notification non lue</p>
      {/if}
      <div class="section-actions">
        <Button variant="ghost" onclick={() => goto("/teacher/notifications")}>
          Voir toutes les notifications
        </Button>
      </div>
    </section>
  </div>

  <div class="quick-actions">
    <h2>Raccourcis</h2>
    <div class="actions-grid">
      <Button
        variant="primary"
        onclick={() => goto("/teacher/proposed-subjects")}
      >
        Mes sujets
      </Button>
      <Button
        variant="primary"
        onclick={() => goto("/teacher/subjects-to-validate")}
      >
        Sujets a valider
      </Button>
      <Button
        variant="primary"
        onclick={() => goto("/teacher/supervised-pfes")}
      >
        Mes encadrements
      </Button>
      <Button variant="primary" onclick={() => goto("/teacher/availability")}>
        Modifier ma disponibilité
      </Button>
    </div>
  </div>
</Page>

<style>
  .stat-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 1rem;
    margin-bottom: var(--spacing-xl);

    @media screen and (max-width: 768px) {
      & {
        grid-template-columns: repeat(2, 1fr);
      }
    }

    @media screen and (max-width: 480px) {
      & {
        grid-template-columns: 1fr;
      }
    }
  }

  .stat-card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: 1.25rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    cursor: pointer;
    transition: box-shadow var(--transition-normal);

    &:hover {
      box-shadow: var(--shadow-md);
    }
  }

  .stat-icon {
    width: 38px;
    height: 38px;
    border-radius: 8px;
    background: var(--color-background);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--color-text-muted);
  }

  .stat-label {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    font-family: var(--font-sans);
    margin: 0;
  }

  .stat-value {
    font-size: 1.9rem;
    font-weight: 700;
    color: var(--color-text);
    font-family: var(--font-sans);
    margin: 0;
    line-height: 1;
  }

  .dashboard-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-lg);
    margin-bottom: var(--spacing-xl);

    @media screen and (max-width: 768px) {
      & {
        grid-template-columns: 1fr;
      }
    }
  }

  section {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: var(--spacing-lg);
  }

  h2 {
    font-size: var(--text-lg);
    font-weight: 600;
    font-family: var(--font-sans);
    color: var(--color-text);
    margin: 0 0 var(--spacing-md);
  }

  .list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .list-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem;
    border-radius: 8px;
    cursor: pointer;
    transition: background var(--transition-fast);

    &:hover {
      background: var(--color-background);
    }
  }

  .list-item-icon {
    width: 32px;
    height: 32px;
    border-radius: 6px;
    background: var(--color-background);
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--color-accent);
    flex-shrink: 0;
  }

  .list-item-title {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text);
    margin: 0;
  }

  .list-item-subtitle {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
    margin: 0.15rem 0 0;
  }

  .empty {
    text-align: center;
    color: var(--color-text-muted);
    font-style: italic;
    font-size: var(--text-sm);
    padding: 1rem;
  }

  .unread-banner {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem;
    background: color-mix(
      in srgb,
      var(--color-accent) 10%,
      var(--color-surface)
    );
    border-radius: 8px;
    color: var(--color-accent);
    font-size: var(--text-sm);
    font-weight: 600;
  }

  .section-actions {
    margin-top: var(--spacing-md);
  }

  .quick-actions {
    margin-bottom: var(--spacing-xl);
  }

  .actions-grid {
    display: flex;
    flex-wrap: wrap;
    gap: var(--spacing-sm);
  }
</style>
