<script lang="ts">
    import AppShell from "$lib/components/ui/AppShell.svelte";
    import {
        LayoutDashboard,
        Users,
        Calendar,
        GraduationCap,
        Gavel,
        FileText,
        ClipboardList,
        BarChart3,
        Download,
        ScrollText,
        Mail,
        Cog,
    } from "lucide-svelte";

    let { children, data } = $props();

    const adminNavLinks = [
        {
            href: "/admin/dashboard",
            label: "Tableau de bord",
            icon: LayoutDashboard,
        },
        { href: "/admin/statistics", label: "Statistiques", icon: BarChart3 },
        { href: "/admin/subjects", label: "Sujets", icon: FileText },
        { href: "/admin/pfe", label: "PFE", icon: GraduationCap },
        { href: "/admin/defenses", label: "Soutenances", icon: Gavel },
        { href: "/admin/users", label: "Utilisateurs", icon: Users },
        { href: "/admin/exports", label: "Exports", icon: Download },
        { href: "/admin/emails", label: "Emails", icon: Mail },
        {
            href: "/admin/audit-log",
            label: "Journal d'activité",
            icon: ScrollText,
        },
        { href: "/admin/settings", label: "Paramètres", icon: Cog },
    ];
</script>

<AppShell
    links={adminNavLinks}
    user={data.profile ?? { full_name: "Administrateur", role: "admin" }}
    notificationsHref="/admin/notifications"
>
    {@render children()}
</AppShell>
