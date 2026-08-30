import { expect, test } from "./fixture";

test("landing is mainnet product", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByText("MAINNET only")).toBeVisible();
  await expect(page.getByRole("button", { name: "TESTNET" })).toHaveCount(0);
  await expect(page.getByText("The laboratory exists for CI and developers, not for the public desk.")).toBeVisible();
});

test("how-it-works stays mainnet", async ({ page }) => {
  await page.goto("/how-it-works");
  await expect(page.getByText("MAINNET only")).toBeVisible();
});
