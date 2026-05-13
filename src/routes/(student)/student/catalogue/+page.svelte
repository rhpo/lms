<script lang="ts">
    import { goto } from "$app/navigation";
    import { Search, Users, Building2 } from "lucide-svelte";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    const subjects = [
        {
            id: "subj-1",
            title: "Application mobile de suivi des stages",
            description:
                "Developper une application mobile permettant le suivi en temps reel des stagiaires par les encadreurs.",
            group_type: "binome",
            proposer_role: "teacher",
            proposer_name: "Dr. Amine Benali",
            specialty: "ISIL",
        },
        {
            id: "subj-2",
            title: "Systeme de recommandation pour e-commerce",
            description:
                "Conception et realisation d'un systeme de recommandation base sur le filtrage collaboratif.",
            group_type: "trinome",
            proposer_role: "company",
            proposer_name: "TechCorp SARL",
            specialty: "ISIL",
        },
        {
            id: "subj-3",
            title: "Analyse des donnees de production industrielle",
            description:
                "Mise en place d'une solution Big Data pour l'analyse des donnees de production.",
            group_type: "monome",
            proposer_role: "teacher",
            proposer_name: "Dr. Sarah Khelif",
            specialty: "ELEC",
        },
    ];

    let searchQuery = $state("");
    let groupTypeFilter = $state("");
    let proposerFilter = $state("");

    const filteredSubjects = $derived.by(() => {
        let list = subjects;

        if (searchQuery) {
            const q = searchQuery.toLowerCase();
            list = list.filter(
                (s: { title: string; description: string }) =>
                    s.title.toLowerCase().includes(q) ||
                    s.description.toLowerCase().includes(q),
            );
        }

        if (groupTypeFilter) {
            list = list.filter(
                (s: { group_type: string }) => s.group_type === groupTypeFilter,
            );
        }

        if (proposerFilter) {
            list = list.filter(
                (s: { proposer_role: string }) =>
                    s.proposer_role === proposerFilter,
            );
        }

        return list;
    });

    const groupTypeLabel = $derived.by(() => {
        const labels: Record<string, string> = {
            monome: "Monome",
            binome: "Binome",
            trinome: "Trinome",
        };
        return labels;
    });
</script>

<Page
    title="Catalogue des sujets"
    subtitle="Consultez les sujets PFE disponibles filtres par votre specialite."
>
    <div class="search-bar">
        <div class="search-input-wrapper">
            <Search size={16} />
            <input
                type="text"
                placeholder="Rechercher par titre ou mot-cle..."
                bind:value={searchQuery}
            />
        </div>

        <select bind:value={groupTypeFilter}>
            <option value="">Tous les groupes</option>
            <option value="monome">Monome</option>
            <option value="binome">Binome</option>
            <option value="trinome">Trinome</option>
        </select>

        <select bind:value={proposerFilter}>
            <option value="">Tous les encadreurs</option>
            <option value="teacher">Enseignant</option>
            <option value="company">Entreprise</option>
        </select>
    </div>

    {#if filteredSubjects.length === 0}
        <p class="empty">
            {subjects.length === 0
                ? "Aucun sujet disponible pour le moment."
                : "Aucun sujet ne correspond a vos filtres."}
        </p>
    {:else}
        <div class="subject-grid">
            {#each filteredSubjects as subject}
                <div
                    class="subject-card"
                    onclick={() => goto(`/student/catalogue/${subject.id}`)}
                    role="button"
                    tabindex="0"
                >
                    <div class="card-header">
                        <h3>{subject.title}</h3>
                        <Badge
                            variant="info"
                            label={groupTypeLabel[subject.group_type] ??
                                subject.group_type}
                        />
                    </div>

                    <p class="card-description">
                        {subject.description.length > 150
                            ? subject.description.slice(0, 150) + "..."
                            : subject.description}
                    </p>

                    <div class="card-meta">
                        <span class="meta-item">
                            {#if subject.proposer_role === "company"}
                                <Building2 size={14} />
                            {:else}
                                <Users size={14} />
                            {/if}
                            {subject.proposer_name ?? subject.proposer_role}
                        </span>
                        <span class="meta-item">
                            {subject.specialty}
                        </span>
                    </div>

                    <div class="card-footer">
                        <Button variant="ghost" size="sm">
                            Voir le detail
                        </Button>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</Page>

<style>
    .search-bar {
        display: flex;
        gap: var(--spacing-sm);
        margin-bottom: var(--spacing-lg);
        flex-wrap: wrap;

        select {
            padding: 0.5rem 0.75rem;
            border: 1px solid var(--color-border);
            border-radius: 8px;
            font-size: var(--text-sm);
            font-family: var(--font-sans);
            background: var(--color-surface);
            color: var(--color-text);
            cursor: pointer;
        }
    }

    .search-input-wrapper {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        flex: 1;
        min-width: 200px;
        padding: 0.5rem 0.75rem;
        border: 1px solid var(--color-border);
        border-radius: 8px;
        background: var(--color-surface);
        color: var(--color-text-muted);

        input {
            border: none;
            outline: none;
            flex: 1;
            font-size: var(--text-sm);
            font-family: var(--font-sans);
            background: transparent;
            color: var(--color-text);
        }
    }

    .empty {
        text-align: center;
        color: var(--color-text-muted);
        font-style: italic;
        font-size: var(--text-sm);
        padding: 2rem;
    }

    .subject-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
        gap: var(--spacing-md);
    }

    .subject-card {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
        display: flex;
        flex-direction: column;
        gap: 0.75rem;
        cursor: pointer;
        transition: box-shadow var(--transition-normal);

        &:hover {
            box-shadow: var(--shadow-md);
        }
    }

    .card-header {
        display: flex;
        align-items: flex-start;
        justify-content: space-between;
        gap: var(--spacing-sm);

        h3 {
            font-size: var(--text-base);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0;
            flex: 1;
        }
    }

    .card-description {
        font-size: var(--text-sm);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
        margin: 0;
        line-height: 1.5;
    }

    .card-meta {
        display: flex;
        gap: var(--spacing-md);
        flex-wrap: wrap;
    }

    .meta-item {
        display: flex;
        align-items: center;
        gap: 0.35rem;
        font-size: var(--text-xs);
        color: var(--color-text-muted);
        font-family: var(--font-sans);
    }

    .card-footer {
        margin-top: auto;
        padding-top: 0.5rem;
    }
</style>
