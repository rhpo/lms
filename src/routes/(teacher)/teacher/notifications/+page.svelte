<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { Check, CheckCheck } from "lucide-svelte";

  import Badge from "$lib/components/ui/Badge.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Page from "$lib/components/ui/Page.svelte";

  let { data } = $props();

  const { notifications } = $derived(data);

  const TYPE_LABELS: Record<string, string> = {
    validation_requise: "Validation requise",
    affectation: "Affectation",
    jury: "Jury",
    disponibilite: "Disponibilite",
  };

  const TYPE_VARIANTS: Record<string, "info" | "warning" | "success" | "danger" | "neutral"> = {
    validation_requise: "warning",
    affectation: "info",
    jury: "neutral",
    disponibilite: "info",
  };

  async function markRead(id: string) {
    try {
      await fetch(`/api/notifications/${id}/read`, { method: "POST" });
      await invalidateAll();
    } catch {
      // silent
    }
  }

  async function markAllRead() {
    try {
      await fetch("/api/notifications/read-all", { method: "POST" });
      await invalidateAll();
    } catch {
      // silent
    }
  }
</script>

<Page title="Notifications" subtitle="Toutes vos notifications.">
  {#if notifications.length > 0}
    <div class="toolbar">
      <Button variant="ghost" Icon={CheckCheck} onclick={markAllRead}>
        Tout marquer comme lu
      </Button>
    </div>
  {/if}

  {#if notifications.length === 0}
    <div class="empty">
      <p>Aucune notification.</p>
    </div>
  {:else}
    <div class="list">
      {#each notifications as notif}
        <div class="item" class:unread={!notif.read_at}>
          <div class="item-header">
            <Badge
              variant={TYPE_VARIANTS[notif.type] || "info"}
              label={TYPE_LABELS[notif.type] || notif.type}
            />
            <span class="date">{new Date(notif.created_at).toLocaleDateString("fr-FR", {
              day: "numeric",
              month: "short",
              year: "numeric",
              hour: "2-digit",
              minute: "2-digit",
            })}</span>
          </div>
          <p class="message">{notif.payload?.message ?? notif.type}</p>
          {#if !notif.read_at}
            <div class="item-actions">
              <Button variant="ghost" Icon={Check} onclick={() => markRead(notif.id)}>
                Marquer comme lu
              </Button>
            </div>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</Page>

<style>
  .toolbar {
    margin-bottom: var(--spacing-md);
    display: flex;
    justify-content: flex-end;
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

  .list {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .item {
    padding: var(--spacing-md);
    border: 1px solid var(--color-border);
    border-radius: 8px;
    background: var(--color-surface);
    transition: background var(--transition-fast);
  }

  .item.unread {
    background: var(--color-accent-50);
    border-color: var(--color-accent-200);
  }

  .item-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.5rem;
  }

  .date {
    font-size: 0.75rem;
    color: var(--color-text-muted);
  }

  .message {
    margin: 0;
    font-size: var(--text-sm);
    color: var(--color-text);
    line-height: 1.5;
  }

  .item-actions {
    margin-top: 0.5rem;
    display: flex;
    justify-content: flex-end;
  }
</style>
