<script lang="ts">
	import { ChevronLeft, ChevronRight } from 'lucide-svelte';

	let {
		events = [],
	}: {
		events?: Array<{
			date: string;
			title: string;
			time?: string;
			type?: 'defense' | 'meeting' | 'deadline';
		}>;
	} = $props();

	const months = [
		'Janvier', 'Fevrier', 'Mars', 'Avril', 'Mai', 'Juin',
		'Juillet', 'Aout', 'Septembre', 'Octobre', 'Novembre', 'Decembre',
	];

	const daysOfWeek = ['Lun', 'Mar', 'Mer', 'Jeu', 'Ven', 'Sam', 'Dim'];

	let currentMonth = $state(new Date().getMonth());
	let currentYear = $state(new Date().getFullYear());

	function daysInMonth(month: number, year: number): number {
		return new Date(year, month + 1, 0).getDate();
	}

	function firstDayOfMonth(month: number, year: number): number {


		const day = new Date(year, month, 1).getDay();
		return day === 0 ? 6 : day - 1;
	}

	function prevMonth() {
		if (currentMonth === 0) {
			currentMonth = 11;
			currentYear--;
		} else {
			currentMonth--;
		}
	}

	function nextMonth() {
		if (currentMonth === 11) {
			currentMonth = 0;
			currentYear++;
		} else {
			currentMonth++;
		}
	}

	function getEventsForDay(day: number) {
		const dateStr = `${currentYear}-${String(currentMonth + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}`;
		return events.filter(e => e.date === dateStr);
	}

	const totalDays = $derived(daysInMonth(currentMonth, currentYear));
	const startOffset = $derived(firstDayOfMonth(currentMonth, currentYear));
	const today = $derived(new Date());
</script>

<div class="calendar">
	<div class="calendar-header">
		<button class="nav-btn" onclick={prevMonth} aria-label="Mois precedent">
			<ChevronLeft size={18} />
		</button>
		<h3 class="calendar-title">{months[currentMonth]} {currentYear}</h3>
		<button class="nav-btn" onclick={nextMonth} aria-label="Mois suivant">
			<ChevronRight size={18} />
		</button>
	</div>

	<div class="calendar-grid">
		{#each daysOfWeek as day}
			<div class="day-header">{day}</div>
		{/each}

		{#each Array(startOffset) as _}
			<div class="day empty"></div>
		{/each}

		{#each Array(totalDays) as _, i}
			{@const day = i + 1}
			{@const dayEvents = getEventsForDay(day)}
			{@const isToday = day === today.getDate() && currentMonth === today.getMonth() && currentYear === today.getFullYear()}
			<div class="day" class:today={isToday} class:has-events={dayEvents.length > 0}>
				<span class="day-number">{day}</span>
				{#if dayEvents.length > 0}
					<div class="day-events">
						{#each dayEvents.slice(0, 2) as event}
							<div class="event-dot" data-type={event.type || 'defense'} title={event.title}></div>
						{/each}
						{#if dayEvents.length > 2}
							<span class="more">+{dayEvents.length - 2}</span>
						{/if}
					</div>
				{/if}
			</div>
		{/each}
	</div>
</div>

<style>
	.calendar {
		width: 100%;
		background: var(--color-background-100);
		border-radius: 12px;
		border: 1px solid var(--color-background-200);
		overflow: hidden;
		font-family: var(--font-sans);
	}

	.calendar-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: 1rem;
		border-bottom: 1px solid var(--color-background-200);
	}

	.calendar-title {
		margin: 0;
		font-size: 1rem;
		font-weight: 600;
		color: var(--color-text);
	}

	.nav-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		border: none;
		background: transparent;
		color: var(--color-text);
		border-radius: 8px;
		cursor: pointer;
		transition: background 0.15s;
	}

	.nav-btn:hover {
		background: var(--color-background-200);
	}

	.calendar-grid {
		display: grid;
		grid-template-columns: repeat(7, 1fr);
		padding: 0.5rem;
		gap: 2px;
	}

	.day-header {
		text-align: center;
		font-size: 0.75rem;
		font-weight: 600;
		color: var(--color-text-muted);
		padding: 0.5rem 0;
		text-transform: uppercase;
	}

	.day {
		aspect-ratio: 1;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: flex-start;
		padding: 0.25rem;
		border-radius: 8px;
		cursor: default;
		position: relative;
		min-height: 60px;
	}

	.day.empty {
		pointer-events: none;
	}

	.day.today {
		background: color-mix(in srgb, var(--color-accent) 10%, transparent);
	}

	.day-number {
		font-size: 0.85rem;
		font-weight: 500;
		color: var(--color-text);
	}

	.day.today .day-number {
		color: var(--color-accent);
		font-weight: 700;
	}

	.day-events {
		display: flex;
		flex-wrap: wrap;
		gap: 2px;
		margin-top: 2px;
		justify-content: center;
	}

	.event-dot {
		width: 6px;
		height: 6px;
		border-radius: 50%;
	}

	.event-dot[data-type="defense"] {
		background: var(--color-accent);
	}

	.event-dot[data-type="meeting"] {
		background: #3b82f6;
	}

	.event-dot[data-type="deadline"] {
		background: #ef4444;
	}

	.more {
		font-size: 0.65rem;
		color: var(--color-text-muted);
		font-weight: 600;
	}
</style>
