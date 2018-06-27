package util

import (
    "os"
    "github.com/tidwall/gjson"
)

type RanzhiServer struct {
    RanzhiAddr  string
    RanzhiToken []byte
    //RanzhiEncrypt bool
}

type ConfigIni struct {
    Ip         string
    ChatPort   string
    CommonPort string
    IsHttps    string

    UploadPath     string
    UploadFileSize int64

    // multiSite or singleSite
    SiteType      string
    DefaultServer string
    RanzhiServer  map[string]RanzhiServer

    LogPath string
    CrtPath string
}

const configPath = "config/xxd.json"

var Config = ConfigIni{SiteType: "singleSite", RanzhiServer: make(map[string]RanzhiServer)}
//异常处理
func check(e error) {
    if e != nil {
        panic(e)
    }
}
func init() {
    //dir, _ := os.Getwd()
    //dat, err := ioutil.ReadFile(dir+"/config/xxd.json")
    //data:=string(dat)
    //if err != nil {

        Config.Ip = "192.168.1.250"//"127.0.0.1"
        Config.ChatPort = "11444"
        Config.CommonPort = "11443"
        Config.IsHttps = "0"

        //Config.UploadPath = "tmpfile"
        //Config.UploadFileSize = 32 * MB
		//
        //Config.SiteType = "singleSite"
        Config.DefaultServer = "ttd"
        //Config.RanzhiServer["xuanxuan"] = RanzhiServer{"serverInfo", []byte("serverInfo")}

        Config.LogPath = GetCurrentPath() + "/log/"
        //Config.CrtPath = dir + "/certificate/"

        //log.Println("config init error，use default conf!")
        //log.Println(Config)
        //return
   // }
}
//
////获取上传目录
//func getUploadPath(config *goconfig.ConfigFile) (err error) {
//    Config.UploadPath, err = config.GetValue("server", "uploadPath")
//    if err != nil {
//        log.Fatal("config: get server upload path error,", err)
//    }
//
//    return
//}
//
////获取上传大小
//func getUploadFileSize(config *goconfig.ConfigFile) error {
//
//    Config.UploadFileSize = 32 * MB
//    var fileSize int64 = 0
//
//    uploadFileSize, err := config.GetValue("server", "uploadFileSize")
//    if err != nil {
//        log.Printf("config: get server upload file size error:%v, default size 32MB.", err)
//        return err
//    }
//
//    switch size, suffix := sizeSuffix(uploadFileSize); suffix {
//    case "K":
//        if fileSize, err = String2Int64(size); err == nil {
//            Config.UploadFileSize = fileSize * KB
//        }
//
//    case "M":
//        if fileSize, err = String2Int64(size); err == nil {
//            Config.UploadFileSize = fileSize * MB
//        }
//
//    case "G":
//        if fileSize, err = String2Int64(size); err == nil {
//            Config.UploadFileSize = fileSize * GB
//        }
//
//    default:
//        if fileSize, err = String2Int64(size); err == nil {
//            Config.UploadFileSize = fileSize
//        } else {
//            log.Println("config: get server upload file size error, default size 32MB.")
//        }
//    }
//
//    if err != nil {
//        log.Println("upload file size parse error:", err)
//    }
//
//    return err
//}
//
////获取服务器列表,conf中[ranzhi]段不能改名.
//func getRanzhi(config *goconfig.ConfigFile) {
//    keyList := config.GetKeyList("ranzhi")
//
//    Config.DefaultServer = ""
//    if len(keyList) > 1 {
//        Config.SiteType = "multiSite"
//    }
//
//    for _, ranzhiName := range keyList {
//        ranzhiServer, err := config.GetValue("ranzhi", ranzhiName)
//        if err != nil {
//            log.Fatal("config: get ranzhi server error,", err)
//        }
//
//        serverInfo := strings.Split(ranzhiServer, ",")
//        //逗号前面是地址，后面是token，token长度固定为32
//        if len(serverInfo) < 2 || len(serverInfo[1]) != 32 {
//            log.Fatal("config: ranzhi server config error")
//        }
//
//        if len(serverInfo) >= 3 && serverInfo[2] == "default" {
//            Config.DefaultServer = ranzhiName
//        }
//
//        Config.RanzhiServer[ranzhiName] = RanzhiServer{serverInfo[0], []byte(serverInfo[1])}
//    }
//}

//获取日志路径
func getLogPath(config string) () {
    dir, _ := os.Getwd()
    logPath :=gjson.Get(config,"log.logPath").String()
    Config.LogPath = dir + "/" + logPath
}

//获取证书路径
func getCrtPath(config string) () {
    dir, _ := os.Getwd()
    crtPath := gjson.Get(config,"certificate.crtPath").String()
    Config.CrtPath = dir + "/" + crtPath
}

//func sizeSuffix(uploadFileSize string) (string, string) {
//    if strings.HasSuffix(uploadFileSize, "K") {
//        return strings.TrimSuffix(uploadFileSize, "K"), "K"
//    }
//
//    if strings.HasSuffix(uploadFileSize, "M") {
//        return strings.TrimSuffix(uploadFileSize, "M"), "M"
//    }
//
//    if strings.HasSuffix(uploadFileSize, "G") {
//        return strings.TrimSuffix(uploadFileSize, "G"), "G"
//    }
//
//    return uploadFileSize, ""
//}
