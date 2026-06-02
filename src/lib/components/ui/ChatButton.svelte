<script lang="ts">
	import { ArrowDown, MessageCircle } from 'lucide-svelte';


	let isOpen = $state(false);

	let chatButton: HTMLButtonElement | null = $state(null);
	let customButton: HTMLButtonElement | null = $state(null);

	var tries = 0;
	var maxTries = 10;

	function tryBind() {
		if (tries >= maxTries) {
			console.warn('Chat button not found after maximum attempts');
			return;
		}


		chatButton = document.querySelector('#chatbase-bubble-button') as HTMLButtonElement;

		if (!chatButton) {
			console.log('Chat button not found, retrying...');
			tries++;
			setTimeout(tryBind, 500);
		}
	}
	$effect(() => {

		tryBind();
	});

	function toggleChat() {
		if (chatButton) {
			chatButton.click();
			isOpen = !isOpen;
		}
	}
</script>

<button
	class="chat-toggle bounce"
	onclick={toggleChat}
	aria-label="Toggle chat"
	bind:this={customButton}
	class:open={isOpen}
>
	{#if isOpen}
		<ArrowDown />
	{:else}
		<MessageCircle />
	{/if}

	IA
</button>

<style>
	.chat-toggle {
		color: var(--color-text);
		position: fixed;
		right: 2rem;
		bottom: 2rem;
		z-index: 999999999;
		background: var(--color-surface);
		border: 2px solid var(--color-border);
		border-radius: 50px;
		height: 50px;
		padding: 0.75rem;
		cursor: pointer;
		transition: all var(--transition-normal);
		box-shadow: var(--shadow-sm);

		display: flex;
		justify-content: center;
		align-items: center;

		gap: 0.5rem;
	}

	.chat-toggle.open {
		background: var(--color-accent);
		color: #fff;
		transform: translateY(15px);
	}

	.chat-toggle:hover {
		box-shadow: var(--shadow-md);
	}

	.chat-toggle:not(.open):hover {
		transform: scale(1.1);
	}

	@media (max-width: 768px) {
		.chat-toggle {
			right: 1rem;
			padding: 0.5rem;
		}
	}
</style>
