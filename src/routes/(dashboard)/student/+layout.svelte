<script lang="ts">
    import AppShell from "$lib/components/ui/AppShell.svelte";
    import { notificationStore } from "$lib/stores/notifications";
    import {
        LayoutDashboard,
        BookOpen,
        Heart,
        FileText,
        Gavel,
        Bell,
    } from "lucide-svelte";

    let { children, data } = $props();

    let unreadCount = $state(0);
    notificationStore.subscribe((n) => (unreadCount = n));

    const studentNavLinks = $derived([
        {
            href: "/student/dashboard",
            label: "Tableau de bord",
            icon: LayoutDashboard,
        },
        { href: "/student/catalogue", label: "Catalogue", icon: BookOpen },
        { href: "/student/voeux", label: "Mes voeux", icon: Heart },
        { href: "/student/my-pfe", label: "Mon PFE", icon: FileText },
        { href: "/student/soutenance", label: "Soutenance", icon: Gavel },
        {
            href: "/student/notifications",
            label: "Notifications",
            icon: Bell,
            count: unreadCount,
        },
    ]);
</script>

<AppShell
    links={studentNavLinks}
    quickAccess={[
        { href: "/student/catalogue", label: "Catalogue", icon: BookOpen },
        { href: "/student/voeux", label: "Mes voeux", icon: Heart },
        { href: "/student/my-pfe", label: "Mon PFE", icon: FileText },
        { href: "/student/notifications", label: "Notifications", icon: Bell, count: unreadCount },
    ]}
    user={data.profile ?? { full_name: "Student", role: "student" }}
>
    {@render children()}
</AppShell>
