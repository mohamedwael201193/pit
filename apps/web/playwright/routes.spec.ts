// MOCK TEST HARNESS — public routes. Never stub a live order.

import { expect, test } from "./fixture";

const ROUTES = ["/", "/radar", "/capital", "/autonomy", "/missions", "/proof", "/agent", "/how-it-works", "/download", "/pair", "/protect"];

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
  await expect(page.getByRole("link", { name: "Pair" }).first()).toBeVisible();
  await expect(page.getByRole("button", { name: "Ask PIT" }).first()).toBeVisible();
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
  await expect(page.getByText("531667200134")).toBeVisible();
  await expect(page.getByText("529167222216")).toBeVisible();
  await expect(page.getByRole("link", { name: /Research 0x1d2113bd/ })).toBeVisible();
  await expect(page.getByRole("link", { name: /Order 0x8c28051b/ })).toBeVisible();
  await expect(page.getByText("RECORDED ROOTS")).toBeVisible();
  await expect(page.getByText("0x9fd42770545ecaacbfff12e3ef7a537b564e31c9ef5515b3a820fd276c22f72e").first()).toBeVisible();
  await expect(page.getByText("0x8c94ec8e643c90fe69276ff20f50a0bc3121f007d611e10e6ab9f24d26f2ff66").first()).toBeVisible();
  await expect(page.getByRole("link", { name: /DOGE job b4ed73ce/ })).toBeVisible();
});

test("agent shows iTransfer not live and desk id", async ({ page }) => {
  await page.goto("/agent");
  await expect(page.getByRole("heading", { name: "PIT-4bbee556" })).toBeVisible();
  await expect(page.getByText("NOT LIVE ON MAINNET")).toBeVisible();
  await expect(page.getByText("Desk ID · ERC-7857")).toBeVisible();
  await expect(page.getByRole("button", { name: /transfer/i })).toHaveCount(0);
});

test("download does not claim authenticode and files the installer", async ({ page }) => {
  await page.goto("/download");
  await expect(page.getByText("unsigned")).toBeVisible();
  const installer = page.getByRole("link", { name: "Download Windows installer" });
  await expect(installer).toBeVisible();
  await expect(installer).toHaveAttribute("href", /\/windows$/);
  await expect(installer).not.toHaveAttribute("href", /releases\/latest/);
  await expect(page.getByRole("link", { name: "SHA256SUMS" })).toHaveAttribute("href", /\/checksums$/);
});

test("unknown replay is honest", async ({ page }) => {
  await page.goto("/missions/8C91/replay");
  await expect(page.getByText("No public-safe mission with that id")).toBeVisible();
});
