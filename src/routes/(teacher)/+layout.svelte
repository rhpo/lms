<script lang="ts">
    import AppShell from "$lib/components/ui/AppShell.svelte";
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

    const teacherNavLinks = [
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
            label: "Disponibilite",
            icon: Calendar,
        },
        { href: "/teacher/notifications", label: "Notifications", icon: Bell },
    ];
</script>

<AppShell
    links={teacherNavLinks}
    user={data.profile ?? { full_name: "Teacher", role: "teacher" }}
    notificationsHref="/teacher/notifications"
>
    {@render children()}
</AppShell>
