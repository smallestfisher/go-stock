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
      ["stock-mobile-add-drawer", "mobile stock add drawer"],
      ["stock-mobile-add-trigger", "mobile stock add trigger"],
      ["stock-mobile-primary-actions", "mobile primary stock actions"],
      ["stock-mobile-more-menu", "mobile secondary stock menu"],
      ["stock-mobile-card-body", "mobile stock card tap target"],
    ],
  },
  {
    file: "frontend/src/components/kline-analysis.vue",
    expectations: [
      ["kline-analysis-page", "K-line page shell class"],
      ["mobile-kline-search", "mobile K-line search layout"],
      ["kline-mobile-chart-stage", "mobile K-line chart stage"],
      ["kline-mobile-search-trigger", "mobile K-line search trigger"],
      ["kline-mobile-search-drawer", "mobile K-line search drawer"],
      ["kline-mobile-recent-list", "mobile K-line recent list"],
    ],
  },
  {
    file: "frontend/src/components/StockLightweightKlineChart.vue",
    expectations: [
      ["lw-kline-top-controls", "mobile compact K-line top controls"],
      ["lw-kline-mobile-compact", "mobile compact K-line rules"],
    ],
  },
  {
    file: "frontend/src/components/market.vue",
    expectations: [
      ["market-page-shell", "market page shell class"],
      ["market-mobile-scroll", "market mobile scroll class"],
      ["market-mobile-heat-panel", "mobile market heat panel"],
      ["market-mobile-news-tabs", "mobile market news tabs"],
      ["market-desktop-news-grid", "desktop market news grid"],
      ["market-mobile-index-card", "mobile market index card"],
      ["market-mobile-summary-action", "mobile market summary action"],
    ],
  },
  {
    file: "frontend/src/components/promptPlaza.vue",
    expectations: [
      ["prompt-plaza-page", "prompt plaza page shell class"],
      ["prompt-plaza-toolbar", "prompt plaza responsive toolbar"],
      ["prompt-plaza-mobile-filter-trigger", "prompt plaza mobile filter trigger"],
      ["prompt-plaza-mobile-filter-drawer", "prompt plaza mobile filter drawer"],
      ["prompt-plaza-mobile-card", "prompt plaza mobile card"],
      ["prompt-plaza-mobile-stats", "prompt plaza mobile stats"],
      ["prompt-plaza-mobile-tags", "prompt plaza mobile tags"],
      ["prompt-plaza-detail-modal", "prompt plaza mobile detail modal"],
      ["prompt-plaza-detail-actions", "prompt plaza mobile detail actions"],
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
