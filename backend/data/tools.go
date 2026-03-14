package data

// @Author spark
// @Date 2026/3/7 18:48
// @Desc AI工具定义，精简描述以提高大模型兼容性
//-----------------------------------------------------------------------------------

func Tools(tools []Tool) []Tool {
	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name: "SearchStockByIndicators",
			Description: "根据自然语言筛选股票。支持K线形态、财务指标、技术指标（MACD, RSI, KDJ等）、资金流向及板块排名等综合选股条件。" +
				"例如：'10点半之前涨停'、'连续2日主力净流入'、'MACD金叉，非ST'。",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"words": map[string]any{
						"type":        "string",
						"description": "选股自然语言描述",
					},
				},
				Required: []string{"words"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name: "SearchBk",
			Description: "查询板块、概念或指数的整体行情数据。" +
				"例如：'今日涨幅前15的概念板块'、'存储芯片板块成分股'、'上证指数'。",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"words": map[string]any{
						"type":        "string",
						"description": "板块/概念/指数查询自然语言",
					},
				},
				Required: []string{"words"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name: "SearchETF",
			Description: "查询ETF数据。" +
				"例如：'机器人ETF'、'溢价率介于0%-10%的ETF'、'今日涨幅前50的ETF'。",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"words": map[string]any{
						"type":        "string",
						"description": "ETF查询自然语言",
					},
				},
				Required: []string{"words"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "GetStockKLine",
			Description: "获取股票日K线数据。",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"days": map[string]any{
						"type":        "string",
						"description": "日K数据条数",
					},
					"stockCode": map[string]any{
						"type":        "string",
						"description": "股票代码（A股：sh,sz开头;港股hk开头,美股：us开头）",
					},
				},
				Required: []string{"days", "stockCode"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "InteractiveAnswer",
			Description: "获取投资者与上市公司互动问答的数据,反映当前投资者关注的热点问题",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"page": map[string]any{
						"type":        "string",
						"description": "分页号",
					},
					"pageSize": map[string]any{
						"type":        "string",
						"description": "分页大小",
					},
					"keyWord": map[string]any{
						"type":        "string",
						"description": "搜索关键词（可输入股票名称或者当前热门板块/行业/概念/标的/事件等）",
					},
				},
				Required: []string{"page", "pageSize"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "GetStockResearchReport",
			Description: "获取市场分析师的股票研究报告",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"stockCode": map[string]any{
						"type":        "string",
						"description": "股票代码",
					},
				},
				Required: []string{"stockCode"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "HotStrategyTable",
			Description: "获取当前热门选股策略",
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "HotStockTable",
			Description: "当前热门股票排名",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"pageSize": map[string]any{
						"type":        "string",
						"description": "分页大小",
					},
				},
				Required: []string{"pageSize"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "GetStockMoneyData",
			Description: "今日股票资金流入排名",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"pageSize": map[string]any{
						"type":        "string",
						"description": "分页大小",
					},
				},
				Required: []string{"pageSize"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "GetStockConceptInfo",
			Description: "获取股票所属概念详细信息",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"code": map[string]any{
						"type":        "string",
						"description": "股票代码",
					},
				},
				Required: []string{"code"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "GetStockFinancialInfo",
			Description: "获取股票财务报表信息",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"stockCode": map[string]any{
						"type":        "string",
						"description": "股票代码",
					},
				},
				Required: []string{"stockCode"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "GetStockHolderNum",
			Description: "获取股票股东人数信息",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"stockCode": map[string]any{
						"type":        "string",
						"description": "股票代码",
					},
				},
				Required: []string{"stockCode"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "GetStockHistoryMoneyData",
			Description: "获取股票历史资金流向数据",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"stockCode": map[string]any{
						"type":        "string",
						"description": "股票代码",
					},
				},
				Required: []string{"stockCode"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "GetStockRZRQInfo",
			Description: "获取股票融资融券信息",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"stockCode": map[string]any{
						"type":        "string",
						"description": "股票代码",
					},
				},
				Required: []string{"stockCode"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "GetIndustryValuation",
			Description: "获取行业/板块平均估值和中值",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"bkName": map[string]any{
						"type":        "string",
						"description": "行业/板块名称",
					},
				},
				Required: []string{"bkName"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "GetNewsListData",
			Description: "获取新闻资讯",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"keyWord": map[string]any{
						"type":        "string",
						"description": "搜索关键词",
					},
					"startTime": map[string]any{
						"type":        "string",
						"description": "开始时间",
					},
					"limit": map[string]any{
						"type":        "string",
						"description": "返回条数",
					},
				},
				Required: []string{"startTime", "limit"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "CreateAiRecommendStocks",
			Description: "保存单条AI推荐股票记录",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"modelName":       map[string]any{"type": "string"},
					"stockCode":       map[string]any{"type": "string"},
					"stockName":       map[string]any{"type": "string"},
					"bkName":          map[string]any{"type": "string"},
					"stockPrice":      map[string]any{"type": "string"},
					"recommendReason": map[string]any{"type": "string"},
				},
				Required: []string{"stockCode", "stockName", "bkName"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "BatchCreateAiRecommendStocks",
			Description: "批量保存AI推荐股票记录",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"stocks": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"stockCode": map[string]any{"type": "string"},
								"stockName": map[string]any{"type": "string"},
								"bkName":    map[string]any{"type": "string"},
							},
							"required": []string{"stockCode", "stockName", "bkName"},
						},
					},
				},
				Required: []string{"stocks"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "AiRecommendStocks",
			Description: "获取近期AI分析记录列表",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"startDate": map[string]any{"type": "string"},
					"endDate":   map[string]any{"type": "string"},
					"page":      map[string]any{"type": "string"},
					"pageSize":  map[string]any{"type": "string"},
				},
				Required: []string{"startDate", "endDate", "page", "pageSize"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "StockNotice",
			Description: "获取上市公司公告列表",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"stock_list": map[string]any{"type": "string", "description": "股票代码列表"},
				},
				Required: []string{"stock_list"},
			},
		},
	})

	tools = append(tools, Tool{
		Type: "function",
		Function: ToolFunction{
			Name:        "GetSecuritiesCompanyOpinion",
			Description: "获取券商机构观点",
			Parameters: &FunctionParameters{
				Type: "object",
				Properties: map[string]any{
					"startDate": map[string]any{"type": "string"},
					"endDate":   map[string]any{"type": "string"},
				},
				Required: []string{"startDate", "endDate"},
			},
		},
	})

	return tools
}
