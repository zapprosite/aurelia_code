import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import { z } from "zod";
import { Stagehand } from "@browserbasehq/stagehand";

const server = new McpServer({
  name: "aurelia-stagehand",
  version: "1.0.0",
});

let stagehand: Stagehand | null = null;

async function getStagehand() {
  if (!stagehand) {
    stagehand = new Stagehand({
      env: "LOCAL",
      verbose: 1,
      model: {
        modelName: "aurelia-smart",
        apiKey: "sk-1234",
        baseURL: "http://localhost:4000/v1",
      }
    });
    await stagehand.init();
  }
  return stagehand;
}

server.tool(
  "navigate",
  { url: z.string().url() },
  async ({ url }) => {
    const sh = await getStagehand();
    const page = await sh.context.awaitActivePage();
    await page.goto(url);
    return {
      content: [{ type: "text", text: `Navegado para ${url}` }],
    };
  }
);

server.tool(
  "act",
  { instruction: z.string() },
  async ({ instruction }) => {
    const sh = await getStagehand();
    await sh.act(instruction);
    return {
      content: [{ type: "text", text: `Ação executada: ${instruction}` }],
    };
  }
);

server.tool(
  "extract",
  { instruction: z.string() },
  async ({ instruction }) => {
    const sh = await getStagehand();
    const data = await sh.extract(instruction);
    return {
      content: [{ type: "text", text: JSON.stringify(data, null, 2) }],
    };
  }
);

async function main() {
  const transport = new StdioServerTransport();
  await server.connect(transport);
  console.error("Aurelia Stagehand MCP Server running on stdio");
}

main().catch((error) => {
  console.error("Fatal error in main():", error);
  process.exit(1);
});
