import fs from "node:fs";
import path from "node:path";

const workspaceRoot = path.resolve(import.meta.dirname, "..");

function mustRead(relativePath) {
  return fs.readFileSync(path.join(workspaceRoot, relativePath), "utf8");
}

function assertIncludes(content, expected, label) {
  if (!content.includes(expected)) {
    throw new Error(`Expected ${label} to include ${JSON.stringify(expected)}`);
  }
}

const reactIndex = mustRead("packages/react/src/index.ts");
const vueIndex = mustRead("packages/vue/src/index.ts");
const vanillaIndex = mustRead("packages/vanilla/src/index.ts");

assertIncludes(reactIndex, "export * from '@servify/core';", "sdk/packages/react/src/index.ts");
assertIncludes(reactIndex, "export { createWebServifySDK } from '@servify/core';", "sdk/packages/react/src/index.ts");

assertIncludes(vueIndex, "export * from '@servify/core';", "sdk/packages/vue/src/index.ts");
assertIncludes(vueIndex, "export { createWebServifySDK } from '@servify/core';", "sdk/packages/vue/src/index.ts");

assertIncludes(vanillaIndex, "createWebServifySDK", "sdk/packages/vanilla/src/index.ts");
assertIncludes(vanillaIndex, "export type { WebServifyConfig as ServifyConfig", "sdk/packages/vanilla/src/index.ts");
assertIncludes(vanillaIndex, "export default VanillaServifySDK;", "sdk/packages/vanilla/src/index.ts");

console.log("SDK surface smoke tests passed.");
