<script lang="ts">
    import AppShell from "$lib/components/ui/AppShell.svelte";
    import { notificationStore } from "$lib/stores/notifications";
    import {
        LayoutDashboard,
        PlusCircle,
        FileText,
        Users,
        Bell,
    } from "lucide-svelte";

    let { children, data } = $props();

    let unreadCount = $state(0);
    notificationStore.subscribe((n) => (unreadCount = n));

    const companyNavLinks = $derived([
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
        {
            href: "/company/notifications",
            label: "Notifications",
            icon: Bell,
            count: unreadCount,
        },
    ]);
</script>

<AppShell
    links={companyNavLinks}
    quickAccess={[
        { href: "/company/propose-subject", label: "Proposer un sujet", icon: PlusCircle },
        { href: "/company/my-subjects", label: "Mes sujets", icon: FileText },
        { href: "/company/supervised-pfes", label: "Encadrements", icon: Users },
        { href: "/company/notifications", label: "Notifications", icon: Bell, count: unreadCount },
    ]}
    user={data.profile ?? {
        full_name: "Company Representative",
        role: "company",
    }}
>
    {@render children()}
</AppShell>
