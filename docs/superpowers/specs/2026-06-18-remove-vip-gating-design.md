# Remove VIP Gating Design

## Objective

Remove paid/VIP gating from go-stock so all local application features are free by default. The change should remove user-visible VIP and sponsor-code workflows while keeping existing RPC names compatible enough that older frontend calls do not break during the transition.

## Current State

The server exposes sponsor-related RPCs from `server/modules/system.go` and `server/modules/system_rest.go`:

- `GetEffectiveSponsorVip` returns `vipLevel` and `active` from `backend/data.EffectiveSponsorVipLevel`.
- `GetSponsorInfo` returns the same effective level plus the saved sponsor code.
- `CheckSponsorCode` validates and stores an encrypted sponsor code.
- `CheckDeviceBinding` checks prompt plaza device limits through the remote community API.

The frontend uses those values to restrict several features:

- AI assistant floating panels.
- K-line and multi-period K-line views.
- Trading record K-line details.
- Fund ranking K-line views.
- Stock and AI recommendation detail enhancements.
- Technical stock filtering.
- Prompt plaza VIP prompt access, VIP login prompts, VIP sync, and device binding checks.

Settings and About pages still expose sponsor code inputs, VIP status, expiration, and paid plan tables.

## Recommended Approach

Use a compatibility-first removal:

1. Make the server report a free, active entitlement for compatibility.
2. Remove frontend gates and paid messaging so features open directly.
3. Remove sponsor-code UI and prompt plaza VIP/device-binding sync.
4. Update documentation to describe all features as free.

This keeps the behavioral source of truth on the server while making the product experience match the new free model.

## Backend Design

`backend/data.EffectiveSponsorVipLevel` will stop depending on saved encrypted sponsor codes for feature access. It should return a stable active level high enough for all existing feature checks, currently `2, true`.

`GetEffectiveSponsorVip` should keep its response shape:

```json
{"vipLevel": 2, "active": true}
```

`GetSponsorInfo` should keep its response shape but no longer expose saved sponsor codes as meaningful access credentials. It should return an active free entitlement and empty sponsor-code fields where practical.

`CheckSponsorCode` should become a compatibility no-op. It should not require, validate, or store a sponsor code. It should return a successful message explaining that sponsor codes are no longer needed.

`CheckDeviceBinding` should not block local free application usage. Prompt plaza can still keep normal account login behavior, but local entitlement and device binding should not be used to gate local features.

## Frontend Design

Feature gates based on VIP level should be removed or made unreachable. Users should be able to open the feature immediately without a warning, auto-close timer, or VIP modal.

Affected local feature areas:

- `FloatingAiAssistant.vue`
- `FloatingAgentAssistant.vue`
- `kline-analysis.vue`
- `stock.vue`
- `TradingRecordManager.vue`
- `SelectStock.vue`
- `FundRanking.vue`
- `allStockList.vue`
- `aiRecommendStocksList.vue`

Settings should remove the sponsor-code form field, validation button, and sponsor-code help copy.

About should remove VIP badges, expiration display, and paid-plan table. It may keep a voluntary donation/support section, but it must not imply feature access depends on payment.

Prompt plaza should stop:

- Prompting local VIP users to log in for VIP-only benefits.
- Syncing local VIP level, expiration, or sponsor code to the remote API.
- Checking device binding as a condition of local usage.
- Blocking prompt import because a prompt is marked VIP-only.

Prompt plaza may still display remote metadata such as author badges or prompt labels if the remote API sends them, but those labels must not block local import.

## Data Flow

After the change:

1. The frontend may still call `GetEffectiveSponsorVip` or `GetSponsorInfo`.
2. The server returns active free access.
3. Frontend components do not branch into paid-access denial flows.
4. Features execute normal data loading and rendering paths.

Saved `sponsorCode` values in existing configs should be harmless. The migration should not delete user config fields unless that is already part of normal config saving.

## Error Handling

Sponsor-code validation errors should disappear from normal user flows because the UI no longer asks for codes.

If old clients call `CheckSponsorCode`, the server should return success and not fail on an empty or malformed code.

If prompt plaza remote calls fail, normal prompt plaza error handling should remain unchanged. Removing VIP sync should not hide errors from unrelated login, prompt loading, or prompt import operations.

## Documentation

Update user-facing docs that currently say features require VIP2, sponsor codes, paid plans, or permission upgrades. The docs should state that all local features are free, while optional donation/support remains voluntary if kept.

Primary files:

- `README.md`
- `docs/go-stock使用手册.md`
- `docs/go-stock帮助问答手册_v2.md`

## Tests And Verification

Run backend tests around server modules and data where feasible:

```bash
go test ./server/... ./backend/data/...
```

Run frontend build:

```bash
cd frontend
npm run build
```

Manual verification should cover:

- Settings no longer shows sponsor-code input.
- About no longer shows VIP status, expiration, or paid plan gating.
- AI assistant and AI agent assistant open without VIP warnings.
- Multi-period K-line opens without auto-close behavior.
- Stock selection, fund ranking, trading record K-line, all-stock filtering, and AI recommendation detail views work without VIP warnings.
- Prompt plaza imports VIP-marked prompts without local VIP blocking and does not perform local VIP sync or device-limit checks.

## Non-Goals

This change does not redesign prompt plaza account authentication or remove normal remote account login. It only removes local paid/VIP gating and sponsor-code/device-binding flows tied to entitlement.

This change does not remove third-party API keys needed for data sources or AI providers. Those are external service credentials, not go-stock VIP gates.
