import { chromium } from "playwright";

const targetUrl = process.argv[2] || "https://example.com";
const screenshotPath = process.argv[3] || "/tmp/jarvis-playwright-smoke.png";

async function main() {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();
  await page.goto(targetUrl, { waitUntil: "domcontentloaded", timeout: 30000 });
  await page.waitForLoadState("networkidle", { timeout: 30000 });

  const result = {
    url: page.url(),
    title: await page.title(),
    screenshot: screenshotPath,
  };

  await page.screenshot({ path: screenshotPath, fullPage: true });
  console.log(JSON.stringify(result, null, 2));
  await browser.close();
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
