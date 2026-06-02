<script lang="ts">
  import { CalendarDays, ChevronLeft, ChevronRight } from "lucide-svelte";

  let {
    value = $bindable(""),
    required = false,
    min,
    onchange,
    placeholder = "Choisir une date",
  }: {
    value?: string;
    required?: boolean;
    min?: string;
    onchange?: (e: Event) => void;
    placeholder?: string;
  } = $props();

  let showCalendar = $state(false);
  let calendarEl = $state<HTMLDivElement>();
  let inputEl = $state<HTMLDivElement>();


  const selectedDate = $derived(value ? new Date(value + "T00:00:00") : null);


  let viewYear = $state(new Date().getFullYear());
  let viewMonth = $state(new Date().getMonth());


  $effect(() => {
    if (showCalendar && selectedDate) {
      viewYear = selectedDate.getFullYear();
      viewMonth = selectedDate.getMonth();
    } else if (showCalendar) {
      viewYear = new Date().getFullYear();
      viewMonth = new Date().getMonth();
    }
  });


  $effect(() => {
    if (!showCalendar) return;
    function handleClick(e: MouseEvent) {
      if (
        calendarEl &&
        !calendarEl.contains(e.target as Node) &&
        inputEl &&
        !inputEl.contains(e.target as Node)
      ) {
        showCalendar = false;
      }
    }
    document.addEventListener("click", handleClick);
    return () => document.removeEventListener("click", handleClick);
  });

  const MONTHS = [
    "Janvier",
    "Fevrier",
    "Mars",
    "Avril",
    "Mai",
    "Juin",
    "Juillet",
    "Aout",
    "Septembre",
    "Octobre",
    "Novembre",
    "Decembre",
  ];

  const DAYS = ["Di", "Lu", "Ma", "Me", "Je", "Ve", "Sa"];

  function daysInMonth(year: number, month: number): number {
    return new Date(year, month + 1, 0).getDate();
  }

  function firstDayOfMonth(year: number, month: number): number {
    return new Date(year, month, 1).getDay();
  }

  function isBeforeMin(date: string): boolean {
    return !!min && date < min;
  }

  const calendarDays = $derived.by(() => {
    const totalDays = daysInMonth(viewYear, viewMonth);
    const startDay = firstDayOfMonth(viewYear, viewMonth);
    const days: Array<{ day: number; currentMonth: boolean; date: string }> =
      [];


    const prevMonthDays = daysInMonth(
      viewYear,
      viewMonth - 1 < 0 ? 11 : viewMonth - 1,
    );
    const prevMonth = viewMonth - 1 < 0 ? 11 : viewMonth - 1;
    const prevYear = viewMonth - 1 < 0 ? viewYear - 1 : viewYear;
    for (let i = startDay - 1; i >= 0; i--) {
      const d = prevMonthDays - i;
      const m = String(prevMonth + 1).padStart(2, "0");
      const day = String(d).padStart(2, "0");
      days.push({
        day: d,
        currentMonth: false,
        date: `${prevYear}-${m}-${day}`,
      });
    }


    const m = String(viewMonth + 1).padStart(2, "0");
    for (let d = 1; d <= totalDays; d++) {
      const day = String(d).padStart(2, "0");
      days.push({
        day: d,
        currentMonth: true,
        date: `${viewYear}-${m}-${day}`,
      });
    }


    const nextMonth = viewMonth + 1 > 11 ? 0 : viewMonth + 1;
    const nextYear = viewMonth + 1 > 11 ? viewYear + 1 : viewYear;
    const remaining = 42 - days.length;
    const nm = String(nextMonth + 1).padStart(2, "0");
    for (let d = 1; d <= remaining; d++) {
      const day = String(d).padStart(2, "0");
      days.push({
        day: d,
        currentMonth: false,
        date: `${nextYear}-${nm}-${day}`,
      });
    }

    return days;
  });

  function selectDate(dateStr: string) {
    if (isBeforeMin(dateStr)) return;

    value = dateStr;
    showCalendar = false;

    if (onchange) {
      const event = new Event("change", { bubbles: true });
      onchange(event);
    }
  }
  function prevMonth() {
    if (viewMonth === 0) {
      viewMonth = 11;
      viewYear--;
    } else {
      viewMonth--;
    }
  }

  function nextMonth() {
    if (viewMonth === 11) {
      viewMonth = 0;
      viewYear++;
    } else {
      viewMonth++;
    }
  }

  function today() {
    const now = new Date();
    const y = now.getFullYear();
    const m = String(now.getMonth() + 1).padStart(2, "0");
    const d = String(now.getDate()).padStart(2, "0");
    selectDate(`${y}-${m}-${d}`);
  }

  function clearDate() {
    value = "";
    showCalendar = false;
    if (onchange) {
      const event = new Event("change", { bubbles: true });
      onchange(event);
    }
  }

  function formatDisplay(iso: string): string {
    if (!iso) return "";
    const [y, m, d] = iso.split("-");
    return `${d}/${m}/${y}`;
  }

  function isSelected(day: { date: string }): boolean {
    return value === day.date;
  }

  function isToday(day: { date: string }): boolean {
    const now = new Date();
    const y = now.getFullYear();
    const m = String(now.getMonth() + 1).padStart(2, "0");
    const d = String(now.getDate()).padStart(2, "0");
    return day.date === `${y}-${m}-${d}`;
  }
</script>

<div class="date-input-wrapper">
  <div
    class="date-input-trigger"
    class:has-value={!!value}
    role="button"
    tabindex="0"
    onclick={() => (showCalendar = !showCalendar)}
    onkeydown={(e) => {
      if (e.key === "Enter") showCalendar = !showCalendar;
    }}
    bind:this={inputEl}
  >
    <CalendarDays size={16} class="icon" />
    <span class="display-text">
      {value ? formatDisplay(value) : placeholder}
    </span>
    {#if value}
      <button
        class="clear-btn"
        onclick={(e) => {
          e.stopPropagation();
          clearDate();
        }}
        aria-label="Effacer la date"
      >
        &times;
      </button>
    {/if}
  </div>

  {#if showCalendar}
    <div class="calendar-dropdown" bind:this={calendarEl}>
      <div class="calendar-header">
        <button class="nav-btn" onclick={prevMonth} aria-label="Mois precedent">
          <ChevronLeft size={16} />
        </button>
        <span class="month-year">{MONTHS[viewMonth]} {viewYear}</span>
        <button class="nav-btn" onclick={nextMonth} aria-label="Mois suivant">
          <ChevronRight size={16} />
        </button>
      </div>

      <div class="calendar-grid">
        {#each DAYS as day}
          <div class="day-header">{day}</div>
        {/each}

        {#each calendarDays as day}
          <button
            class="day-cell"
            class:other-month={!day.currentMonth}
            class:selected={isSelected(day)}
            class:today={isToday(day)}
            onclick={() => selectDate(day.date)}
          >
            {day.day}
          </button>
        {/each}
      </div>

      <div class="calendar-footer">
        <button class="today-btn" onclick={today}>Aujourd'hui</button>
      </div>
    </div>
  {/if}
</div>

<!-- Hidden native input for form submission -->
<input type="hidden" {value} {required} />

<style>
  .date-input-wrapper {
    position: relative;
    width: 100%;
  }

  .date-input-trigger {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 0.75rem;
    border: 1px solid var(--color-border);
    border-radius: 8px;
    background: var(--color-surface);
    color: var(--color-text-muted);
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    cursor: pointer;
    transition:
      border-color var(--transition-fast),
      box-shadow var(--transition-fast);
    user-select: none;
    min-height: 38px;
  }

  .date-input-trigger:hover {
    border-color: var(--color-accent);
  }

  .date-input-trigger:focus-visible {
    outline: none;
    border-color: var(--color-accent);
    box-shadow: 0 0 0 3px
      color-mix(in srgb, var(--color-accent) 20%, transparent);
  }

  .date-input-trigger.has-value {
    color: var(--color-text);
  }

  .icon {
    flex-shrink: 0;
    color: var(--color-text-muted);
  }

  .display-text {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .clear-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    border: none;
    background: none;
    color: var(--color-text-muted);
    font-size: 1.1rem;
    cursor: pointer;
    border-radius: 50%;
    padding: 0;
    line-height: 1;
    flex-shrink: 0;
  }

  .clear-btn:hover {
    background: var(--color-background-100);
    color: var(--color-text);
  }

  .calendar-dropdown {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    z-index: 100;
    background: var(--color-surface);
    border: 1px solid var(--color-border);
    border-radius: 12px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
    padding: 0.75rem;
    width: 280px;
    animation: fadeIn 0.15s ease-out;
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(-4px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .calendar-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 0.75rem;
  }

  .nav-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    background: none;
    color: var(--color-text-muted);
    cursor: pointer;
    border-radius: 8px;
    transition: background var(--transition-fast);
  }

  .nav-btn:hover {
    background: var(--color-background-100);
    color: var(--color-text);
  }

  .month-year {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text);
    font-family: var(--font-sans);
  }

  .calendar-grid {
    display: grid;
    grid-template-columns: repeat(7, 1fr);
    gap: 2px;
  }

  .day-header {
    text-align: center;
    font-size: 0.7rem;
    font-weight: 600;
    color: var(--color-text-muted);
    padding: 0.35rem 0;
    font-family: var(--font-sans);
    text-transform: uppercase;
  }

  .day-cell {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 100%;
    aspect-ratio: 1;
    border: none;
    background: none;
    color: var(--color-text);
    font-size: var(--text-sm);
    font-family: var(--font-sans);
    cursor: pointer;
    border-radius: 8px;
    transition:
      background var(--transition-fast),
      color var(--transition-fast);
  }

  .day-cell:hover {
    background: var(--color-background-100);
  }

  .day-cell.other-month {
    color: var(--color-text-muted);
    opacity: 0.4;
  }

  .day-cell.today {
    font-weight: 700;
    color: var(--color-accent);
    border: 1px solid var(--color-accent);
  }

  .day-cell.selected {
    background: var(--color-accent);
    color: #fff;
    font-weight: 600;
  }

  .day-cell.selected:hover {
    background: color-mix(in srgb, var(--color-accent) 85%, black);
  }

  .calendar-footer {
    display: flex;
    justify-content: center;
    margin-top: 0.5rem;
    padding-top: 0.5rem;
    border-top: 1px solid var(--color-border);
  }

  .today-btn {
    border: none;
    background: none;
    color: var(--color-accent);
    font-size: var(--text-xs);
    font-weight: 600;
    font-family: var(--font-sans);
    cursor: pointer;
    padding: 0.35rem 0.75rem;
    border-radius: 6px;
    transition: background var(--transition-fast);
  }

  .today-btn:hover {
    background: color-mix(in srgb, var(--color-accent) 10%, transparent);
  }
</style>
