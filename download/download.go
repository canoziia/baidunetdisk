package download

import (
	cfg "baidunetdisk/config"
	"baidunetdisk/request"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/robertkrimen/otto"
)

var (
	// Config 配置文件总成
	Config cfg.Config
	// Appid 百度网盘appid
	Appid string
	// Channel 百度网盘channel
	Channel string
	// BDUSS 百度网盘的BDUSS(SVIP)
	BDUSS string
	// STOKEN 百度网盘的STOKEN(SVIP)
	STOKEN string
	// MyBDUSS 百度网盘的BDUSS(自己的)
	MyBDUSS string
	// MySTOKEN 百度网盘的STOKEN(自己的)
	MySTOKEN string
	// Timeout 请求的Timeout,单位为s
	Timeout int
)

// Init 初始化
func Init(configPath string) {
	Config = cfg.GetConfig(configPath)
	Appid = Config["baidunetdisk"]["appid"]
	Channel = Config["baidunetdisk"]["channel"]
	BDUSS = Config["baidunetdisk"]["BDUSS"]
	STOKEN = Config["baidunetdisk"]["STOKEN"]
	MyBDUSS = Config["baidunetdisk"]["MyBDUSS"]
	MySTOKEN = Config["baidunetdisk"]["MySTOKEN"]
	Timeout, _ = strconv.Atoi(Config["request"]["timeout"]) // 这里错了就默认0
}

// VerifyPwd 测试提取码
func VerifyPwd(surl, pwd string) (randsk string, err error) {
	link := fmt.Sprintf("https://pan.baidu.com/share/verify?channel=%s&clienttype=0&web=1&app_id=%s&surl=%s", Channel, Appid, surl)
	//url := "https://pan.baidu.com/share/verify?channel="+Channel+"&clienttype=0&web=1&app_id=250528&surl=" + surl
	headers := make(http.Header, 0)
	headers.Set("User-Agent", "netdisk")
	headers.Set("Referer", "https://pan.baidu.com/disk/home")
	params := map[string]string{"pwd": pwd}
	_, content, err := request.Request("POST", link, headers, params, Timeout)
	if err != nil {
		return "", err
	}
	var contentMap map[string]interface{}
	err = json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		return "", err
	}
	code := contentMap["errno"].(float64)
	if code != 0 {
		return "", errors.New("百度返回码错误")
	}
	randsk = contentMap["randsk"].(string)
	return randsk, nil
}

// GetSignFromShare 获取Sign(我也不知道是什么)
func GetSignFromShare(surl, randsk string) (sign map[string]interface{}, err error) {
	link := fmt.Sprintf("https://pan.baidu.com/s/1%s", surl)
	cookie := fmt.Sprintf("BDUSS=%s; STOKEN=%s; BDCLND=%s", BDUSS, STOKEN, randsk)
	headers := make(http.Header, 0)
	headers.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36`)
	headers.Set("Cookie", cookie)
	_, content, err := request.Request("GET", link, headers, make(map[string]string, 0), Timeout)
	if err != nil {
		return sign, err
	}
	reg, _ := regexp.Compile(`yunData.setData\(({.*?)\);`)
	resultSlice := reg.FindStringSubmatch(content)
	if len(resultSlice) == 0 {
		return sign, errors.New("匹配不到sign")
	}
	signStr := resultSlice[1]
	err = json.Unmarshal([]byte(signStr), &sign)
	if err != nil {
		return sign, err
	}
	return sign, err
}

// GetFileListFromShare 从shareid获取文件信息
func GetFileListFromShare(shareid, uk, randsk string) (fileList map[string]interface{}, err error) {
	link := fmt.Sprintf("https://pan.baidu.com/share/list?app_id=250528&channel=chunlei&clienttype=0&desc=0&num=100&order=name&page=1&root=1&shareid=%s&showempty=0&uk=%s&web=1", shareid, uk)
	cookie := fmt.Sprintf("BDUSS=%s; STOKEN=%s; BDCLND=%s", BDUSS, STOKEN, randsk)
	headers := make(http.Header, 0)
	headers.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36`)
	headers.Set("Cookie", cookie)
	_, content, err := request.Request("GET", link, headers, make(map[string]string, 0), Timeout)
	if err != nil {
		return fileList, err
	}
	err = json.Unmarshal([]byte(content), &fileList)
	if err != nil {
		return fileList, err
	}
	code := fileList["errno"].(float64)
	if code != 0 {
		return fileList, errors.New("百度返回码错误")
	}
	return fileList, nil
}

// GetParamsFromShare 获取全部参数
func GetParamsFromShare(surl, pwd string) (timestamp, sign, randsk, shareid, uk string, err error) {
	randsk, err = VerifyPwd(surl, pwd)
	if err != nil {
		return "", "", "", "", "", err
	}

	signMap, err := GetSignFromShare(surl, randsk)
	if err != nil {
		return "", "", "", "", "", err
	}

	sign = signMap["sign"].(string)
	timestamp = strconv.FormatFloat(signMap["timestamp"].(float64), 'f', -1, 64)
	shareid = strconv.FormatFloat(signMap["shareid"].(float64), 'f', -1, 64)
	uk = strconv.FormatFloat(signMap["uk"].(float64), 'f', -1, 64)
	return timestamp, sign, randsk, shareid, uk, nil
}

// GetDownloadLink 从文件信息获取初始下载链接
func GetDownloadLink(fsid, timestamp, sign, randsk, shareid, uk string) (dlink string, err error) {
	link := fmt.Sprintf("https://pan.baidu.com/api/sharedownload?app_id=%s&channel=%s&clienttype=12&sign=%s&timestamp=%s&web=1", Appid, Channel, sign, timestamp)
	cookie := fmt.Sprintf("BDUSS=%s; STOKEN=%s; BDCLND=%s", BDUSS, STOKEN, randsk)
	headers := make(http.Header, 0)
	headers.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36`)
	headers.Set("Cookie", cookie)
	params := map[string]string{
		"encrypt":   "0",
		"extra":     `{"sekey":"` + DecodeURIComponent(randsk) + `"}`,
		"fid_list":  `[` + fsid + `]`,
		"primaryid": shareid,
		"uk":        uk,
		"product":   "share",
		"type":      "nolimit",
	}
	_, content, err := request.Request("POST", link, headers, params, Timeout)
	if err != nil {
		return "", err
	}
	rawlinkMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(content), &rawlinkMap)
	if err != nil {
		return "", err
	}
	code := rawlinkMap["errno"].(float64)
	if code != 0 {
		return "", errors.New("百度返回码错误")
	}

	if tmpdlink := rawlinkMap["list"].([]interface{})[0].(map[string]interface{})["dlink"]; tmpdlink == nil {
		dlink, err = "", errors.New("是文件夹")
	} else {
		dlink, err = tmpdlink.(string), nil
	}
	return dlink, err
}

// GetRealLink 从初始链接获取真实链接
func GetRealLink(downloadLink string) (realLink string, err error) {
	noRedirectRequest := func(link string, headers http.Header, params map[string]string, timeout int) (*(http.Response), string, error) {

		// 超时时间:
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Timeout: time.Duration(timeout) * time.Second,
		}

		// var body *(bytes.Buffer)
		body := new(bytes.Buffer) // 这里相当于nil?
		req, err := http.NewRequest("GET", link, body)

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
		return resp, string(content), nil
	}
	cookie := fmt.Sprintf("BDUSS=%s;", BDUSS)
	headers := make(http.Header, 0)
	headers.Set("User-Agent", `LogStatistic`)
	headers.Set("Cookie", cookie)
	resp, _, err := noRedirectRequest(downloadLink, headers, make(map[string]string, 0), Timeout)
	if err != nil {
		return "", err
	}
	tmpLink, err := resp.Location()
	if err != nil {
		return "", err
	}
	realLink = tmpLink.String()
	return realLink, nil
}

// DecodeURIComponent URL解码
func DecodeURIComponent(encoded string) (decoded string) {
	tmp := strings.Replace(encoded, "%20", "+", -1) // 有和没有都一样
	decoded, _ = url.QueryUnescape(tmp)
	decoded = strings.Replace(encoded, "+", "%20", -1) // 不这样吧+换成%20的话js理解不了
	return decoded
}

// GetDirFileList 从文件夹fs_id获取下属的文件列表
func GetDirFileList(randsk, uk, shareid, path string) (fileInfoSlice []map[string]string, err error) {
	dirStr := url.QueryEscape(path)
	link := fmt.Sprintf("https://pan.baidu.com/share/list?&order=other&desc=1&showempty=0&web=1&page=1&num=100&channel=%s&app_id=%s&uk=%s&shareid=%s&dir=%s", Channel, Appid, uk, shareid, dirStr)
	cookie := fmt.Sprintf("BDUSS=%s; STOKEN=%s; BDCLND=%s", BDUSS, STOKEN, randsk)
	headers := make(http.Header, 0)
	headers.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36`)
	headers.Set("Cookie", cookie)
	_, content, err := request.Request("GET", link, headers, make(map[string]string, 0), Timeout)
	if err != nil {
		return
	}
	fileListMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(content), &fileListMap)
	if err != nil {
		return
	}

	fileInfoSlice = CUFromFileListMap(fileListMap, 1)
	return fileInfoSlice, nil
}

// GetShareLinkFileList 从分享链接获取第一次文件列表
func GetShareLinkFileList(surl, pwd string) (fileInfoSlice []map[string]string, err error) {

	randsk, err := VerifyPwd(surl, pwd)
	if err != nil {
		return nil, err
	}

	signMap, err := GetSignFromShare(surl, randsk)
	if err != nil {
		return nil, err
	}

	shareid := strconv.FormatFloat(signMap["shareid"].(float64), 'f', -1, 64)
	uk := strconv.FormatFloat(signMap["uk"].(float64), 'f', -1, 64)
	fileListMap, err := GetFileListFromShare(shareid, uk, randsk)
	if err != nil {
		return nil, err
	}
	// fmt.Println(fileListMap)
	fileInfoSlice = CUFromFileListMap(fileListMap, 0)
	// fmt.Println(fileListSlice)
	return fileInfoSlice, nil
}

// GetFileRealLink 从参数获取真实链接
func GetFileRealLink(fsid, timestamp, sign, randsk, shareid, uk string) (realLink string, err error) {
	dlink, err := GetDownloadLink(fsid, timestamp, sign, randsk, shareid, uk) //  (rawlinkMap map[string]interface{}, err error)
	if err != nil {
		return "", err
	}
	// fmt.Println(rawlinkMap)
	realLink, err = GetRealLink(dlink)
	if err != nil {
		return "", err
	}
	return realLink, nil
}

// GetFileListFromParams 从分享链接获取第一次文件列表
func GetFileListFromParams(shareid, uk, randsk string) (fileInfoSlice []map[string]string, err error) {

	fileListMap, err := GetFileListFromShare(shareid, uk, randsk)
	if err != nil {
		return nil, err
	}
	fileInfoSlice = CUFromFileListMap(fileListMap, 0)
	// fmt.Println(fileListSlice)
	return fileInfoSlice, nil
}

// CUFromFileListMap ChooseUsefulInfoFromFileListMap
func CUFromFileListMap(fileListMap map[string]interface{}, caseInt int) (fileInfoSlice []map[string]string) {
	// fmt.Println(fileListMap)
	switch caseInt {
	case 0:
		fileListSlice := fileListMap["list"].([]interface{})
		for _, fileInfo := range fileListSlice {
			oldInfoMap := fileInfo.(map[string]interface{})
			newInfoMap := make(map[string]string, 0)
			newInfoMap["fsid"] = oldInfoMap["fs_id"].(string)
			newInfoMap["isdir"] = oldInfoMap["isdir"].(string)
			newInfoMap["filename"] = oldInfoMap["server_filename"].(string)
			newInfoMap["filesize"] = oldInfoMap["size"].(string) // 这边从shareid获取和后面不一样
			newInfoMap["path"] = oldInfoMap["path"].(string)
			fileInfoSlice = append(fileInfoSlice, newInfoMap)
		}
	case 1:
		fileSlice := fileListMap["list"].([]interface{})
		// fileInfoSlice = make([]map[string]string, 0)
		for _, fileInfo := range fileSlice {
			oldInfoMap := fileInfo.(map[string]interface{})
			newInfoMap := make(map[string]string, 0)
			newInfoMap["fsid"] = strconv.FormatFloat(oldInfoMap["fs_id"].(float64), 'f', -1, 64)
			newInfoMap["isdir"] = strconv.FormatFloat(oldInfoMap["isdir"].(float64), 'f', -1, 64)
			newInfoMap["filename"] = oldInfoMap["server_filename"].(string)
			newInfoMap["filesize"] = strconv.FormatFloat(oldInfoMap["size"].(float64), 'f', -1, 64)
			newInfoMap["path"] = oldInfoMap["path"].(string)
			fileInfoSlice = append(fileInfoSlice, newInfoMap)
		}
	}
	return fileInfoSlice
}

// GetMyParams 获取我的参数
func GetMyParams() (sign string, timestamp string) {
	link := "https://pan.baidu.com/api/gettemplatevariable?fields=[%22sign1%22,%22sign2%22,%22sign3%22,%22timestamp%22]&channel=" + Channel + "&web=1&app_id=" + Appid + "&clienttype=0"
	cookie := fmt.Sprintf("BDUSS=%s; STOKEN=%s;", MyBDUSS, MySTOKEN)
	headers := make(http.Header, 0)
	headers.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36`)
	headers.Set("Cookie", cookie)
	_, content, err := request.Request("GET", link, headers, make(map[string]string, 0), Timeout)
	if err != nil {
		return
	}
	//fmt.Println(resp, content)
	respMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(content), &respMap)
	tmpParamsMap := respMap["result"].(map[string]interface{})
	sign1, sign2, sign3 := tmpParamsMap["sign1"].(string), tmpParamsMap["sign2"].(string), tmpParamsMap["sign3"].(string)
	timestamp = strconv.FormatFloat(tmpParamsMap["timestamp"].(float64), 'f', -1, 64)
	vm := otto.New()
	code := `var sign1 = "` + sign1 + `";` + `var sign3 = "` + sign3 + `";` + sign2 + `;`
	code += `
	var base64EncodeChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/";
	function base64encode(str){
		var out, i, len;
		var c1, c2, c3;
		len = str.length;
		i = 0;
		out = "";
		while (i < len) {
			c1 = str.charCodeAt(i++) & 0xff;
			if (i == len) {
				out += base64EncodeChars.charAt(c1 >> 2);
				out += base64EncodeChars.charAt((c1 & 0x3) << 4);
				out += "==";
				break;
			}
			c2 = str.charCodeAt(i++);
			if (i == len) {
				out += base64EncodeChars.charAt(c1 >> 2);
				out += base64EncodeChars.charAt(((c1 & 0x3) << 4) | ((c2 & 0xF0) >> 4));
				out += base64EncodeChars.charAt((c2 & 0xF) << 2);
				out += "=";
				break;
			}
			c3 = str.charCodeAt(i++);
			out += base64EncodeChars.charAt(c1 >> 2);
			out += base64EncodeChars.charAt(((c1 & 0x3) << 4) | ((c2 & 0xF0) >> 4));
			out += base64EncodeChars.charAt(((c2 & 0xF) << 2) | ((c3 & 0xC0) >> 6));
			out += base64EncodeChars.charAt(c3 & 0x3F);
		}
		return out;
	}`
	vm.Run(code)
	value, _ := vm.Run(`base64encode(s(sign3,sign1))`)
	sign, _ = value.ToString()
	return sign, timestamp
}

// GetMyFileList 获取网盘内文件链接
func GetMyFileList(dirPath string) (fileInfoSlice []map[string]string, err error) {
	dirStr := url.QueryEscape(dirPath)
	link := fmt.Sprintf("https://pan.baidu.com/api/list?order=time&desc=1&showempty=0&web=1&clienttype=0&page=1&num=100&dir=%s&channel=%s&app_id=%s", dirStr, Channel, Appid)
	cookie := fmt.Sprintf("BDUSS=%s; STOKEN=%s;", MyBDUSS, MySTOKEN)
	headers := make(http.Header, 0)
	headers.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36`)
	headers.Set("Cookie", cookie)
	_, content, err := request.Request("GET", link, headers, make(map[string]string, 0), Timeout)
	if err != nil {
		return
	}
	fileListMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(content), &fileListMap)
	if err != nil {
		return
	}
	fileInfoSlice = CUFromFileListMap(fileListMap, 1)
	return fileInfoSlice, nil
}

// GetMyDownloadLink 获取网盘内文件下载链接
func GetMyDownloadLink(fsid, sign, timestamp string) (dlink string, err error) {
	signStr := url.QueryEscape(sign)
	link := "https://pan.baidu.com/api/download?sign=" + signStr + "&timestamp=" + timestamp + "&fidlist=%5B" + fsid + "%5D&type=dlink&vip=0&channel=" + Channel + "&web=1&app_id=" + Appid + "&clienttype=0"
	cookie := fmt.Sprintf("BDUSS=%s; STOKEN=%s;", MyBDUSS, MySTOKEN)
	headers := make(http.Header, 0)
	headers.Set("User-Agent", `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36`)
	headers.Set("Cookie", cookie)
	_, content, err := request.Request("GET", link, headers, make(map[string]string, 0), Timeout)
	if err != nil {
		return
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(content), &respMap)
	if err != nil {
		return
	}
	dlinkMap := respMap["dlink"].([]interface{})[0].(map[string]interface{})
	dlink = dlinkMap["dlink"].(string)
	return
}
