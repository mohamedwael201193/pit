// MOCK TEST HARNESS — public UI copy. Never mix venues or stub a live order.

import { expect, test } from "@playwright/test";

test("testnet lab copy is distinct from production", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("button", { name: "TESTNET" }).click();
  await expect(page.getByText("TESTNET is the full integration lab.")).toBeVisible();
  await expect(page.getByText("Different model catalog than production")).toBeVisible();
});
