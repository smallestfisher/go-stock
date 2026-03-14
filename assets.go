package main

import "embed"

//go:embed build/appicon.png
var icon []byte

//go:embed build/app.ico
var icon2 []byte

//go:embed build/screenshot/alipay.jpg
var alipay []byte

//go:embed build/screenshot/wxpay.jpg
var wxpay []byte

//go:embed build/screenshot/扫码_搜索联合传播样式-白色版.png
var wxgzh []byte

//go:embed build/stock_basic.json
var stocksBin []byte

//go:embed build/stock_base_info_hk.json
var stocksBinHK []byte

//go:embed build/stock_base_info_us.json
var stocksBinUS []byte

//go:embed frontend/dist/*
var webAssets embed.FS
