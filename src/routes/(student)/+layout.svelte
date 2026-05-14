<script lang="ts">
    import AppShell from "$lib/components/ui/AppShell.svelte";
    import {
        LayoutDashboard,
        BookOpen,
        Heart,
        FileText,
        Gavel,
        Bell,
    } from "lucide-svelte";

    let { children, data } = $props();

    const studentNavLinks = [
        {
            href: "/student/dashboard",
            label: "Tableau de bord",
            icon: LayoutDashboard,
        },
        { href: "/student/catalogue", label: "Catalogue", icon: BookOpen },
        { href: "/student/voeux", label: "Mes voeux", icon: Heart },
        { href: "/student/my-pfe", label: "Mon PFE", icon: FileText },
        { href: "/student/soutenance", label: "Soutenance", icon: Gavel },
        { href: "/student/notifications", label: "Notifications", icon: Bell },
    ];
</script>

<AppShell
    links={studentNavLinks}
    user={data.profile ?? { full_name: "Student", role: "student" }}
    notificationsHref="/student/notifications"
>
    {@render children()}
</AppShell>
