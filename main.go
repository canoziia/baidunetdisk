package main

import (
	cfg "baidunetdisk/config"
	dl "baidunetdisk/download"
	page "baidunetdisk/page"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var (
	configPath string
	config     cfg.Config
	root       string
	port       string
)

func init() {
	fmt.Println("初始化程序---------------------")
	flag.StringVar(&configPath, "c", "conf.ini", "配置文件路径")
	flag.Parse()
	dl.Init(configPath)
	page.Init(configPath)
	config = cfg.GetConfig(configPath)
	root = config["httpserver"]["root"]
	port = config["httpserver"]["port"]
}
func main() {
	http.HandleFunc(root, netdiskpageV1)
	// http.HandleFunc("/baidu/netdisk/v1/file", netdiskdirpageV1)
	http.HandleFunc(root+"download", netdiskpageDownloadV1)
	http.HandleFunc(root+"help/", netdiskpageHelpV1)
	if port == "" {
		port = "80"
	}
	fmt.Println("开始监听服务: 127.0.0.1:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func netdiskpageV1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Index--------------------")
	fmt.Println("from: ", GetIP(r))

	// 检查是否为post请求
	if r.Method != http.MethodPost {
		fmt.Println("Method: GET")
		w.Write([]byte(page.Landing))
		return
	}
	fmt.Println("Method: POST")
	// dl.Init(configPath)
	postMap, err := postToMap(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	var (
		fileInfoSlice []map[string]string
		timestamp     string
		sign          string
		randsk        string
		shareid       string
		uk            string
		path          string
	)
	switch len(postMap) {
	case 2:
		surl := postMap["surl"]
		pwd := postMap["pwd"]
		timestamp, sign, randsk, shareid, uk, err = dl.GetParamsFromShare(surl, pwd)
		if err != nil {
			// w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(page.Error + `
				<div class="alert alert-danger" role="alert">
		  		<h5 class="alert-heading">提示</h5>
		  		<hr>
		  		<p class="card-text">提取码错误或文件失效,也可能是暂时性错误</p>
		  		</div>` + page.Errordiv))
			fmt.Println(err)
			return
		}
		fileInfoSlice, err = dl.GetFileListFromParams(shareid, uk, randsk)
		if err != nil {
			// w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(page.Error + `
				<div class="alert alert-danger" role="alert">
		  		<h5 class="alert-heading">提示</h5>
		  		<hr>
		  		<p class="card-text">提取码错误或文件失效,也可能是暂时性错误</p>
		  		</div>` + page.Errordiv))
			fmt.Println(err)
			return
		}
	case 6:
		randsk, uk, shareid, path, timestamp, sign = postMap["randsk"], postMap["uk"], postMap["shareid"], postMap["path"], postMap["timestamp"], postMap["sign"]
		fileInfoSlice, err = dl.GetDirFileList(randsk, uk, shareid, path)
		if err != nil {
			// w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(page.Error + `
				<div class="alert alert-danger" role="alert">
		  		<h5 class="alert-heading">提示</h5>
		  		<hr>
		  		<p class="card-text">可能是暂时性错误,请重试</p>
		  		</div>` + page.Errordiv))
			fmt.Println(err)
			return
		}
	}

	filecontent := ""
	for _, fileInfo := range fileInfoSlice {
		switch fileInfo["isdir"] {
		case "0":
			filecontent += `<li class="list-group-item border-muted rounded text-muted py-2">
			<i class="far fa-file mr-2"></i>
			<a href="javascript:void(0)" onclick="dl('` + fileInfo["fsid"] + `',` + timestamp + `,'` + sign + `','` + randsk + `','` + shareid + `','` + uk + `')">` + fileInfo["filename"] + `</a>
			<span class="float-right">` + fileInfo["filesize"] + `</span>
			</li>`
		default:
			filecontent += `<li class="list-group-item border-muted rounded text-muted py-2">
			<i class="far fa-folder mr-2"></i>
			<a href="javascript:void(0)" onclick="getdirfilelist('` + randsk + `',` + uk + `,'` + shareid + `','` + fileInfo["path"] + `','` + timestamp + `','` + sign + `')">` + fileInfo["filename"] + `</a>
			<span class="float-right">` + fileInfo["filesize"] + `</span>
			</li>`
		}

	}
	w.Write([]byte(page.Filebody + filecontent + page.Filefoot))
	fmt.Println("success")
}

func netdiskpageDownloadV1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Download-----------------")
	fmt.Println("from: ", GetIP(r))

	// 检查是否为post请求
	if r.Method != http.MethodPost {
		fmt.Println("Method: GET; Response: 301")
		w.Header().Set("Location", root)
		w.WriteHeader(http.StatusMovedPermanently)
		return
	}
	fmt.Println("Method: POST")
	postMap, err := postToMap(r)
	if err != nil {
		w.Write([]byte(page.Error + `
			<div class="alert alert-danger" role="alert">
		  	<h5 class="alert-heading">提示</h5>
		  	<hr>
		  	<p class="card-text">可能是暂时性错误,请重试</p>
		  	</div>` + page.Errordiv))
		fmt.Println(err)
		return
	}
	dresult := ""
	fsid := postMap["fsid"]
	timestamp := postMap["time"]
	sign := postMap["sign"]
	randsk := postMap["randsk"]
	shareid := postMap["shareid"]
	uk := postMap["uk"]
	// dl.Init(configPath)
	realLink, err := dl.GetFileRealLink(fsid, timestamp, sign, randsk, shareid, uk)
	if err != nil {
		w.Write([]byte(page.Error + `
			<div class="alert alert-danger" role="alert">
		  	<h5 class="alert-heading">提示</h5>
		  	<hr>
		  	<p class="card-text">可能是暂时性错误,请重试</p>
		  	</div>` + page.Errordiv))
		fmt.Println(err)
		return
	}
	realLink = realLink[7:] // 去掉前面的http://
	dresult += `<div class="alert alert-primary" role="alert">
	<h5 class="alert-heading">获取下载链接成功</h5>
	<hr>
	<p class="card-text"><a href="http://` + realLink + `" onclick= target=_blank>下载链接(http)</a> <a href="https://` + realLink + `" target=_blank>下载链接(https)</a><br><br><a>推送到Aria2(即将支持)</a><br><br><a href="./help">下载链接使用方法（必读）</a></p>
	</div>`
	w.Write([]byte(page.Dbody + dresult + page.Dfooter))
	fmt.Println("success")
}

func netdiskpageHelpV1(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Help---------------------")
	fmt.Println("from: ", GetIP(r))

	w.Write([]byte(page.Helpbody + page.Helpcontent + page.Dfooter))
	fmt.Println("success")
}

func postToMap(r *http.Request) (postMap map[string]string, err error) {
	err = nil
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("InvalidPostData")
		}
	}()
	postByte, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	// fmt.Println(string(postByte))
	postMap = make(map[string]string)
	json.Unmarshal(postByte, &postMap)
	if len(postMap) == 0 {
		postValues, _ := url.ParseQuery(string(postByte))
		for postKey, postValueSlice := range postValues {
			postMap[postKey] = postValueSlice[0]
		}
		return
	}
	return

}

// GetIP 获取请求IP地址
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
