package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	neturl "net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CaptchaConfig struct {
	RenderTo    string `json:"renderTo"`
	CustomImage string `json:"customImage"`
	Nctokenstr  string `json:"NCTOKENSTR"`
	Action      string `json:"action"`
	Host        string `json:"HOST"`
	Path        string `json:"PATH"`
	Formactioin string `json:"FORMACTIOIN"`
	Bxstep      string `json:"BXSTEP"`
	Secdata     string `json:"SECDATA"`
	Ncappkey    string `json:"NCAPPKEY"`
	IsUpgrade   string `json:"isUpgrade"`
	CrossSite   string `json:"crossSite"`
	Qrcode      string `json:"qrcode"`
	Pp          struct {
		Enc string `json:"enc"`
		F   string `json:"f"`
		I   int    `json:"i"`
		K   string `json:"k"`
		L   int    `json:"l"`
		Mt  int    `json:"mt"`
		Q   string `json:"q"`
		T   string `json:"t"`
	} `json:"pp"`
	CaptchaConfigInfo string `json:"captchaConfigInfo"`
}

type FruitConfig struct {
	EncryptToken string `json:"encryptToken"`
	ImageData    string `json:"imageData"`
	Ques         string `json:"ques"`
}

type Data struct {
	CaptchaConfig CaptchaConfig `json:"captcha_config"`
	Key           string        `json:"key"`
	EncryptToken  string        `json:"encryptToken"`
	ImageData     string        `json:"imageData"`
	Ques          string        `json:"ques"`
}

type ResponseData struct {
	UrlEncode string `json:"urlEncode"`
	BxEt      string `json:"bx_et"`
	BxPp      string `json:"bx_pp"`
	Referer   string `json:"referer"`
	Count     int    `json:"count"`
}

type Cookie struct {
	Name         string `json:"name"`
	Value        string `json:"value"`
	Domain       string `json:"domain"`
	Path         string `json:"path"`
	Expires      int64  `json:"expires"`
	Size         int    `json:"size"`
	HttpOnly     bool   `json:"httpOnly"`
	Secure       bool   `json:"secure"`
	Session      bool   `json:"session"`
	Priority     string `json:"priority"`
	SourceScheme string `json:"sourceScheme"`
	SourcePort   int    `json:"sourcePort"`
}

var logger *log.Logger

func init() {
	logFile, err := os.OpenFile("fcklog.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}
	logger = log.New(logFile, "", log.LstdFlags)
}

func getCaptchaConfig(wconfigJson string) (CaptchaConfig, error) {
	var config CaptchaConfig
	err := json.Unmarshal([]byte(wconfigJson), &config)
	if err != nil {
		return CaptchaConfig{}, err
	}
	logger.Printf("获取验证码配置信息config：%+v\n", config)
	return config, nil
}

func getFruitConfig(punishURL string, config CaptchaConfig) (FruitConfig, error) {
	var url string
	if config.Action == "captchacapslidev2" || config.Action == "captchascene" {
		url = punishURL + "/_____tmd_____/newslidecaptcha?"
	} else if config.Action == "captchacappuzzle" {
		url = punishURL + "/_____tmd_____/puzzleCaptchaGet?"
	}

	params := neturl.Values{}
	params.Add("token", config.Nctokenstr)
	params.Add("appKey", config.Ncappkey)
	params.Add("x5secdata", config.Secdata)
	params.Add("language", "cn")
	params.Add("v", fmt.Sprintf("%v", randomString()))

	client := &http.Client{}
	req, err := http.NewRequest("GET", url+params.Encode(), nil)
	if err != nil {
		return FruitConfig{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return FruitConfig{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return FruitConfig{}, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return FruitConfig{}, err
	}

	data := result["data"].(map[string]interface{})
	return FruitConfig{
		EncryptToken: data["encryptToken"].(string),
		ImageData:    data["imageData"].(string),
		Ques:         data["ques"].(string),
	}, nil
}

func postConfigToServer(wconfigJson string) (string, error) {
	captchaConfig, err := getCaptchaConfig(wconfigJson)
	if err != nil {
		return "", err
	}

	captchaURL := "https://" + captchaConfig.Host + captchaConfig.Path
	captchaURL = strings.Replace(captchaURL, ":443", "", -1)
	fruitConfig, err := getFruitConfig(captchaURL, captchaConfig)
	if err != nil {
		logger.Println(err)
		return "", err
	}

	data := Data{
		CaptchaConfig: captchaConfig,
		Key:           "UserTest",
		EncryptToken:  fruitConfig.EncryptToken,
		ImageData:     fruitConfig.ImageData,
		Ques:          fruitConfig.Ques,
	}

	client := &http.Client{}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "http://121.36.8.43:3200/tb_n", strings.NewReader(string(jsonData)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		logger.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	logger.Println("API接口返回：", string(body))
	var responseData ResponseData
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		logger.Println(err)
		return "", err
	}
	if len(responseData.UrlEncode) == 0 || len(responseData.BxEt) == 0 || len(responseData.BxPp) == 0 {
		return "", fmt.Errorf("ck接口返回异常")
	}

	headers := map[string]string{
		// "authority":          "login.taobao.com",
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "zh-CN,zh;q=0.9",
		"bx-v":               "2.5.3",
		"cache-control":      "no-cache",
		"content-type":       "application/x-www-form-urlencoded",
		"pragma":             "no-cache",
		"sec-ch-ua":          `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": `"Windows"`,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
	headers["bx_et"] = responseData.BxEt
	headers["bx-pp"] = responseData.BxPp
	headers["referer"] = responseData.Referer

	time.Sleep(time.Duration(rand.Float64()*0.5+1) * time.Second)

	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err = http.NewRequest("GET", responseData.UrlEncode, nil)
	if err != nil {
		return "", err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err = client.Do(req)
	if err != nil {
		logger.Println(err)
		return "", err
	} else {
		cookies := resp.Cookies()
		logger.Println("cookies返回：", cookies)
		for _, cookie := range cookies {
			if cookie.Name == "x5sec" {
				// logger.Println("成功获取 x5sec cookie", cookie.Value)
				// 解析 Max-Age
				maxAgePattern := regexp.MustCompile(`Max-Age=(\d+);`)
				maxAgeMatch := maxAgePattern.FindStringSubmatch(cookie.String())
				if len(maxAgeMatch) < 2 {
					return "", fmt.Errorf("max-Age not found in cookie")
				}
				maxAge, err := strconv.Atoi(maxAgeMatch[1])
				if err != nil {
					return "", fmt.Errorf("invalid Max-Age value")
				}
				// 计算新的 expires 值
				expires := time.Now().Add(time.Duration(maxAge-30) * time.Second).Unix()
				// 组合 JSON
				cookieData := []Cookie{
					{
						Name:   "x5sec",
						Value:  cookie.Value,
						Domain: ".taobao.com",
						// Domain:       ".tmall.com",
						Path:         "/",
						Expires:      expires,
						Size:         len(cookie.Value),
						HttpOnly:     false,
						Secure:       false,
						Session:      false,
						Priority:     "Medium",
						SourceScheme: "Secure",
						SourcePort:   443,
					},
				}
				cookieJson, err := json.Marshal(cookieData)
				if err != nil {
					return "", err
				}
				return string(cookieJson), nil
			}
		}
	}
	return "", fmt.Errorf("no match x5sec cookie")
}

func Wconfig2Ck(jsonString string) (string, error) {
	// 请求打码ck
	ck, err := postConfigToServer(jsonString)
	if err != nil {
		return "", err
	}
	if len(ck) == 0 {
		return "", fmt.Errorf("no Wconfig2Ck found")
	}
	return ck, nil
}

func randomString() string {
	/// 定义字符集 08806067884483191
	digits := "0123456789"
	// 生成15位随机数
	length := 15
	result := make([]byte, length)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	// 将结果转换为字符串
	randomString := string(result)
	return "08" + randomString
}
