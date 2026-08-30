import { test as base } from "@playwright/test";

export const test = base.extend({
  page: async ({ page }, use) => {
    const original = page.goto.bind(page);
    page.goto = (url, options) => original(url, { waitUntil: "domcontentloaded", ...options });
    await use(page);
  },
});

export { expect } from "@playwright/test";
