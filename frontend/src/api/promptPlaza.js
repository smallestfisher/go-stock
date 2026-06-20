import { getToken, serverURL } from "./transport.js";

const PROXY_PREFIX = "/api/prompt-plaza";

export function promptPlazaURL(path, params = {}) {
  const normalizedPath = path.startsWith("/") ? path : "/" + path;
  const url = new URL(serverURL(PROXY_PREFIX + normalizedPath));
  Object.entries(params).forEach(([key, value]) => {
    if (value !== null && value !== undefined && value !== "") {
      url.searchParams.set(key, value);
    }
  });
  return url.toString();
}

export function promptPlazaHeaders(promptPlazaToken = "") {
  const headers = { "Content-Type": "application/json" };
  const goStockToken = getToken();
  if (goStockToken) {
    headers["X-Go-Stock-Token"] = goStockToken;
  }
  if (promptPlazaToken) {
    headers.Authorization = `Bearer ${promptPlazaToken}`;
  }
  return headers;
}

export async function parsePromptPlazaResponse(resp) {
  const text = await resp.text();
  let json;
  try {
    json = text ? JSON.parse(text) : {};
  } catch (_) {
    throw new Error(`接口返回非JSON (HTTP ${resp.status}): ${text.substring(0, 200)}`);
  }
  if (!resp.ok) {
    throw new Error(json.message || `请求失败 (HTTP ${resp.status})`);
  }
  if (json.code !== 0) {
    throw new Error(json.message || "请求失败");
  }
  return json.data;
}
