const http = require("http");
const fs = require("fs");
const path = require("path");

async function run() {
  const { server, url } = await startServer();
  const puppeteer = await import("puppeteer");
  const browser = await puppeteer.default.launch({ args: ["--no-sandbox"] });
  const page = await browser.newPage();
  await page.goto(url + "#Capital");
  await page.type("#start-system", "Start");
  await page.type("#end-system", "End");
  await page.click("#find-route");
  await page.waitForSelector(".leaflet-interactive");
  const count = await page.$$eval(".leaflet-interactive", (els) => els.length);
  if (count === 0) {
    throw new Error("Polyline not rendered");
  }
  await browser.close();
  server.close();
}

function startServer() {
  const distDir = path.resolve(__dirname, "../../dist");
  return new Promise((resolve) => {
    const server = http.createServer((req, res) => {
      const url = req.url ? req.url.split("?")[0] : "/";
      if (url.startsWith("/api/capital")) {
        res.writeHead(200, { "Content-Type": "application/json" });
        res.end(
          JSON.stringify({
            route: [
              { id: 1, name: "Start", x: 0, y: 0, z: 0 },
              { id: 2, name: "End", x: 1e16, y: 0, z: 0 },
            ],
          }),
        );
        return;
      }
      if (url.startsWith("/api/auth/user")) {
        res.writeHead(200, { "Content-Type": "application/json" });
        res.end(
          JSON.stringify({
            name: "Test",
            allianceName: "",
            allianceTicker: "",
            roles: [],
            csrfHeaderKey: "X-CSRF-Token",
            csrfToken: "1",
          }),
        );
        return;
      }
      if (url.startsWith("/api/route/map-connections")) {
        res.writeHead(200, { "Content-Type": "application/json" });
        res.end(JSON.stringify({ code: "", ansiblexes: [], temporary: [] }));
        return;
      }
      const filePath =
        url === "/"
          ? path.join(distDir, "index.html")
          : path.join(distDir, url);
      fs.readFile(filePath, (err, data) => {
        if (err) {
          res.writeHead(404);
          res.end("not found");
          return;
        }
        let type = "text/html";
        if (filePath.endsWith(".js")) type = "application/javascript";
        if (filePath.endsWith(".css")) type = "text/css";
        res.writeHead(200, { "Content-Type": type });
        res.end(data);
      });
    });
    server.listen(0, () => {
      const { port } = server.address();
      resolve({ server, url: `http://localhost:${port}` });
    });
  });
}

run().catch((err) => {
  console.error(err);
  process.exit(1);
});
