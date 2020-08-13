package page

import (
	cfg "baidunetdisk/config"
)

var (
	config cfg.Config
	root   string
	link   string
	name   string
	// 以下为页面
	Error       string
	Filebody    string
	Errordiv    string
	Landing     string
	Helpbody    string
	Dbody       string
	Dfooter     string
	Helpcontent string
	Filefoot    string
)

func Init(configPath string) {
	config = cfg.GetConfig(configPath)
	root = config["httpserver"]["root"]
	link = config["httpserver"]["url"]
	name = config["httpserver"]["name"]

	Error = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"/>
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="author" content="Ling Macker"/>
<meta name="description" content="PanDownload网页版,百度网盘分享链接在线解析工具"/>
<meta name="keywords" content="PanDownload,百度网盘,分享链接,下载,不限速"/>
<link rel="icon" href="https://pandownload.com/favicon.ico" type="image/x-icon"/>
<link rel="stylesheet" href="https://cdn.staticfile.org/twitter-bootstrap/4.1.2/css/bootstrap.min.css">
<script src="https://cdn.staticfile.org/jquery/3.2.1/jquery.min.js"></script>
<script src="https://cdn.staticfile.org/popper.js/1.12.5/umd/popper.min.js"></script>
<script src="https://cdn.staticfile.org/twitter-bootstrap/4.1.2/js/bootstrap.min.js"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag() { dataLayer.push(arguments); }
  gtag('js', new Date());
  gtag('config', 'UA-112166098-2');
</script>
<style>
  body {
    background-image: url("https://pandownload.com/img/baiduwp/bg.png");
  }
  .logo-img {
    width: 1.1em;
    position: relative;
    top: -3px;
  }
</style>
<meta name="referrer" content="never">
<title>提示</title>
<style>
    .alert {
      position: relative;
      top: 5em;
    }
    .alert-heading {
      height: 0.8em;
    }
  </style>
</head>
<body>
<nav class="navbar navbar-expand-sm bg-dark navbar-dark">
<div class="container">
<a class="navbar-brand" href="` + link + `">
<img src="https://pandownload.com/img/baiduwp/logo.png" class="img-fluid rounded logo-img mr-2" alt="LOGO">` + name + `
</a>
<button class="navbar-toggler border-0" type="button" data-toggle="collapse" data-target="#collpase-bar">
<span class="navbar-toggler-icon"></span>
</button>
<div class="collapse navbar-collapse" id="collpase-bar">
<ul class="navbar-nav">
<li class="nav-item">
<a class="nav-link" href="` + link + `">主页</a>
</li>
<li class="nav-item">
<a class="nav-link" href="https://pandownload.com/" target="_blank">网盘下载器</a>
</li>
<li class="nav-item">
<a class="nav-link" href="https://www.pandownload.com/donate.html">捐助</a>
</li>
</ul>
</div>
</div>
</nav>
<div class="container">
<div class="row justify-content-center">
<div class="col-md-7 col-sm-8 col-11">`

	Filebody = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"/>
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="author" content="Ling Macker"/>
<meta name="description" content="PanDownload网页版,百度网盘分享链接在线解析工具"/>
<meta name="keywords" content="PanDownload,百度网盘,分享链接,下载,不限速"/>
<link rel="icon" href="https://pandownload.com/favicon.ico" type="image/x-icon"/>
<link rel="stylesheet" href="https://cdn.staticfile.org/twitter-bootstrap/4.1.2/css/bootstrap.min.css">
<script src="https://cdn.staticfile.org/jquery/3.2.1/jquery.min.js"></script>
<script src="https://cdn.staticfile.org/popper.js/1.12.5/umd/popper.min.js"></script>
<script src="https://cdn.staticfile.org/twitter-bootstrap/4.1.2/js/bootstrap.min.js"></script>
<script src="https://cdn.staticfile.org/limonte-sweetalert2/8.11.8/sweetalert2.all.min.js"></script>
<style>
  body {
    background-image: url("https://pandownload.com/img/baiduwp/bg.png");
  }
  .logo-img {
    width: 1.1em;
    position: relative;
    top: -3px;
  }
</style>
<meta name="referrer" content="never">
<link href="https://cdn.staticfile.org/font-awesome/5.8.1/css/all.min.css" rel="stylesheet">
<title>文件列表</title>
<script>
  function dl(fsid, timestamp, sign, randsk, shareid, uk) {
    var form = $('<form method="post" action="./download" target="_blank"></form>');
    form.append('<input type="hidden" name="fsid" value="'+fsid+'">');
    form.append('<input type="hidden" name="time" value="'+timestamp+'">');
    form.append('<input type="hidden" name="sign" value="'+sign+'">');
    form.append('<input type="hidden" name="randsk" value="'+randsk+'">');
    form.append('<input type="hidden" name="shareid" value="'+shareid+'">');
    form.append('<input type="hidden" name="uk" value="'+uk+'">');
    $(document.body).append(form);
    form.submit();
  }
  function getdirfilelist(randsk,uk,shareid,path,timestamp,sign) {
    var form = $('<form method="post" action="` + root + `" target="_parent"></form>');
    form.append('<input type="hidden" name="randsk" value="'+randsk+'">');
    form.append('<input type="hidden" name="uk" value="'+uk+'">');
    form.append('<input type="hidden" name="shareid" value="'+shareid+'">');
    form.append('<input type="hidden" name="path" value="'+path+'">');
    form.append('<input type="hidden" name="timestamp" value="'+timestamp+'">');
    form.append('<input type="hidden" name="sign" value="'+sign+'">');
    $(document.body).append(form);
    form.submit();
  }
  function getIconClass(filename){
    var filetype = {
      file_video: ["wmv", "rmvb", "mpeg4", "mpeg2", "flv", "avi", "3gp", "mpga", "qt", "rm", "wmz", "wmd", "wvx", "wmx", "wm", "mpg", "mp4", "mkv", "mpeg", "mov", "asf", "m4v", "m3u8", "swf"],
      file_audio: ["wma", "wav", "mp3", "aac", "ra", "ram", "mp2", "ogg", "aif", "mpega", "amr", "mid", "midi", "m4a", "flac"],
      file_image: ["jpg", "jpeg", "gif", "bmp", "png", "jpe", "cur", "svgz", "ico"],
      file_archive: ["rar", "zip", "7z", "iso"],
      windows: ["exe"],
      apple: ["ipa"],
      android: ["apk"],
      file_alt: ["txt", "rtf"],
      file_excel: ["xls", "xlsx"],
      file_word: ["doc", "docx"],
      file_powerpoint: ["ppt", "pptx"],
      file_pdf: ["pdf"],
    };
    var point = filename.lastIndexOf(".");
    var t = filename.substr(point+1);
    if (t == ""){
      return "";
    }
    t = t.toLowerCase();
    for(var icon in filetype) {
      for(var type in filetype[icon]) {
        if (t == filetype[icon][type])
        {
          return "fa-"+icon.replace('_', '-');
        }
      }
    }
    return "";
  }
  $(document).ready(function(){
    $(".fa-file").each(function(){
      var icon = getIconClass($(this).next().text());
      if (icon != "")
      {
        if (icon == "fa-windows" || icon == "fa-android" || icon == "fa-apple")
        {
          $(this).removeClass("far").addClass("fab");
        }
        $(this).removeClass("fa-file").addClass(icon);
      }
    });
  });
</script>
</head>
<body>
<nav class="navbar navbar-expand-sm bg-dark navbar-dark">
<div class="container">
<a class="navbar-brand" href="` + link + `">
<img src="https://pandownload.com/img/baiduwp/logo.png" class="img-fluid rounded logo-img mr-2" alt="LOGO">` + name + `
</a>
<button class="navbar-toggler border-0" type="button" data-toggle="collapse" data-target="#collpase-bar">
<span class="navbar-toggler-icon"></span>
</button>
<div class="collapse navbar-collapse" id="collpase-bar">
<ul class="navbar-nav">
<li class="nav-item">
<a class="nav-link" href="` + link + `">主页</a>
</li>
<li class="nav-item">
<a class="nav-link" href="https://pandownload.com/" target="_blank">网盘下载器</a>
</li>
<li class="nav-item">
<a class="nav-link" href="https://www.pandownload.com/donate.html">捐助</a>
</li>
</ul>
</div>
</div>
</nav>
<div class="container">
<ol class="breadcrumb my-4">
文件列表 </ol>
<div>
<ul class="list-group ">`

	Errordiv = `</div>
</div>
</div>
<div style="display:none">
</div>
</body>
</html>`

	Landing = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"/>
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="author" content="Ling Macker"/>
<meta name="description" content="PanDownload网页版,百度网盘分享链接在线解析工具"/>
<meta name="keywords" content="PanDownload,百度网盘,分享链接,下载,不限速"/>
<link rel="icon" href="https://pandownload.com/favicon.ico" type="image/x-icon"/>
<link rel="stylesheet" href="https://cdn.staticfile.org/twitter-bootstrap/4.1.2/css/bootstrap.min.css">
<script src="https://cdn.staticfile.org/jquery/3.2.1/jquery.min.js"></script>
<script src="https://cdn.staticfile.org/popper.js/1.12.5/umd/popper.min.js"></script>
<script src="https://cdn.staticfile.org/twitter-bootstrap/4.1.2/js/bootstrap.min.js"></script>
<script>
window.dataLayer = window.dataLayer || [];
function gtag() { dataLayer.push(arguments); }
gtag('js', new Date());
gtag('config', 'UA-112166098-2');
</script>
<style>
body {
background-image: url("https://pandownload.com/img/baiduwp/bg.png");
}
.logo-img {
width: 1.1em;
position: relative;
top: -3px;
}
</style>
<title>` + name + `网页版</title>
<style>
.form-inline input {
width: 500px;
}
.input-card {
position: relative;
top: 7.0em;
}
.card-header {
height: 3.2em;
font-size: 20px;
line-height: 2.0em;
}
form input,
form button {
height: 3em;
}
</style>
<script>
function validateForm() {
var link = document.forms["form1"]["surl"].value;
if (link == null || link == "") {
document.forms["form1"]["surl"].focus();
return false;
}
var surl = link.match(/surl=([A-Za-z0-9-_]+)/);
if (surl == null) {
surl = link.match(/1[A-Za-z0-9-_]+/);
if (surl == null) {
document.forms["form1"]["surl"].focus();
return false;
}
else {
surl = surl[0].substring(1);
}
}
else {
surl = surl[1];
}
document.forms["form1"]["surl"].value = surl;
return true;
}
</script>
</head>
<body>
<nav class="navbar navbar-expand-sm bg-dark navbar-dark">
<div class="container">
<a class="navbar-brand" href="` + link + `">
<img src="https://pandownload.com/img/baiduwp/logo.png" class="img-fluid rounded logo-img mr-2" alt="LOGO">` + name + `
</a>
<button class="navbar-toggler border-0" type="button" data-toggle="collapse" data-target="#collpase-bar">
<span class="navbar-toggler-icon"></span>
</button>
<div class="collapse navbar-collapse" id="collpase-bar">
<ul class="navbar-nav">
<li class="nav-item">
<a class="nav-link" href="` + link + `">主页</a>
</li>
<li class="nav-item">
<a class="nav-link" href="https://pandownload.com/" target="_blank">网盘下载器</a>
</li>
</ul>
</div>
</div>
</nav>
<div class="container">
<div class="col-lg-6 col-md-9 mx-auto mb-5 input-card">
<div class="card">
<div class="card-header bg-dark text-light">分享链接在线解析</div>
<div class="card-body">
<form name="form1" method="post" action="` + root + `" onsubmit="return validateForm()">
<div class="form-group my-2">
<input type="text" class="form-control" name="surl" placeholder="分享链接">
</div>
<div class="form-group my-4">
<input type="text" class="form-control" name="pwd" placeholder="提取码">
</div>
<button type="submit" class="mt-4 mb-3 form-control btn btn-success btn-block">打开</button>
</form>
</div>
</div>
</div>
</div>
<div style="display:none">
</div>
</body>
</html>
`

	Helpbody = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"/>
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="author" content="Ling Macker"/>
<meta name="description" content="PanDownload网页版,百度网盘分享链接在线解析工具"/>
<meta name="keywords" content="PanDownload,百度网盘,分享链接,下载,不限速"/>
<link rel="icon" href="https://pandownload.com/favicon.ico" type="image/x-icon"/>
<link rel="stylesheet" href="https://cdn.staticfile.org/twitter-bootstrap/4.1.2/css/bootstrap.min.css">
<script src="https://cdn.staticfile.org/jquery/3.2.1/jquery.min.js"></script>
<script src="https://cdn.staticfile.org/popper.js/1.12.5/umd/popper.min.js"></script>
<script src="https://cdn.staticfile.org/twitter-bootstrap/4.1.2/js/bootstrap.min.js"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag() { dataLayer.push(arguments); }
  gtag('js', new Date());
  gtag('config', 'UA-112166098-2');
</script>
<style>
  body {
    background-image: url("https://pandownload.com/img/baiduwp/bg.png");
  }
  .logo-img {
    width: 1.1em;
    position: relative;
    top: -3px;
  }
</style>
<meta name="referrer" content="never">
<title>下载链接使用方法</title>
<style>
    .alert {
      position: relative;
      top: 5em;
    }
    .alert-heading {
      height: 0.8em;
    }
  </style>
</head>
<body>
<nav class="navbar navbar-expand-sm bg-dark navbar-dark">
<div class="container">
<a class="navbar-brand" href="` + link + `">
<img src="https://pandownload.com/img/baiduwp/logo.png" class="img-fluid rounded logo-img mr-2" alt="LOGO">` + name + `
</a>
<button class="navbar-toggler border-0" type="button" data-toggle="collapse" data-target="#collpase-bar">
<span class="navbar-toggler-icon"></span>
</button>
<div class="collapse navbar-collapse" id="collpase-bar">
<ul class="navbar-nav">
<li class="nav-item">
<a class="nav-link" href="` + link + `">主页</a>
</li>
<li class="nav-item">
<a class="nav-link" href="https://pandownload.com/" target="_blank">网盘下载器</a>
</li>
<li class="nav-item">
<a class="nav-link" href="https://www.pandownload.com/donate.html">捐助</a>
</li>
</ul>
</div>
</div>
</nav>
<div class="container">
<div class="row justify-content-center">
<div class="col-md-7 col-sm-8 col-11">`

	Dbody = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"/>
<meta name="viewport" content="width=device-width, initial-scale=1">
<meta name="author" content="Ling Macker"/>
<meta name="description" content="PanDownload网页版,百度网盘分享链接在线解析工具"/>
<meta name="keywords" content="PanDownload,百度网盘,分享链接,下载,不限速"/>
<link rel="icon" href="https://pandownload.com/favicon.ico" type="image/x-icon"/>
<link rel="stylesheet" href="https://cdn.staticfile.org/twitter-bootstrap/4.1.2/css/bootstrap.min.css">
<script src="https://cdn.staticfile.org/jquery/3.2.1/jquery.min.js"></script>
<script src="https://cdn.staticfile.org/popper.js/1.12.5/umd/popper.min.js"></script>
<script src="https://cdn.staticfile.org/twitter-bootstrap/4.1.2/js/bootstrap.min.js"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag() { dataLayer.push(arguments); }
  gtag('js', new Date());
  gtag('config', 'UA-112166098-2');
</script>
<style>
  body {
    background-image: url("https://pandownload.com/img/baiduwp/bg.png");
  }
  .logo-img {
    width: 1.1em;
    position: relative;
    top: -3px;
  }
</style>
<meta name="referrer" content="never">
<title>提示</title>
<style>
    .alert {
      position: relative;
      top: 5em;
    }
    .alert-heading {
      height: 0.8em;
    }
  </style>
</head>
<body>
<nav class="navbar navbar-expand-sm bg-dark navbar-dark">
<div class="container">
<a class="navbar-brand" href="` + link + `">
<img src="https://pandownload.com/img/baiduwp/logo.png" class="img-fluid rounded logo-img mr-2" alt="LOGO">` + name + `
</a>
<button class="navbar-toggler border-0" type="button" data-toggle="collapse" data-target="#collpase-bar">
<span class="navbar-toggler-icon"></span>
</button>
<div class="collapse navbar-collapse" id="collpase-bar">
<ul class="navbar-nav">
<li class="nav-item">
<a class="nav-link" href="` + link + `">主页</a>
</li>
<li class="nav-item">
<a class="nav-link" href="https://pandownload.com/" target="_blank">网盘下载器</a>
</li>
<li class="nav-item">
<a class="nav-link" href="https://www.pandownload.com/donate.html">捐助</a>
</li>
</ul>
</div>
</div>
</nav>
<div class="container">
<div class="row justify-content-center">
<div class="col-md-7 col-sm-8 col-11">`

	Dfooter = `
</div>
</div>
</div>
</body>
</html>`

	Helpcontent = `
<div class="alert alert-primary" role="alert">
<h5 class="alert-heading">提示</h5>
<hr>
<p class="card-text">因百度限制，需修改浏览器UA后下载。<br>
<div class="page-inner">
<section class="normal" id="section-">
<h4>IDM（推荐）</h4>
<ol>
<li>选项 -> 下载 -> 手动添加任务时使用的用户代理（UA）-> 填入 <b>LogStatistic</b></li>
<li>右键复制下载链接，在 IDM 新建任务，粘贴链接即可下载。</li>
</ol>
<h4>Chrome浏览器</h4>
<ol>
<li>安装浏览器扩展程序 <a href="https://chrome.google.com/webstore/detail/user-agent-switcher-for-c/djflhoibgkdhkhhcedjiklpkjnoahfmg" target="_blank">User-Agent Switcher for Chrome</a></li>
<li>右键点击扩展图标 -> 选项</li>
<li>New User-agent name 填入 百度网盘分享下载</li>
<li>New User-Agent String 填入 LogStatistic</li>
<li>Group 填入 百度网盘</li>
<li>Append? 选择 Replace</li>
<li>Indicator Flag 填入 Log，点击 Add 保存</li>
<li>保存后点击扩展图标，出现"百度网盘"，进入并选择"百度网盘分享下载"。</li>
</ol>
<blockquote>
<p>Chrome应用商店打不开或者其他Chromium内核的浏览器，<a href="http://pandownload.com/static/user_agent_switcher_1_0_43_0.crx" target="_blank">请点此下载</a></p>
<p><a href="https://appcenter.browser.qq.com/search/detail?key=User-Agent%20Switcher%20for%20Chrome&amp;id=djflhoibgkdhkhhcedjiklpkjnoahfmg%20&amp;title=User-Agent%20Switcher%20for%20Chrome" target="_blank">QQ浏览器插件下载</a></p>
</blockquote>
<h4>Pure浏览器（Android）</h4>
<ol>
<li>设置 –&gt; 浏览设置 -&gt; 浏览器标识(UA)</li>
<li>添加自定义UA：LogStatistic</li>
</ol>
<h4>Alook浏览器（IOS）</h4>
<ol>
<li>设置 -&gt; 通用设置 -&gt; 浏览器标识 -&gt; 移动版浏览器标识 -&gt; 自定义 -><br> 填入 <b>LogStatistic</b></li>
</ol>
</section>
</div>
</p>
</div>
`

	Filefoot = `</ul>
</div>
</div>
</body>
</html>`

}
