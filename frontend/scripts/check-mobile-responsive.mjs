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
      ["mobile-ai-modal__reader", "mobile AI reader area"],
      ["mobile-ai-modal__footer", "mobile AI footer metadata"],
      ["mobile-ai-modal__control-panel", "mobile AI control panel"],
      ["mobile-ai-modal__bottom-actions", "mobile AI bottom actions"],
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
      ["lw-kline-mobile-indicators", "mobile K-line indicator controls"],
      ["lw-kline-mobile-indicator-trigger", "mobile K-line indicator trigger"],
      ["lw-kline-mobile-indicator-panel", "mobile K-line indicator panel"],
      ["lw-kline-mobile-signal-strip", "mobile K-line signal summary"],
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
  {
    file: "frontend/src/components/researchReport.vue",
    expectations: [
      ["research-report-page", "research report page shell"],
      ["research-report-search", "research report mobile search"],
      ["research-report-table", "research report mobile table"],
      ["research-ai-modal", "research AI modal"],
      ["research-ai-reader", "research AI reader area"],
      ["research-ai-footer", "research AI footer metadata"],
      ["research-ai-actions", "research AI actions"],
    ],
  },
  {
    file: "frontend/src/components/fund.vue",
    expectations: [
      ["fund-page-shell", "fund page shell"],
      ["fund-mobile-tabs", "fund mobile tabs"],
    ],
  },
  {
    file: "frontend/src/components/FundFollow.vue",
    expectations: [
      ["fund-follow-page", "fund follow page shell"],
      ["fund-follow-search", "fund follow mobile search"],
      ["fund-follow-grid", "fund follow mobile grid"],
      ["fund-follow-card", "fund follow mobile card"],
      ["fund-follow-card-body", "fund follow card body"],
      ["fund-follow-actions", "fund follow actions"],
      ["fund-follow-chart-modal", "fund follow chart modal"],
      ["fund-follow-net-table", "fund follow net value table"],
      ["fund-mobile-add-bar", "fund mobile add bar"],
    ],
  },
  {
    file: "frontend/src/components/FundRanking.vue",
    expectations: [
      ["fund-ranking-page", "fund ranking page shell"],
      ["fund-ranking-toolbar", "fund ranking mobile toolbar"],
      ["fund-ranking-table", "fund ranking mobile table"],
      ["fund-ranking-holdings-modal", "fund ranking holdings modal"],
      ["fund-ranking-kline-modal", "fund ranking K-line modal"],
    ],
  },
  {
    file: "frontend/src/components/FundKlineChart.vue",
    expectations: [
      ["fund-kline-chart", "fund K-line chart shell"],
      ["fund-kline-controls", "fund K-line mobile controls"],
      ["fund-kline-source", "fund K-line source label"],
    ],
  },
  {
    file: "frontend/src/components/agent-chat.vue",
    expectations: [
      ["agent-chat-page", "agent chat page shell"],
      ["agent-chat-thread", "agent chat message thread"],
      ["agent-chat-config-bar", "agent chat config controls"],
      ["agent-chat-sender", "agent chat sender"],
      ["agent-chat-send-button", "agent chat send button"],
      ["agent-chat-bottom-button", "agent chat bottom button"],
      ["agent-chat-report", "agent chat report block"],
      ["agent-chat-steps", "agent chat steps block"],
    ],
  },
  {
    file: "frontend/src/components/settings.vue",
    expectations: [
      ["settings-page-shell", "settings page shell"],
      ["settings-section-card", "settings section cards"],
      ["settings-form-grid", "settings mobile form grid"],
      ["settings-ai-config-grid", "settings mobile AI config grid"],
      ["settings-mobile-action-bar", "settings mobile action bar"],
      ["settings-prompt-modal", "settings prompt modal"],
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
