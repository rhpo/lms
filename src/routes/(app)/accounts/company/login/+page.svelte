<script lang="ts">
  import { tippy } from "$lib/scripts/tippy";
  import { BRAND } from "$lib/constants/branding";
  import { GraduationCap, User } from "lucide-svelte";

  import Logo from "$lib/components/Logo.svelte";
  import View from "$lib/components/ui/View.svelte";
  import Switch from "$lib/components/ui/Switch.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import InputEmail from "$lib/components/ui/Input/InputEmail.svelte";
  import InputPassword from "$lib/components/ui/Input/InputPassword.svelte";

  import { authStore } from "$lib/stores/auth";

  let checked = $state(false);
  let loading = $state(false);
  let email = $state("");
  let password = $state("");
  let error = $state("");

  async function handleLogin() {
    loading = true;
    error = "";
    try {
      await authStore.devLogin(email);
    } catch (err: any) {
      error = err.message || "Failed to login";
    } finally {
      loading = false;
    }
  }
</script>

<View fullScreen center>
  <div class="main" data-aos="flip-left">
    <div class="card-wrapper">
      <Logo width="180px" />

      <div class="header">
        <h1 class="title">Company</h1>
        <p>Connect as a company to manage external PFE's.</p>
      </div>

      {#if error}
        <div class="error-banner">
          <p>{error}</p>
        </div>
      {/if}

      <form
        onsubmit={(e) => {
          e.preventDefault();
          handleLogin();
        }}
        class="form-wrapper"
      >
        <div class="form">
          <InputEmail
            name="email"
            label="Enter your email"
            placeholder="email@company.com"
            bind:value={email}
          />

          <InputPassword
            name="password"
            label="Enter your password"
            placeholder="********"
            bind:value={password}
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
          </div>

          <Button
            type="submit"
            disabledLabel={"Please accept the TOS."}
            disabled={!checked || loading}
            Icon={User}
          >
            {loading ? "Logging in..." : "Login"}
          </Button>
        </div>
      </form>

      <div class="questions">
        <p
          class="question"
          use:tippy={{
            content: `Student or supervisor at the ${BRAND.university.name}`,
            followCursor: "horizontal",
          }}
        >
          <a href="../login" class="inline">
            <GraduationCap size="1.5rem" />
            I am an {BRAND.university.name} Intern
          </a>
        </p>

        <p class="question">
          <span class="content">
            Don't have an account?
            <a href="./create">Create account</a>.
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

  .inline {
    display: flex;
    gap: 0.5rem;
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
