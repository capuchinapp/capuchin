declare module "tom-select" {
  class TomSelect {
    constructor(selector: string | HTMLElement, options?: Record<string, any>);
    setValue(value: string, silent?: boolean): void;
    getValue(): string;
    clear(silent?: boolean): void;
    clearOptions(): void;
    addOptions(options: any[], silent?: boolean): void;
    refreshOptions(silent?: boolean): void;
    destroy(): void;
  }
  export default TomSelect;
}

declare module "bootstrap/dist/js/bootstrap.esm" {
  export class Modal {
    constructor(element: HTMLElement, options?: Record<string, any>);
    show(): void;
    hide(): void;
    toggle(): void;
  }
}
