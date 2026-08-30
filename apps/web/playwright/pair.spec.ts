// MOCK TEST HARNESS — public UI copy. Never stub VerifyE2EE, AUTHORIZE, or a live order.

import { expect, test } from "./fixture";

test("pair page never asks for a seed or authorize", async ({ page }) => {
  await page.goto("/pair");
  await expect(page.getByText("PIT never asks for a seed phrase.")).toBeVisible();
  await expect(page.getByRole("link", { name: "Protect my strategy" }).first()).toHaveAttribute("href", "/signin");
  await expect(page.getByRole("button", { name: "Authorize" })).toHaveCount(0);
  await expect(page.locator('input[type="password"]')).toHaveCount(0);
  await expect(page.getByRole("link", { name: "Download PIT" })).toHaveAttribute(
    "href",
    "https://github.com/mohamedwael201193/pit/releases/latest",
  );
  await expect(page.getByRole("button", { name: "Open PIT Desktop" })).toBeVisible();
  await expect(page.getByLabel("pairing code")).toBeVisible();
});

test("protect my strategy opens wallet sign-in", async ({ page }) => {
  await page.goto("/pair");
  await page.getByRole("link", { name: "Protect my strategy" }).first().click();
  await expect(page).toHaveURL(/\/signin$/);
  await expect(page.getByRole("heading", { name: /Sign in with your wallet|Wallet connected/ })).toBeVisible();
  await expect(page.getByRole("button", { name: "Connect your wallet" }).or(page.getByRole("link", { name: "Overview" }))).toBeVisible();
});
