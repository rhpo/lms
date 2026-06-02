<script lang="ts">
  import {
    ArrowLeft,
    Calendar,
    MapPin,
    User,
    GraduationCap,
    CheckCircle,
  } from "lucide-svelte";
  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();
  const defense = $derived(data.defense);

  const statusVariant: Record<
    string,
    "success" | "warning" | "danger" | "info" | "neutral"
  > = {
    scheduled: "warning",
    done: "success",
    postponed: "danger",
  };

  const statusLabel: Record<string, string> = {
    scheduled: "Planifiée",
    done: "Terminée",
    postponed: "Reportée",
  };

  const resultVariant: Record<
    string,
    "success" | "warning" | "danger" | "info" | "neutral"
  > = {
    admitted: "success",
    corrections_required: "warning",
    not_admitted: "danger",
  };

  const resultLabel: Record<string, string> = {
    admitted: "Admis",
    corrections_required: "Corrections requises",
    not_admitted: "Ajourné",
  };

  function formatDate(d: string | Date | null | undefined) {
    if (!d) return "-";
    return new Date(d).toLocaleString("fr-FR", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }
</script>

{#if defense}
  <Page
    title="Détails de la Soutenance"
    subtitle={`PFE: ${defense.assignment?.subject?.title ?? "Sujet inconnu"}`}
  >
    {#snippet actions()}
      <a href="/admin/pfe">
        <Button variant="ghost" Icon={ArrowLeft}>Retour</Button>
      </a>
    {/snippet}

    <div class="grid-layout">
      <!-- Info Principales -->
      <section class="card">
        <h2>Plannification</h2>
        <div class="details-list">
          <div class="detail-item">
            <Calendar class="icon" size={18} />
            <div>
              <span class="label">Date et heure</span>
              <span class="value">{formatDate(defense.scheduled_at)}</span>
            </div>
          </div>
          <div class="detail-item">
            <MapPin class="icon" size={18} />
            <div>
              <span class="label">Salle</span>
              <span class="value">{defense.room || "-"}</span>
            </div>
          </div>
          <div class="detail-item">
            <CheckCircle class="icon" size={18} />
            <div>
              <span class="label">Statut</span>
              <span class="value">
                <Badge
                  variant={statusVariant[defense.status] ?? "neutral"}
                  label={statusLabel[defense.status] ?? defense.status}
                />
              </span>
            </div>
          </div>
        </div>
      </section>

      <!-- Assignation PFE -->
      <section class="card">
        <h2>Informations du PFE</h2>
        <div class="details-list">
          <div class="detail-item">
            <GraduationCap class="icon" size={18} />
            <div>
              <span class="label">Étudiant</span>
              <span class="value">
                {defense.assignment?.student?.profile?.full_name ?? "-"}
                {#if defense.assignment?.student2}
                  , {defense.assignment.student2.profile?.full_name}
                {/if}
                {#if defense.assignment?.student3}
                  , {defense.assignment.student3.profile?.full_name}
                {/if}
              </span>
            </div>
          </div>
          <div class="detail-item">
            <User class="icon" size={18} />
            <div>
              <span class="label">Encadrant</span>
              <span class="value"
                >{defense.assignment?.supervisor?.profile?.full_name ??
                  "-"}</span
              >
            </div>
          </div>
          {#if defense.assignment?.co_supervisor_id}
            <div class="detail-item">
              <User class="icon" size={18} />
              <div>
                <span class="label">Co-encadrant</span>
                <span class="value"
                  >{defense.assignment?.co_supervisor?.profile?.full_name ??
                    "-"}</span
                >
              </div>
            </div>
          {/if}
        </div>
      </section>

      <!-- Informations Jury -->
      <section class="card full-width">
        <h2>Jury de Soutenance</h2>
        {#if defense.jury}
          <div class="grid-2">
            <div class="jury-member">
              <h3>Président du Jury</h3>
              <p>{defense.jury.president?.profile?.full_name ?? "-"}</p>
              <Badge
                variant={defense.jury.president_confirmed
                  ? "success"
                  : "warning"}
                label={defense.jury.president_confirmed
                  ? "Confirmé"
                  : "En attente"}
              />
            </div>

            <div class="jury-member">
              <h3>Examinateur</h3>
              <p>{defense.jury.member?.profile?.full_name ?? "-"}</p>
              <Badge
                variant={defense.jury.member_confirmed ? "success" : "warning"}
                label={defense.jury.member_confirmed
                  ? "Confirmé"
                  : "En attente"}
              />
            </div>
          </div>
        {:else}
          <p class="text-muted">Aucun jury n'est assigné à cette soutenance.</p>
        {/if}
      </section>

      <!-- Résultats -->
      {#if defense.status === "done" && defense.result}
        <section class="card full-width outcome">
          <h2>Résultat de la Soutenance</h2>
          <div class="result-display">
            <div class="decision">
              <span class="label">Décision du jury :</span>
              <Badge
                variant={resultVariant[defense.result] ?? "neutral"}
                label={resultLabel[defense.result] ?? defense.result}
              />
            </div>
            {#if defense.final_grade}
              <div class="grade">
                <span class="label">Note finale :</span>
                <span class="grade-value"
                  >{defense.final_grade.toFixed(2)} / 20</span
                >
              </div>
            {/if}
          </div>
        </section>
      {/if}
    </div>
  </Page>
{:else}
  <Page title="Soutenance introuvable" subtitle="">
    <p>Cette soutenance n'existe pas ou a été supprimée.</p>
    <a href="/admin/pfe">
      <Button variant="ghost" Icon={ArrowLeft}>Retour</Button>
    </a>
  </Page>
{/if}

<style>
  .grid-layout {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-lg);
  }

  .card {
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    padding: var(--spacing-lg);
    box-shadow: 0 1px 3px rgb(0 0 0 / 0.1);
  }

  .full-width {
    grid-column: 1 / -1;
  }

  h2 {
    font-family: var(--font-sans);
    font-size: var(--text-lg);
    font-weight: 700;
    color: var(--color-text);
    margin: 0 0 var(--spacing-md);
    padding-bottom: var(--spacing-sm);
    border-bottom: 1px solid var(--color-border);
  }

  .details-list {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-md);
  }

  .detail-item {
    display: flex;
    align-items: flex-start;
    gap: var(--spacing-md);
  }

  .icon {
    margin-top: 2px;
    color: var(--color-primary);
  }

  .label {
    display: block;
    font-family: var(--font-sans);
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-muted);
    margin-bottom: 4px;
  }

  .value {
    display: block;
    font-family: var(--font-sans);
    font-size: var(--text-md);
    color: var(--color-text);
  }

  .grid-2 {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--spacing-lg);
  }

  .jury-member {
    background: var(--color-background-100);
    padding: var(--spacing-md);
    border-radius: 8px;
    border: 1px solid var(--color-border);
  }

  .jury-member h3 {
    font-size: var(--text-sm);
    color: var(--color-text-muted);
    margin: 0 0 var(--spacing-sm);
  }

  .jury-member p {
    font-size: var(--text-md);
    font-weight: 500;
    margin: 0 0 var(--spacing-sm);
  }

  .text-muted {
    color: var(--color-text-muted);
    font-style: italic;
  }

  .outcome {
    background: color-mix(
      in srgb,
      var(--color-success) 5%,
      var(--color-surface)
    );
    border-color: color-mix(in srgb, var(--color-success) 30%, transparent);
  }

  .result-display {
    display: flex;
    justify-content: space-around;
    align-items: center;
    padding: var(--spacing-sm) 0;
  }

  .decision,
  .grade {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: var(--spacing-sm);
  }

  .grade-value {
    font-size: 2rem;
    font-weight: 800;
    color: var(--color-primary);
  }

  @media screen and (max-width: 768px) {
    .grid-layout {
      grid-template-columns: 1fr;
    }
    .grid-2 {
      grid-template-columns: 1fr;
    }
    .result-display {
      flex-direction: column;
      gap: var(--spacing-lg);
    }
  }
</style>
