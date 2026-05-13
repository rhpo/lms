<script lang="ts">
  import { goto } from "$app/navigation";
  import { Plus, Eye, RefreshCw } from "lucide-svelte";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { subjects } = $derived(data);

  const STATUS_LABELS: Record<string, string> = {
    en_attente: "En attente",
    valide: "Valide",
    refuse: "Refuse",
    expire: "Expire",
  };

  function statusVariant(
    status: string,
  ): "success" | "warning" | "danger" | "info" | "neutral" {
    switch (status) {
      case "valide":
        return "success";
      case "en_attente":
        return "warning";
      case "refuse":
        return "danger";
      case "expire":
        return "neutral";
      default:
        return "neutral";
    }
  }

  const GROUP_LABELS: Record<string, string> = {
    monome: "Monome",
    binome: "Binome",
    trinome: "Trinome",
  };
</script>

<Page
  title="Mes sujets proposes"
  subtitle="Gérer les sujets que vous avez proposes."
>
  {#snippet actions()}
    <Button
      variant="primary"
      Icon={Plus}
      onclick={() => goto("/teacher/proposed-subjects/new")}
    >
      Proposer un nouveau sujet
    </Button>
  {/snippet}

  {#if subjects.length === 0}
    <div class="empty">
      <p>Aucun sujet propose pour le moment.</p>
      <Button
        variant="primary"
        onclick={() => goto("/teacher/proposed-subjects/new")}
      >
        Proposer mon premier sujet
      </Button>
    </div>
  {:else}
    <div class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>Titre</th>
            <th>Specialite</th>
            <th>Groupe</th>
            <th>Statut</th>
            <th>Date</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each subjects as subject}
            <tr>
              <td class="title">{subject.title}</td>
              <td>{subject.specialty}</td>
              <td>{GROUP_LABELS[subject.group_type] ?? subject.group_type}</td>
              <td>
                <Badge
                  variant={statusVariant(subject.status)}
                  label={STATUS_LABELS[subject.status] ?? subject.status}
                />
              </td>
              <td>{subject.created_at}</td>
              <td class="actions-cell">
                {#if subject.status === "valide"}
                  <Button
                    variant="ghost"
                    Icon={Eye}
                    onclick={() =>
                      goto(
                        `/teacher/proposed-subjects/${subject.id}/candidats`,
                      )}
                  >
                    Voir les candidats
                  </Button>
                {/if}
                {#if subject.status === "refuse"}
                  <Button
                    variant="ghost"
                    Icon={RefreshCw}
                    onclick={() => goto(`/teacher/proposed-subjects/new`)}
                  >
                    Modifier et resoumettre
                  </Button>
                {/if}
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Page>

<style>
  .empty {
    text-align: center;
    padding: 3rem 1rem;
    color: var(--color-text-muted);

    & p {
      font-size: var(--text-sm);
      margin: 0 0 var(--spacing-md);
    }
  }

  .table-wrapper {
    overflow-x: auto;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    font-family: var(--font-sans);
    font-size: 0.85rem;
  }

  th {
    text-align: left;
    padding: 0.75rem;
    font-weight: 600;
    color: var(--color-text);
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.04em;
    border-bottom: 1px solid var(--color-border);
  }

  td {
    padding: 0.75rem;
    color: var(--color-text);
    border-bottom: 1px solid var(--color-border);
  }

  tr:last-child td {
    border-bottom: none;
  }

  .title {
    font-weight: 600;
    max-width: 280px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .actions-cell {
    display: flex;
    gap: 0.5rem;
    align-items: center;
  }
</style>
