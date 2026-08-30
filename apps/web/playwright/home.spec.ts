// MOCK TEST HARNESS — public UI copy. Never stub VerifyE2EE, AUTHORIZE, or a live order.

import { expect, test } from "./fixture";

test("home is intelligence, not a wallet-first landing", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: /It hunts while you sleep/ })).toBeVisible();
  await expect(page.getByText("Your keys never leave your machine.")).toBeVisible();
  await expect(page.getByRole("link", { name: "Explore PIT" }).first()).toBeVisible();
  await expect(page.getByRole("link", { name: "Download PIT Desktop" }).first()).toBeVisible();
  await expect(page.getByRole("link", { name: "See Sleep Missions" }).first()).toBeVisible();
  await expect(page.getByText("PIT never asks for a seed phrase.")).toHaveCount(0);
  await expect(page.locator('input[type="password"]')).toHaveCount(0);
  await expect(page.getByRole("button", { name: "Authorize" })).toHaveCount(0);
  await expect(page.getByRole("button", { name: "Connect your wallet" })).toHaveCount(0);
  await expect(page.getByText("MAINNET only")).toBeVisible();
  await expect(page.getByText("The laboratory exists for CI and developers, not for the public desk.")).toBeVisible();
  await expect(page.getByText("This computer must stay awake for the bound.")).toBeVisible();
  await expect(page.getByText("Let PIT watch. Keep execution on your machine.")).toBeVisible();
});
