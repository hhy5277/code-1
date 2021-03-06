package agents

import (
	"bytes"
	"fmt"
	"github.com/iwind/TeaGo/maps"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// WebHook
type WebHookSource struct {
	URL        string           `yaml:"url" json:"url"`
	Timeout    string           `yaml:"timeout" json:"timeout"`
	Method     string           `yaml:"method" json:"method"`         // 请求方法
	DataFormat SourceDataFormat `yaml:"dataFormat" json:"dataFormat"` // 数据格式

	timeoutDuration time.Duration
}

// 获取新对象
func NewWebHookSource() *WebHookSource {
	return &WebHookSource{}
}

// 校验
func (this *WebHookSource) Validate() error {
	this.timeoutDuration, _ = time.ParseDuration(this.Timeout)
	if len(this.Method) == 0 {
		this.Method = http.MethodPost
	} else {
		this.Method = strings.ToUpper(this.Method)
	}

	if len(this.URL) == 0 {
		return errors.New("url should not be empty")
	}

	return nil
}

// 名称
func (this *WebHookSource) Name() string {
	return "WebHook"
}

// 代号
func (this *WebHookSource) Code() string {
	return "webhook"
}

// 描述
func (this *WebHookSource) Description() string {
	return "通过HTTP或者HTTPS接口获取数据"
}

// 数据格式
func (this *WebHookSource) DataFormatCode() SourceDataFormat {
	return this.DataFormat
}

// 执行
func (this *WebHookSource) Execute(params map[string]string) (value interface{}, err error) {
	if this.timeoutDuration.Seconds() <= 0 {
		this.timeoutDuration = 10 * time.Second
	}

	client := http.Client{
		Timeout: this.timeoutDuration,
	}

	query := url.Values{}
	for name, value := range params {
		query[name] = []string{value}
	}
	rawQuery := query.Encode()

	urlString := this.URL
	var body io.Reader = nil
	if len(rawQuery) > 0 {
		if this.Method == http.MethodGet {
			if strings.Index(this.URL, "?") > 0 {
				urlString += "&" + rawQuery
			} else {
				urlString += "?" + rawQuery
			}
		} else {
			body = bytes.NewReader([]byte(rawQuery))
		}
	}

	req, err := http.NewRequest(this.Method, urlString, body)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response status code should be 200, now is " + fmt.Sprintf("%d", resp.StatusCode))
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return DecodeSource(respBytes, this.DataFormat)
}

// 获取简要信息
func (this *WebHookSource) Summary() maps.Map {
	return maps.Map{
		"name":        this.Name(),
		"code":        this.Code(),
		"description": this.Description(),
	}
}
