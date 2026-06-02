<!-- Input.svelte -->
<script module lang="ts">
  import type {
    HTMLInputAttributes,
    HTMLTextareaAttributes,
  } from "svelte/elements";

  type SharedProps = {
    label?: string;
    value?: string;
    error?: string;
    required?: boolean;
    validation?: ValidationFunc;
    onInput?: (e: Event) => void;
  };

  type InputType = NonNullable<HTMLInputAttributes["type"]>;

  const validationMatch: Partial<Record<InputType, ValidationFunc>> = {
    email: validationCollection.email,
    text: (value: string) => "",
  };

  export type InputProps =
    | ({
        category?: "input";
        type?: HTMLInputAttributes["type"];
        validation?: ValidationFunc;
      } & SharedProps &
        HTMLInputAttributes)
    | ({
        category: "textarea";
        type?: HTMLInputAttributes["type"];
        validation?: ValidationFunc;
      } & SharedProps &
        HTMLTextareaAttributes);
</script>

<script lang="ts">
  import { tippy, type TooltipParams } from "$lib/scripts/tippy";
  import {
    validate,
    validateValue,
    validation as validationCollection,
    type ValidationFunc,
  } from "$lib/utils/validation";

  let {
    category = "input",
    label = "",
    value = $bindable(""),
    type = "text",
    error = $bindable(""),
    placeholder,
    required = false,
    validation = validationMatch[type || "text"],
    onInput = () => {},
    ...rest
  }: InputProps = $props();

  let inputElement: HTMLInputElement | HTMLTextAreaElement = $state() as any;

  export function focus() {
    inputElement?.focus();
  }

  function handleInput(e: Event) {
    error = validateValue(value, validation);
    onInput(e);
  }

  function handleBlur() {
    error = validateValue(value, validation);
  }

  let tippyProps = $state({
    content: error || undefined,
    enabled: !!error,
    followCursor: "horizontal",
    maxWidth: 200,
  }) as TooltipParams;
</script>

<main class:margin-bottom={!!label}>
  {#if label}
    <label for={label.toLowerCase().replace(/\s+/g, "-")}>
      {label}
      {#if required}<span class="required">*</span>{/if}
    </label>
  {/if}

  {#if category === "textarea"}
    <textarea
      bind:this={inputElement}
      class="input"
      id={label.toLowerCase().replace(/\s+/g, "-")}
      bind:value
      use:tippy={tippyProps}
      {placeholder}
      class:error={!!error}
      oninput={handleInput}
      onblur={handleBlur}
      {required}
      {...rest as HTMLTextareaAttributes}
    ></textarea>
  {:else}
    <input
      bind:this={inputElement}
      use:tippy={tippyProps}
      {placeholder}
      class="input"
      id={label.toLowerCase().replace(/\s+/g, "-")}
      bind:value
      class:error={!!error}
      oninput={handleInput}
      onblur={handleBlur}
      {required}
      {...rest as HTMLInputAttributes}
    />
  {/if}
</main>

<style>
  main {
    overflow: hidden;

    &.margin-bottom {
      margin-bottom: 0.5rem;
    }
  }

  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    font-family: var(--font-sans);
    color: var(--color-text);
    font-size: 13px;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .required {
    color: var(--color-danger);
    margin-left: 0.25rem;
  }

  .input {
    width: 100%;
    padding: 0.8rem 1rem;
    border: none;
    font-size: 1rem;
    transition: all var(--transition-normal);
    background: transparent;
    font-family: var(--font-sans);
    color: var(--color-text);
    border-radius: 1rem;
    background-color: var(--color-background-100);
    border: 1px solid var(--color-background-200);

    &:focus {
      outline: none;
      border-bottom-color: var(--color-text);
    }

    &::placeholder {
      color: var(--color-text-muted);
    }

    &.error {
      border: 1px solid var(--color-danger);
    }
  }

  :global(.error-message) {
    font-family: var(--font-sans);
  }
</style>
