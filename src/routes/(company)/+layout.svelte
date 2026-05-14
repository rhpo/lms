<script lang="ts">
    import AppShell from "$lib/components/ui/AppShell.svelte";
    import {
        LayoutDashboard,
        PlusCircle,
        FileText,
        Users,
        Bell,
    } from "lucide-svelte";

    let { children, data } = $props();

    const companyNavLinks = [
        {
            href: "/company/dashboard",
            label: "Tableau de bord",
            icon: LayoutDashboard,
        },
        {
            href: "/company/propose-subject",
            label: "Proposer un sujet",
            icon: PlusCircle,
        },
        { href: "/company/my-subjects", label: "Mes sujets", icon: FileText },
        {
            href: "/company/supervised-pfes",
            label: "Mes encadrements",
            icon: Users,
        },
        { href: "/company/notifications", label: "Notifications", icon: Bell },
    ];
</script>

<AppShell
    links={companyNavLinks}
    user={data.profile ?? {
        full_name: "Company Representative",
        role: "company",
    }}
    notificationsHref="/company/notifications"
>
    {@render children()}
</AppShell>
