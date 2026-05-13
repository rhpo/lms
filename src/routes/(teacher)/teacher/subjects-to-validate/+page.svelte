<script lang="ts">
  import { goto } from "$app/navigation";
  import { Eye } from "lucide-svelte";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { subjects } = $derived(data);

  const GROUP_LABELS: Record<string, string> = {
    monome: "Monome",
    binome: "Binome",
    trinome: "Trinome",
  };

  const PROPOSER_LABELS: Record<string, string> = {
    teacher: "Enseignant",
    company: "Entreprise",
    student: "Etudiant",
  };
</script>

<Page title="Sujets a valider" subtitle="Sujets en attente de votre validation.">
  {#if subjects.length === 0}
    <div class="empty">
      <p>Aucun sujet en attente de validation pour le moment.</p>
    </div>
  {:else}
    <div class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th>Titre</th>
            <th>Auteur</th>
            <th>Specialite</th>
            <th>Groupe</th>
            <th>Date de soumission</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          {#each subjects as subject}
            <tr>
              <td class="title">{subject.title}</td>
              <td>{PROPOSER_LABELS[subject.proposer_role] ?? subject.proposer_role}</td>
              <td>{subject.specialty}</td>
              <td>{GROUP_LABELS[subject.group_type] ?? subject.group_type}</td>
              <td>{subject.created_at}</td>
              <td>
                <Button
                  variant="ghost"
                  Icon={Eye}
                  onclick={() => goto(`/teacher/subjects-to-validate/${subject.id}`)}
                >
                  Examiner
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
