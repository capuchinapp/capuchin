import "./../dayjs";

import { describe, it, expect, vi, beforeEach } from "vitest";
import MockDate from "mockdate";
import {
  nowDate,
  nowTime,
  getDateTime,
  getDate,
  getTime,
  getDiffSeconds,
  hoursFromSeconds,
  calculateTime,
  formatDurationSmart,
} from "./DateTime";

describe("DateTime Service", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    MockDate.reset();
  });

  describe("nowDate", () => {
    it("должен вернуть текущую дату", () => {
      MockDate.set("2023-12-15");
      expect(nowDate()).toBe("15.12.2023");
    });
  });

  describe("nowTime", () => {
    it("должен вернуть текущее время без секунд, когда withSeconds имеет значение false", () => {
      MockDate.set("2023-12-15 14:30:45");
      expect(nowTime(false)).toBe("14:30");
    });

    it("должен вернуть текущее время с секундами, когда withSeconds имеет значение true", () => {
      MockDate.set("2023-12-15 14:30:45");
      expect(nowTime(true)).toBe("14:30:45");
    });

    it("должен вернуть текущее время без секунд, когда withSeconds не указан", () => {
      MockDate.set("2023-12-15 14:30:45");
      expect(nowTime(undefined)).toBe("14:30");
    });
  });

  describe("getDateTime", () => {
    it("должен вернуть дату и время без секунд", () => {
      MockDate.set("2023-12-15 14:30:45");
      expect(getDateTime("2023-12-15T14:30:45")).toBe("15.12.2023 14:30");
    });

    it("должен вернуть тире для неверной даты", () => {
      MockDate.set("2023-12-15 14:30:45");
      expect(getDateTime("invalid-date")).toBe("-");
    });
  });

  describe("getDate", () => {
    it("должен вернуть дату в формате из настроек", () => {
      MockDate.set("2023-12-15");
      expect(getDate("2023-12-15T14:30:00")).toBe("15.12.2023");
    });

    it("должен вернуть тире для неверной даты", () => {
      MockDate.set("2023-12-15");
      expect(getDate("invalid-date")).toBe("-");
    });
  });

  describe("getTime", () => {
    it("должен вернуть время без секунд", () => {
      MockDate.set("2023-12-15 14:30:45");
      expect(getTime("2023-12-15T14:30:00")).toBe("14:30");
    });

    it("должен вернуть тире для неверного времени", () => {
      MockDate.set("2023-12-15 14:30:45");
      expect(getTime("invalid-date")).toBe("-");
    });
  });

  describe("getDiffSeconds", () => {
    it("должен рассчитать разницу в секундах между двумя временами", () => {
      expect(getDiffSeconds("10:00", "11:00")).toBe(3600);
    });
  });

  describe("hoursFromSeconds", () => {
    it("должен преобразовать секунды в часы с точностью", () => {
      expect(hoursFromSeconds(5400, 2)).toBe(1.5); // 1.5 hours
      expect(hoursFromSeconds(3600, 0)).toBe(1); // 1 hour
    });

    it("должен вернуть необработанное значение, когда точность отрицательна", () => {
      expect(hoursFromSeconds(5400, -1)).toBe(1.5);
    });
  });

  describe("calculateTime", () => {
    it("должен форматировать продолжительность как время без секунд", () => {
      expect(calculateTime(5400, false)).toBe("01:30");
    });

    it("должен форматировать продолжительность как время с секундами", () => {
      expect(calculateTime(5415, true)).toBe("01:30:15");
    });
  });

  describe("formatDurationSmart", () => {
    it("должен правильно форматировать секунды", () => {
      expect(formatDurationSmart(30)).toBe("30с");
      expect(formatDurationSmart(90)).toBe("1м");
      expect(formatDurationSmart(3599)).toBe("59м");
      expect(formatDurationSmart(3600)).toBe("1ч");
      expect(formatDurationSmart(5400)).toBe("1ч 30м");
      expect(formatDurationSmart(7200)).toBe("2ч");
      expect(formatDurationSmart(90000)).toBe("1д 1ч");
      expect(formatDurationSmart(86400)).toBe("1д");
    });
  });
});
