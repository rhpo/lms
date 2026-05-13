<script lang="ts">
	import {
		X,
		CheckCircle,
		AlertCircle,
		Info,
		AlertTriangle,
	} from "lucide-svelte";

	let {
		type = "info",
		message = "",
		visible = $bindable(true),
		duration = 0,
		onClose,
	}: {
		type?: "success" | "error" | "info" | "warning";
		message?: string;
		visible?: boolean;
		duration?: number;
		onClose?: () => void;
	} = $props();

	const icons = {
		success: CheckCircle,
		error: AlertCircle,
		info: Info,
		warning: AlertTriangle,
	};

	$effect(() => {
		if (duration > 0 && visible) {
			const timer = setTimeout(() => {
				visible = false;
				onClose?.();
			}, duration);
			return () => clearTimeout(timer);
		}
	});

	function close() {
		visible = false;
		onClose?.();
	}

	let IconComponent = $derived(icons[type]);
</script>

{#if visible}
	<div class="notification" data-type={type}>
		<div class="notification-icon">
			<IconComponent size={18} />
		</div>
		<span class="notification-message">{message}</span>
		<button class="notification-close" onclick={close} aria-label="Fermer">
			<X size={14} />
		</button>
	</div>
{/if}

<style>
	.notification {
		display: flex;
		align-items: center;
		gap: 0.75rem;
		padding: 0.75rem 1rem;
		border-radius: 10px;
		font-family: var(--font-sans);
		font-size: 0.9rem;
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
		animation: slideIn 0.2s ease-out;
	}

	@keyframes slideIn {
		from {
			opacity: 0;
			transform: translateY(-8px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	.notification[data-type="success"] {
		background: color-mix(
			in srgb,
			var(--color-accent) 10%,
			var(--color-background)
		);
		color: var(--color-accent);
		border: 1px solid
			color-mix(in srgb, var(--color-accent) 20%, transparent);
	}

	.notification[data-type="error"] {
		background: color-mix(in srgb, #ef4444 10%, var(--color-background));
		color: #ef4444;
		border: 1px solid color-mix(in srgb, #ef4444 20%, transparent);
	}

	.notification[data-type="info"] {
		background: color-mix(in srgb, #3b82f6 10%, var(--color-background));
		color: #3b82f6;
		border: 1px solid color-mix(in srgb, #3b82f6 20%, transparent);
	}

	.notification[data-type="warning"] {
		background: color-mix(in srgb, #f59e0b 10%, var(--color-background));
		color: #f59e0b;
		border: 1px solid color-mix(in srgb, #f59e0b 20%, transparent);
	}

	.notification-icon {
		display: flex;
		align-items: center;
		flex-shrink: 0;
	}

	.notification-message {
		flex: 1;
	}

	.notification-close {
		display: flex;
		align-items: center;
		justify-content: center;
		background: none;
		border: none;
		color: inherit;
		opacity: 0.6;
		cursor: pointer;
		padding: 0.2rem;
		border-radius: 4px;
		flex-shrink: 0;
		transition: opacity 0.15s;
	}

	.notification-close:hover {
		opacity: 1;
	}
</style>
