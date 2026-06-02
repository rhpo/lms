<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { goto } from "$app/navigation";
    import {
        GraduationCap,
        FileText,
        UserPlus,
        X,
        Star,
        Loader,
    } from "lucide-svelte";
    import { PFE_ASSIGNMENT_STATUS_LABELS } from "$lib/constants/labels";
    import { admin } from "$lib/api";
    import { showToast } from "$lib/utils/toast";
    import type { PfeAssignment } from "$lib/types";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";
    import Modal from "$lib/components/ui/Modal.svelte";

    let { data } = $props();
    const { assignments } = $derived(data);

    const statusVariant: Record<
        string,
        "success" | "warning" | "danger" | "info" | "neutral"
    > = {
        en_cours: "info",
        soutenance_planifiee: "warning",
        memoire_soumis: "warning",
        valide: "success",
        refuse: "danger",
    };

    let showCoSupModal = $state(false);
    let selectedAssignment = $state<PfeAssignment | null>(null);
    let recommendations = $state<
        { teacher: any; score: number; matching_domains: any[] }[]
    >([]);
    let subjectDomains = $state<any[]>([]);
    let loadingRec = $state(false);
    let selectedTeacherId = $state<number | null>(null);
    let saving = $state(false);

    async function openCoSupModal(a: PfeAssignment) {
        selectedAssignment = a;
        selectedTeacherId = null;
        recommendations = [];
        subjectDomains = [];
        showCoSupModal = true;
        loadingRec = true;
        try {
            const res = await admin.recommendCoSupervisor(a.id);
            recommendations = res.recommended ?? [];
            subjectDomains = res.subject_domains ?? [];
        } catch {
            showToast.error("Impossible de charger les recommandations");
        } finally {
            loadingRec = false;
        }
    }

    async function assignCoSupervisor() {
        if (!selectedAssignment || !selectedTeacherId) return;
        saving = true;
        try {
            await admin.assignmentAction(
                selectedAssignment.id,
                "assign-co-supervisor",
                {
                    teacher_id: selectedTeacherId,
                },
            );
            showToast.success("Co-encadrant assigné avec succès");
            showCoSupModal = false;
            await invalidateAll();
        } catch (err) {
            showToast.error(
                err instanceof Error ? err.message : "Erreur réseau",
            );
        } finally {
            saving = false;
        }
    }

    function selectRecommended(teacherId: number) {
        selectedTeacherId = teacherId;
    }
</script>

<Page title="PFE" subtitle="Tous les PFE actifs et archivés">
    {#if assignments && assignments.length > 0}
        <div class="table-wrapper">
            <table>
                <thead>
                    <tr>
                        <th>Code PFE</th>
                        <th>Sujet</th>
                        <th>Étudiants</th>
                        <th>Encadrant</th>
                        <th>Co-encadrant</th>
                        <th>Spécialité</th>
                        <th>Statut</th>
                        <th>Mémoire</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#each assignments as a}
                        <tr>
                            <td class="code">{a.pfe_code}</td>
                            <td>{a.subject?.title ?? "-"}</td>
                            <td>
                                {#if a.student}
                                    <ul class="student-list">
                                        <li>
                                            {a.student.profile?.full_name ??
                                                "-"}
                                        </li>
                                        {#if a.student2}<li>
                                                {a.student2.profile
                                                    ?.full_name ?? "-"}
                                            </li>{/if}
                                        {#if a.student3}<li>
                                                {a.student3.profile
                                                    ?.full_name ?? "-"}
                                            </li>{/if}
                                    </ul>
                                {:else}-{/if}
                            </td>
                            <td>{a.supervisor?.profile?.full_name ?? "-"}</td>
                            <td>
                                {#if a.co_supervisor}
                                    <span class="co-sup-name"
                                        >{a.co_supervisor.profile?.full_name ??
                                            "-"}</span
                                    >
                                {:else}
                                    <span class="none">-</span>
                                {/if}
                            </td>
                            <td>
                                {[
                                    a.student?.speciality?.code,
                                    a.student2?.speciality?.code,
                                    a.student3?.speciality?.code,
                                ]
                                    .filter(Boolean)
                                    .join(", ") || "-"}
                            </td>
                            <td>
                                <Badge
                                    variant={statusVariant[a.status] ??
                                        "neutral"}
                                    label={PFE_ASSIGNMENT_STATUS_LABELS[
                                        a.status
                                    ] ?? a.status}
                                />
                            </td>
                            <td>
                                {#if a.memoire_url}
                                    <a
                                        href={a.memoire_url}
                                        target="_blank"
                                        class="a-icon"
                                        rel="noopener noreferrer"
                                    >
                                        <FileText size={16} />Voir
                                    </a>
                                {:else}-{/if}
                            </td>
                            <td>
                                <div class="action-btns">
                                    {#if a.co_supervisor}
                                        <!-- <button class="icon-btn danger" title="Retirer le co-encadrant" onclick={() => removeCoSupervisor(a)}>
                                            <X size={14} />
                                        </button> -->
                                    {:else}
                                        <button
                                            class="icon-btn"
                                            title="Assigner un co-encadrant"
                                            onclick={() => openCoSupModal(a)}
                                        >
                                            <UserPlus size={14} />
                                        </button>
                                    {/if}
                                    <Button
                                        variant="secondary"
                                        onclick={() =>
                                            goto(
                                                `/admin/pfe/defenses/${a.id}/plan`,
                                            )}
                                    >
                                        Planifier
                                    </Button>
                                </div>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {:else}
        <p class="empty">Aucun PFE enregistré.</p>
    {/if}
</Page>

<!-- Co-supervisor Modal -->
{#if showCoSupModal && selectedAssignment}
    <Modal
        open={true}
        title="Assigner un co-encadrant"
        description="Choisissez un enseignant disponible. Les recommandations sont triées par correspondance de domaines avec le sujet."
        onClose={() => (showCoSupModal = false)}
        width="640px"
    >
        <div class="cosup-modal">
            <!-- Subject domains chip row -->
            {#if subjectDomains.length > 0}
                <div class="domain-chips">
                    <span class="chip-label">Domaines du sujet :</span>
                    {#each subjectDomains as d}
                        <span class="chip">{d.name}</span>
                    {/each}
                </div>
            {/if}

            <!-- Recommendations list -->
            {#if loadingRec}
                <div class="loading-state">
                    <Loader size={20} class="spin" />
                    <span>Chargement des recommandations…</span>
                </div>
            {:else if recommendations.length === 0}
                <p class="empty-rec">Aucun enseignant disponible trouvé.</p>
            {:else}
                <div class="rec-list">
                    {#each recommendations as rec}
                        {@const isSelected =
                            selectedTeacherId === rec.teacher.id}
                        <button
                            class="rec-card"
                            class:selected={isSelected}
                            onclick={() => selectRecommended(rec.teacher.id)}
                            type="button"
                        >
                            <div class="rec-left">
                                <span class="rec-name"
                                    >{rec.teacher.profile?.full_name ??
                                        "-"}</span
                                >
                                <span class="rec-grade"
                                    >{rec.teacher.grade ?? ""}</span
                                >
                            </div>
                            <div class="rec-right">
                                {#if rec.score > 0}
                                    <span class="score-badge">
                                        <Star size={11} />
                                        {rec.score} domaine{rec.score > 1
                                            ? "s"
                                            : ""} en commun
                                    </span>
                                    <div class="match-domains">
                                        {#each rec.matching_domains as d}
                                            <span class="chip small"
                                                >{d.name}</span
                                            >
                                        {/each}
                                    </div>
                                {:else}
                                    <span class="no-match"
                                        >Aucun domaine en commun</span
                                    >
                                {/if}
                            </div>
                        </button>
                    {/each}
                </div>
            {/if}

            <!-- Or pick from full dropdown -->
            <div class="separator-row">
                <span>ou choisir manuellement</span>
            </div>
            <div class="field">
                <label for="cosup-select">Enseignant</label>
                <select
                    id="cosup-select"
                    class="input"
                    bind:value={selectedTeacherId}
                >
                    <option value={null}>Sélectionner un enseignant</option>
                    {#each recommendations as rec}
                        {@const isSupervisor =
                            rec.teacher.id === selectedAssignment.supervisor_id}
                        <option value={rec.teacher.id} disabled={isSupervisor}>
                            {rec.teacher.profile?.full_name ?? "-"}{isSupervisor
                                ? " - (Encadrant principal)"
                                : ""}
                        </option>
                    {/each}
                </select>
            </div>

            <div class="modal-actions">
                <Button
                    variant="ghost"
                    onclick={() => (showCoSupModal = false)}
                    disabled={saving}
                >
                    Annuler
                </Button>
                <Button
                    variant="primary"
                    onclick={assignCoSupervisor}
                    disabled={saving || !selectedTeacherId}
                >
                    {saving ? "Assignation…" : "Assigner"}
                </Button>
            </div>
        </div>
    </Modal>
{/if}

<style>
    .a-icon {
        display: flex;
        align-items: center;
        gap: var(--spacing-xs);
        color: var(--color-accent);
        font-size: var(--text-sm);
    }

    .table-wrapper {
        overflow-x: auto;
    }

    table {
        width: 100%;
        border-collapse: collapse;

        thead th {
            text-align: left;
            padding: var(--spacing-sm) var(--spacing-md);
            font-size: var(--text-sm);
            font-weight: 600;
            color: var(--color-text-muted);
            border-bottom: 2px solid var(--color-border);
            white-space: nowrap;
        }

        tbody tr:hover {
            background: var(--color-background-100);
        }

        tbody td {
            padding: var(--spacing-sm) var(--spacing-md);
            font-size: var(--text-sm);
            border-bottom: 1px solid var(--color-border);
            vertical-align: middle;

            &.code {
                font-family: monospace;
                font-weight: 600;
            }
        }
    }

    .student-list {
        margin: 0;
        padding: 0;
        list-style: none;
        li {
            font-size: var(--text-sm);
            line-height: 1.5;
        }
    }

    .co-sup-name {
        font-size: var(--text-sm);
        color: var(--color-text);
    }
    .none {
        color: var(--color-text-muted);
    }

    .action-btns {
        display: flex;
        align-items: center;
        gap: var(--spacing-xs);
    }

    .icon-btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 2rem;
        height: 2rem;
        border-radius: 6px;
        border: 1px solid var(--color-border);
        background: var(--color-surface);
        color: var(--color-text-muted);
        cursor: pointer;
        transition: all 0.15s;
        flex-shrink: 0;

        &:hover {
            border-color: var(--color-accent);
            color: var(--color-accent);
            background: color-mix(
                in srgb,
                var(--color-accent) 8%,
                var(--color-surface)
            );
        }

        &.danger:hover {
            border-color: var(--color-danger);
            color: var(--color-danger);
            background: color-mix(
                in srgb,
                var(--color-danger) 8%,
                var(--color-surface)
            );
        }
    }

    .empty {
        text-align: center;
        color: var(--color-text-muted);
        padding: var(--spacing-xl);
    }

    /* Modal styles */
    .cosup-modal {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-md);
    }

    .domain-chips {
        display: flex;
        flex-wrap: wrap;
        align-items: center;
        gap: 0.4rem;
    }

    .chip-label {
        font-size: var(--text-xs);
        font-weight: 600;
        color: var(--color-text-muted);
        white-space: nowrap;
    }

    .chip {
        font-size: var(--text-xs);
        padding: 0.2rem 0.55rem;
        background: color-mix(
            in srgb,
            var(--color-accent) 12%,
            var(--color-surface)
        );
        color: var(--color-accent);
        border-radius: 99px;
        border: 1px solid
            color-mix(in srgb, var(--color-accent) 25%, transparent);
        font-weight: 500;

        &.small {
            font-size: 10px;
            padding: 0.1rem 0.4rem;
        }
    }

    .loading-state {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        color: var(--color-text-muted);
        font-size: var(--text-sm);
        padding: var(--spacing-md) 0;
    }

    .empty-rec {
        color: var(--color-text-muted);
        font-size: var(--text-sm);
        font-style: italic;
        margin: 0;
    }

    .rec-list {
        display: flex;
        flex-direction: column;
        gap: 0.4rem;
        max-height: 280px;
        overflow-y: auto;
        padding-right: 2px;
    }

    .rec-card {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        gap: 1rem;
        padding: 0.65rem 0.9rem;
        border-radius: 8px;
        border: 1.5px solid var(--color-border);
        background: var(--color-surface);
        cursor: pointer;
        text-align: left;
        transition:
            border-color 0.15s,
            background 0.15s;

        &:hover {
            border-color: var(--color-accent);
            background: color-mix(
                in srgb,
                var(--color-accent) 4%,
                var(--color-surface)
            );
        }
        &.selected {
            border-color: var(--color-accent);
            background: color-mix(
                in srgb,
                var(--color-accent) 8%,
                var(--color-surface)
            );
        }
    }

    .rec-left {
        display: flex;
        flex-direction: column;
        gap: 0.15rem;
    }

    .rec-name {
        font-size: var(--text-sm);
        font-weight: 600;
        color: var(--color-text);
    }

    .rec-grade {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        text-transform: uppercase;
    }

    .rec-right {
        display: flex;
        flex-direction: column;
        align-items: flex-end;
        gap: 0.3rem;
        flex-shrink: 0;
    }

    .score-badge {
        display: inline-flex;
        align-items: center;
        gap: 0.25rem;
        font-size: var(--text-xs);
        font-weight: 600;
        color: var(--color-success);
        background: color-mix(
            in srgb,
            var(--color-success) 10%,
            var(--color-surface)
        );
        border: 1px solid
            color-mix(in srgb, var(--color-success) 25%, transparent);
        border-radius: 99px;
        padding: 0.15rem 0.5rem;
    }

    .match-domains {
        display: flex;
        flex-wrap: wrap;
        gap: 0.25rem;
        justify-content: flex-end;
    }

    .no-match {
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        font-style: italic;
    }

    .separator-row {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        color: var(--color-text-muted);
        font-size: var(--text-xs);

        &::before,
        &::after {
            content: "";
            flex: 1;
            height: 1px;
            background: var(--color-border);
        }
    }

    .field {
        display: flex;
        flex-direction: column;
        gap: 0.35rem;

        label {
            font-size: var(--text-sm);
            font-weight: 600;
            color: var(--color-text-muted);
        }
    }

    .input {
        width: 100%;
        padding: 0.5rem 0.75rem;
        border: 1px solid var(--color-border);
        border-radius: 8px;
        font-size: var(--text-sm);
        background: var(--color-surface);
        color: var(--color-text);

        &:focus {
            border-color: var(--color-accent);
            outline: none;
        }
    }

    .modal-actions {
        display: flex;
        justify-content: flex-end;
        gap: var(--spacing-sm);
        margin-top: var(--spacing-xs);
    }
</style>
