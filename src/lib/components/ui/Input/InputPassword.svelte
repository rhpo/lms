<script lang="ts">
  import type { InputProps } from "./Input.svelte";

  import { validation } from "$lib/utils/validation";

  import Input from "./Input.svelte";

  let { value = $bindable(""), ...rest }: InputProps = $props();
</script>

<Input
  type="password"
  bind:value
  validation={validation.password}
  {...rest as any}
/>
