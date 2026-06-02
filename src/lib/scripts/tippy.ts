import _tippy, { followCursor, animateFill, sticky } from "tippy.js";
import type { Props, Plugin } from "tippy.js";

import "tippy.js/dist/tippy.css";
import "tippy.js/dist/backdrop.css";
import "tippy.js/animations/shift-away.css";

export type TooltipParams = Partial<Props> & {
    enabled?: boolean;
    plugins?: Plugin[];
    content?: Props["content"] | string | null;
};

const DEFAULT_PROPS: Partial<Props> = {
    animation: "shift-away",
    duration: [200, 150],
    arrow: true,
};

const DEFAULT_PLUGINS: Plugin[] = [followCursor, animateFill, sticky];

export function tippy(node: HTMLElement, params: TooltipParams = {}) {
    const { enabled = true, plugins = [], ...tippyProps } = params;


    const instance = _tippy(node, {
        ...DEFAULT_PROPS,
        ...tippyProps,
        plugins: [...DEFAULT_PLUGINS, ...plugins],
    });


    if (!enabled || !tippyProps.content) {
        instance.disable();
    }

    return {
        update(newParams: TooltipParams) {
            const { enabled: newEnabled = true, plugins: newPlugins = [], ...newProps } = newParams;

            instance.setProps({
                ...DEFAULT_PROPS,
                ...newProps,
                plugins: [...DEFAULT_PLUGINS, ...newPlugins],
            });

            if (newEnabled && newProps.content) {
                instance.enable();
                instance.show();
            } else {
                instance.hide();
                instance.disable();
            }
        },

        destroy() {
            instance.destroy();
        },
    };
}
