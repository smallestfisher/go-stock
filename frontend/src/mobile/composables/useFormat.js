/**
 * 移动端公共格式化工具
 * 统一金额、百分比、热度等数据的展示格式，避免各页面重复实现
 */

/**
 * 格式化金额：自动转 万/亿
 * @param {number} value 金额数值
 * @param {number} digits 小数位
 * @returns {string}
 */
export function formatMoney(value, digits = 2) {
  if (value == null || value === '' || isNaN(value)) return '--'
  const n = Number(value)
  const sign = n < 0 ? '-' : ''
  const abs = Math.abs(n)
  if (abs >= 100000000) return `${sign}${(abs / 100000000).toFixed(digits)}亿`
  if (abs >= 10000) return `${sign}${(abs / 10000).toFixed(digits)}万`
  return `${sign}${abs}`
}

/**
 * 格式化成交量/额（与 formatMoney 同义，保留语义化命名）
 */
export const formatVolume = formatMoney

/**
 * 格式化百分比：加正负号，保留两位小数
 * @param {number} value 百分比数值（如 2.34）
 * @returns {string} 如 +2.34% / -1.20%
 */
export function formatPercent(value) {
  const n = Number(value)
  if (!Number.isFinite(n)) return '--'
  return `${n > 0 ? '+' : ''}${n.toFixed(2)}%`
}

/**
 * 根据涨跌返回颜色类名（m-rise/m-fall/m-flat）
 * @param {number} value 涨跌值或百分比
 * @returns {string}
 */
export function trendClass(value) {
  const n = Number(value)
  if (!Number.isFinite(n) || n === 0) return 'm-flat'
  return n > 0 ? 'm-rise' : 'm-fall'
}

/**
 * 格式化热度值：大数转 K/M
 * @param {number|string} heat 热度值
 * @returns {string}
 */
export function formatHeat(heat) {
  if (typeof heat === 'string') return heat
  heat = Number(heat) || 0
  if (heat >= 1000000) return `${(heat / 1000000).toFixed(1)}M`
  if (heat >= 1000) return `${(heat / 1000).toFixed(0)}K`
  return String(heat)
}
