package main

import (
    "net/http"
    "encoding/json"
    "log"
)

func main()  {

    http.HandleFunc("/login1", login1)
    http.HandleFunc("/login2", login2)
    http.ListenAndServe("0.0.0.0:8082", nil)
}

type Resp struct {
    Code    string `json:"code"`
    Msg     string `json:"msg"`
}

type  Auth struct {
    Username string `json:"username"`
    Pwd      string   `json:"password"`
}

//post接口接收json數據
func login1(writer http.ResponseWriter,  request *http.Request)  {
    var auth Auth
    if err := json.NewDecoder(request.Body).Decode(&auth); err != nil {
        request.Body.Close()
        log.Fatal(err)
    }
    var result  Resp
    if auth.Username == "admin" && auth.Pwd == "123456" {
        result.Code = "200"
        result.Msg = "登錄成功"
    } else {
        result.Code = "401"
        result.Msg = "賬戶名或密碼錯誤"
    }
    if err := json.NewEncoder(writer).Encode(result); err != nil {
        log.Fatal(err)
    }
}

//接收x-www-form-urlencoded類型的post請求或者普通get請求
func login2(writer http.ResponseWriter,  request *http.Request)  {
    request.ParseForm()
    username, uError :=  request.Form["username"]
    pwd, pError :=  request.Form["password"]

    var result  Resp
    if !uError || !pError {
        result.Code = "401"
        result.Msg = "登錄失敗"
    } else if username[0] == "admin" && pwd[0] == "123456" {
        result.Code = "200"
        result.Msg = "登錄成功"
    } else {
        result.Code = "203"
        result.Msg = "賬戶名或密碼錯誤"
    }
    if err := json.NewEncoder(writer).Encode(result); err != nil {
        log.Fatal(err)
    }
}