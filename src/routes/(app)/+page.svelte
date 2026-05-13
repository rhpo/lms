<script lang="ts">
    import Logo from "$lib/components/Logo.svelte";
    import Page from "$lib/components/ui/Page.svelte";
    import View from "$lib/components/ui/View.svelte";
    import Button from "$lib/components/ui/Button.svelte";

    import { BRAND } from "$lib/constants/branding";
    import { Grid2x2Icon, UserPlus } from "lucide-svelte";
</script>

<Page
    title="Main Page"
    description="Digitize and automate the full Projet de Fin d'Études lifecycle - from subject submission to defense archiving."
>
    <View center fullScreen>
        <View center>
            <Logo width="10rem" />
            <h1>
                {BRAND.name} Management Software
            </h1>
        </View>

        <div class="actions">
            <Button Icon={UserPlus} href="/accounts/login" variant="secondary">
                Sign-In
            </Button>

            <Button Icon={Grid2x2Icon} href="/dashboard" variant="primary">
                Get Started
            </Button>
        </div>
    </View>
</Page>

<style>
    .actions {
        margin-top: 1rem;

        display: flex;
        gap: 1rem;
    }
</style>
