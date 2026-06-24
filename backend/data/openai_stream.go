package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"go-stock/backend/logger"
	"go-stock/backend/models"
	"go-stock/backend/util"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/tidwall/gjson"
)

func (o *OpenAi) NewSummaryStockNewsStreamWithTools(userQuestion string, sysPromptId *int, tools []Tool, thinking bool, history []map[string]interface{}) <-chan map[string]any {
	ch := make(chan map[string]any, 512)
	defer func() {
		if err := recover(); err != nil {
			logger.SugaredLogger.Error("NewSummaryStockNewsStream panic", err)
		}
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.SugaredLogger.Errorf("NewSummaryStockNewsStream goroutine panic: %s", err)
				logger.SugaredLogger.Errorf("NewSummaryStockNewsStream goroutine panic config: %s", o.String())
			}
		}()
		defer close(ch)

		sysPrompt := ""
		if sysPromptId == nil || *sysPromptId == 0 {
			sysPrompt = o.Prompt
		} else {
			sysPrompt = NewPromptTemplateApi().GetPromptTemplateByID(*sysPromptId)
		}
		if sysPrompt == "" {
			sysPrompt = o.Prompt
		}

		sysPrompt += `

【强制规则】你必须通过工具调用获取实时数据，严禁凭记忆编造或使用过时数据。以下场景必须调用工具：
1. 股票/指数行情数据（价格、涨跌幅、成交量等）——必须调用工具获取最新实时数据
2. 财务数据（营收、利润、市盈率等）——必须调用工具获取最新财报数据
3. 新闻资讯——必须调用工具获取最新新闻
4. 宏观经济数据——必须调用工具获取最新数据
任何涉及具体数字的回答，都必须先通过工具查询确认，不得使用训练数据中的过时信息。如果你没有获取到最新数据，必须明确告知用户"当前未能获取到最新数据"，绝不能编造数据。`

		sysPrompt += "最后必须调用CreateAiRecommendStocks工具函数保存ai股票推荐记录。"

		msg := []map[string]interface{}{
			{
				"role":    "system",
				"content": sysPrompt,
			},
		}
		msg = append(msg, map[string]interface{}{
			"role":    "user",
			"content": "当前时间",
		})
		msg = append(msg, map[string]interface{}{
			"role":              "assistant",
			"reasoning_content": "使用工具查询",
			"content":           "当前本地时间是:" + time.Now().Format("2006-01-02 15:04:05"),
		})

		if userQuestion == "" {
			userQuestion = "请根据当前时间，总结和分析股票市场新闻中的投资机会"
		}
		msg = append(msg, map[string]interface{}{
			"role":    "user",
			"content": userQuestion,
		})
		AskAiWithTools(o, errors.New(""), msg, ch, userQuestion, tools, thinking)
	}()
	return ch
}

func (o *OpenAi) NewSummaryStockNewsStream(userQuestion string, sysPromptId *int, think bool, history []map[string]interface{}) <-chan map[string]any {
	ch := make(chan map[string]any, 512)
	defer func() {
		if err := recover(); err != nil {
			logger.SugaredLogger.Error("NewSummaryStockNewsStream panic", err)
		}
	}()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.SugaredLogger.Errorf("NewSummaryStockNewsStream goroutine  panic :%s", err)
				logger.SugaredLogger.Errorf("NewSummaryStockNewsStream goroutine  panic  config:%s", o.String())
			}
		}()
		defer close(ch)

		sysPrompt := ""
		if sysPromptId == nil || *sysPromptId == 0 {
			sysPrompt = o.Prompt
		} else {
			sysPrompt = NewPromptTemplateApi().GetPromptTemplateByID(*sysPromptId)
		}
		if sysPrompt == "" {
			sysPrompt = o.Prompt
		}
		msg := []map[string]interface{}{
			{
				"role":    "system",
				"content": sysPrompt,
			},
		}
		msg = append(msg, map[string]interface{}{
			"role":    "user",
			"content": "当前时间",
		})
		msg = append(msg, map[string]interface{}{
			"role":    "assistant",
			"content": "当前本地时间是:" + time.Now().Format("2006-01-02 15:04:05"),
		})
		wg := &sync.WaitGroup{}
		// 同 NewChatStream：用 pairCh 收集各 goroutine 产出的消息对，wg.Wait 后串行 append，
		// 消除并发 append 同一 slice 的数据竞争。
		pairCh := make(chan []map[string]interface{}, 3)
		wg.Add(3)

		go func() {
			defer wg.Done()
			md := strings.Builder{}
			res := NewMarketNewsApi().ClsCalendar()
			for _, a := range res {
				bytes, err := json.Marshal(a)
				if err != nil {
					continue
				}
				date := gjson.Get(string(bytes), "calendar_day")
				md.WriteString("\n### 事件/会议日期：" + date.String())
				list := gjson.Get(string(bytes), "items")
				list.ForEach(func(key, value gjson.Result) bool {
					md.WriteString("\n- " + gjson.Get(value.String(), "title").String())
					return true
				})
			}
			pairCh <- []map[string]interface{}{
				{"role": "user", "content": "近期重大事件/会议"},
				{"role": "assistant", "reasoning_content": "使用工具查询", "content": "近期重大事件/会议如下：\n" + md.String()},
			}
		}()

		go func() {
			defer wg.Done()
			datas := NewMarketNewsApi().InteractiveAnswer(1, 100, "")
			content := util.MarkdownTableWithTitle("当前最新投资者互动数据", datas.Results)
			pairCh <- []map[string]interface{}{
				{"role": "user", "content": "投资者互动数据"},
				{"role": "assistant", "content": content},
			}
		}()

		go func() {
			defer wg.Done()
			res := NewSearchStockApi("").HotStrategy()
			bytes, _ := json.Marshal(res)
			strategy := &models.HotStrategy{}
			json.Unmarshal(bytes, strategy)
			for _, data := range strategy.Data {
				data.Chg = mathutil.RoundToFloat(100*data.Chg, 2)
			}
			markdownTable := util.MarkdownTableWithTitle("当前热门选股策略", strategy.Data)
			pairCh <- []map[string]interface{}{
				{"role": "user", "content": "当前热门选股策略"},
				{"role": "assistant", "content": markdownTable},
			}
		}()

		go func() {
			wg.Wait()
			close(pairCh)
		}()

		for pairs := range pairCh {
			msg = append(msg, pairs...)
		}

		// 资讯条数过多会单独占满 context；限制在合理范围（原先 200–1000 极易撑爆上下文）
		news := NewMarketNewsApi().GetNews24HoursList("", random.RandInt(20, 40))
		messageText := strings.Builder{}
		for _, telegraph := range *news {
			messageText.WriteString("## " + telegraph.Time + ":" + "\n")
			messageText.WriteString("### " + telegraph.Content + "\n")
		}

		msg = append(msg, map[string]interface{}{
			"role":    "user",
			"content": "市场资讯",
		})
		msg = append(msg, map[string]interface{}{
			"role":    "assistant",
			"content": messageText.String(),
		})

		//for _, m := range TrimAiAssistantHistoryForAPI(history) {
		//	msg = append(msg, m)
		//}
		if userQuestion == "" {
			userQuestion = "请根据当前时间，总结和分析股票市场新闻中的投资机会"
		}
		msg = append(msg, map[string]interface{}{
			"role":    "user",
			"content": userQuestion,
		})
		AskAi(o, errors.New(""), msg, ch, userQuestion, think)
	}()
	return ch
}

func buildStockAnalysisContext(stock, stockCode string, followedStock FollowedStock) string {
	stock = RemoveAllBlankChar(stock)
	stockCode = RemoveAllBlankChar(stockCode)

	var context strings.Builder
	context.WriteString("本次AI诊股的股票上下文如下，后续所有分析都必须围绕该股票展开：\n")
	if stock != "" {
		context.WriteString("- 股票名称：" + stock + "\n")
	}
	if stockCode != "" {
		context.WriteString("- 股票代码：" + stockCode + "\n")
	}
	if followedStock.Name != "" && followedStock.Name != stock {
		context.WriteString("- 自选股名称：" + followedStock.Name + "\n")
	}
	if followedStock.StockCode != "" && followedStock.StockCode != stockCode {
		context.WriteString("- 自选股代码：" + followedStock.StockCode + "\n")
	}
	if followedStock.CostPrice > 0 {
		context.WriteString(fmt.Sprintf("- 用户持仓成本价：%.4f\n", followedStock.CostPrice))
	}
	if followedStock.Volume > 0 {
		context.WriteString(fmt.Sprintf("- 用户持仓数量：%d\n", followedStock.Volume))
	}
	if followedStock.EntryPrice > 0 {
		context.WriteString(fmt.Sprintf("- 开仓价：%.4f\n", followedStock.EntryPrice))
	}
	if followedStock.TakeProfitPrice > 0 {
		context.WriteString(fmt.Sprintf("- 止盈价：%.4f\n", followedStock.TakeProfitPrice))
	}
	if followedStock.StopLossPrice > 0 {
		context.WriteString(fmt.Sprintf("- 止损价：%.4f\n", followedStock.StopLossPrice))
	}
	context.WriteString("如果用户只提到成本价、持仓、后续操作等问题，默认都是针对上述股票；不要反问用户是哪只股票。实时行情、K线、财务、新闻、公告等数据必须优先通过工具获取。")
	return context.String()
}

// ---- 技术指标本地计算 ----
// 目的：避免让 AI 凭记忆从原始 OHLCV 心算 MACD/KDJ/RSI（必然不准），
// 由后端算好最新指标值 + 信号判断，连同 K 线一起喂给 AI。
// 公式均采用主流默认周期，与东方财富/同花顺一致。

// computeEMA 序列的指数移动平均（返回与输入等长的数组，前 period-1 个为 0）。
func computeEMA(values []float64, period int) []float64 {
	if period <= 0 || len(values) == 0 {
		return nil
	}
	ema := make([]float64, len(values))
	mult := 2.0 / float64(period+1)
	// 首个有效点取前 period 个的简单平均作为种子
	if len(values) >= period {
		sum := 0.0
		for i := 0; i < period; i++ {
			sum += values[i]
		}
		ema[period-1] = sum / float64(period)
		for i := period; i < len(values); i++ {
			ema[i] = (values[i]-ema[i-1])*mult + ema[i-1]
		}
	}
	return ema
}

// computeMACD 返回最新一根的 (dif, dea, macd柱, 是否金叉, 是否死叉)。
// 金叉/死叉依据最近两根 dif 与 dea 的穿越判断。
func computeMACD(closes []float64) (dif, dea, macd float64, golden, death bool, ok bool) {
	if len(closes) < 30 {
		return 0, 0, 0, false, false, false
	}
	ema12 := computeEMA(closes, 12)
	ema26 := computeEMA(closes, 26)
	n := len(closes)
	difs := make([]float64, n)
	for i := 0; i < n; i++ {
		difs[i] = ema12[i] - ema26[i]
	}
	deaArr := computeEMA(difs, 9)
	dif = difs[n-1]
	dea = deaArr[n-1]
	macd = (dif - dea) * 2
	// 信号：dif 上穿 dea = 金叉，下穿 = 死叉
	prevDiff := difs[n-2] - deaArr[n-2]
	curDiff := dif - dea
	golden = prevDiff <= 0 && curDiff > 0
	death = prevDiff >= 0 && curDiff < 0
	return dif, dea, macd, golden, death, true
}

// computeKDJ 返回最新一根的 (k, d, j)。经典 9 日 KDJ。
func computeKDJ(klines []KLineData) (k, d, j float64, ok bool) {
	n := len(klines)
	if n < 9 {
		return 0, 0, 0, false
	}
	rsv := make([]float64, n)
	for i := 0; i < n; i++ {
		lo, _ := parseFloatToFloat(klines[i].Low)
		hi, _ := parseFloatToFloat(klines[i].High)
		cl, _ := parseFloatToFloat(klines[i].Close)
		hn := hi
		ln := lo
		for j := i - 8; j <= i; j++ {
			if j < 0 {
				continue
			}
			h, _ := parseFloatToFloat(klines[j].High)
			l, _ := parseFloatToFloat(klines[j].Low)
			if h > hn {
				hn = h
			}
			if l < ln {
				ln = l
			}
		}
		if hn-ln != 0 {
			rsv[i] = (cl - ln) / (hn - ln) * 100
		} else {
			rsv[i] = 50
		}
	}
	// K = 2/3 * 前K + 1/3 * RSV；D = 2/3 * 前D + 1/3 * K；初始 K=D=50
	kArr := make([]float64, n)
	dArr := make([]float64, n)
	pk, pd := 50.0, 50.0
	for i := 0; i < n; i++ {
		pk = 2.0/3.0*pk + 1.0/3.0*rsv[i]
		pd = 2.0/3.0*pd + 1.0/3.0*pk
		kArr[i] = pk
		dArr[i] = pd
	}
	k = kArr[n-1]
	d = dArr[n-1]
	j = 3*k - 2*d
	return k, d, j, true
}

// computeRSI 返回最新一根的 RSI(period)。采用 Wilder 平滑。
func computeRSI(closes []float64, period int) (float64, bool) {
	n := len(closes)
	if n < period+1 {
		return 0, false
	}
	var gainSum, lossSum float64
	for i := 1; i <= period; i++ {
		ch := closes[i] - closes[i-1]
		if ch >= 0 {
			gainSum += ch
		} else {
			lossSum -= ch
		}
	}
	avgGain := gainSum / float64(period)
	avgLoss := lossSum / float64(period)
	for i := period + 1; i < n; i++ {
		ch := closes[i] - closes[i-1]
		gain, loss := 0.0, 0.0
		if ch >= 0 {
			gain = ch
		} else {
			loss = -ch
		}
		avgGain = (avgGain*float64(period-1) + gain) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + loss) / float64(period)
	}
	if avgLoss == 0 {
		return 100, true
	}
	rs := avgGain / avgLoss
	return 100 - 100/(1+rs), true
}

// buildKLineIndicatorSummary 本地计算并汇总技术指标，返回给 AI 的中文描述。
func buildKLineIndicatorSummary(klines []KLineData) string {
	n := len(klines)
	if n < 20 {
		return ""
	}
	closes := make([]float64, n)
	for i, k := range klines {
		v, _ := parseFloatToFloat(k.Close)
		closes[i] = v
	}
	var b strings.Builder
	b.WriteString("\n## 技术指标（后端本地计算，请直接参考，勿重新心算）：\n")

	// MA
	for _, p := range []int{5, 10, 20, 60} {
		if ma := computeSMA(closes, n-1, p); ma >= 0 {
			b.WriteString(fmt.Sprintf("- MA%d: %.2f；", p, ma))
		}
	}
	b.WriteString("\n")

	// MACD
	if dif, dea, macd, golden, death, ok := computeMACD(closes); ok {
		b.WriteString(fmt.Sprintf("- MACD(12,26,9): DIF=%.2f, DEA=%.2f, MACD柱=%.2f", dif, dea, macd))
		if golden {
			b.WriteString("，【金叉信号（DIF 上穿 DEA）】")
		} else if death {
			b.WriteString("，【死叉信号（DIF 下穿 DEA）】")
		} else if dif > dea {
			b.WriteString("，DIF>DEA（多头）")
		} else {
			b.WriteString("，DIF<DEA（空头）")
		}
		b.WriteString("\n")
	}

	// KDJ
	if k, d, j, ok := computeKDJ(klines); ok {
		b.WriteString(fmt.Sprintf("- KDJ(9,3,3): K=%.2f, D=%.2f, J=%.2f", k, d, j))
		if j < 0 {
			b.WriteString("，J<0（超卖）")
		} else if j > 100 {
			b.WriteString("，J>100（超买）")
		} else if k > d {
			b.WriteString("，K>D")
		} else {
			b.WriteString("，K<D")
		}
		b.WriteString("\n")
	}

	// RSI
	if rsi, ok := computeRSI(closes, 14); ok {
		b.WriteString(fmt.Sprintf("- RSI(14): %.2f", rsi))
		if rsi >= 70 {
			b.WriteString("（超买）")
		} else if rsi <= 30 {
			b.WriteString("（超卖）")
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (o *OpenAi) NewChatStream(stock, stockCode, userQuestion string, sysPromptId *int, tools []Tool, thinking bool) <-chan map[string]any {
	ch := make(chan map[string]any, 512)

	defer func() {
		if err := recover(); err != nil {
			logger.SugaredLogger.Error("NewChatStream panic", err)
		}
	}()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.SugaredLogger.Errorf("NewChatStream goroutine  panic :%s", err)
				logger.SugaredLogger.Errorf("NewChatStream goroutine  panic  stock:%s stockCode:%s", stock, stockCode)
				logger.SugaredLogger.Errorf("NewChatStream goroutine  panic  config:%s", o.String())
			}
		}()
		defer close(ch)

		sysPrompt := ""
		if sysPromptId == nil || *sysPromptId == 0 {
			sysPrompt = o.Prompt
		} else {
			sysPrompt = NewPromptTemplateApi().GetPromptTemplateByID(*sysPromptId)
		}
		if sysPrompt == "" {
			sysPrompt = o.Prompt
		}
		if len(tools) > 0 {
			sysPrompt += `

【实时数据规则】分析具体股票时，必须优先通过工具获取实时行情、K线、财务、新闻、公告等数据；如工具未能获取到数据，必须明确说明数据缺失，不能凭记忆编造具体数字。`
		}

		msg := []map[string]interface{}{
			{
				"role":    "system",
				"content": sysPrompt,
			},
		}

		msg = append(msg, map[string]interface{}{
			"role":    "user",
			"content": "当前时间",
		})
		msg = append(msg, map[string]interface{}{
			"role":    "assistant",
			"content": "当前本地时间是:" + time.Now().Format("2006-01-02 15:04:05"),
		})

		replaceTemplates := map[string]string{
			"{{stockName}}": RemoveAllBlankChar(stock),
			"{{stockCode}}": RemoveAllBlankChar(stockCode),
			"{stockName}":   RemoveAllBlankChar(stock),
			"{stockCode}":   RemoveAllBlankChar(stockCode),
			"stockName":     RemoveAllBlankChar(stock),
			"stockCode":     RemoveAllBlankChar(stockCode),
		}
		followedStock := NewStockDataApi().GetFollowedStockByStockCode(stockCode)
		if followedStock.CostPrice > 0 {
			replaceTemplates["{{costPrice}}"] = convertor.ToString(followedStock.CostPrice)
			replaceTemplates["{costPrice}"] = convertor.ToString(followedStock.CostPrice)
			replaceTemplates["costPrice"] = convertor.ToString(followedStock.CostPrice)
		}

		question := ""
		if userQuestion == "" {
			question = strutil.ReplaceWithMap(o.QuestionTemplate, replaceTemplates)
		} else {
			question = strutil.ReplaceWithMap(userQuestion, replaceTemplates)
		}

		if len(tools) > 0 {
			msg = append(msg, map[string]interface{}{
				"role":    "user",
				"content": buildStockAnalysisContext(stock, stockCode, followedStock),
			})
			msg = append(msg, map[string]interface{}{
				"role":    "user",
				"content": question,
			})
			AskAiWithTools(o, errors.New(""), msg, ch, question, tools, thinking)
			return
		}

		stockData, err := NewStockDataApi().GetStockCodeRealTimeData(stockCode)
		if err == nil && len(*stockData) > 0 {
			msg = append(msg, map[string]interface{}{
				"role":    "user",
				"content": fmt.Sprintf("当前%s[%s]价格是多少？", stock, stockCode),
			})
			msg = append(msg, map[string]interface{}{
				"role":    "assistant",
				"content": fmt.Sprintf("截止到%s,当前%s[%s]价格是%s", (*stockData)[0].Date+" "+(*stockData)[0].Time, stock, stockCode, (*stockData)[0].Price),
			})
		}

		wg := &sync.WaitGroup{}
		// 每个 goroutine 把自己产出的若干条消息（成对的 user/assistant）作为一个 []map 发到 pairCh，
		// 在 wg.Wait() 之后由主 goroutine 串行收齐并 append 到 msg。
		// 这样彻底消除原来 8 个 goroutine 并发 append 同一 slice 的数据竞争，
		// 且保证每对 user/assistant 不会被打乱、财报的多条也不会与其他数据交错。
		pairCh := make(chan []map[string]interface{}, 8)
		wg.Add(8)

		go func() {
			defer wg.Done()
			datas := NewMarketNewsApi().InteractiveAnswer(1, 100, stock)
			content := util.MarkdownTableWithTitle("当前最新投资者互动数据", datas.Results)
			pairCh <- []map[string]interface{}{
				{"role": "user", "content": "投资者互动数据"},
				{"role": "assistant", "reasoning_content": "使用工具查询", "content": content},
			}
		}()

		go func() {
			defer wg.Done()
			var market strings.Builder
			res := NewMarketNewsApi().GetGDP()
			md := util.MarkdownTableWithTitle("国内生产总值(GDP)", res.GDPResult.Data)
			market.WriteString(md)
			res2 := NewMarketNewsApi().GetCPI()
			md2 := util.MarkdownTableWithTitle("居民消费价格指数(CPI)", res2.CPIResult.Data)
			market.WriteString(md2)
			res3 := NewMarketNewsApi().GetPPI()
			md3 := util.MarkdownTableWithTitle("工业品出厂价格指数(PPI)", res3.PPIResult.Data)
			market.WriteString(md3)
			res4 := NewMarketNewsApi().GetPMI()
			md4 := util.MarkdownTableWithTitle("采购经理人指数(PMI)", res4.PMIResult.Data)
			market.WriteString(md4)

			pairCh <- []map[string]interface{}{
				{"role": "user", "content": "国内宏观经济数据"},
				{"role": "assistant", "reasoning_content": "使用工具查询", "content": "\n# 国内宏观经济数据：\n" + market.String()},
			}
		}()

		go func() {
			defer wg.Done()
			md := strings.Builder{}
			res := NewMarketNewsApi().ClsCalendar()
			for _, a := range res {
				bytes, err := json.Marshal(a)
				if err != nil {
					continue
				}
				date := gjson.Get(string(bytes), "calendar_day")
				md.WriteString("\n### 事件/会议日期：" + date.String())
				list := gjson.Get(string(bytes), "items")
				list.ForEach(func(key, value gjson.Result) bool {
					md.WriteString("\n- " + gjson.Get(value.String(), "title").String())
					return true
				})
			}
			pairCh <- []map[string]interface{}{
				{"role": "user", "content": "近期重大事件/会议"},
				{"role": "assistant", "reasoning_content": "使用工具查询", "content": "近期重大事件/会议如下：\n" + md.String()},
			}
		}()

		go func() {
			defer wg.Done()
			//logger.SugaredLogger.Infof("NewChatStream getKLineData stock:%s stockCode:%s", stock, stockCode)
			if strutil.HasPrefixAny(stockCode, []string{"sz", "sh", "hk", "us", "gb_"}) {
				K := &[]KLineData{}
				if strutil.HasPrefixAny(stockCode, []string{"sz", "sh"}) {
					K = NewStockDataApi().GetKLineData(stockCode, "240", o.KDays)
				}
				if strutil.HasPrefixAny(stockCode, []string{"hk", "us", "gb_"}) {
					K = NewStockDataApi().GetHK_KLineData(stockCode, "day", o.KDays)
				}
				Kmap := &[]map[string]any{}
				for _, kline := range *K {
					mapk := make(map[string]any, 6)
					mapk["日期"] = kline.Day
					mapk["开盘价"] = kline.Open
					mapk["最高价"] = kline.High
					mapk["最低价"] = kline.Low
					mapk["收盘价"] = kline.Close
					Volume, _ := convertor.ToFloat(kline.Volume)
					mapk["成交量(万手)"] = Volume / 10000.00 / 100.00
					*Kmap = append(*Kmap, mapk)
				}
				jsonData, _ := json.Marshal(Kmap)
				markdownTable, _ := JSONToMarkdownTable(jsonData)
				content := "## " + stock + "日K数据如下：\n" + markdownTable
				// 后端本地算好的技术指标，避免 AI 凭记忆心算（提升准确度）。
				if ind := buildKLineIndicatorSummary(*K); ind != "" {
					content += ind
				}
				pairCh <- []map[string]interface{}{
					{"role": "user", "content": stock + "日K数据"},
					{"role": "assistant", "content": content},
				}
			}
		}()

		go func() {
			defer wg.Done()
			messages := SearchStockPriceInfo(stock, stockCode, o.CrawlTimeOut)
			if messages == nil || len(*messages) == 0 {
				ch <- map[string]any{
					"code":         1,
					"question":     question,
					"extraContent": "***❗获取股票价格失败,分析结果可能不准确***<hr>",
				}
				logger.SugaredLogger.Warn("获取股票价格失败,分析结果可能不准确")
				return
			}
			price := ""
			for _, message := range *messages {
				price += message + ";"
			}
			pairCh <- []map[string]interface{}{
				{"role": "user", "content": stock + "股价数据"},
				{"role": "assistant", "content": "\n## " + stock + "股价数据：\n" + price},
			}
		}()

		go func() {
			defer wg.Done()
			if tools != nil && len(tools) > 0 {
				return
			}
			if checkIsIndexBasic(stock) {
				return
			}
			messages := GetFinancialReportsByXUEQIU(stockCode, o.CrawlTimeOut)
			if messages == nil || len(*messages) == 0 {
				ch <- map[string]any{
					"code":         1,
					"question":     question,
					"extraContent": "***❗获取股票财报失败,分析结果可能不准确***<hr>",
				}
				logger.SugaredLogger.Warn("获取股票财报失败,分析结果可能不准确")
				return
			}
			pairs := []map[string]interface{}{
				{"role": "user", "content": stock + "财报数据"},
			}
			for _, message := range *messages {
				pairs = append(pairs, map[string]interface{}{
					"role":    "assistant",
					"content": stock + message,
				})
			}
			pairCh <- pairs
		}()

		go func() {
			defer wg.Done()
			messages := NewMarketNewsApi().GetNews24HoursList("", random.RandInt(80, 200))
			if messages == nil || len(*messages) == 0 {
				return
			}
			var messageText strings.Builder
			for _, telegraph := range *messages {
				messageText.WriteString("## " + telegraph.Time + ":" + "\n")
				messageText.WriteString("### " + telegraph.Content + "\n")
			}
			pairCh <- []map[string]interface{}{
				{"role": "user", "content": "市场资讯"},
				{"role": "assistant", "content": messageText.String()},
			}
		}()

		go func() {
			defer wg.Done()
			messages := SearchStockInfo(stock, "telegram", o.CrawlTimeOut)
			if messages == nil || len(*messages) == 0 {
				return
			}
			var newsText strings.Builder
			for _, message := range *messages {
				newsText.WriteString(message + "\n")
			}
			pairCh <- []map[string]interface{}{
				{"role": "user", "content": stock + "相关新闻资讯"},
				{"role": "assistant", "content": newsText.String()},
			}
		}()

		// 等待全部抓取完成，再关闭 pairCh，随后串行收齐（无并发写 msg）。
		go func() {
			wg.Wait()
			close(pairCh)
		}()

		for pairs := range pairCh {
			msg = append(msg, pairs...)
		}

		msg = append(msg, map[string]interface{}{
			"role":    "user",
			"content": question,
		})

		if tools != nil && len(tools) > 0 {
			AskAiWithTools(o, errors.New(""), msg, ch, question, tools, thinking)
		} else {
			AskAi(o, errors.New(""), msg, ch, question, thinking)
		}
	}()
	return ch
}
