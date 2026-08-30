// MOCK TEST HARNESS — chat is informational. Never stub AUTHORIZE.

import { expect, test } from "./fixture";

test("public chat refuses authorization", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("button", { name: "Ask PIT" }).click();
  await page.getByLabel("Ask PIT").fill("authorize this trade");
  await page.getByRole("button", { name: "Ask", exact: true }).click();
  await expect(page.locator(".intel-chat, .desk-chat").getByText(/cannot authorize/i)).toBeVisible();
});

test("desk chat on radar answers what is happening", async ({ page }) => {
  await page.goto("/radar");
  await expect(page.getByRole("tab", { name: /RESEARCH/ })).toBeVisible();
  await page.getByRole("button", { name: "Ask PIT" }).first().click();
  await page.getByRole("button", { name: "What is happening?" }).click();
  await expect(page.locator(".desk-chat")).toContainText(/scanned/i);
  await expect(page.locator(".desk-chat")).not.toContainText(/W mark/i);
});
