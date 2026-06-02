<script lang="ts">
  interface Props {
    children?: import("svelte").Snippet;
    align?: "row" | "column" | "default";
    center?: boolean;
    fullScreen?: boolean;
    maxWidth?: string;
    editor?: boolean;
    main?: boolean;
    [key: string]: any;
  }

  let {
    children,
    align = "column",
    center = false,
    fullScreen = false,
    maxWidth,
    editor = false,
    main = false,
    ...rest
  }: Props = $props();
</script>

<div
  class="view"
  class:center
  class:full-screen={fullScreen}
  class:editor
  class:main
  style:max-width={maxWidth}
  data-align={align}
  style:height={fullScreen ? "unset" : "100%"}
  {...rest}
>
  {@render children?.()}
</div>

<style>
  /* Layout utilities */
  .view {
    width: 100%;
    max-width: var(--max-view);
    margin: 0 auto;
    padding: 0 24px;

    transition: height var(--transition-normal);
  }

  .view.editor {
    max-width: var(--max-view-large);
  }

  .view.main {
    padding: 2rem 0;
  }

  @media (max-width: 768px) {
    .view {
      padding: 0rem 1rem;
    }
  }

  div[data-align] {
    display: flex;
  }

  div[data-align="row"] {
    flex-direction: row;
  }

  div[data-align="column"] {
    flex-direction: column;
  }
  .center {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .full-screen {
    min-height: var(--screen-height);
  }

  .full-screen.editor {
    min-height: var(--screen-height-editor);
  }
</style>
