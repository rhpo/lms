<script lang="ts">
  import { goto } from "$app/navigation";
  import { Eye } from "lucide-svelte";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";
  import {
    GROUP_TYPE_LABELS,
    ROLE_LABELS,
    formatDate,
  } from "$lib/constants/labels";
  import type { PfeSubject } from "$lib/types";

  let { data } = $props();

  const { subjects } = $derived(data);

  function otherValidatorDecision(subject: PfeSubject): string | null {
    const d1 = subject.validator1_decision;
    const d2 = subject.validator2_decision;

    if (d1 && !d2) return d1;
    if (d2 && !d1) return d2;
    return null;
  }

  const DECISION_LABELS: Record<string, string> = {
    valide: "Validé",
    accepte_sous_reserve: "Sous réserve",
    refuse: "Refusé",
  };
  const DECISION_VARIANTS: Record<string, "success" | "warning" | "danger"> = {
    valide: "success",
    accepte_sous_reserve: "warning",
    refuse: "danger",
  };
</script>

<Page title="Sujets à valider" subtitle="Sujets en attente de votre décision.">
  {#if subjects.length === 0}
    <div class="empty">
      <p>Aucun sujet en attente de votre décision pour le moment.</p>
    </div>
  {:else}
    <div class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>Titre</th>
            <th>Auteur</th>
            <th>Domaines</th>
            <th>Groupe</th>
            <th>Autre validateur</th>
            <th>Soumis le</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each subjects as subject}
            {@const other = otherValidatorDecision(subject)}
            <tr>
              <td class="title">{subject.title}</td>
              <td
                >{ROLE_LABELS[subject.proposer_role] ??
                  subject.proposer_role}</td
              >
              <td>{subject.domains?.map((d) => d.name).join(", ") || "-"}</td>
              <td
                >{GROUP_TYPE_LABELS[subject.group_type] ??
                  subject.group_type}</td
              >
              <td>
                {#if other}
                  <Badge
                    variant={DECISION_VARIANTS[other] ?? "info"}
                    label={DECISION_LABELS[other] ?? other}
                  />
                {:else}
                  <span class="muted">En attente</span>
                {/if}
              </td>
              <td>{formatDate(subject.created_at)}</td>
              <td>
                <Button
                  variant="ghost"
                  Icon={Eye}
                  onclick={() =>
                    goto(`/teacher/subjects-to-validate/${subject.id}`)}
                >
                  Décider
                </Button>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  {/if}
</Page>

<style>
  .muted {
    font-size: var(--text-xs);
    color: var(--color-text-muted);
  }

  .empty {
    text-align: center;
    padding: 3rem 1rem;
    color: var(--color-text-muted);

    & p {
      font-size: var(--text-sm);
      margin: 0;
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
</style>
