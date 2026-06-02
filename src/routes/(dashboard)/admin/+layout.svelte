<script lang="ts">
  import AppShell from "$lib/components/ui/AppShell.svelte";
  import { notificationStore } from "$lib/stores/notifications";
  import {
    LayoutDashboard,
    Users,
    GraduationCap,
    Gavel,
    FileText,
    Mail,
    Cog,
    Bell,
  } from "lucide-svelte";

  let { children, data } = $props();

  let unreadCount = $state(0);
  notificationStore.subscribe((n) => (unreadCount = n));

  const adminNavLinks = $derived([
    {
      href: "/admin/dashboard",
      label: "Tableau de bord",
      icon: LayoutDashboard,
    },

    { href: "/admin/users", label: "Utilisateurs", icon: Users },
    { href: "/admin/subjects", label: "Sujets", icon: FileText },
    { href: "/admin/pfe", label: "PFE", icon: GraduationCap },


    { href: "/admin/emails", label: "Emails", icon: Mail },

    {
      href: "/admin/notifications",
      label: "Notifications",
      icon: Bell,
      count: unreadCount,
    },
    { href: "/admin/settings", label: "Paramètres", icon: Cog },
  ]);
</script>

<AppShell
  links={adminNavLinks}
  quickAccess={[

    {
      href: "/admin/subjects?status=pending",
      label: "Sujets à valider",
      icon: FileText,
    },
    {
      href: "/admin/notifications",
      label: "Notifications",
      icon: Bell,
    },
  ]}
  user={data.profile ?? { full_name: "Administrateur", role: "admin" }}
>
  {@render children()}
</AppShell>
