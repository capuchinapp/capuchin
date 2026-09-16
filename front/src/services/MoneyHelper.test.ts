import { describe, it, expect, vi, beforeEach } from "vitest";
import { moneyInteger, moneyFormatted } from "./MoneyHelper";

describe("MoneyHelper Service", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("moneyInteger", () => {
    it("должен вернуть целое значение когда на входе строка", () => {
      expect(moneyInteger("100.55")).toBe(10055);
      expect(moneyInteger("100")).toBe(10000);
    });

    it("должен вернуть целое значение когда на входе число", () => {
      expect(moneyInteger(100.55)).toBe(10055);
      expect(moneyInteger(100)).toBe(10000);
    });
  });

  describe("moneyFormatted", () => {
    it("должен вернуть форматированное значение когда на входе строка", () => {
      expect(moneyFormatted("10 000,55")).toBe("10 000.55");
      expect(moneyFormatted("10000.55")).toBe("10 000.55");
      expect(moneyFormatted("10000")).toBe("10 000");
    });

    it("должен вернуть форматированное значение когда на входе число", () => {
      expect(moneyFormatted(10000.55)).toBe("10 000.55");
      expect(moneyFormatted(10000)).toBe("10 000");
    });
  });
});
