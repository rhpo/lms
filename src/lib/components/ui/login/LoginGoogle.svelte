<script lang="ts">
</script>

<a class="button" href="/dev/login">
  <div class="icon">
    <img width="22px" alt="Google logo" src="/media/icons/google.png" />
  </div>

  <div class="content">
    Continue with <span class="google">Google</span>
  </div>
</a>

<style>
  :root {
    --padding: 1rem;
  }

  .button {
    border: none;
    background-color: var(--color-background-200);
    color: var(--color-text);
    display: flex;
    align-items: center;
    gap: calc(var(--padding) / 1.5);
    width: fit-content;
    border-radius: 999rem;
    padding: calc(var(--padding) / 1.7) var(--padding);
    font-size: 1.2rem;
    font-family: var(--font-sans);
    cursor: pointer;
    transition: all var(--transition-normal);

    &:hover {
      text-decoration: none;
      transform: scale(1.05) translateY(-2px);
    }

    &:hover .google {
      text-decoration: underline;
    }

    &:active {
      transform: none;
    }

    .google {
      font-weight: 600;
    }
  }
</style>
