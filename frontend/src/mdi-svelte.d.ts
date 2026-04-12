declare module 'mdi-svelte' {
  import { SvelteComponentTyped } from 'svelte';
  export interface IconProps {
    path: string;
    size?: number | string;
    color?: string;
    flip?: boolean | 'h' | 'v';
    rotate?: number | string;
    spin?: boolean | number;
  }
  export default class Icon extends SvelteComponentTyped<IconProps> {}
}