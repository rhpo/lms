<script lang="ts">
	import { Search } from 'lucide-svelte';

	let {
		placeholder = 'Rechercher...',
		onSearch = () => {},
		value = $bindable(''),
		searchInput = $bindable()
	} = $props();
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
<main onclick={() => searchInput?.focus()}>
	<input
		type="text"
		{placeholder}
		bind:this={searchInput}
		bind:value
		onkeydown={(e) => {
			if (e.key === 'Enter') {
				onSearch(value);
			}
		}}
	/>

	<button onclick={() => onSearch(value)}>
		<Search />
	</button>
</main>

<style>
	main {
		width: 100%;
		max-width: 56rem;
		display: inline-flex;
		align-items: center;
		background-color: var(--color-background-100);
		border-radius: 50px;
		box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
		cursor: text;
		height: 100%;
	}

	main input {
		flex: 1;
		height: 100%;
		width: 100%;
		padding-left: 2rem;
		font-size: 1.2rem;
		background: transparent;
		border: none;
		color: var(--color-text);
		font-family: var(--font-sans);
	}

	main input:focus {
		outline: none;
	}

	main input::placeholder {
		font-family: var(--font-sans);
		font-style: italic;
		color: var(--color-text-muted);
	}

	button {
		background-color: var(--color-accent);
		color: var(--color-background);
		border: none;
		border-radius: 50%;
		padding: 0.75rem;
		margin: 0.25rem;
		display: grid;
		place-items: center;
		cursor: pointer;
		transition: all var(--transition-normal);
	}

	button:hover {
		transform: scale(1.05) rotate(5deg);
	}

	@media screen and (max-width: 600px) {
		main input {
			padding-left: 1rem;
			font-size: 0.875rem;
		}
	}
</style>
