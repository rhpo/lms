<script>
    import Logo from "$lib/components/Logo.svelte";
    import View from "$lib/components/ui/View.svelte";
    import { BRAND } from "$lib/constants/branding";

    import { Building2 } from "lucide-svelte";
    import LoginGoogle from "$lib/components/ui/login/LoginGoogle.svelte";
</script>

<View fullScreen center>
    <main data-aos="flip-left">
        <div class="card-wrapper">
            <Logo width="180px" />

            <div class="header">
                <h1>Connexion</h1>
                <p>
                    Connectez-vous à votre compte Google {BRAND.university
                        .name}.
                </p>
            </div>

            <div class="cta">
                <LoginGoogle />
            </div>

            <p class="question">
                <span class="icon">
                    <Building2 size="1.5rem" />
                </span>

                <span class="content">
                    Vous êtes une entreprise ?
                    <a href="./company/login">Cliquez ici</a>.
                </span>
            </p>
        </div>
    </main>
</View>

<style>
    .card-wrapper {
        display: flex;
        flex-direction: column;
        align-items: center;

        gap: 1rem;
    }

    .cta {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .header {
        text-align: center;

        margin-bottom: 4rem;
    }

    .question {
        font-size: 1.2rem;

        display: flex;
        align-items: center;
        gap: 0.5rem;

        .content {
            transform: translateY(2px);
        }
    }
</style>
