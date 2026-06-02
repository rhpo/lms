<script lang="ts">
    import { invalidateAll } from "$app/navigation";
    import {
        Building2,
        Globe,
        Mail,
        MapPin,
        Phone,
        Edit3,
        CheckCircle,
        XCircle,
    } from "lucide-svelte";
    import { admin } from "$lib/api";

    import Badge from "$lib/components/ui/Badge.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import Modal from "$lib/components/ui/Modal.svelte";
    import FormField from "$lib/components/ui/FormField.svelte";
    import Page from "$lib/components/ui/Page.svelte";

    let { data } = $props();
    const { companies } = $derived(data);

    let editingCompany = $state<(typeof companies)[0] | null>(null);
    let editForm = $state({
        company_name: "",
        contact_email: "",
        sector: "",
        contact_phone: "",
        website: "",
        description: "",
    });
    let editError = $state("");
    let editLoading = $state(false);

    function openEdit(company: (typeof companies)[0]) {
        editingCompany = company;
        editForm = {
            company_name: company.company_name ?? "",
            contact_email: company.contact_email ?? "",
            sector: company.sector ?? "",
            contact_phone: company.contact_phone ?? "",
            website: company.website ?? "",
            description: company.description ?? "",
        };
    }

    function closeEdit() {
        editingCompany = null;
        editError = "";
    }

    async function toggleVerification(profileId: number, current: boolean) {
        try {
            await admin.companyAction(profileId, "verify-company", {
                is_verified: !current,
            });
            await invalidateAll();
        } catch {

        }
    }

    async function saveEdit() {
        if (!editingCompany) return;
        editError = "";
        editLoading = true;
        try {
            await admin.companyAction(
                editingCompany.profile_id,
                "update",
                editForm,
            );
            await invalidateAll();
            closeEdit();
        } catch (err) {
            editError = err instanceof Error ? err.message : "Erreur reseau";
        } finally {
            editLoading = false;
        }
    }
</script>

<Page title="Entreprises" subtitle="Gestion des comptes entreprises.">
    {#if companies.length === 0}
        <div class="empty-state">
            <Building2 size={48} />
            <h2>Aucune entreprise</h2>
            <p>Aucune entreprise n'est enregistree pour le moment.</p>
        </div>
    {:else}
        <div class="table-container">
            <table>
                <thead>
                    <tr>
                        <th>Entreprise</th>
                        <th>Contact</th>
                        <th>Secteur</th>
                        <th>Statut</th>
                        <th>Date d'inscription</th>
                        <th>Actions</th>
                    </tr>
                </thead>
                <tbody>
                    {#each companies as company (company.profile_id)}
                        <tr>
                            <td>
                                <div class="company-name">
                                    {#if company.logo_url}
                                        <img
                                            src={company.logo_url}
                                            alt={company.company_name}
                                            class="company-logo"
                                        />
                                    {/if}
                                    <span>{company.company_name}</span>
                                </div>
                            </td>
                            <td>
                                <div class="contact-info">
                                    {#if company.contact_email}
                                        <span class="contact-item"
                                            ><Mail size={12} />
                                            {company.contact_email}</span
                                        >
                                    {/if}
                                    {#if company.contact_phone}
                                        <span class="contact-item"
                                            ><Phone size={12} />
                                            {company.contact_phone}</span
                                        >
                                    {/if}
                                </div>
                            </td>
                            <td>{company.sector}</td>
                            <td>
                                <Badge
                                    variant={company.is_verified
                                        ? "success"
                                        : "warning"}
                                    label={company.is_verified
                                        ? "Validee"
                                        : "En attente"}
                                />
                            </td>
                            <td>{new Date(company.created_at).toLocaleDateString("fr-FR")}</td>
                            <td>
                                <div class="actions">
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        Icon={Edit3}
                                        onclick={() => openEdit(company)}
                                        title="Modifier"
                                    />
                                    <Button
                                        variant="ghost"
                                        size="sm"
                                        Icon={company.is_verified
                                            ? XCircle
                                            : CheckCircle}
                                        onclick={() =>
                                            toggleVerification(
                                                company.profile_id,
                                                company.is_verified,
                                            )}
                                        title={company.is_verified
                                            ? "Retirer validation"
                                            : "Valider"}
                                    />
                                </div>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    {/if}
</Page>

<Modal
    open={editingCompany !== null}
    title="Modifier l'entreprise"
    onClose={closeEdit}
>
    {#if editError}
        <div class="error-banner">{editError}</div>
    {/if}
    <div class="edit-form">
        <FormField label="Nom de l'entreprise" required>
            <input
                type="text"
                bind:value={editForm.company_name}
                class="input"
                required
            />
        </FormField>
        <FormField label="Email" required>
            <input
                type="email"
                bind:value={editForm.contact_email}
                class="input"
                required
            />
        </FormField>
        <FormField label="Secteur" required>
            <input
                type="text"
                bind:value={editForm.sector}
                class="input"
                required
            />
        </FormField>
        <FormField label="Telephone">
            <input
                type="text"
                bind:value={editForm.contact_phone}
                class="input"
            />
        </FormField>
        <FormField label="Site web">
            <input
                type="text"
                bind:value={editForm.website}
                class="input"
                placeholder="https://..."
            />
        </FormField>
        <FormField label="Description">
            <textarea bind:value={editForm.description} rows="3" class="input"
            ></textarea>
        </FormField>
        <div class="form-actions">
            <Button variant="ghost" onclick={closeEdit}>Annuler</Button>
            <Button variant="primary" onclick={saveEdit} disabled={editLoading}>
                {editLoading ? "Enregistrement..." : "Enregistrer"}
            </Button>
        </div>
    </div>
</Modal>

<style>
    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        gap: 1rem;
        padding: 3rem;
        text-align: center;
        color: var(--color-text-muted);

        h2 {
            font-size: var(--text-lg);
            font-weight: 600;
            font-family: var(--font-sans);
            color: var(--color-text);
            margin: 0;
        }

        p {
            font-size: var(--text-sm);
            font-family: var(--font-sans);
            margin: 0;
        }
    }

    .table-container {
        overflow-x: auto;
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 12px;
    }

    table {
        width: 100%;
        border-collapse: collapse;

        thead {
            th {
                padding: var(--spacing-sm) var(--spacing-md);
                font-size: var(--text-xs);
                font-weight: 600;
                text-transform: uppercase;
                letter-spacing: 0.05em;
                color: var(--color-text-muted);
                font-family: var(--font-sans);
                text-align: left;
                border-bottom: 1px solid var(--color-border);
                background: var(--color-background);
                white-space: nowrap;
            }
        }

        tbody {
            tr {
                border-bottom: 1px solid var(--color-border);
                transition: background var(--transition-fast);

                &:last-child {
                    border-bottom: none;
                }

                &:hover {
                    background: var(--color-background-100);
                }
            }

            td {
                padding: var(--spacing-sm) var(--spacing-md);
                font-size: var(--text-sm);
                font-family: var(--font-sans);
                color: var(--color-text);
                vertical-align: middle;
            }
        }
    }

    .company-name {
        display: flex;
        align-items: center;
        gap: var(--spacing-sm);
        font-weight: 600;
    }

    .company-logo {
        width: 28px;
        height: 28px;
        border-radius: 6px;
        object-fit: cover;
    }

    .contact-info {
        display: flex;
        flex-direction: column;
        gap: 0.25rem;
    }

    .contact-item {
        display: flex;
        align-items: center;
        gap: 0.25rem;
        font-size: var(--text-xs);
        color: var(--color-text-muted);
    }

    .actions {
        display: flex;
        gap: var(--spacing-xs);
    }

    .edit-form {
        display: flex;
        flex-direction: column;
        gap: var(--spacing-md);

        .input {
            width: 100%;
            padding: 0.5rem 0.75rem;
            border: 1px solid var(--color-border);
            border-radius: 8px;
            font-size: var(--text-sm);
            font-family: var(--font-sans);
            background: var(--color-background);
            color: var(--color-text);
            box-sizing: border-box;

            &:focus {
                outline: none;
                border-color: var(--color-accent);
            }
        }

        textarea.input {
            resize: vertical;
        }
    }

    .error-banner {
        padding: 0.75rem 1rem;
        background: color-mix(in srgb, var(--color-danger) 10%, transparent);
        color: var(--color-danger);
        border: 1px solid
            color-mix(in srgb, var(--color-danger) 20%, transparent);
        border-radius: 8px;
        font-size: var(--text-sm);
        font-family: var(--font-sans);
    }

    .form-actions {
        display: flex;
        justify-content: flex-end;
        gap: var(--spacing-sm);
        margin-top: var(--spacing-sm);
    }
</style>
