// MOCK TEST HARNESS — chat is informational. Never stub AUTHORIZE.

import { expect, test } from "./fixture";

test("public chat refuses authorization", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("button", { name: "Ask PIT" }).click();
  await page.getByLabel("Ask PIT").fill("authorize this trade");
  await page.getByRole("button", { name: "Ask", exact: true }).click();
  await expect(page.getByText(/cannot authorize/i)).toBeVisible();
});
