<script lang="ts">
    import { MetaTags } from "svelte-meta-tags";
    import { getContext } from "svelte";


    import { BRAND } from "$lib/constants/branding";

    let {
        children,

        title,
        description,
        image = "",
    }: {
        children: any;
        title: string;
        description: string;
        image?: string;
    } = $props();

    let defaultData = {
        keywords:
            "pfe, esst, esst-sup, esst-sup.com, pfe.esst-sup.com, algérie, pfe algérie, pfe en ligne, pfe en ligne".split(
                ", ",
            ),

        image: "/media/hero.mp4",
    };
</script>

<MetaTags
    title={title + " | " + BRAND.name}
    {description}
    keywords={defaultData.keywords}
    additionalMetaTags={[
        { property: "og:image", content: image || defaultData.image },
        { name: "twitter:image", content: image || defaultData.image },
    ]}
/>

{@render children?.()}
