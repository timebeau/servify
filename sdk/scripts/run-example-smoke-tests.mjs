import fs from "node:fs";
import path from "node:path";

const workspaceRoot = path.resolve(import.meta.dirname, "..");

function mustRead(filePath) {
  return fs.readFileSync(filePath, "utf8");
}

function assertIncludes(content, expected, label) {
  if (!content.includes(expected)) {
    throw new Error(`Expected ${label} to include ${JSON.stringify(expected)}`);
  }
}

const reactPackage = mustRead(path.join(workspaceRoot, "examples", "react", "package.json"));
const reactMain = mustRead(path.join(workspaceRoot, "examples", "react", "src", "main.tsx"));
const vanillaHtml = mustRead(path.join(workspaceRoot, "examples", "vanilla", "index.html"));

assertIncludes(reactPackage, "\"@servify/react\"", "sdk/examples/react/package.json");
assertIncludes(reactPackage, "\"build\": \"vite build\"", "sdk/examples/react/package.json");
assertIncludes(reactMain, "ReactDOM.createRoot", "sdk/examples/react/src/main.tsx");
assertIncludes(reactMain, "App", "sdk/examples/react/src/main.tsx");

assertIncludes(vanillaHtml, "../../packages/vanilla/dist/index.js", "sdk/examples/vanilla/index.html");
assertIncludes(vanillaHtml, "new Servify(", "sdk/examples/vanilla/index.html");
assertIncludes(vanillaHtml, "startChat()", "sdk/examples/vanilla/index.html");

console.log("SDK example smoke tests passed.");
