<script lang="ts">
  import { validation } from "$lib/utils/validation";
  import type { InputProps } from "./Input.svelte";
  import Input from "./Input.svelte";

  let { value = $bindable(""), ...rest }: InputProps = $props();
</script>

<Input type="email" bind:value validation={validation.email} {...rest as any} />
