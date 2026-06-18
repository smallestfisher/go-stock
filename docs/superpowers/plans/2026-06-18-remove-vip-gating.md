# Remove VIP Gating Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make all go-stock local features free by default and remove user-visible VIP, sponsor-code, and entitlement-sync flows.

**Architecture:** Keep existing RPC names for compatibility, but make entitlement RPCs return active free access and make sponsor-code validation a no-op. Remove frontend denial branches and paid UI so the free behavior is visible, not only implied by backend responses.

**Tech Stack:** Go 1.26, net/http RPC handlers, Vue 3 Composition API, Vite, Naive UI, markdown documentation.

---

## File Structure

- Modify `backend/data/sponsor_vip.go`: change effective entitlement source to free active level.
- Create `backend/data/sponsor_vip_test.go`: prove entitlement no longer depends on sponsor code.
- Modify `server/modules/system_rest.go`: make sponsor-code validation and device binding compatibility no-ops.
- Create `server/modules/system_rest_test.go`: prove sponsor RPCs return free access and accept malformed/empty sponsor codes.
- Modify `frontend/src/components/settings.vue`: remove sponsor-code import, form state, validation on save, import hydration, and form item.
- Modify `frontend/src/components/about.vue`: remove VIP badge/expiry display and paid plan table; keep voluntary donation/support copy.
- Modify local feature components: remove VIP denial branches from assistants, K-line flows, stock filtering, AI recommendation detail.
- Modify `frontend/src/components/promptPlaza.vue`: remove local VIP/device-binding imports and behavior; make VIP-marked prompt import local-free.
- Modify docs: `README.md`, `docs/go-stock使用手册.md`, `docs/go-stock帮助问答手册_v2.md`.

---

### Task 1: Backend Free Entitlement Tests

**Files:**
- Create: `backend/data/sponsor_vip_test.go`
- Create: `server/modules/system_rest_test.go`

- [ ] **Step 1: Write the failing data-layer test**

Create `backend/data/sponsor_vip_test.go`:

```go
package data

import "testing"

func TestEffectiveSponsorVipLevelIsFreeAndActive(t *testing.T) {
	SponsorDecryptKeyHex = ""

	level, active := EffectiveSponsorVipLevel()

	if level != 2 {
		t.Fatalf("level = %d, want 2", level)
	}
	if !active {
		t.Fatal("active = false, want true")
	}
}
```

- [ ] **Step 2: Write failing RPC compatibility tests**

Create `server/modules/system_rest_test.go`:

```go
package modules

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-stock/server"
)

func callRPC(t *testing.T, method string, args []any) map[string]any {
	t.Helper()
	body, err := json.Marshal(map[string]any{"args": args})
	if err != nil {
		t.Fatalf("marshal rpc body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/rpc/"+method, bytes.NewReader(body))
	rec := httptest.NewRecorder()
	server.RPCDispatcher(&server.Core{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("%s status = %d, body = %s", method, rec.Code, rec.Body.String())
	}
	var payload map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode %s response: %v", method, err)
	}
	return payload
}

func TestGetSponsorInfoReportsFreeEntitlement(t *testing.T) {
	payload := callRPC(t, "GetSponsorInfo", nil)

	if got := int(payload["vipLevel"].(float64)); got != 2 {
		t.Fatalf("vipLevel = %d, want 2", got)
	}
	if active, _ := payload["active"].(bool); !active {
		t.Fatal("active = false, want true")
	}
	if got, _ := payload["sponsorCode"].(string); got != "" {
		t.Fatalf("sponsorCode = %q, want empty", got)
	}
}

func TestCheckSponsorCodeNoLongerValidatesCodes(t *testing.T) {
	for _, args := range [][]any{nil, {""}, {"not-hex"}} {
		payload := callRPC(t, "CheckSponsorCode", args)
		if got := int(payload["code"].(float64)); got != 1 {
			t.Fatalf("code = %d, want 1 for args %#v", got, args)
		}
	}
}

func TestCheckDeviceBindingDoesNotBlockLocalFreeUse(t *testing.T) {
	payload := callRPC(t, "CheckDeviceBinding", []any{"token", "https://example.invalid"})

	if bound, _ := payload["bound"].(bool); !bound {
		t.Fatal("bound = false, want true")
	}
	if deviceCount := int(payload["deviceCount"].(float64)); deviceCount != 0 {
		t.Fatalf("deviceCount = %d, want 0", deviceCount)
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run:

```bash
go test ./backend/data -run TestEffectiveSponsorVipLevelIsFreeAndActive -v
go test ./server/modules -run 'TestGetSponsorInfoReportsFreeEntitlement|TestCheckSponsorCodeNoLongerValidatesCodes|TestCheckDeviceBindingDoesNotBlockLocalFreeUse' -v
```

Expected: tests fail because current code depends on sponsor code parsing and rejects malformed codes.

---

### Task 2: Backend Free Entitlement Implementation

**Files:**
- Modify: `backend/data/sponsor_vip.go`
- Modify: `server/modules/system_rest.go`
- Test: `backend/data/sponsor_vip_test.go`
- Test: `server/modules/system_rest_test.go`

- [ ] **Step 1: Replace `EffectiveSponsorVipLevel` with free entitlement**

In `backend/data/sponsor_vip.go`, reduce imports to none and replace the function body with:

```go
package data

// DefaultSponsorAESKeyHex is kept for compatibility with older builds and configs.
const DefaultSponsorAESKeyHex = ""

// SponsorDecryptKeyHex is kept for compatibility with older startup code.
var SponsorDecryptKeyHex string

// EffectiveSponsorVipLevel reports the free entitlement used by all local features.
func EffectiveSponsorVipLevel() (level int, active bool) {
	return 2, true
}
```

- [ ] **Step 2: Replace sponsor RPC implementations**

In `server/modules/system_rest.go`, reduce imports to:

```go
import (
	"context"

	"go-stock/backend/data"
	"go-stock/server"
)
```

Replace the `init()` body with:

```go
func init() {
	server.Register("GetSponsorInfo", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		level, active := data.EffectiveSponsorVipLevel()
		return map[string]any{
			"vipLevel":       level,
			"active":         active,
			"sponsorCode":    "",
			"vipStartTime":   "",
			"vipEndTime":     "",
			"vipAuthTime":    "",
			"freeEntitlement": true,
		}, nil
	})

	server.Register("CheckSponsorCode", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return map[string]any{
			"code": 1,
			"msg":  "全部本地功能已免费开放，无需赞助码。",
		}, nil
	})

	server.Register("CheckDeviceBinding", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		return map[string]any{
			"bound":       true,
			"deviceCount": 0,
			"maxDevices":  0,
		}, nil
	})
}
```

- [ ] **Step 3: Run backend tests**

Run:

```bash
gofmt -w backend/data/sponsor_vip.go backend/data/sponsor_vip_test.go server/modules/system_rest.go server/modules/system_rest_test.go
go test ./backend/data -run TestEffectiveSponsorVipLevelIsFreeAndActive -v
go test ./server/modules -run 'TestGetSponsorInfoReportsFreeEntitlement|TestCheckSponsorCodeNoLongerValidatesCodes|TestCheckDeviceBindingDoesNotBlockLocalFreeUse' -v
```

Expected: all selected tests pass.

- [ ] **Step 4: Commit backend compatibility change**

Run:

```bash
git add backend/data/sponsor_vip.go backend/data/sponsor_vip_test.go server/modules/system_rest.go server/modules/system_rest_test.go
git commit -m "feat: make VIP entitlement free"
```

---

### Task 3: Remove Sponsor-Code UI

**Files:**
- Modify: `frontend/src/components/settings.vue`
- Modify: `frontend/src/components/about.vue`

- [ ] **Step 1: Remove sponsor-code import and state from settings**

In `frontend/src/components/settings.vue`, remove `CheckSponsorCode` from the app import list.

Remove this field from `formValue`:

```js
sponsorCode: "",
```

Remove this property from the `new data.SettingConfig` payload:

```js
sponsorCode: formValue.value.sponsorCode,
```

Remove this block from save handling:

```js
if (config.sponsorCode) {
  CheckSponsorCode(config.sponsorCode).then(res => {
    if (!res.code) {
      message.warning(res.msg || '赞助码验证失败')
    }
  })
}
```

Remove this assignment from import/config hydration:

```js
formValue.value.sponsorCode = config.sponsorCode
```

- [ ] **Step 2: Remove sponsor-code form item from settings template**

In `frontend/src/components/settings.vue`, delete the entire form item:

```vue
<n-form-item-gi :span="11" label="赞助码：" path="sponsorCode">
  <n-input-group>
    <n-input :show-count="true" placeholder="联系作者QQ或微信获取，激活VIP功能" v-model:value="formValue.sponsorCode">
    </n-input>
    <n-button type="success" secondary strong
              @click="CheckSponsorCode(formValue.sponsorCode).then((res) => {message.warning(res.msg)})">验证
    </n-button>
    <n-popover trigger="hover" placement="top">
      <template #trigger>
        <n-icon color="#0e7a0d" size="20">
          <HelpCircleFilledIcon />
        </n-icon>
      </template>
      <n-gradient-text :type="'warning'">
        <div style="max-width: 400px;text-align: left">
          赞助码获取方式：<br>
          联系作者获取赞助码，激活VIP功能<br>
          享受更多高级功能和优先支持
        </div>
      </n-gradient-text>
    </n-popover>
  </n-input-group>
</n-form-item-gi>
```

- [ ] **Step 3: Remove VIP status and paid plan table from About**

In `frontend/src/components/about.vue`, replace the `<h1>` block that conditionally shows `VIP{{vipLevel}}` with:

```vue
<h1>
  <n-badge :value="versionInfo" :offset="[80,10]" type="success">
    <n-gradient-text type="info" :size="50">go-stock</n-gradient-text>
  </n-badge>
</h1>
```

Delete this line:

```vue
<n-gradient-text  :type="expired?'error':'warning'" v-if="vipLevel" >vip到期时间：{{vipEndTime}}</n-gradient-text>
```

Replace the sponsor plan divider/table section with:

```vue
<n-divider title-placement="center">支持开源</n-divider>
<n-space vertical align="center">
  <p>go-stock 本地功能已全部免费开放。如果觉得好用，可以自愿支持项目持续维护。</p>
</n-space>
```

Replace the technical support sentence:

```vue
<p>
  开源不易，本人精力和时间有限，如确实需要一对一技术支持，<i style="color: crimson">请先赞助！</i>联系微信(备注 技术支持)：ArvinLovegood
</p>
```

with:

```vue
<p>
  开源不易，本人精力和时间有限，如确实需要一对一技术支持，请联系微信(备注 技术支持)：ArvinLovegood
</p>
```

- [ ] **Step 4: Build-check frontend after UI removal**

Run:

```bash
cd frontend
npm run build
```

Expected: Vite build succeeds without `CheckSponsorCode` reference errors.

---

### Task 4: Remove VIP Denial From Local Feature Entry Points

**Files:**
- Modify: `frontend/src/components/FloatingAiAssistant.vue`
- Modify: `frontend/src/components/FloatingAgentAssistant.vue`
- Modify: `frontend/src/components/kline-analysis.vue`
- Modify: `frontend/src/components/stock.vue`
- Modify: `frontend/src/components/TradingRecordManager.vue`
- Modify: `frontend/src/components/SelectStock.vue`
- Modify: `frontend/src/components/FundRanking.vue`

- [ ] **Step 1: Open assistant panels without VIP checks**

In `FloatingAiAssistant.vue`, replace `togglePanel()` with:

```js
async function togglePanel() {
  if (!panelVisible.value) {
    ensureSummaryEvent()
    openPanel()
  } else {
    closePanel()
  }
}
```

In `FloatingAgentAssistant.vue`, replace `togglePanel()` with:

```js
async function togglePanel() {
  if (!panelVisible.value) {
    openPanel()
  } else {
    closePanel()
  }
}
```

- [ ] **Step 2: Prevent K-line VIP modal from opening**

In `kline-analysis.vue`, replace `startVipCheck()` with:

```js
function startVipCheck() {
  if (vipTimer) {
    clearInterval(vipTimer)
    vipTimer = null
  }
  showVipModal.value = false
}
```

Delete the `<n-modal v-model:show="showVipModal">...</n-modal>` VIP dialog block from the template. If `NCard` is no longer used after that deletion, remove `NCard` from the Naive UI import list.

- [ ] **Step 3: Open stock lightweight K-line directly**

In `stock.vue`, replace the VIP denial block inside `showLightweightKline`:

```js
await refreshEffectiveVip()
// 检查 VIP 权限：有效期内 VIP2 及以上（与 AI 助手 Web 端校验一致）
if (vipLevel.value < 2) {
  message.warning('多周期 K 线仅限 VIP2 及以上用户使用，您当前权限不足，将在 10 秒后自动关闭')
  lwKlineCode.value = em
  lwKlineName.value = name || ''
  modalShow6.value = true
  // 10 秒后自动关闭
  klineAutoCloseTimer.value = setTimeout(() => {
    modalShow6.value = false
    message.info('权限不足，多周期 K 线已自动关闭')
  }, 10000)
  return
}
```

with:

```js
if (klineAutoCloseTimer.value) {
  clearTimeout(klineAutoCloseTimer.value)
  klineAutoCloseTimer.value = null
}
```

- [ ] **Step 4: Open trading-record K-line directly**

In `TradingRecordManager.vue`, replace `openKlineChart(row)` with:

```js
function openKlineChart(row) {
  klineStockCode.value = toEastMoneyCode(row.StockCode)
  klineStockName.value = row.StockName || ''
  showKlineModal.value = true
  longStopLossPrice.value = row.StopLossPrice || 0
  longTakeProfitPrice.value = row.TakeProfitPrice || 0
  costPrice.value = row.Price || 0
}
```

- [ ] **Step 5: Open select-stock K-line directly**

In `SelectStock.vue`, replace the `refreshEffectiveVip().then(() => { ... })` body in `showStockKline(row)` with:

```js
klineStockCode.value = em
klineStockName.value = stockName || ''
klineModalShow.value = true
if (klineAutoCloseTimer) {
  clearTimeout(klineAutoCloseTimer)
  klineAutoCloseTimer = null
}
```

- [ ] **Step 6: Open fund K-line directly**

In `FundRanking.vue`, replace `showStockKline(stockCode, stockName, market)` with:

```js
function showStockKline(stockCode, stockName, market) {
  klineStockCode.value = toEastMoneyCode(stockCode, market)
  klineStockName.value = stockName
  klineModalShow.value = true
  if (klineAutoCloseTimer) {
    clearTimeout(klineAutoCloseTimer)
    klineAutoCloseTimer = null
  }
}
```

- [ ] **Step 7: Build-check feature entry changes**

Run:

```bash
cd frontend
npm run build
```

Expected: build succeeds and no remaining error references to removed timers/functions appear.

---

### Task 5: Remove VIP-Dependent List And Detail Behavior

**Files:**
- Modify: `frontend/src/components/allStockList.vue`
- Modify: `frontend/src/components/aiRecommendStocksList.vue`

- [ ] **Step 1: Let all-stock technical filters work for everyone**

In `allStockList.vue`, replace:

```js
const vipLevel=ref("");
const vipStartTime=ref("");
const vipEndTime=ref("");
const expired=ref(false)
const isValidVip=ref(false) // 是否是会员
```

with:

```js
const isValidVip=ref(true)
```

In the `GetSponsorInfo().then((res) => { ... })` block, replace its body with:

```js
isValidVip.value = true
```

In `loadStocks`, delete:

```js
if((vipLevel.value===""|| Number(vipLevel.value) <=0)){
  handleReset()
}
```

Replace `handleCheckedChange(checked)` with:

```js
function handleCheckedChange(checked) {
}
```

- [ ] **Step 2: Let AI recommendation details work for everyone**

In `aiRecommendStocksList.vue`, replace:

```js
const vipLevel=ref("");
const vipStartTime=ref("");
const vipEndTime=ref("");
const expired=ref(false)
const isValidVip=ref(false) // 是否是会员
```

with:

```js
const vipLevel=ref(2)
const isValidVip=ref(true)
```

In the `GetSponsorInfo().then((res) => { ... })` block, replace its body with:

```js
vipLevel.value = 2
isValidVip.value = true
```

Replace `showDetail(row)` guard:

```js
if(vipLevel.value===""|| Number(vipLevel.value) <=0){
  notify.warning({content: '未开通VIP或者已经过期'})
  return
}
```

with no guard; `showDetail(row)` should start with:

```js
function showDetail(row) {
  modalDataRef.title = row.stockName
```

- [ ] **Step 3: Scan for remaining local denial copy**

Run:

```bash
rg -n "未开通VIP|VIP2|权限不足|自动关闭|赞助码|激活VIP|VIP专属功能" frontend/src/components
```

Expected: only prompt plaza remote labels may remain before Task 6; local feature denial copy should be gone.

---

### Task 6: Remove Prompt Plaza VIP Sync And Blocking

**Files:**
- Modify: `frontend/src/components/promptPlaza.vue`

- [ ] **Step 1: Remove local entitlement imports and dialog use**

Change the import:

```js
import {GetConfig, GetSponsorInfo, GetMachineId, CheckDeviceBinding, GetEffectiveSponsorVip, AddPromptTemplate} from "../api/app";
import {useMessage, useDialog} from "naive-ui";
```

to:

```js
import {GetConfig, AddPromptTemplate} from "../api/app";
import {useMessage} from "naive-ui";
```

Remove:

```js
const dialog = useDialog()
const vipRequireLogin = ref(false)
```

- [ ] **Step 2: Stop VIP login prompts and device checks**

In `onMounted`, replace:

```js
if (token.value) {
  fetchCurrentUser()
} else {
  checkVipAndPromptLogin()
}
```

with:

```js
if (token.value) {
  fetchCurrentUser()
}
```

Replace `fetchCurrentUser()` with:

```js
async function fetchCurrentUser() {
  try {
    const data = await apiGet('/user/me')
    currentUser.value = data
  } catch (e) {
    token.value = ''
    localStorage.removeItem('promptPlazaToken')
    currentUser.value = null
  }
}
```

Delete `checkVipAndPromptLogin`, `checkDeviceLimit`, and `syncVipInfo`.

- [ ] **Step 3: Allow importing VIP-marked prompts locally**

In `addPromptToTemplate(prompt)`, delete:

```js
if (prompt.needVip) {
  const vipInfo = await GetEffectiveSponsorVip()
  if (!vipInfo || vipInfo.vipLevel <= 0 || !vipInfo.active) {
    message.warning('该提示词为VIP专属，请先开通VIP')
    return
  }
}
```

- [ ] **Step 4: Remove VIP-only create/edit controls from local UI**

In `handleCreate`, force the outgoing body to non-VIP:

```js
vipOnly: false
```

In `handleEdit`, force the outgoing body to non-VIP:

```js
vipOnly: false
```

Remove the create modal controls:

```vue
<n-divider vertical />
<n-text>VIP专属</n-text>
<n-switch v-model:value="createModal.vipOnly" />
<n-text depth="3" style="font-size: 12px">仅VIP用户可查看完整内容</n-text>
```

Remove the edit modal controls:

```vue
<n-divider vertical />
<n-text>VIP专属</n-text>
<n-switch v-model:value="editModal.vipOnly" />
<n-text depth="3" style="font-size: 12px">仅VIP用户可查看完整内容</n-text>
```

Change the login modal tag:

```vue
<n-modal v-model:show="loginModal.show" preset="card" style="width: 400px" :title="vipRequireLogin ? '🎉 VIP专属福利' : '账号'" :closable="!vipRequireLogin" :maskClosable="!vipRequireLogin" :closeOnEsc="!vipRequireLogin">
```

to:

```vue
<n-modal v-model:show="loginModal.show" preset="card" style="width: 400px" title="账号">
```

Delete the `<div v-if="vipRequireLogin">...</div>` benefits panel.

- [ ] **Step 5: Build-check prompt plaza**

Run:

```bash
cd frontend
npm run build
```

Expected: build succeeds with no `vipRequireLogin`, `GetSponsorInfo`, `GetMachineId`, `CheckDeviceBinding`, or `GetEffectiveSponsorVip` unresolved references in `promptPlaza.vue`.

---

### Task 7: Documentation Update

**Files:**
- Modify: `README.md`
- Modify: `docs/go-stock使用手册.md`
- Modify: `docs/go-stock帮助问答手册_v2.md`

- [ ] **Step 1: Replace README paid plan table with free statement**

In `README.md`, replace the "支持开源计划" table with text:

```markdown
### 支持开源

go-stock 本地功能已全部免费开放。项目仍欢迎自愿赞助，用于持续维护、文档、社区答疑和新功能开发。
```

- [ ] **Step 2: Remove VIP restriction language from manual**

Run:

```bash
rg -n "VIP|赞助码|赞助计划|权限不足|开通VIP|VIP2" docs/go-stock使用手册.md
```

For each hit, rewrite the sentence so the feature is described as free. Use these concrete replacements:

```markdown
> 多周期 K 线图已免费开放。美股代码（`gb_` 前缀）不支持多周期 K 线。
```

```markdown
> 浮动 AI 助手已免费开放。
```

```markdown
全市场股票形态筛选工具，支持按各种指标条件筛选，支持排序和搜索，可将筛选结果添加到自选，并可查看 K 线图。
```

```markdown
| 保存设置 | 保存所有配置 |
```

- [ ] **Step 3: Remove VIP restriction language from Q&A**

Run:

```bash
rg -n "VIP|赞助码|赞助计划|权限不足|开通VIP|VIP2" docs/go-stock帮助问答手册_v2.md
```

Rewrite Q&A entries so they state:

```markdown
**A:** 多周期 K 线功能已免费开放，可直接使用。
```

```markdown
**A:** 浮动 AI 助手已免费开放，可直接使用。
```

```markdown
**A:** 交易记录管理中的 K 线查看已免费开放，可叠加显示持仓成本线、止损线、止盈线。
```

- [ ] **Step 4: Final paid-copy scan**

Run:

```bash
rg -n "VIP|赞助码|开通VIP|权限不足|VIP2|赞助计划" README.md docs/go-stock使用手册.md docs/go-stock帮助问答手册_v2.md frontend/src/components backend server
```

Expected: no paid gating language remains. Mentions of third-party URLs containing `vip.stock.finance.sina.com.cn` are allowed because they are Sina API hosts, not go-stock entitlement.

---

### Task 8: Full Verification And Final Commit

**Files:**
- Verify all modified files.

- [ ] **Step 1: Run backend focused tests**

Run:

```bash
go test ./backend/data -run TestEffectiveSponsorVipLevelIsFreeAndActive -v
go test ./server/modules -run 'TestGetSponsorInfoReportsFreeEntitlement|TestCheckSponsorCodeNoLongerValidatesCodes|TestCheckDeviceBindingDoesNotBlockLocalFreeUse' -v
```

Expected: all selected tests pass.

- [ ] **Step 2: Run broader backend tests if feasible**

Run:

```bash
go test ./server/... ./backend/data/...
```

Expected: tests pass, or failures are unrelated external-network/data-source tests and are documented in the final response.

- [ ] **Step 3: Run frontend build**

Run:

```bash
cd frontend
npm run build
```

Expected: Vite build succeeds.

- [ ] **Step 4: Inspect git diff for accidental unrelated churn**

Run:

```bash
git status --short
git diff --stat
```

Expected: only planned files are modified.

- [ ] **Step 5: Commit implementation**

Run:

```bash
git add backend/data/sponsor_vip.go backend/data/sponsor_vip_test.go server/modules/system_rest.go server/modules/system_rest_test.go frontend/src/components/settings.vue frontend/src/components/about.vue frontend/src/components/FloatingAiAssistant.vue frontend/src/components/FloatingAgentAssistant.vue frontend/src/components/kline-analysis.vue frontend/src/components/stock.vue frontend/src/components/TradingRecordManager.vue frontend/src/components/SelectStock.vue frontend/src/components/FundRanking.vue frontend/src/components/allStockList.vue frontend/src/components/aiRecommendStocksList.vue frontend/src/components/promptPlaza.vue README.md docs/go-stock使用手册.md docs/go-stock帮助问答手册_v2.md
git commit -m "feat: remove VIP gating"
```

Expected: implementation commit is created after tests/build.
