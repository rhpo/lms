<script lang="ts">
  import { page } from "$app/state";
  import { CircleQuestionMark, TriangleAlert } from "lucide-svelte";

  import View from "$lib/components/ui/View.svelte";
</script>

<View center fullScreen>
  <div class="error" data-aos="flip-up">
    <div class="icon">
      {#if page.status === 404}
        <CircleQuestionMark size="8rem" />
      {:else}
        <TriangleAlert size="8rem" />

        <h1>
          {page.status}
        </h1>
      {/if}
    </div>

    <div class="divider"></div>

    <div class="info">
      <div class="header">
        <h1>Oops...</h1>

        <p>
          {#if page.status === 404}
            This page was not found, please <a href="/">Go back</a> to the site.
          {:else}
            Something unexpected happened, please contact the administrator
            <a href="mailto:ramyhadid.js@gmail.com">ramyhadid.js@gmail.com</a>.
          {/if}
        </p>
      </div>

      <div class="message">
        <p>Error message: <code>{page?.error?.message}</code></p>
      </div>
    </div>
  </div>
</View>

<style>
  .error {
    display: flex;
    gap: 2.5rem;

    .icon,
    .info {
      padding: 2rem;
    }

    .icon {
      display: flex;
      flex-direction: column;
      align-items: center;

      gap: 0.5rem;
    }

    .divider {
      border-right: 1px solid var(--color-text);
    }

    .info {
      flex: 1;

      display: flex;
      flex-direction: column;
      align-content: center;
      justify-content: space-around;
    }
  }

  @media screen and (max-width: 500px) {
    .error {
      flex-direction: column;
      margin: 1rem 2rem;

      .divider {
        border-right: none;
        border-bottom: 1px solid var(--color-text);
      }

      .info {
        text-align: center;

        gap: 2rem;
      }
    }
  }
</style>
