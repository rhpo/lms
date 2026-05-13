import { Book, CircleQuestionMark, House, LayoutDashboard } from "lucide-svelte"

type NavigationItem = {
    url: string;
    name: string;
    icon: typeof Book;
}

export const NAV_LINKS: NavigationItem[] = [
    { name: "Home", icon: House, url: "/" },
    { name: "About", icon: Book, url: "/about" },
    { name: "Guide", icon: CircleQuestionMark, url: "/guide" },
    { name: "Dashboard", icon: LayoutDashboard, url: "/dashboard" },
];
