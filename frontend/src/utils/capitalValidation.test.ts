import { validateSystems } from "./capitalValidation";

describe("capital form validation", () => {
  test("fails on empty systems", () => {
    expect(validateSystems("", "")).toBe("capital.validation-required");
  });

  test("fails on same systems", () => {
    expect(validateSystems("Jita", "Jita")).toBe("capital.validation-same");
  });

  test("passes on different systems", () => {
    expect(validateSystems("Jita", "Amarr")).toBeNull();
  });
});
