<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import { Plus, Info } from "lucide-svelte";

    import Button from "$lib/components/ui/Button.svelte";
    import FormField from "$lib/components/ui/FormField.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    type Speciality = {
        id: string;
        name: string;
        code: string;
        year_type: "licence" | "master" | "ingenieur";
    };

    const specialities: Speciality[] = [
        {
            id: "spec-1",
            name: "Ingenierie des Systemes d'Information",
            code: "ISIL",
            year_type: "ingenieur",
        },
        {
            id: "spec-2",
            name: "Chimie des Materiaux",
            code: "CHIM",
            year_type: "master",
        },
        {
            id: "spec-3",
            name: "Electrotechnique",
            code: "ELEC",
            year_type: "ingenieur",
        },
    ];

    const niveaux = [
        { value: "licence", label: "Licence" },
        { value: "master", label: "Master" },
        { value: "ingenieur", label: "Ingenieur" },
    ] as const;

    let selectedNiveau = $state<string>("");
    let selectedSpecialtyId = $state<string>("");
    let title = $state("");
    let description = $state("");
    let groupType = $state("monome");

    let error = $state("");
    let success = $state(false);

    const filteredSpecialities = $derived(
        selectedNiveau
            ? specialities.filter((s) => s.year_type === selectedNiveau)
            : [],
    );

    async function submitAction() {
        error = "";
        success = false;

        if (!title || !description || !selectedSpecialtyId) {
            error = "Veuillez remplir tous les champs obligatoires.";
            return;
        }

        try {
            const res = await fetch("/api/subjects", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    title,
                    description,
                    specialty: selectedSpecialtyId,
                    level: selectedNiveau,
                    group_type: groupType,
                    proposer_role: "company",
                }),
            });
            if (res.ok) {
                success = true;
                title = "";
                description = "";
                selectedSpecialtyId = "";
                selectedNiveau = "";
                groupType = "monome";
                await invalidateAll();
            } else {
                const err = await res.json();
                error = err.message ?? "Erreur lors de la soumission";
            }
        } catch {
            error = "Erreur reseau";
        }
    }
</script>

<Page
    title="Proposer un sujet"
    subtitle="Soumettez un nouveau sujet PFE pour validation."
>
    <div class="info-banner">
        <Info size={16} />
        <span>
            Le sujet sera visible dans le catalogue apres validation par deux
            enseignants de la specialite concernee.
        </span>
    </div>

    {#if success}
        <div class="success-banner">
            Votre sujet a ete soumis avec succes. Il sera examine par les
            validateurs.
        </div>
    {/if}

    {#if error}
        <div class="error-banner">{error}</div>
    {/if}

    <form
        class="subject-form"
        onsubmit={(e) => {
            e.preventDefault();
            submitAction();
        }}
    >
        <FormField label="Titre du sujet" required>
            <input
                type="text"
                bind:value={title}
                placeholder="Ex: Developpement d'une application mobile..."
                required
                class="input"
            />
        </FormField>

        <FormField label="Description" required>
            <textarea
                bind:value={description}
                rows="6"
                placeholder="Decrivez le sujet en detail, les objectifs, les technologies attendues..."
                required
                class="input"
            ></textarea>
        </FormField>

        <div class="form-row">
            <FormField label="Niveau" required>
                <select
                    bind:value={selectedNiveau}
                    required
                    class="input"
                    onchange={() => (selectedSpecialtyId = "")}
                >
                    <option value="">Selectionner un niveau</option>
                    {#each niveaux as n}
                        <option value={n.value}>{n.label}</option>
                    {/each}
                </select>
            </FormField>

            <FormField label="Specialite" required>
                <select
                    bind:value={selectedSpecialtyId}
                    required
                    disabled={!selectedNiveau}
                    class="input"
                >
                    <option value="">
                        {selectedNiveau
                            ? "Selectionner une specialite"
                            : "Choisissez d'abord un niveau"}
                    </option>
                    {#each filteredSpecialities as spec}
                        <option value={spec.id}
                            >{spec.code} — {spec.name}</option
                        >
                    {/each}
                </select>
            </FormField>
        </div>

        <FormField label="Type de groupe" required>
            <select bind:value={groupType} required class="input">
                <option value="monome">Monome (1 etudiant)</option>
                <option value="binome">Binome (2 etudiants)</option>
                <option value="trinome">Trinome (3 etudiants)</option>
            </select>
        </FormField>

        <div class="form-actions">
            <Button variant="primary" type="submit">
                <Plus size={14} />
                Soumettre le sujet
            </Button>
        </div>
    </form>
</Page>

<style>
    .info-banner {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1rem;
        background: color-mix(
            in srgb,
            var(--color-info) 10%,
            var(--color-surface)
        );
        border: 1px solid color-mix(in srgb, var(--color-info) 20%, transparent);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
        color: var(--color-info);
        margin-bottom: var(--spacing-lg);
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

    .subject-form {
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
        padding: var(--spacing-lg);
        display: flex;
        flex-direction: column;
        gap: var(--spacing-md);

        textarea {
            width: 100%;
            padding: 0.5rem 0.75rem;
            border: 1px solid var(--color-border);
            border-radius: 8px;
            font-size: var(--text-sm);
            font-family: var(--font-sans);
            background: var(--color-background);
            color: var(--color-text);
            resize: vertical;
        }

        select {
            padding: 0.5rem 0.75rem;
            border: 1px solid var(--color-border);
            border-radius: 8px;
            font-size: var(--text-sm);
            font-family: var(--font-sans);
            background: var(--color-background);
            color: var(--color-text);
            cursor: pointer;
        }
    }

    .form-row {
        display: grid;
        grid-template-columns: 1fr 1fr;
        gap: var(--spacing-md);

        @media screen and (max-width: 480px) {
            & {
                grid-template-columns: 1fr;
            }
        }
    }

    .form-actions {
        display: flex;
        justify-content: flex-end;
        padding-top: var(--spacing-sm);
    }
</style>
