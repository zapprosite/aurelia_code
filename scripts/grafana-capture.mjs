import { chromium } from "playwright";
import fs from "fs";
import path from "path";

const GRAFANA_URL = process.env.GRAFANA_URL || "https://monitor.zappro.site/login";
const DASHBOARD_URL = process.env.DASHBOARD_URL || "https://monitor.zappro.site/d/40977c27-62f1-4a7d-a455-2f4d698538d8/nvidia-gpu-metrics?orgId=1&from=now-5m&to=now&timezone=browser&var-job=nvidia-gpu&var-node=will-zappro&var-gpu=bc42e64f-64d5-4711-e976-6141787b60a2&refresh=2s&kiosk";
const GRAFANA_USER = process.env.GRAFANA_USER || "admin";
const GRAFANA_PASS = process.env.GRAFANA_PASS || "2LCwzksQxxF7PhFnbgB5dF1G";
const OUTPUT_PATH = process.env.OUTPUT_PATH || "/tmp/grafana_snapshot.png";

async function capture() {
  console.log(`🚀 Iniciando captura do Grafana em: ${GRAFANA_URL}`);
  const browser = await chromium.launch({ headless: true });
  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 }
  });
  const page = await context.newPage();

  try {
    // 1. Login
    await page.goto(GRAFANA_URL, { waitUntil: "networkidle" });
    
    // Verificar se já estamos logados ou se precisamos preencher o form
    if (await page.isVisible('input[name="user"]')) {
      await page.fill('input[name="user"]', GRAFANA_USER);
      await page.fill('input[name="password"]', GRAFANA_PASS);
      await page.click('button[type="submit"]');
      await page.waitForNavigation({ waitUntil: "networkidle" });
    }

    // 2. Ir para o Dashboard
    console.log(`📊 Carregando dashboard: ${DASHBOARD_URL}`);
    await page.goto(DASHBOARD_URL, { waitUntil: "networkidle" });
    
    // Aguardar o carregamento dos painéis (ajustar conforme necessário)
    await page.waitForTimeout(5000); 

    // 3. Screenshot
    await page.screenshot({ path: OUTPUT_PATH, fullPage: true });
    console.log(`✅ Screenshot salva em: ${OUTPUT_PATH}`);

  } catch (error) {
    console.error("❌ Erro durante a captura:", error);
    process.exit(1);
  } finally {
    await browser.close();
  }
}

capture();
