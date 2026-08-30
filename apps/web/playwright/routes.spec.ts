// MOCK TEST HARNESS — public routes. Never stub a live order.

import { expect, test } from "./fixture";

const ROUTES = ["/", "/radar", "/capital", "/missions", "/proof", "/agent", "/how-it-works", "/download", "/pair"];

for (const path of ROUTES) {
  test(`${path} renders without authorize`, async ({ page }) => {
    const res = await page.goto(path);
    expect(res?.ok()).toBeTruthy();
    await expect(page.getByRole("button", { name: "Authorize" })).toHaveCount(0);
    await expect(page.locator('input[type="password"]')).toHaveCount(0);
  });
}

test("radar copy is public-safe", async ({ page }) => {
  await page.goto("/radar");
  await expect(page.getByRole("heading", { name: "What is happening right now?" })).toBeVisible();
  await expect(page.getByRole("tab", { name: /RESEARCH/ })).toBeVisible();
});

test("capital is labeled simulation", async ({ page }) => {
  await page.goto("/capital");
  await expect(page.getByText("SIMULATION", { exact: true })).toBeVisible();
  await expect(page.getByRole("heading", { name: "What could PIT do with this capital?" })).toBeVisible();
  await expect(page.getByText("No live PIT LP/swap route. Hyperliquid perps only.")).toBeVisible();
});

test("missions do not invent a live stream", async ({ page }) => {
  await page.goto("/missions");
  await expect(page.getByText("No live public mission")).toBeVisible();
  await expect(page.getByText("HISTORICAL", { exact: true })).toBeVisible();
});

test("proof does not badge verified without a check", async ({ page }) => {
  await page.goto("/proof");
  await expect(page.getByRole("heading", { name: "What was verified, and how" })).toBeVisible();
  await expect(page.getByText("NO LIVE RECEIPT")).toBeVisible();
});

test("agent shows iTransfer not live and desk id", async ({ page }) => {
  await page.goto("/agent");
  await expect(page.getByRole("heading", { name: "PIT-4bbee556" })).toBeVisible();
  await expect(page.getByText("NOT LIVE ON MAINNET")).toBeVisible();
  await expect(page.getByText("Desk ID · ERC-7857")).toBeVisible();
  await expect(page.getByRole("button", { name: /transfer/i })).toHaveCount(0);
});

test("download does not claim authenticode", async ({ page }) => {
  await page.goto("/download");
  await expect(page.getByText("unsigned")).toBeVisible();
  await expect(page.getByRole("link", { name: "GitHub release" })).toBeVisible();
});

test("unknown replay is honest", async ({ page }) => {
  await page.goto("/missions/8C91/replay");
  await expect(page.getByText("No public-safe mission with that id")).toBeVisible();
});
