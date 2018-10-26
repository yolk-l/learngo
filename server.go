package main

import (
    "bytes"
    "errors"
    "fmt"
    "html/template"
    "io"
    "net/http"
    "runtime"
    "os"
    "os/exec"
    "strconv"
    "io/ioutil"
)

// 端口
const (
    HTTP_PORT  string = "8080"
)

// 目录
const (
    CSS_CLIENT_PATH   = "/css/"
    DART_CLIENT_PATH  = "/js/"
    IMAGE_CLIENT_PATH = "/image/"

    CSS_SVR_PATH   = "web"
    DART_SVR_PATH  = "web"
    IMAGE_SVR_PATH = "web"
)

func init() {
    runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
    // 先把css和脚本服务上去
    http.Handle(CSS_CLIENT_PATH, http.FileServer(http.Dir(CSS_SVR_PATH)))
    http.Handle(DART_CLIENT_PATH, http.FileServer(http.Dir(DART_SVR_PATH)))

    // 网址与处理逻辑对应起来
    http.HandleFunc("/", HomePage)
    http.HandleFunc("/ajax", OnAjax)
    http.HandleFunc("/update", OnUpdate)

    // 开始服务
    err := http.ListenAndServe(":"+HTTP_PORT, nil)
    if err != nil {
        fmt.Println("服务失败 /// ", err)
    }
}

func WriteTemplateToHttpResponse(res http.ResponseWriter, t *template.Template) error {
    if t == nil || res == nil {
        return errors.New("WriteTemplateToHttpResponse: t must not be nil.")
    }
    var buf bytes.Buffer
    err := t.Execute(&buf, nil)
    if err != nil {
        return err
    }
    res.Header().Set("Content-Type", "text/html; charset=utf-8")
    _, err = res.Write(buf.Bytes())
    return err
}

func HomePage(res http.ResponseWriter, req *http.Request) {
    t, err := template.ParseFiles("web/index.html")
    if err != nil {
        fmt.Println(err)
        return
    }
    err = WriteTemplateToHttpResponse(res, t)
    if err != nil {
        fmt.Println(err)
        return
    }
}

func OnAjax(res http.ResponseWriter, req *http.Request) {
    req.ParseForm()
    fmt.Println("on ajax", req.Form["begin"])
    var begin = req.Form["begin"][0];
    var offset = req.Form["offset"][0];
    var begin_num, _ = strconv.Atoi(begin);
    var offset_num, _ = strconv.Atoi(offset);

    var file_path = "~/learn/go/update/server.go";
    var cmd_str = "sed -n '"+begin+", "+strconv.Itoa(begin_num+offset_num)+"p' "+file_path;
    var cmd *exec.Cmd;
    cmd = exec.Command("/bin/bash", "-c", cmd_str);
    stdout, err := cmd.StdoutPipe();
    if err != nil {
        fmt.Println(err);
        os.Exit(1);
    }

    err = cmd.Start()
    if err != nil {
        fmt.Println(err);
        os.Exit(2);
    }
    result, err := ioutil.ReadAll(stdout)
    io.WriteString(res, string(result));
}

func OnUpdate(res http.ResponseWriter, req *http.Request) {
    fmt.Println("on update")
    var cmd *exec.Cmd;
    cmd = exec.Command("/bin/bash", "-c", "git pull");
    cmd.Dir = "/home/lhh/work/s3server";
    stdout, err := cmd.StdoutPipe();
    if err != nil {
        fmt.Println(err);
        os.Exit(1);
    }

    err = cmd.Start()
    if err != nil {
        fmt.Println(err);
        os.Exit(2);
    }
    result, err := ioutil.ReadAll(stdout)
    if err != nil {
        fmt.Println(result, err);
        return;
    }
    fmt.Println(string(result));
    io.WriteString(res, string(result));
}
