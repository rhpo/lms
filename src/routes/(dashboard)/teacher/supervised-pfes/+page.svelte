<script lang="ts">
  import { goto } from "$app/navigation";
  import { Eye } from "lucide-svelte";
  import { authStore } from "$lib/stores/auth";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { supervisedPfes } = $derived(data);
  const myProfileId = $derived(authStore.profile?.id);

  const STATUS_LABELS: Record<string, string> = {
    en_cours: "En cours",
    soutenance_planifiee: "Soutenance planifiee",
    valide: "Valide",
    refuse: "Refuse",
  };

  const STATUS_VARIANTS: Record<
    string,
    "info" | "warning" | "success" | "danger"
  > = {
    en_cours: "info",
    soutenance_planifiee: "warning",
    valide: "success",
    refuse: "danger",
  };
</script>

<Page
  title="Mes encadrements"
  subtitle="Liste des PFE que vous encadrez cette annee."
>
  {#if supervisedPfes.length === 0}
    <div class="empty">
      <p>Aucun PFE encadre pour le moment.</p>
    </div>
  {:else}
    <div class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>Sujet</th>
            <th>Etudiant(s)</th>
            <th>Rôle</th>
            <th>Statut</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each supervisedPfes as pfe}
            <tr>
              <td class="title"
                >{pfe.subject?.title ?? pfe.subject_title ?? "-"}</td
              >
              <td>
                {pfe.student?.profile?.full_name ??
                  "Étudiant #" + pfe.student_id}
                {#if pfe.student2_id}, {pfe.student2?.profile?.full_name ??
                    "Étudiant #" + pfe.student2_id}{/if}
                {#if pfe.student3_id}, {pfe.student3?.profile?.full_name ??
                    "Étudiant #" + pfe.student3_id}{/if}
              </td>
              <td>
                {#if pfe.supervisor?.profile_id === myProfileId}
                  <Badge variant="info" label="Encadrant" />
                {:else}
                  <Badge variant="warning" label="Co-encadrant" />
                {/if}
              </td>
              <td>
                <Badge
                  variant={STATUS_VARIANTS[pfe.status] || "info"}
                  label={STATUS_LABELS[pfe.status] || pfe.status}
                />
              </td>
              <td>
                <Button
                  variant="ghost"
                  Icon={Eye}
                  onclick={() => goto(`/teacher/supervised-pfes/${pfe.id}`)}
                >
                  Voir le detail
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
