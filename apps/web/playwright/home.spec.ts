// MOCK TEST HARNESS — public UI copy. Never stub VerifyE2EE, AUTHORIZE, or a live order.

import { expect, test } from "./fixture";

test("home is intelligence, not a wallet-first landing", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: /0G seals the book/ })).toBeVisible();
  await expect(page.getByRole("heading", { name: "Watch PIT in action" })).toBeVisible();
  await expect(page.getByText("Hyperliquid fills the order.")).toBeVisible();
  await expect(page.getByText(/OID 531667200134/)).toBeVisible();
  await expect(page.getByRole("link", { name: "Watch PIT in action" }).first()).toHaveAttribute("href", /#watch$/);
  await expect(page.getByRole("link", { name: "Explore PIT" }).first()).toBeVisible();
  await expect(page.getByRole("link", { name: "Download PIT Desktop" }).first()).toBeVisible();
  await expect(page.getByRole("link", { name: "Download", exact: true }).first()).toHaveAttribute("href", /\/windows$/);
  await expect(page.getByText("PIT never asks for a seed phrase.")).toHaveCount(0);
  await expect(page.locator('input[type="password"]')).toHaveCount(0);
  await expect(page.getByRole("button", { name: "Authorize" })).toHaveCount(0);
  await expect(page.getByRole("button", { name: "Connect your wallet" })).toHaveCount(0);
  await expect(page.getByRole("heading", { name: "MAINNET", exact: true })).toBeVisible();
  await expect(page.getByText("The laboratory exists for CI and developers, not for the public desk.")).toBeVisible();
  await expect(page.getByText("Sleep Mission is bounded host automation on this computer.")).toBeVisible();
  await expect(page.getByText("The web discovers. The desktop acts.")).toBeVisible();
  await expect(page.getByRole("heading", { name: "Authorized." })).toBeVisible();
  await expect(page.getByText("Let PIT watch. Keep execution on your machine.")).toBeVisible();
});
