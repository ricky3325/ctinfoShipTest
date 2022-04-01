package main

import (
    "net/http"
    "encoding/json"
    "log"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "gopkg.in/redis.v3"
    "fmt"
    "time"
    "io/ioutil"
    "math"
    "strings"
    "errors"
    "github.com/gorilla/mux"
    "strconv"
)

/*func main()  {

    http.HandleFunc("/login1", login1)
    http.HandleFunc("/login2", login2)
    http.ListenAndServe("0.0.0.0:8080", nil)
}*/

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/login1/{id}", login1)
    r.HandleFunc("/login1", login1)
    r.HandleFunc("/login2", login2)
    r.HandleFunc("/login3", login3)
    http.ListenAndServe("0.0.0.0:8080", r)
}

type FullData struct {
	ID  int
	Url string
	ExpireAt string
}

type Cumstomer struct {
	ID  int
	Username string
	Password string
}

type CretUrsho struct {
    Url         string `json:"url"`
    ExpireAt    string `json:"expireAt"`
}

type Resp struct {
    Code    string `json:"code"`
    Msg     string `json:"msg"`
}

type  Auth struct {
    Username string `json:"username"`
    Pwd      string   `json:"password"`
}

type creatUrlShortnerType struct {
    ResID    string `json:"id"`
    ResUrl     string `json:"shortUrl":`
}

const (
	UserName     string = "root"
	Password     string = "12345"
	Addr         string = "mysql"
	Port         int    = 3306
	Database     string = "mydb"
	MaxLifetime  int    = 10
	MaxOpenConns int    = 10
	MaxIdleConns int    = 10
)

/*type App struct {
	MyDB    *sql.DB
}*/
var (
	MyDB         *sql.DB
    RedisDB      *redis.Client
)

const (
    alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    length   = uint64(len(alphabet))
)

//post接口接收json數據
func login1(writer http.ResponseWriter,  request *http.Request)  {
    /*var auth Auth
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
    }*/

    /*var result  Resp
    result.Code = "401"
    result.Msg = "登錄失敗"

    if err := json.NewEncoder(writer).Encode(result); err != nil {
        log.Fatal(err)
    }*/
    //var results []string
    if request.Method == "POST" {
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "Error reading request body",
				http.StatusInternalServerError)
		}
		//results = append(results, string(body))

        var postData CretUrsho
        if err := json.Unmarshal(body, &postData); err != nil {   // Parse []byte to go struct pointer
            fmt.Println("Can not unmarshal JSON", err)
        }
        fmt.Println(postData.Url)

        x, y := autoAdd(postData.Url, postData.ExpireAt)
        //z := x + y + "POST done"
		//fmt.Fprint(writer, z)

        res2D := &creatUrlShortnerType{
            ResID:  x,
            ResUrl: y,
        }
        res2B, _ := json.Marshal(res2D)
        writer.Write(res2B)
        //fmt.Println(string(res2B))

	} else if request.Method == "GET"{
        vars := mux.Vars(request)
        id, ok := vars["id"]
        if !ok {
            fmt.Println("id is missing in parameters")
        }
        B, C := Decode(id)
        if C != nil {
            //fmt.Fprint(writer, "Decode Fail:", C)
            fmt.Println("Decode Fail:", C)
            //return
        }else{
            fmt.Println("Decode(A)B:", B)

            if getDataToRedis(strconv.FormatUint(B, 10)) != "Fail" {
                fmt.Println("get from redis")
                http.Redirect(writer, request, getDataToRedis(strconv.FormatUint(B, 10)), http.StatusSeeOther)
            }else{
                fmt.Println("get from mysql")
                FullUrl := returnUrl(strconv.FormatUint(B, 10))
                fmt.Println(FullUrl)
                if FullUrl == "overDue" {
                    addDataToRedis(strconv.FormatUint(B, 10), "http://localhost/NotFound/")
                    fmt.Println("overDue overDue overDue")
                    http.Redirect(writer, request, "http://localhost/NotFound/", http.StatusSeeOther)
                }else if FullUrl == "Scan Failed" {
                    addDataToRedis(strconv.FormatUint(B, 10), "http://localhost/NotFound/")
                    fmt.Println("Scan Failed Scan Failed Scan Failed")
                    http.Redirect(writer, request, "http://localhost/NotFound/", http.StatusSeeOther)
                }else{
                    addDataToRedis(strconv.FormatUint(B, 10), FullUrl)
                    http.Redirect(writer, request, FullUrl, http.StatusSeeOther)
                }
            }
        }
        
        //fmt.Println(`id := `, id)
        //fmt.Fprint(writer, "GET done:",FullUrl)

    }else {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func NotFound(writer http.ResponseWriter,  request *http.Request)  {
    var result  Resp
    result.Code = "404"
    result.Msg = "NotFound"
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
    } else if username[0] == "admin" && pwd[0] == "0" {
        result.Code = "200"
        result.Msg = "0Connect"
        DbConnectSQL()
        DbConnectRedis()
    }else if username[0] == "admin" && pwd[0] == "1" {
        result.Code = "200"
        result.Msg = "1CreateTable"
        CreateTable()
    }else if username[0] == "3"{
        result.Code = "200"
        result.Msg = "3ReadFullData"
        ReadFullData(pwd[0])
    }else if username[0] == "4" {
        result.Code = "200"
        result.Msg = "4getDataToRedis"
        getDataToRedis(pwd[0])
    }else {
        result.Code = "203"
        result.Msg = "賬戶名或密碼錯誤"
    }
    if err := json.NewEncoder(writer).Encode(result); err != nil {
        log.Fatal(err)
    }
}

func login3(writer http.ResponseWriter,  request *http.Request)  {
    request.ParseForm()
    id, uError :=  request.Form["id"]
    info, iError :=  request.Form["info"]
    numberId, err := strconv.ParseUint(id[0], 10, 64)
    if request.Method == "POST"{
        var result  Resp
        if !uError || !iError || err != nil {
            result.Code = "401"
            result.Msg = "失敗"
        } else {
            result.Code = "200"
            result.Msg = "Good>> ID: " + id[0] + "INFO: " + info[0]
            addDataToRedis(strconv.FormatUint(numberId, 10), info[0])
        }
        if err := json.NewEncoder(writer).Encode(result); err != nil {
            log.Fatal(err)
        }
    } else if request.Method == "GET"{
        var result  Resp
        if !uError || !iError || err != nil {
            result.Code = "401"
            result.Msg = "失敗"
        } else {
            writer.WriteHeader(http.StatusOK)
            writer.Header().Set("Content-Type", "application/text")
            writer.Write([]byte(getDataToRedis(strconv.FormatUint(numberId, 10))))
            return 
        }
    }
}

func DbConnectSQL(){
    //組合sql連線字串
    conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", UserName, Password, Addr, Port, Database)
    //連接MySQL
    //DB, err := sql.Open("mysql", conn)
    //MyDB = DB
    err := error(nil)
    MyDB, err = sql.Open("mysql", conn)
    if err != nil {
        fmt.Println("connection to mysql failed:", err)
        return
    }else {
        fmt.Println("connected to mysql")
    }
    MyDB.SetConnMaxLifetime(time.Duration(MaxLifetime) * time.Second)
    MyDB.SetMaxOpenConns(MaxOpenConns)
    MyDB.SetMaxIdleConns(MaxIdleConns)
    //fmt.Println("connected to mysql")
}

func DbConnectRedis() {
    fmt.Println("golang連接redis")
    RedisDB = redis.NewClient(&redis.Options{
        Addr: "redis:6379",
        Password: "",
        DB: 0,
    })
    pong, err := RedisDB.Ping().Result()
    fmt.Println(pong, err)
}

func addDataToRedis(id string, url string) {
    // 第三个参数是过期时间, 如果是0, 则表示没有过期时间. 这里设置过期时间.
    fmt.Println("addDataToRedis")
    err := RedisDB.Set(id, url, 600 * time.Second).Err()
    if err != nil {
        fmt.Println("add Data To Redis failed:", err)
    }
}

func getDataToRedis(id string) (string) {
    val, err := RedisDB.Get(id).Result()
    if err != nil {
        fmt.Println("get Data To Redis failed:", err)
        return "Fail"
    }
    fmt.Println(id, ":get: ", val)
    return val
}

func CreateTable() {
	sql := `CREATE TABLE urshoner(
	id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
	url NVARCHAR(2084) NOT NULL,
	ExpireAt timestamp NOT NULL
	); `

	if _, err := MyDB.Exec(sql); err != nil {
		fmt.Println("create table failed:", err)
		return
	}
	fmt.Println("create table successd")
}


func ReadFullData(Num string) {
    var fullData FullData
    //單筆資料
	row := MyDB.QueryRow("select id,url,ExpireAt from urshoner where id=?", Num)
    //Scan對應的欄位與select語法的欄位順序一致
	if err := row.Scan(&fullData.ID, &fullData.Url, &fullData.ExpireAt); err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return
	}
    fmt.Println("fullData.ID:", fullData.ID)
    fmt.Println("fullData.Url:", fullData.Url)
    fmt.Println("fullData.ExpireAt:", fullData.ExpireAt)
    A := Encode(uint64(fullData.ID))
    fmt.Println("Encode(fullData.ID):", A)
    B, C := Decode(A)
    fmt.Println("Decode(A)B:", B)
    fmt.Println("Decode(A)C:", C)
	//fmt.Println("fullData:%+v:", fullData)
}

func returnUrl(Num string) (string) {
    fmt.Println("returnUrl")
    var fullData FullData
    //單筆資料
	row := MyDB.QueryRow("select id,url,ExpireAt from urshoner where id=?", Num)
    //Scan對應的欄位與select語法的欄位順序一致
	if err := row.Scan(&fullData.ID, &fullData.Url, &fullData.ExpireAt); err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return "Scan Failed"
	}
    fmt.Println("fullData.Url:", fullData.Url)

    //local_location, err := time.LoadLocation("Asia/Taipei")
    /*if err != nil {
        fmt.Println(err)
    }*/
    NowTimeStr := time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05") //.In(local_location)
    //NowTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().In(local_location))
    //NowTime := time.Now().In(time.FixedZone("CST", 8*3600))
    NowTime := time.Now().Unix()
    
    /*str_array := strings.Split(fullData.ExpireAt, "UTC")
    var year, mon, day, hh, mm, ss int
    fmt.Sscanf(str_array[0], "%d-%d-%d %d:%d:%d ", &year, &mon, &day, &hh, &mm, &ss)
    time_string_to_parse := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d+00:00", year, mon, day, hh, mm, ss)*/
    //NowTime, _ := time.Parse("2006-01-02 15:04:05", NowTimeStr)
    url_create_time, _ := time.Parse("2006-01-02 15:04:05", fullData.ExpireAt)
    Url_ExTime := url_create_time.Unix()
    fmt.Println(NowTimeStr,"NOOOOOW",NowTime,"WWTTTIII",fullData.ExpireAt,"EEEXXX",Url_ExTime,"XTTTT")
    
    /*if err == nil && NowTime.Before(url_create_time) {
        //处理逻辑
        fmt.Println("true")
    }*/

    //if time.Since(url_create_time) < 0{
    if (Url_ExTime-NowTime) > 28800 {
        fmt.Println("return fullData.Url",Url_ExTime-NowTime)
        return fullData.Url
    }else{
        Re := fullData.ExpireAt + ":==:" + NowTimeStr + "overDue"
        fmt.Println(Re)
        return "overDue"
    }
} 

func SHOW_TABLES() {
    sql := `SHOW TABLES;`
    
        if _, err := MyDB.Exec(sql); err != nil {
            fmt.Println("SHOW TABLES failed:", err)
            return
        }
        fmt.Println("SHOW TABLES successd")
} 

func autoAdd(Url string, ExpireAt string) (string, string) {
	//sql := `insert INTO cumstomer(username,password) values('test','123456'); `
    result, err := MyDB.Exec("insert INTO urshoner(url,ExpireAt) values(?,?)", Url, ExpireAt)
    //result, err := MyDB.Exec(sql)
	if err != nil {
		fmt.Printf("Insert data failed,err:%v", err)
		return "fail", "fail"
	}
    //sql.Result 的LastInsertId()可取得AUTO_INCREMENT的值
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Get insert id failed,err:%v", err)
		return "fail", "fail"
	}
	fmt.Println("Insert data id:", lastInsertID)

    //RowsAffected() 影響的資料筆數，如果很嚴謹的寫法會判斷RowsAffected()是否與新增的資料筆數一致
	rowsaffected, err := result.RowsAffected() 
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return "fail", "fail"
	}
	fmt.Println("Affected rows:", rowsaffected)
    x := Encode(uint64(lastInsertID))
    y := "http://localhost/login1/" + x
    return x, y
}

func Encode(number uint64) string {
    var encodedBuilder strings.Builder
    encodedBuilder.Grow(11)
  
    for ; number > 0; number = number / length {
       encodedBuilder.WriteByte(alphabet[(number % length)])
    }
  
    return encodedBuilder.String()
}

func Decode(encoded string) (uint64, error) {
    var number uint64
  
    for i, symbol := range encoded {
       alphabeticPosition := strings.IndexRune(alphabet, symbol)
  
       if alphabeticPosition == -1 {
          return uint64(alphabeticPosition), errors.New("invalid character: " + string(symbol))
       }
       number += uint64(alphabeticPosition) * uint64(math.Pow(float64(length), float64(i)))
    }
  
    return number, nil
}