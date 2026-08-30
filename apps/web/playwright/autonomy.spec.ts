// MOCK TEST HARNESS — public copy. Never stub ARM, AUTHORIZE, or a live order.

import { expect, test } from "./fixture";

test("autonomy page never arms", async ({ page }) => {
  await page.goto("/autonomy");
  await expect(page.getByRole("heading", { name: /Sleep Missions/ })).toBeVisible();
  await expect(page.getByRole("button", { name: /arm/i })).toHaveCount(0);
  await expect(page.getByRole("button", { name: "Authorize" })).toHaveCount(0);
  await expect(page.getByRole("link", { name: "Open PIT Desktop" })).toBeVisible();
  await expect(page.getByText("Cannot arm, authorize, pin, or execute")).toBeVisible();
});

test("mission detail stays redacted", async ({ page }) => {
  await page.goto("/missions/does-not-exist");
  await expect(page.getByText("Private strategy remains on desktop.")).toBeVisible();
  await expect(page.getByRole("button", { name: /arm/i })).toHaveCount(0);
});
