<script>
  import { BRAND } from "$lib/constants/branding";
  import { NAV_LINKS } from "$lib/constants/navigation";

  import Logo from "./Logo.svelte";
  import View from "./ui/View.svelte";
  import ThemeToggle from "./ui/ThemeToggle.svelte";

  import { Hamburger } from "svelte-hamburgers";

  let menuOpen = $state(false);
</script>

<menu class="mobile" class:open={menuOpen}>
  <View center style="gap:40px">
    <ul class="links">
      {#each NAV_LINKS as link}
        <a href={link.url}>
          <div class="icon">
            <link.icon size="32" />
          </div>
          {link.name}
        </a>
      {/each}
    </ul>

    <div class="theme">
      <ThemeToggle />
    </div>
  </View>
</menu>

<nav>
  <View align="row">
    <main>
      <a class="brand" href="/">
        <div class="logo">
          <Logo height="65%" textColor="transparent" logoColor="white" />
        </div>

        <h1>
          {BRAND.name}
        </h1>
      </a>

      <div class="actions">
        <ul class="links">
          {#each NAV_LINKS as link}
            <a href={link.url}>
              <div class="icon">
                <link.icon size="16" />
              </div>
              {link.name}
            </a>
          {/each}
        </ul>

        <div class="theme">
          <ThemeToggle />
        </div>

        <!-- display user data. -->
        <!-- {#if } -->
      </div>

      <div class="hamburger">
        <Hamburger bind:open={menuOpen} />
      </div>
    </main>
  </View>
</nav>

<style>
  menu {
    z-index: 3;
    position: absolute;

    top: var(--nav-height);
    left: 0;
    width: 100%;
    height: calc(100% - var(--nav-height));

    background-color: var(--color-background);

    transition: all var(--transition-normal);

    transform: translateX(-100%);
    &.open {
      transform: translateX(0);
    }

    .links {
      display: flex;
      flex-direction: column;
      gap: 1rem;

      a {
        display: flex;
        gap: 0.7rem;
        font-size: 2rem;

        font-weight: 600;
        letter-spacing: 0.5px;
        text-transform: uppercase;
        font-family: var(--font-sans);

        .icon {
          display: flex;
          justify-self: center;
          align-items: center;
        }
      }
    }

    .theme {
      transform: scale(1.2);
    }
  }

  nav {
    height: var(--nav-height);
    background-color: var(--color-accent-fixed);
    animation: logo 3s linear infinite;

    main {
      width: 100%;

      display: flex;
      align-items: center;
      justify-content: space-between;

      > * {
        height: 100%;
      }

      .brand {
        display: flex;
        gap: 1rem;

        font-family: var(--font-serif);

        align-items: center;
        text-decoration: none;

        &,
        .logo {
          transition: all var(--transition-fast);
        }

        .logo {
          height: 100%;

          display: flex;
          align-items: center;
          justify-content: center;
        }

        h1 {
          font-weight: 500;
        }

        &:hover .logo {
          transform: scale(1.1);
        }

        &:active {
          transform: scale(0.9);
        }

        &:active .logo {
          transform: none;
        }
      }
    }

    .links {
      display: flex;
      align-items: center;
      gap: 1rem;

      a {
        gap: 0.5rem;
        display: flex;
        align-items: center;

        font-weight: 600;
        font-size: 1rem;
        letter-spacing: 0.5px;
        text-transform: uppercase;
        font-family: var(--font-sans);

        .icon {
          transition: all var(--transition-normal);
        }

        &:hover .icon {
          transform: scale(1.2) rotate(10deg);
        }
      }
    }

    .actions {
      display: flex;
      align-items: center;
      gap: 0.5rem;

      @media screen and (max-width: 665px) {
        & {
          display: none;
        }
      }

      color: var(--color-gray-50);
    }

    .hamburger {
      display: flex;
      justify-content: center;
      align-items: center;

      display: none;

      @media screen and (max-width: 665px) {
        & {
          display: flex;
        }
      }

      :global(button) {
        padding: 0;

        :global(.hamburger-inner) {
          &,
          &::before,
          &::after {
            background-color: var(--color-gray-50) !important;
          }
        }
      }
    }

    .brand *,
    .links * {
      color: var(--color-gray-50);
    }
  }

  /* @property --color {
    syntax: "<color>";
    inherits: true;
    initial-value: red;
  } */
</style>
