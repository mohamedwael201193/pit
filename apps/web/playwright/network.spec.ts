// MOCK TEST HARNESS — public UI copy. Never mix venues or stub a live order.

import { expect, test } from "@playwright/test";

test("landing is mainnet product", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByText("MAINNET only")).toBeVisible();
  await expect(page.getByRole("button", { name: "TESTNET" })).toHaveCount(0);
  await expect(page.getByText("The laboratory exists for CI and developers, not for the public desk.")).toBeVisible();
});
