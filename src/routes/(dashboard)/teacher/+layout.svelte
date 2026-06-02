<script lang="ts">
    import AppShell from "$lib/components/ui/AppShell.svelte";
    import { notificationStore } from "$lib/stores/notifications";
    import {
        LayoutDashboard,
        FileText,
        CheckSquare,
        Users,
        Gavel,
        Calendar,
        Bell,
    } from "lucide-svelte";

    let { children, data } = $props();

    let unreadCount = $state(0);
    notificationStore.subscribe((n) => (unreadCount = n));

    const teacherNavLinks = $derived([
        {
            href: "/teacher/dashboard",
            label: "Tableau de bord",
            icon: LayoutDashboard,
        },
        {
            href: "/teacher/proposed-subjects",
            label: "Mes sujets",
            icon: FileText,
        },
        {
            href: "/teacher/subjects-to-validate",
            label: "Sujets a valider",
            icon: CheckSquare,
        },
        {
            href: "/teacher/supervised-pfes",
            label: "Mes encadrements",
            icon: Users,
        },
        { href: "/teacher/jury-duties", label: "Jury", icon: Gavel },
        {
            href: "/teacher/availability",
            label: "Disponibilité",
            icon: Calendar,
        },
        {
            href: "/teacher/notifications",
            label: "Notifications",
            icon: Bell,
            count: unreadCount,
        },
    ]);
</script>

<AppShell
    links={teacherNavLinks}
    quickAccess={[
        { href: "/teacher/subjects-to-validate", label: "Sujets à valider", icon: CheckSquare },
        { href: "/teacher/supervised-pfes", label: "Mes encadrements", icon: Users },
        { href: "/teacher/jury-duties", label: "Jury", icon: Gavel },
        { href: "/teacher/notifications", label: "Notifications", icon: Bell, count: unreadCount },
    ]}
    user={data.profile ?? { full_name: "Teacher", role: "teacher" }}
>
    {@render children()}
</AppShell>
