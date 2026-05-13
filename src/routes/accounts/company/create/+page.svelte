<script>
  import { UserPlus } from "lucide-svelte";

  import Logo from "$lib/components/Logo.svelte";
  import View from "$lib/components/ui/View.svelte";
  import Input from "$lib/components/ui/Input/Input.svelte";
  import Switch from "$lib/components/ui/Switch.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import EmailInput from "$lib/components/ui/Input/InputEmail.svelte";

  import { enhance } from "$app/forms";

  let { form } = $props();
  let checked = $state(false);
  let loading = $state(false);
</script>

<View fullScreen center>
  <div class="main" data-aos="flip-left">
    <div class="card-wrapper">
      <Logo width="180px" />

      <div class="header">
        <h1 class="title">Company</h1>
        <p>Create a PFE Management company account.</p>
      </div>

      {#if form?.error}
        <div class="error-banner">
          <p>{form.error}</p>
        </div>
      {/if}

      <form
        method="POST"
        use:enhance={() => {
          loading = true;
          return async ({ update }) => {
            loading = false;
            await update();
          };
        }}
        class="form-wrapper"
      >
        <div class="form">
          <EmailInput
            name="email"
            placeholder="email@company.com"
            value={form?.email ?? ''}
          />

          <Input
            name="password"
            autocomplete="off"
            type="password"
            placeholder="password *"
          />

          <Input
            name="confirmPassword"
            autocomplete="off"
            type="password"
            placeholder="confirm password *"
          />
        </div>

        <div class="actions">
          <div class="condition">
            <Switch size={25} bind:checked>
              {#snippet label()}
                <p class="content">
                  I have read and accepted the
                  <a href="/conditions" target="_blank">Terms & Conditions</a>.
                </p>
              {/snippet}
            </Switch>
            <input type="hidden" name="acceptedTos" value={checked ? 'on' : ''} />
          </div>

          <Button
            type="submit"
            disabledLabel={"Please accept the TOS."}
            disabled={!checked || loading}
            Icon={UserPlus}
          >
            {loading ? 'Creating...' : 'Create account'}
          </Button>
        </div>
      </form>

      <div class="questions">
        <p class="question">
          <span class="content">
            Already have an account?
            <a href="./login">Log-in</a>.
          </span>
        </p>
      </div>
    </div>
  </div>
</View>

<style>
  .main {
    width: 100%;
    max-width: 520px;
  }

  .questions {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .header {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    text-align: center;
    margin-bottom: 2rem;

    .title {
      display: flex;
      justify-content: center;
      width: fit-content;
      align-items: center;
      gap: 0.5rem;
    }
  }

  .error-banner {
    background: color-mix(in srgb, #ef4444 15%, transparent);
    border: 1px solid #ef4444;
    border-radius: 0.5rem;
    padding: 0.75rem 1rem;
    width: 100%;
    text-align: center;

    p {
      color: #ef4444;
      font-size: 0.95rem;
      margin: 0;
    }
  }

  .form-wrapper {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    width: 100%;

    .form {
      display: flex;
      flex-direction: column;
      width: 100%;
      gap: 0.5rem;
      margin-bottom: 3rem;
    }

    .actions {
      display: flex;
      flex-direction: column;
      width: 100%;
      gap: 0.7rem;

      margin-bottom: 2rem;
    }
  }

  .condition {
    font-family: var(--font-sans);

    display: flex;
    justify-content: center;
    align-items: center;
    gap: 0.5rem;

    .content,
    .content * {
      font-size: 1rem;
    }

    .content {
      color: var(--color-text-muted);
    }
  }

  .card-wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
  }

  .question {
    font-size: 1.2rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;

    .content {
      transform: translateY(2px);
    }
  }
</style>
