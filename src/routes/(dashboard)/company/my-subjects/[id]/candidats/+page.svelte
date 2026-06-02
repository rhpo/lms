<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { Users, Check, ArrowLeft } from "lucide-svelte";
    import { company } from "$lib/api";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();

    const { subject, wishes, isLocked } = $derived(data);

    let selectedIds = $state<number[]>([]);
    let assignError = $state("");
    let assignSuccess = $state(false);

    const maxStudents = $derived.by(() => {
        if (!subject) return 1;
        const limits: Record<string, number> = {
            monome: 1,
            binome: 2,
            trinome: 3,
        };
        return limits[subject.group_type] ?? 1;
    });

    function toggleSelection(studentId: number) {
        if (selectedIds.includes(studentId)) {
            selectedIds = selectedIds.filter((id) => id !== studentId);
        } else {
            if (selectedIds.length < maxStudents) {
                selectedIds = [...selectedIds, studentId];
            }
        }
    }

    async function confirmAssignAction() {
        assignError = "";
        assignSuccess = false;

        if (!subject) return;

        if (selectedIds.length === 0) {
            assignError = "Veuillez selectionner au moins un etudiant.";
            return;
        }

        try {
            await company.acceptCandidat(subject.id, {
                student_ids: selectedIds,
            });
            assignSuccess = true;
            selectedIds = [];
            await invalidateAll();
        } catch (err: unknown) {
            assignError = err instanceof Error ? err.message : "Erreur reseau";
        }
    }
</script>

<Page
    title="Candidats"
    subtitle="Selectionnez les etudiants a affecter a ce sujet."
>
    <div class="back-link">
        <a href="/company/my-subjects">
            <ArrowLeft size={14} />
            Retour a mes sujets
        </a>
    </div>

    {#if subject}
        <div class="subject-info">
            <h2>{subject.title}</h2>
            <div class="subject-meta">
                <Badge label={subject.group_type} variant="info" />
                <span class="max-info">
                    Ce sujet accepté {maxStudents} etudiant{maxStudents !== 1
                        ? "s"
                        : ""}
                </span>
            </div>
        </div>

        {#if assignSuccess}
            <div class="success-banner">Affectation confirmee avec succes.</div>
        {/if}

        {#if assignError}
            <div class="error-banner">{assignError}</div>
        {/if}

        {#if isLocked}
            <div class="locked-banner">
                <Check size={16} />
                <span>Affectation déjà effectuée. Les étudiants suivants ont été acceptés :</span>
            </div>
            <div class="student-list">
                {#each wishes.filter(w => w.status === 'accepte') as wish}
                    {@const studentName = wish.student?.profile?.full_name ?? "Étudiant inconnu"}
                    {@const speciality = wish.student?.speciality?.name ?? null}
                    {@const level = wish.student?.level ?? null}
                    {@const detail = [speciality, level].filter(Boolean).join(" · ") || "Spécialité non renseignée"}
                    <div class="student-card is-selected is-locked">
                        <div class="student-checkbox">
                            <Check size={16} />
                        </div>
                        <div class="student-info">
                            <span class="student-name">{studentName}</span>
                            <span class="student-detail">{detail}</span>
                        </div>
                    </div>
                {/each}
            </div>
        {:else if wishes.length === 0}
            <p class="empty">
                Aucun etudiant n'a mis ce sujet dans ses voeux pour le moment.
            </p>
        {:else}
            <div class="selection-info">
                <Users size={16} />
                <span>
                    {selectedIds.length} selectionne{selectedIds.length !== 1
                        ? "s"
                        : ""} sur {maxStudents} maximum
                </span>
            </div>

            <div class="student-list">
                {#each wishes as wish}
                    {@const studentName = wish.student?.profile?.full_name ?? "Étudiant inconnu"}
                    {@const speciality = wish.student?.speciality?.name ?? null}
                    {@const level = wish.student?.level ?? null}
                    {@const detail = [speciality, level].filter(Boolean).join(" · ") || "Spécialité non renseignée"}
                    <div
                        class="student-card"
                        class:is-selected={selectedIds.includes(wish.student_id)}
                        onclick={() => toggleSelection(wish.student_id)}
                        role="button"
                        tabindex="0"
                    >
                        <div class="student-checkbox">
                            {#if selectedIds.includes(wish.student_id)}
                                <Check size={16} />
                            {/if}
                        </div>
                        <div class="student-info">
                            <span class="student-name">{studentName}</span>
                            <span class="student-detail">{detail}</span>
                        </div>
                    </div>
                {/each}
            </div>

            <div class="form-actions">
                <Button
                    variant="primary"
                    onclick={confirmAssignAction}
                    disabled={selectedIds.length === 0}
                >
                    <Check size={14} />
                    Confirmer l'affectation
                </Button>
            </div>
        {/if}
    {:else}
        <p class="empty">Sujet introuvable.</p>
    {/if}
</Page>

<style>
    .back-link {
        margin-bottom: var(--spacing-md);

        a {
            display: inline-flex;
            align-items: center;
            gap: 0.35rem;
            font-size: var(--text-sm);
            color: var(--color-accent);
            text-decoration: none;
            font-family: var(--font-sans);

            &:hover {
                text-decoration: underline;
            }
        }
    }

    .subject-info {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
        margin-bottom: var(--spacing-lg);

        h2 {
            font-size: var(--text-lg);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0 0 var(--spacing-sm);
        }
    }

    .subject-meta {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
    }

    .max-info {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
    }

    .success-banner {
        padding: 0.75rem 1rem;
        background: color-mix(
            in srgb,
            var(--color-success) 10%,
            var(--color-surface)
        );
        border: 1px solid
            color-mix(in srgb, var(--color-success) 20%, transparent);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-success);
        margin-bottom: var(--spacing-md);
    }

    .error-banner {
        padding: 0.75rem 1rem;
        background: color-mix(
            in srgb,
            var(--color-error) 10%,
            var(--color-surface)
        );
        border: 1px solid
            color-mix(in srgb, var(--color-error) 20%, transparent);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-error);
        margin-bottom: var(--spacing-md);
    }

    .empty {
        text-align: center;
        color: var(--color-text-muted);
        font-style: italic;
        font-size: var(--text-sm);
        padding: 2rem;
        font-family: var(--font-sans);
    }

    .selection-info {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1rem;
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-text-muted);
        margin-bottom: var(--spacing-md);
    }

    .student-list {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-xs);
        margin-bottom: var(--spacing-lg);
    }

    .student-card {
        display: flex;
        align-items: center;
        gap: var(--spacing-md);
        padding: var(--spacing-md);
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        cursor: pointer;
        transition: all var(--transition-fast);
    }

    .student-card:hover {
        border-color: var(--color-accent);
    }

    .student-card.is-selected {
        border-color: var(--color-accent);
        background: color-mix(
            in srgb,
            var(--color-accent) 5%,
            var(--color-surface)
        );
    }

    .student-checkbox {
        width: 24px;
        height: 24px;
        border: 2px solid var(--color-border);
        border-radius: 6px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;
        color: var(--color-accent);
    }

    .student-card.is-selected .student-checkbox {
        border-color: var(--color-accent);
        background: var(--color-accent);
        color: white;
    }

    .student-info {
        display: flex;
        flex-direction: column;
        gap: 0.15rem;
        flex: 1;
        min-width: 0;
    }

    .student-name {
        font-size: var(--text-sm);
        font-weight: 600;
        color: var(--color-text);
        font-family: var(--font-sans);
    }

    .student-detail {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
    }

    .form-actions {
        display: flex;
        justify-content: flex-end;
    }

    .locked-banner {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1rem;
        background: color-mix(in srgb, var(--color-success) 10%, var(--color-surface));
        border: 1px solid color-mix(in srgb, var(--color-success) 25%, transparent);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-success);
        margin-bottom: var(--spacing-md);
    }

    .student-card.is-locked {
        cursor: default;
        pointer-events: none;
    }
</style>
