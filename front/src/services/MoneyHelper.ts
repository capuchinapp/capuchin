const precision = 2;

export const moneyFormat = (
  node: HTMLInputElement,
): { destroy: () => void } => {
  function onFocus() {
    node.value = moneyFormatted(node.value).replace(/ /g, "");
  }

  function onBlur() {
    node.value = moneyFormatted(node.value);
  }

  node.addEventListener("focus", onFocus);
  node.addEventListener("blur", onBlur);

  return {
    destroy() {
      node.removeEventListener("focus", onFocus);
      node.removeEventListener("blur", onBlur);
    },
  };
};

export const moneyInteger = (value: number | string): number => {
  return getCleanNumber(value) * 100;
};

export const moneyFormatted = (value: number | string): string => {
  const re = /(\d)(?=(\d\d\d)+([^\d]|$))/g;

  return String(getCleanNumber(value)).replace(re, "$1 ");
};

const getCleanNumber = (value: number | string): number => {
  if (!value) {
    value = 0;
  }

  if (typeof value !== "number") {
    value = value.replace(/ /g, "");
    value = value.replace(/,/g, ".");

    value = parseFloat(value);

    if (isNaN(value)) {
      value = 0;
    }
  }

  return Number(value.toFixed(precision + 1).slice(0, -1));
};
