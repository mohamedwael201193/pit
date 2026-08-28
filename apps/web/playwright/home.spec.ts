// MOCK TEST HARNESS — public UI copy. Never stub VerifyE2EE, AUTHORIZE, or a live order.

import { expect, test } from "@playwright/test";

test("home asks to connect and never collects a seed", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("button", { name: "Connect your wallet" })).toBeVisible();
  await expect(page.getByText("PIT never asks for a seed phrase.")).toBeVisible();
  await expect(page.getByText("Your browser watches. PIT Desktop acts.")).toBeVisible();
  await expect(page.locator('input[type="password"]')).toHaveCount(0);
  await expect(page.getByRole("button", { name: "Authorize" })).toHaveCount(0);
  await expect(page.getByRole("link", { name: "Download PIT" }).first()).toHaveAttribute(
    "href",
    "https://github.com/mohamedwael201193/pit/releases/latest",
  );
});
