import { toast } from "@zerodevx/svelte-toast";

const show = (
  className: string,
  content: string,
  autohide: boolean,
  isError: boolean,
) => {
  toast.push(content, {
    initial: autohide ? 1 : 0,
    classes: ["toast-item", className],
    target: isError ? "errors" : "general",
  });
};

export const toastSuccess = (content: string, autohide: boolean) =>
  show("toast-success", content, autohide, false);

export const toastInfo = (content: string, autohide: boolean) =>
  show("toast-info", content, autohide, false);

export const toastWarning = (content: string, autohide: boolean) =>
  show("toast-warning", content, autohide, true);

export const toastDanger = (content: string, autohide: boolean) =>
  show("toast-danger", content, autohide, true);
