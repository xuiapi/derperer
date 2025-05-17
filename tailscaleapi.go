package derperer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// 顶层配置
type Config struct {
	ACLs    []ACL   `json:"acls"`
	SSH     []SSH   `json:"ssh"`
	DerpMap DerpMap `json:"derpMap"`
}

// ACL 规则
type ACL struct {
	Action string   `json:"action"`
	Src    []string `json:"src"`
	Dst    []string `json:"dst"`
}

// SSH 规则
type SSH struct {
	Action string   `json:"action"`
	Src    []string `json:"src"`
	Dst    []string `json:"dst"`
	Users  []string `json:"users"`
}

// DERP 映射
type DerpMap struct {
	OmitDefaultRegions bool                 `json:"OmitDefaultRegions"`
	Regions            map[int]*DERPRegionR `json:"Regions"`
}

func generateACLConfig(regions map[int]*DERPRegionR) Config {
	// 直接使用变量赋值创建 Config 实例
	cfg := Config{
		ACLs: []ACL{
			{
				Action: "accept",
				Src:    []string{"*"},
				Dst:    []string{"*:*"},
			},
		},

		SSH: []SSH{
			{
				Action: "check",
				Src:    []string{"autogroup:member"},
				Dst:    []string{"autogroup:self"},
				Users:  []string{"autogroup:nonroot", "root"},
			},
		},

		DerpMap: DerpMap{
			OmitDefaultRegions: true,
			Regions:            regions,
		},
	}
	return cfg
}

func UpdateACL(regions map[int]*DERPRegionR, account string, api_key string) []byte {
	url := "https://api.tailscale.com/api/v2/tailnet/" + account + "/acl"

	cfg := generateACLConfig(regions)

	// 将结构体编码为 JSON 并输出
	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		fmt.Println("JSON Marshal 失败: %v", err)
	}

	fmt.Println(string(out))

	payload := strings.NewReader(string(out))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Authorization", "Bearer "+api_key)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	return body
}
