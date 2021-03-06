/**
 * The httpserver file of hyperttp current module of xxd.
 *
 * @copyright   Copyright 2009-2017 青岛易软天创网络科技有限公司(QingDao Nature Easy Soft Network Technology Co,LTD, www.cnezsoft.com)
 * @license     ZPL (http://zpl.pub/page/zplv12.html)
 * @author      Archer Peng <pengjiangxiu@cnezsoft.com>
 * @package     server
 * @link        http://www.zentao.net
 */
package server

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    //"strings"
    "yin/xxd/api"
    "yin/xxd/util"
    "math/rand"
    "os/exec"
    "strings"
    "image"
    "image/png"
    "image/color"
    "yin/xxd/orm"
    "io/ioutil"
    "bytes"
    "image/jpeg"
)

type retCInfo struct {
    // server version
    Version string `json:"version"`

    // encrypt key
    Token string `json:"token"`

    // multiSite or singleSite
    SiteType string `json:"siteType"`

    UploadFileSize int64 `json:"uploadFileSize"`

    ChatPort  int  `json:"chatPort"`
    TestModel bool `json:"testModel"`
}

// route
const (
    download = "/download"
    upload   = "/upload"
    sInfo    = "/serverInfo"
)

// 获取文件大小的接口
type Size interface {
    Size() int64
}

// 获取文件信息的接口
type Stat interface {
    Stat() (os.FileInfo, error)
}

// 启动 http server
func InitHttp() {
    //crt, key, err := CreateSignedCertKey()
    //if err != nil {
    //    util.LogError().Println("https server start error!")
    //    return
    //}

    //err = api.StartXXD()
    //if err != nil {
    //    util.Exit("ranzhi server login error")
    //}

    mux := http.NewServeMux()
    fs := http.FileServer(http.Dir(GetCurrentPath()+"./public"))
    mux.Handle("/", fs)
    //mux.HandleFunc(download, fileDownload)
    //mux.HandleFunc(upload, fileUpload)
    mux.HandleFunc(sInfo, serverInfo)
    mux.HandleFunc("/avatar",userAvatar)
    mux.HandleFunc("/chatlog",getChatLog)
    addr := util.Config.Ip + ":" + util.Config.CommonPort

    //util.Println("file server start,listen addr:", addr, download)
    //util.Println("file server start,listen addr:", addr, upload)
	//
    //util.LogInfo().Println("file server start,listen addr:", addr, download)
    //util.LogInfo().Println("file server start,listen addr:", addr, upload)

    if util.Config.IsHttps != "1" {
        util.Println("http server start,listen addr:http://", addr)
        util.LogInfo().Println("http server start,listen addr:http://", addr, sInfo)

        if err := http.ListenAndServe(addr, mux); err != nil {
            util.LogError().Println("http server listen err:", err)
            util.Exit("http server listen err")
        }
    }else{
        util.Println("https server start,listen addr:https://", addr)
        util.LogInfo().Println("https server start,listen addr:https://", addr, sInfo)

        //if err := http.ListenAndServeTLS(addr, crt, key, mux); err != nil {
        //    util.LogError().Println("https server listen err:", err)
        //    util.Exit("https server listen err")
        //}
    }
}
func getChatLog(rw http.ResponseWriter, req *http.Request)  {

    //rs := []rune(req.URL.RawQuery)
    req.ParseForm()
    id:=req.Form.Get("id")
    //id:=string(rs[3:])
    isLogin:=api.CheckLogin(id)
    if isLogin{
        json,err:=orm.GetChatLog(id)
        util.Check(err)
        rw.Write(json)
    }else {
        rw.Write([]byte(`用户未登陆`))
    }
}
func userAvatar(rw http.ResponseWriter, req *http.Request) {
    const (
        dx = 48
        dy = 48
    )
    var img *image.NRGBA
    rs := []rune(req.URL.RawQuery)
    id:=string(rs[3:])
    var imgByte []byte
    var err error
    imgByte,err=orm.GetAvatar(id)
    if err!=nil{
        ff, _ := ioutil.ReadFile(GetCurrentPath()+"public/img/none.jpg")
        bbb := bytes.NewBuffer(ff)
        //util.Check(er)
        bt,_,_:=image.Decode(bbb)
        rw.Header().Set("Content-Type", "image/jpeg")
        jpeg.Encode(rw, bt,nil)
    }else {
        // 新建一个 指定大小的 RGBA位图
        img = image.NewNRGBA(image.Rect(0, 0, dx, dy))
        site := 4
        if len(imgByte)>0 {
            for y := 0; y < dy; y++ {
                for x := 0; x < dx; x++ {
                    b := imgByte[site - 4]
                    g := imgByte[site - 3]
                    r := imgByte[site - 2]
                    // 设置某个点的颜色，依次是 RGBA
                    img.Set(x, y, color.RGBA{r, g, b, 255})
                    site = site + 4
                }
            }
        }
        // 图片流方式输出
        rw.Header().Set("Content-Type", "image/png")
        png.Encode(rw, img)
    }

}
func GetCurrentPath() string {
    s, err := exec.LookPath(os.Args[0])
    if err != nil {
        fmt.Println(err.Error())
    }
    s = strings.Replace(s, "\\", "/", -1)
    s = strings.Replace(s, "\\\\", "/", -1)
    i := strings.LastIndex(s, "/")
    path := string(s[0 : i+1])
    return path
}

//文件下载
func fileDownload(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        fmt.Fprintln(w, "not supported request")
        return
    }

    r.ParseForm()
    reqFileName := r.Form["fileName"][0]
    reqFileTime := r.Form["time"][0]
    reqFileID := r.Form["id"][0]

    serverName := r.Form["ServerName"][0]
    if serverName == "" {
        serverName = util.Config.DefaultServer
    }

    //新增加验证方式
    reqSid := r.Form["sid"][0]
    reqGid := r.Form["gid"][0]
    session,err :=util.GetUid(serverName, reqGid)
    util.Println("file_session:",session)
    if err!=nil {
        fmt.Fprintln(w, "not supported request")
        return
    }
    if reqSid != string(util.GetMD5( session  + reqFileName )) {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    fileTime, err := util.String2Int64(reqFileTime)
    if err != nil {
        util.LogError().Println("file download,time undefined:", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // new file name = md5(old filename + fileID + fileTime)
    fileName := util.Config.UploadPath + serverName + "/" + util.GetYmdPath(fileTime) + util.GetMD5(reqFileName+reqFileID+reqFileTime)
    //util.Println(fileName)
    if util.IsNotExist(fileName) || util.IsDir(fileName) {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    http.ServeFile(w, r, fileName)
}

//文件上传
func fileUpload(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Access-Control-Allow-Origin", "*")
    w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-FILENAME, ServerName")
    w.Header().Add("Access-Control-Allow-Credentials", "true")

    if r.Method != "POST" {
        fmt.Fprintln(w, "not supported request")
        return
    }

    //util.Println(r.Header)
    serverName := r.Header.Get("ServerName")
    if serverName == "" {
        serverName = util.Config.DefaultServer
    }

    authorization := r.Header.Get("Authorization")
    if authorization != string(util.Token) {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    r.ParseMultipartForm(32 << 20)

    file, handler, err := r.FormFile("file")
    if err != nil {
        util.LogError().Println("form file error:", err)
        fmt.Fprintln(w, "form file error")
        return
    }
    defer file.Close()

    nowTime := util.GetUnixTime()
    savePath := util.Config.UploadPath + serverName + "/" + util.GetYmdPath(nowTime)
    if err := util.Mkdir(savePath); err != nil {
        util.LogError().Println("mkdir error:", err)
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintln(w, "mkdir error")
        return
    }

    var fileSize int64 = 0
    if statInterface, ok := file.(Stat); ok {
        fileInfo, _ := statInterface.Stat()
        fileSize = fileInfo.Size()
    }

    if sizeInterface, ok := file.(Size); ok {
        fileSize = sizeInterface.Size()
    }

    if fileSize <= 0 {
        util.LogError().Println("get file size error")
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintln(w, "get file size error")
        return
    }

    if fileSize > util.Config.UploadFileSize {
        // 400
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintln(w, "file is too large")
        return
    }

    //util.Println(r.Form)
    fileName := util.FileBaseName(handler.Filename)
    nowTimeStr := util.Int642String(nowTime)
    gid := r.Form["gid"][0]
    userID := r.Form["userID"][0]

    x2rJson := `{"userID":` + userID + `,"module":"chat","method":"uploadFile","params":["` + fileName + `","` + savePath + `",` + util.Int642String(fileSize) + `,` + nowTimeStr + `,"` + gid + `"]}`

    //util.Println(x2rJson)
    fileID, err := api.UploadFileInfo(serverName, []byte(x2rJson))
    if err != nil {
        util.LogError().Println("Upload file info error:", err)
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintln(w, "Upload file info error")
        return
    }
    if fileID == "" {
        fileID = fmt.Sprintf("%d", rand.Intn(999999) + 1)
    }

    // new file name = md5(old filename + fileID + nowTime)
    saveFile := savePath + util.GetMD5(fileName+fileID+nowTimeStr)
    //util.Println(saveFile)
    f, err := os.OpenFile(saveFile, os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        util.LogError().Println("open file error:", err)
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintln(w, "open file error")
        return
    }
    defer f.Close()
    io.Copy(f, file)

    x2cJson := `{"result":"success","data":{"time":` + nowTimeStr + `,"id":` + fileID + `,"name":"` + fileName + `"}}`
    //fmt.Fprintln(w, handler.Header)
    //util.Println(x2cJson)
    fmt.Fprintln(w, x2cJson)
}

//服务配置信息
func serverInfo(w http.ResponseWriter, r *http.Request) {

    w.Header().Add("Access-Control-Allow-Origin", "*")
    w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
    w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    w.Header().Add("Access-Control-Allow-Credentials", "true")

    //if r.Method != "POST" {
    //    fmt.Fprintln(w, "not supported request")
    //    return
    //}

    r.ParseForm()

    ok, err := api.VerifyLogin([]byte(r.Form["data"][0]))
    if err != nil {
        util.LogError().Println("verify login error:", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if !ok {
        //util.Println("auth error")
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    chatPort, err := util.String2Int(util.Config.ChatPort)
    if err != nil {
        util.LogError().Println("string to int error:", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    info := retCInfo{
       // Version:        util.Version,
        Token:          string(util.Token),
        SiteType:       util.Config.SiteType,
        UploadFileSize: util.Config.UploadFileSize,
        ChatPort:       chatPort,
        TestModel:      util.IsTest}

    jsonData, err := json.Marshal(info)
    if err != nil {
        util.LogError().Println("json unmarshal error:", err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    fmt.Fprintln(w, string(jsonData))
}
