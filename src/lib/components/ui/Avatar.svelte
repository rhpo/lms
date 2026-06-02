<script lang="ts">
	interface Props {
		user: { full_name: string; avatar_url?: string | null };
		size?: number;
	}

	let { user, size = 34 }: Props = $props();

	const initials = $derived(
		user.full_name
			.split(" ")
			.map((n) => n[0])
			.slice(0, 2)
			.join("")
			.toUpperCase() ?? "U",
	);

	function stringToAvatar(str: string): {
		background: string;
		text: "black" | "white";
	} {
		let hash = 0;

		for (let i = 0; i < str.length; i++) {
			hash = str.charCodeAt(i) + ((hash << 5) - hash);
		}

		const hue = Math.abs(hash) % 360;
		const background = `hsl(${hue}, 70%, 50%)`;


		const lightness = 50;

		const text = lightness > 50 ? "black" : "white";

		return {
			background,
			text,
		};
	}

	let colorPalette = $derived(stringToAvatar(user.full_name));
	let backColor = $derived(colorPalette.background);
	let textColor = $derived(colorPalette.text);
</script>

{#if user?.avatar_url}
	<img
		class="avatar"
		src={user.avatar_url}
		alt={user.full_name}
		style:width="{size}px"
		style:height="{size}px"
	/>
{:else}
	<div
		class="avatar placeholder"
		style:height="{size}px"
		style:width="{size}px"
		style:--font-size="{size * 0.35}px"
		style:--text-color={textColor}
		style:background={backColor}
		aria-label="Avatar utilisateur"
		title={user?.full_name ?? "Utilisateur"}
	>
		<p>
			{initials}
		</p>
	</div>
{/if}

<style>
	.avatar {
		border-radius: 50%;
		object-fit: cover;
		flex-shrink: 0;
	}

	.placeholder {
		background: var(--color-accent);

		display: flex;
		align-items: center;
		justify-content: center;
		cursor: pointer;

		p {
			color: var(--text-color);

			font-weight: 700;
			font-size: var(--font-size);
			font-family: var(--font-sans);
		}
	}
</style>
