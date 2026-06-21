import { readFileSync } from "node:fs";
import { join } from "node:path";

const root = process.cwd();

function read(path) {
  return readFileSync(join(root, path), "utf8");
}

function assertIncludes(file, content, needle, description) {
  if (!content.includes(needle)) {
    throw new Error(`${file}: missing ${description} (${needle})`);
  }
}

const checks = [
  {
    file: "frontend/src/App.vue",
    expectations: [
      ["isMobile", "mobile viewport state"],
      ["mobileMenuVisible", "mobile drawer state"],
      ["mobile-bottom-nav", "bottom navigation"],
      ["mobile-menu-drawer", "mobile menu drawer"],
    ],
  },
  {
    file: "frontend/src/style.css",
    expectations: [
      ["--mobile-bottom-nav-height", "mobile bottom nav CSS variable"],
      ["@media (max-width: 768px)", "mobile breakpoint"],
      [".mobile-only", "mobile-only utility"],
      [".desktop-only", "desktop-only utility"],
    ],
  },
  {
    file: "frontend/src/components/stock.vue",
    expectations: [
      ["stock-page-shell", "stock page shell class"],
      ["mobile-ai-modal", "mobile AI modal class"],
      ["stock-mobile-actions", "mobile stock action layout"],
    ],
  },
  {
    file: "frontend/src/components/kline-analysis.vue",
    expectations: [
      ["kline-analysis-page", "K-line page shell class"],
      ["mobile-kline-search", "mobile K-line search layout"],
    ],
  },
  {
    file: "frontend/src/components/market.vue",
    expectations: [
      ["market-page-shell", "market page shell class"],
      ["market-mobile-scroll", "market mobile scroll class"],
    ],
  },
  {
    file: "frontend/src/components/promptPlaza.vue",
    expectations: [
      ["prompt-plaza-page", "prompt plaza page shell class"],
      ["prompt-plaza-toolbar", "prompt plaza responsive toolbar"],
    ],
  },
];

for (const check of checks) {
  const content = read(check.file);
  for (const [needle, description] of check.expectations) {
    assertIncludes(check.file, content, needle, description);
  }
}

console.log("mobile responsive static checks passed");
