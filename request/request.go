package request

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// Request Request
func Request(method string, url string, headers http.Header, params map[string]string, timeout int) (*(http.Response), string, error) {

	// 超时时间:
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	var body *(bytes.Buffer)
	switch method {
	case "GET":
		body = new(bytes.Buffer) // 这里相当于nil?
	case "POST":
		paramsStr := func() string {
			pstr := ""
			for key, value := range params {
				pstr += key + "=" + value + "&"
			}
			if len(pstr) > 0 {
				pstr = pstr[:len(pstr)-1]
			}
			return pstr
		}() // 使用url 参数方式POST
		// jsonStr, _ := json.Marshal(params) 如果要用map方式的POST
		body = bytes.NewBuffer([]byte(paramsStr))
	}
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return new(http.Response), "", err
	}
	req.Header = headers
	resp, err := client.Do(req)
	if err != nil {
		return new(http.Response), "", err
	}
	// resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	defer resp.Body.Close()
	// var respHeader map[string][]string = resp.Header
	content, _ := ioutil.ReadAll(resp.Body)
	content = bytes.TrimPrefix(content, []byte("\xef\xbb\xbf")) // Or []byte{239, 187, 191} // 为了防止奇怪的无法json.unma错误
	return resp, string(content), nil
}

// 上传文件
// url                请求地址
// params        post form里数据
// nameField  请求地址上传文件对应field
// fileName     文件名
// file               文件

// UploadFile 抄来的Upload
func UploadFile(url string, headers http.Header, params map[string]string, nameField, fileName string, filepath string, timeout int) (resp *(http.Response), content string, err error) {

	HTTPClient := &http.Client{Timeout: time.Duration(timeout) * time.Second}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	formFile, err := writer.CreateFormFile(nameField, fileName)
	if err != nil {
		return nil, "", err
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return nil, "", err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, "", err
	}
	//req.Header.Set("Content-Type","multipart/form-data")
	req.Header = headers
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err = HTTPClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	contentByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	content = string(contentByte)
	return resp, content, nil
}
