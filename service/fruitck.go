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
	"strings"
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

	return string(body), nil
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
