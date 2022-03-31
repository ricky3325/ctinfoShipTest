# 嘉堂測試專用
### 啟動服務
```
docker-compose up -d
```
### 第一次使用<目前無需此操作>
連接資料庫: <http://localhost/myproject?username=admin&password=0>  
建立資料表：<http://localhost/myproject?username=admin&password=1>

### 每次開啟<每次啟動服務後，都要連接資料庫>
連接資料庫:<http://localhost/myproject?username=admin&password=0>  

### RUN Unit Tests
```
docker exec -it myproject /bin/sh
go test -v
```

### API 操作

----
##### POST http://localhost/login3

**Request methods**

| Request methods/headers | Value |
| ------------- | ------------------------------ |
| Method      | POST       |
| Content-Type   | Text/Plain     |

**Request parameters**

| Parameter name | Required/optional | Type | Description |
| --------- | ------------ |------ |------------ |
| id      | Required    |	string    |要輸入的網址    |

**Response**

| Response header | Value |
| ------------- | ------------------------------ |
| Status  | 200: Success       |
| Content-Type   | Application/Json     |

----

##### GET http://localhost/login3

**Request methods**

| Request methods/headers | Value |
| ------------- | ------------------------------ |
| Method      | GET       |


**Response**

| Response header | Value |
| ------------- | ------------------------------ |
| Status  | 200: Success<br>401: Fail|
| Content-Type   | Text/Plain     |

**Response body**

| 狀態 |動作|
| ------------- | ------------------------------ |
| 200  | 顯示ID|


