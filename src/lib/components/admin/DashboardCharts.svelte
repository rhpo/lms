<script lang="ts">
  import { onMount } from "svelte";
  import Chart from "chart.js/auto";

  let { stats = {}, timeline = {}, theme = "light" } = $props();

  let lineCanvas: HTMLCanvasElement;
  let donutCanvas: HTMLCanvasElement;

  let lineChart: Chart | null = null;
  let donutChart: Chart | null = null;

  function getThemeColors() {
    if (typeof document === "undefined") {
      return {
        text: "#111827",
        grid: "rgba(0, 0, 0, 0.1)",
        surface: "#ffffff",
        border: "#e5e7eb",
      };
    }


    const headerEl = document.querySelector(".chart-header h3") as HTMLElement;
    let textRaw = "#111827";
    if (headerEl) {
      const oldTrans = headerEl.style.transition;
      headerEl.style.transition = "none";
      textRaw = getComputedStyle(headerEl).color;
      headerEl.style.transition = oldTrans;
    }


    const dashCard = document.querySelector(".chart-card") as HTMLElement;
    let borderRaw = "#e5e7eb";
    let surfaceRaw = "#ffffff";
    if (dashCard) {
      const oldTrans = dashCard.style.transition;
      dashCard.style.transition = "none";
      const styles = getComputedStyle(dashCard);
      borderRaw = styles.borderColor;
      surfaceRaw = styles.backgroundColor;
      dashCard.style.transition = oldTrans;
    }


    let gridRaw = borderRaw;
    if (gridRaw.startsWith("rgb(")) {
      gridRaw = gridRaw.replace("rgb(", "rgba(").replace(")", ", 0.3)");
    } else if (gridRaw.startsWith("rgba(")) {

      gridRaw = gridRaw.replace(/[\d\.]+\)$/g, "0.3)");
    }

    return {
      text: textRaw,
      grid: gridRaw,
      surface: surfaceRaw,
      border: borderRaw,
    };
  }

  onMount(() => {
    const colors = getThemeColors();


    lineChart = new Chart(lineCanvas, {
      type: "line",
      data: {
        labels: timeline.labels || [
          "Sep",
          "Oct",
          "Nov",
          "Dec",
          "Jan",
          "Feb",
          "Mar",
          "Apr",
          "Mai",
          "Juin",
        ],
        datasets: [
          {
            label: "Soumis le mémoire",
            data: timeline.soumis_memoire || [0, 0, 0, 0, 0, 0, 5, 20, 60, 115],
            borderColor: "#10b981", // green
            backgroundColor: "color-mix(in srgb, #10b981 20%, transparent)",
            fill: true,
            tension: 0.4,
            borderWidth: 3,
            pointBackgroundColor: "#10b981",
          },
          {
            label: "Avec sujet PFE",
            data: timeline.avec_sujet || [
              0, 10, 30, 60, 85, 100, 110, 118, 120, 120,
            ],
            borderColor: "#3b82f6", // blue
            backgroundColor: "color-mix(in srgb, #3b82f6 20%, transparent)",
            fill: true,
            tension: 0.4,
            borderWidth: 3,
            pointBackgroundColor: "#3b82f6",
          },
          {
            label: "Sans sujet PFE",
            data: timeline.sans_sujet || [
              120, 110, 90, 60, 35, 20, 10, 2, 0, 0,
            ],
            borderColor: "#f43f5e",
            backgroundColor: "color-mix(in srgb, #f43f5e 20%, transparent)",
            fill: true,
            tension: 0.4,
            borderWidth: 3,
            pointBackgroundColor: "#f43f5e",
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
          mode: "index",
          intersect: false,
        },
        plugins: {
          legend: {
            position: "bottom",
            labels: {
              usePointStyle: true,
              padding: 20,
              color: colors.text,
              font: {
                family: "ui-sans-serif, system-ui, sans-serif",
                size: 12,
                weight: "normal",
              },
            },
          },
          tooltip: {
            backgroundColor: colors.surface,
            titleColor: colors.text,
            bodyColor: colors.text,
            borderColor: colors.border,
            borderWidth: 1,
            padding: 12,
            boxPadding: 4,
            usePointStyle: true,
            titleFont: {
              size: 14,
              family: "ui-sans-serif, system-ui, sans-serif",
            },
            bodyFont: {
              size: 13,
              family: "ui-sans-serif, system-ui, sans-serif",
            },
          },
        },
        scales: {
          x: {
            grid: { color: colors.grid },
            ticks: {
              color: colors.text,
              font: { family: "ui-sans-serif, system-ui, sans-serif" },
            },
          },
          y: {
            grid: { color: colors.grid },
            ticks: {
              color: colors.text,
              font: { family: "ui-sans-serif, system-ui, sans-serif" },
            },
          },
        },
      },
    });


    donutChart = new Chart(donutCanvas, {
      type: "doughnut",
      data: {
        labels: [
          "En attente validation",
          "Validés & Dispos",
          "Affectés",
          "Rejetés",
        ],
        datasets: [
          {
            data: [
              stats.pendingSubjects || 0,
              (stats.validatedSubjects || 0) - (stats.assignedSubjects || 0),
              stats.assignedSubjects || 0,
              stats.rejectedSubjects || 0,
            ],
            backgroundColor: [
              "#f59e0b",
              "#8b5cf6",
              "#3b82f6",
              "#ef4444",
            ],
            borderWidth: 0,
            hoverOffset: 8,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        cutout: "75%",
        plugins: {
          legend: {
            display: false,
          },
          tooltip: {
            backgroundColor: colors.surface,
            titleColor: colors.text,
            bodyColor: colors.text,
            borderColor: colors.border,
            borderWidth: 1,
            padding: 12,
            boxPadding: 4,
            usePointStyle: true,
            titleFont: {
              size: 14,
              family: "ui-sans-serif, system-ui, sans-serif",
            },
            bodyFont: {
              size: 13,
              family: "ui-sans-serif, system-ui, sans-serif",
            },
          },
        },
      },
    });


    Chart.defaults.color = colors.text;
    Chart.defaults.borderColor = colors.grid;
  });


  $effect(() => {


    const currentTheme = theme;

    if (lineChart && timeline) {
      if (timeline.labels) lineChart.data.labels = timeline.labels;
      if (timeline.soumis_memoire)
        lineChart.data.datasets[0].data = timeline.soumis_memoire;
      if (timeline.avec_sujet)
        lineChart.data.datasets[1].data = timeline.avec_sujet;
      if (timeline.sans_sujet)
        lineChart.data.datasets[2].data = timeline.sans_sujet;
      const colors = getThemeColors();


      Chart.defaults.color = colors.text;
      Chart.defaults.borderColor = colors.grid;

      if (lineChart.options.plugins?.legend?.labels) {
        lineChart.options.plugins.legend.labels.color = colors.text;
      }
      if (lineChart.options.scales?.x?.ticks) {
        lineChart.options.scales.x.ticks.color = colors.text;
      }
      if (lineChart.options.scales?.y?.ticks) {
        lineChart.options.scales.y.ticks.color = colors.text;
      }
      if (lineChart.options.scales?.x?.grid) {
        lineChart.options.scales.x.grid.color = colors.grid;
      }
      if (lineChart.options.scales?.y?.grid) {
        lineChart.options.scales.y.grid.color = colors.grid;
      }

      lineChart.update();
    }

    if (donutChart && stats) {
      const pending = stats.pendingSubjects || 0;
      const assigned = stats.assignedSubjects || 0;
      const validated = stats.validatedSubjects || 0;
      const rejected = stats.rejectedSubjects || 0;

      const valdisp = Math.max(0, validated - assigned);

      donutChart.data.datasets[0].data = [pending, valdisp, assigned, rejected];

      const colors = getThemeColors();
      if (donutChart.options.plugins?.tooltip) {
        donutChart.options.plugins.tooltip.backgroundColor = colors.surface;
        donutChart.options.plugins.tooltip.titleColor = colors.text;
        donutChart.options.plugins.tooltip.bodyColor = colors.text;
        donutChart.options.plugins.tooltip.borderColor = colors.border;
      }

      donutChart.update();
    }
  });
</script>

<div class="charts-wrapper" data-aos="fade-up" data-aos-delay="100">
  <div class="chart-card line-chart">
    <div class="chart-header">
      <h3>Progression des Étudiants</h3>
      <span class="badge">Année Académique Courante</span>
    </div>
    <div class="canvas-container">
      <canvas bind:this={lineCanvas}></canvas>
    </div>
  </div>

  <div class="chart-card donut-chart">
    <div class="chart-header">
      <h3>État Global des Sujets PFE</h3>
    </div>
    <div class="canvas-container inner-donut-wrapper">
      <div class="inner-donut-text">
        <span class="total">{stats.totalSubjects || 0}</span>
        <span class="label">Total Sujets</span>
      </div>
      <canvas bind:this={donutCanvas}></canvas>
    </div>

    <div class="custom-legend">
      <div class="legend-item">
        <span class="dot" style="background:#f59e0b;"></span> En attente validation
      </div>
      <div class="legend-item">
        <span class="dot" style="background:#8b5cf6;"></span> Validés & Dispos
      </div>
      <div class="legend-item">
        <span class="dot" style="background:#3b82f6;"></span> Affectés
      </div>
      <div class="legend-item">
        <span class="dot" style="background:#ef4444;"></span> Rejetés
      </div>
    </div>
  </div>
</div>

<style>
  .charts-wrapper {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 1.5rem;
  }

  .chart-card {
    position: relative;
    background: var(--color-surface, white);
    border: 1px solid var(--color-border, #e5e7eb);
    border-radius: 12px;
    padding: 1.5rem;
    box-shadow: 0 4px 15px -3px color-mix(in srgb, var(--color-text, #000) 5%, transparent);
    display: flex;
    flex-direction: column;
    min-height: 400px;
  }

  .chart-card:hover {
    box-shadow: 0 10px 25px -5px color-mix(in srgb, var(--color-text, #000) 8%, transparent);
  }

  .chart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
  }

  h3 {
    font-size: 1.15rem;
    font-weight: 600;
    color: var(--color-text, #111827);
    margin: 0;
    font-family: var(--font-sans, system-ui);
  }

  .badge {
    font-size: 0.75rem;
    padding: 4px 10px;
    background: color-mix(
      in srgb,
      var(--color-accent, #3b82f6) 15%,
      transparent
    );
    color: var(--color-accent, #3b82f6);
    border-radius: 9999px;
    font-weight: 600;
    font-family: var(--font-sans, system-ui);
  }

  .canvas-container {
    position: relative;
    flex: 1;
    width: 100%;
    min-height: 300px;

    :global(canvas) {
      position: relative;
      z-index: 10;
    }
  }

  .inner-donut-text {
    position: absolute;
    top: 50%; /* now exactly 50% because HTML legend doesn't skew canvas */
    left: 50%;
    transform: translate(-50%, -50%);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    pointer-events: none;
    z-index: 10;
  }

  .inner-donut-text .total {
    font-size: 2.5rem;
    font-weight: 800;
    line-height: 1;
    color: var(--color-text, #111827);
    font-family: var(--font-sans, system-ui);
  }

  .inner-donut-text .label {
    font-size: 0.85rem;
    font-weight: 500;
    color: color-mix(in srgb, var(--color-text, #4b5563) 70%, transparent);
    font-family: var(--font-sans, system-ui);
    margin-top: 4px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .custom-legend {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 0.5rem;
    margin-top: 1.5rem;
  }

  .legend-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 0.85rem;
    font-family: var(--font-sans, system-ui);
    color: var(--color-text, #4b5563);
    font-weight: 300;
  }

  .dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    display: inline-block;
  }

  @media (max-width: 1100px) {
    .charts-wrapper {
      grid-template-columns: 1fr;
    }
    .chart-card {
      min-height: 350px;
    }
  }
</style>
