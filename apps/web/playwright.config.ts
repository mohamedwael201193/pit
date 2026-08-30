import { defineConfig, devices } from "@playwright/test";

const baseURL = process.env.PLAYWRIGHT_BASE_URL ?? "http://127.0.0.1:4173";

export default defineConfig({
  testDir: "./playwright",
  fullyParallel: true,
  workers: 4,
  forbidOnly: !!process.env.CI,
  retries: 0,
  reporter: "list",
  use: {
    baseURL,
    trace: "off",
    screenshot: "off",
    navigationTimeout: 15_000,
    actionTimeout: 10_000,
  },
  projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
  webServer: process.env.PLAYWRIGHT_BASE_URL
    ? undefined
    : {
        command: "npx vite --host 127.0.0.1 --port 4173",
        url: "http://127.0.0.1:4173",
        reuseExistingServer: false,
        timeout: 120_000,
        env: {
          ...process.env,
          VITE_PRIVY_APP_ID: process.env.VITE_PRIVY_APP_ID ?? "cmtafcijw02av0cl1ay81om7m",
          VITE_HEALTH_URL: process.env.VITE_HEALTH_URL ?? "https://pit-health.onrender.com",
        },
      },
});
