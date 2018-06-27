/**
 * The test file of test current module of xxd.
 *
 * @copyright   Copyright 2009-2017 青岛易软天创网络科技有限公司(QingDao Nature Easy Soft Network Technology Co,LTD, www.cnezsoft.com)
 * @license     ZPL (http://zpl.pub/page/zplv12.html)
 * @author      Archer Peng <pengjiangxiu@cnezsoft.com>
 * @package     main
 * @link        http://www.zentao.net
 */

package main

import (
    "net/http"
    "html/template"
	"os"
	"os/exec"
	"fmt"
	"strings"
)


func main() {
    tmpl := template.Must(template.ParseFiles(GetCurrentPath()+"/Test.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        tmpl.Execute(w,nil)
    })

    http.ListenAndServe(":2061", nil)

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

