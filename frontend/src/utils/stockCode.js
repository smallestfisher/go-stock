/**
 * 应用内部股票代码（如 sh000001 / hk00700 / gb_AAPL）与东方财富格式（如 000001.SZ）互转。
 * 抽离自 stock.vue 供移动端详情抽屉与桌面共用。
 */

/** 内部代码转东方财富格式（如 sh000001 → 000001.SZ） */
export function toEastMoneyCode(code) {
    if (!code) return ''
    const c = String(code).trim()
    if (/\.(SH|SZ|BJ|HK|US|SS)$/i.test(c)) return c.toUpperCase()
    const lower = c.toLowerCase()
    if (lower.startsWith('sh')) return lower.slice(2) + '.SH'
    if (lower.startsWith('sz')) return lower.slice(2) + '.SZ'
    if (lower.startsWith('bj')) return lower.slice(2) + '.BJ'
    if (lower.startsWith('hk')) return lower.slice(2).toUpperCase() + '.HK'
    if (lower.startsWith('us')) return lower.slice(2).toUpperCase() + '.US'
    if (lower.startsWith('gb_')) return lower.slice(3).toUpperCase() + '.US'
    if (/^\d+$/.test(c)) {
        const d = c[0]
        if (d === '6') return c + '.SH'
        if (d === '0' || d === '3') return c + '.SZ'
        if (d === '8' || d === '9') return c + '.BJ'
        return c + '.SZ'
    }
    // 纯字母代码视为美股（如 AAPL → AAPL.US）
    if (/^[a-zA-Z]+$/.test(c)) return c.toUpperCase() + '.US'
    return ''
}

/** 东方财富格式转回应用内部代码格式（如 000001.SZ → sh000001） */
export function fromEastMoneyCode(emCode) {
    if (!emCode) return ''
    const c = String(emCode).trim().toUpperCase()
    if (c.endsWith('.SH')) return 'sh' + c.slice(0, -3)
    if (c.endsWith('.SZ')) return 'sz' + c.slice(0, -3)
    if (c.endsWith('.BJ')) return 'bj' + c.slice(0, -3)
    if (c.endsWith('.HK')) return 'hk' + c.slice(0, -3).toLowerCase()
    if (c.endsWith('.US')) return 'us' + c.slice(0, -3).toLowerCase()
    return c.toLowerCase()
}
