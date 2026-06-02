<script lang="ts">
    import Logo from "$lib/components/Logo.svelte";
    import View from "$lib/components/ui/View.svelte";
    import Button from "$lib/components/ui/Button.svelte";
    import PageMeta from "$lib/components/ui/PageMeta.svelte";

    import { Flag, UserPlus } from "lucide-svelte";
</script>

<PageMeta
    title="Main Page"
    description="Digitize and automate the full Projet de Fin d'Études lifecycle - from subject submission to defense archiving."
>
    <div class="hero">
        <div class="image">
            <div class="overlay"></div>

            <video src="/media/hero.mp4" autoplay muted loop>
                Your browser does not support the video tag.
            </video>
        </div>

        <View center>
            <main data-aos="fade-up" data-aos-duration="1000">
                <div
                    class="left"
                    data-aos="fade-right"
                    data-aos-duration="1500"
                >
                    <h1>
                        Simplifiez la gestion de vos projets de fin d'études,
                        avec la nouvelle solution digitale à l'<a
                            href="https://www.esst-sup.com"
                            target="_blank">ESST</a
                        >.
                    </h1>

                    <div class="actions">
                        <Button
                            Icon={UserPlus}
                            href="/accounts/login"
                            variant="secondary"
                        >
                            Se connecter
                        </Button>

                        <Button
                            Icon={Flag}
                            href="/accounts/login"
                            variant="primary"
                        >
                            Commencer
                        </Button>
                    </div>
                </div>

                <div
                    class="right"
                    data-aos="flip-left"
                    data-aos-duration="1500"
                >
                    <View center>
                        <Logo width="25rem" />
                    </View>
                </div>
            </main>
        </View>
    </div>
</PageMeta>

<style>
    .hero {
        flex: 1;

        min-height: calc(100vh - var(--nav-height));

        position: relative;

        display: flex;
        justify-content: center;
        align-items: center;

        h1 {
            color: var(--color-gray-50);
            font-family: var(--font-serif);

            a {
                color: var(--color-accent-fixed);
            }
        }

        main {
            display: flex;
            justify-content: space-between;
            align-items: center;

            width: 100%;

            .left {
                max-width: 600px;
            }
        }
    }

    .image {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        z-index: -1;

        .overlay {
            position: absolute;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;

            background-color: rgba(8, 22, 38, 0.746);

            z-index: 1;
        }

        video {
            width: 100%;
            height: 100%;
            object-fit: cover;
        }
    }

    .actions {
        margin-top: 1rem;

        display: flex;
        gap: 1rem;
    }

    @media (max-width: 1000px) {
        main {
            flex-direction: column-reverse;
            gap: 2rem;

            padding: 0 1rem;
            padding-bottom: var(--nav-height);
            text-align: center;

            h1 {
                font-size: 1.5rem;
            }

            .actions {
                justify-content: center;
            }

            .right {
                width: 80%;
            }
        }
    }
</style>
