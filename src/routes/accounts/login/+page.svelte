<script>
    import Logo from "$lib/components/Logo.svelte";
    import View from "$lib/components/ui/View.svelte";
    import LoginGoogle from "$lib/components/ui/custom/LoginGoogle.svelte";
    import LoginDev from "$lib/components/ui/custom/LoginDev.svelte";
    import { BRAND } from "$lib/constants/branding";

    import { Building2 } from "lucide-svelte";
</script>

<View fullScreen center>
    <main data-aos="flip-left">
        <div class="card-wrapper">
            <Logo width="180px" />

            <div class="header">
                <h1>Connection</h1>
                <p>Connect with your {BRAND.university.name} Google account.</p>
            </div>

            <div class="cta">
                <LoginGoogle />

                <LoginDev autoWidth />
            </div>

            <p class="question">
                <span class="icon">
                    <Building2 size="1.5rem" />
                </span>

                <span class="content">
                    Are you a company?
                    <a href="./company/login">Click here</a>.
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
