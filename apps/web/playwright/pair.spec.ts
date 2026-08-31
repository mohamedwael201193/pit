// MOCK TEST HARNESS — public UI copy. Never stub VerifyE2EE, AUTHORIZE, or a live order.

import { expect, test } from "./fixture";

test("pair page never asks for a seed or authorize", async ({ page }) => {
  await page.goto("/pair");
  await expect(page.getByText("PIT never asks for a seed phrase.")).toBeVisible();
  await expect(page.getByText("Pairing is step 1. Protect my strategy is step 2.")).toBeVisible();
  await expect(page.getByRole("link", { name: "Protect my strategy" }).first()).toHaveAttribute("href", "/protect");
  await expect(page.getByRole("button", { name: "Authorize" })).toHaveCount(0);
  await expect(page.locator('input[type="password"]')).toHaveCount(0);
  await expect(page.getByRole("link", { name: "Download PIT" })).toHaveAttribute(
    "href",
    /\/windows$/,
  );
  await expect(page.getByRole("button", { name: "Open PIT Desktop" })).toBeVisible();
  await expect(page.getByLabel("pairing code")).toBeVisible();
});

test("protect stays locked until this browser is paired", async ({ page }) => {
  await page.goto("/pair");
  await page.getByRole("link", { name: "Protect my strategy" }).first().click();
  await expect(page).toHaveURL(/\/protect$/);
  await expect(page.getByRole("heading", { name: "Protect my strategy" })).toBeVisible();
  await expect(page.getByText("Connect wallet to link with desktop")).toBeVisible();
  await expect(page.getByRole("button", { name: "Connect your wallet" })).toBeVisible();
  await expect(page.getByText("Protect my strategy stays locked until this browser is paired.")).toBeVisible();
  await expect(page.getByRole("button", { name: "Protect my strategy" })).toHaveCount(0);
  await expect(page.getByRole("link", { name: "Overview" })).toHaveCount(0);
});

test("explore desk shows connect wallet in the header", async ({ page }) => {
  await page.goto("/radar");
  await expect(page.getByRole("heading", { name: "What is happening right now?" })).toBeVisible();
  await expect(page.getByRole("button", { name: "Connect wallet" })).toBeVisible();
  await expect(page.getByRole("link", { name: "Sign in" })).toHaveCount(0);
});

test("protect unlocks wallet connect after a pairing token exists", async ({ page }) => {
  await page.addInitScript(() => sessionStorage.setItem("pit_device", "playwright-device"));
  await page.goto("/protect");
  await expect(page.getByRole("button", { name: "Connect your wallet" })).toBeVisible();
});

test("legacy sign-in and get-started land on pair", async ({ page }) => {
  await page.goto("/signin");
  await expect(page).toHaveURL(/\/pair$/);
  await page.goto("/app/start");
  await expect(page).toHaveURL(/\/pair$/);
});
