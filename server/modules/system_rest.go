package modules

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/duke-git/lancet/v2/cryptor"
	"go-stock/backend/data"
	"go-stock/backend/machineid"
	"go-stock/server"
)

func init() {
	// GetSponsorInfo：返回启动时缓存的赞助信息；Web 版直接返回当前有效赞助等级。
	server.Register("GetSponsorInfo", func(_ context.Context, _ *server.Core, _ server.Args) (any, error) {
		level, active := data.EffectiveSponsorVipLevel()
		return map[string]any{
			"vipLevel":    level,
			"active":      active,
			"sponsorCode": data.GetSettingConfig().SponsorCode,
		}, nil
	})

	// CheckSponsorCode：AES-ECB 解密校验赞助码并持久化。
	server.Register("CheckSponsorCode", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		sponsorCode := strings.TrimSpace(server.ArgString(args, 0))
		if sponsorCode == "" {
			return map[string]any{"code": 0, "message": "赞助码不能为空,请输入正确的赞助码!"}, nil
		}
		encrypted, err := hex.DecodeString(sponsorCode)
		if err != nil {
			return map[string]any{"code": 0, "msg": "赞助码格式错误,请输入正确的赞助码!"}, nil
		}
		key, err := hex.DecodeString(server.BuildKey)
		if err != nil {
			return map[string]any{"code": 0, "msg": "版本错误，不支持赞助码!"}, nil
		}
		decrypt := cryptor.AesEcbDecrypt(encrypted, key)
		if len(decrypt) == 0 {
			return map[string]any{"code": 0, "msg": "赞助码错误，请输入正确的赞助码!"}, nil
		}
		config := data.GetSettingConfig()
		if config.SponsorCode != sponsorCode {
			config.SponsorCode = sponsorCode
			data.UpdateConfig(config)
		}
		return map[string]any{"code": 1, "msg": "赞助码校验成功，感谢您的支持!"}, nil
	})

	// CheckDeviceBinding：向作者服务端查询设备绑定情况（透传）。
	server.Register("CheckDeviceBinding", func(_ context.Context, _ *server.Core, args server.Args) (any, error) {
		token := server.ArgString(args, 0)
		apiBase := server.ArgString(args, 1)
		result := map[string]any{"bound": false, "deviceCount": 0, "maxDevices": 5}
		if token == "" || apiBase == "" {
			return result, nil
		}
		uuid := machineid.GetMachineId()
		resp, err := data.SharedHTTPClient.R().SetHeader("Authorization", "Bearer "+token).
			Get(fmt.Sprintf("%s/user/device-check?uuid=%s", apiBase, uuid))
		if err != nil {
			return result, nil
		}
		var respData struct {
			Code int `json:"code"`
			Data struct {
				Bound       bool `json:"bound"`
				DeviceCount int  `json:"deviceCount"`
				MaxDevices  int  `json:"maxDevices"`
			} `json:"data"`
		}
		if err := json.Unmarshal(resp.Body(), &respData); err != nil || respData.Code != 0 {
			return result, nil
		}
		result["bound"] = respData.Data.Bound
		result["deviceCount"] = respData.Data.DeviceCount
		result["maxDevices"] = respData.Data.MaxDevices
		return result, nil
	})

}
