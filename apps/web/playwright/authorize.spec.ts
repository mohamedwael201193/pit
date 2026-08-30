// MOCK TEST HARNESS — public UI copy. Never type AUTHORIZE against a live venue.

import { expect, test } from "./fixture";

test("web cannot present an authorize control", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("button", { name: "Authorize" })).toHaveCount(0);
  await expect(page.getByRole("button", { name: /enable guarded/i })).toHaveCount(0);
  await page.goto("/radar");
  await expect(page.getByRole("button", { name: "Authorize" })).toHaveCount(0);
  await page.goto("/pair");
  await expect(page.getByRole("button", { name: "Authorize" })).toHaveCount(0);
});
