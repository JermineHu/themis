# themis
用于监控和安保程序

## 实现方案介绍

### 关于agent服务的注册
pc客户端运行后，弹出填写注册服务器地址和token，对话框，服务器验证token后，pc端程序将自己的主机信息，mac，主机名称，ip等，发给服务器，服务端在加rtsp的时候，下拉选择对应的主机。这样就能将rtsp与键盘监听服务关联，监听服务作为一个agent服务部署于对应的windows主机上，每次监听到键盘数据就将数据通过webassembly进行发送到后台，后台收到数据就存储数据和时间，如果同时有人进入websoket，就将接收数据显示到页面!

### 实时显示键盘的数据
根据服务端所添加的rtsp数据ID作为参数建立ws，摄像头与键盘采用一对一的关系。

### 隐含的实现细节
1、PC端的token管理（管理后台生成，类似GitHub的token机制），用于配置远程服务，不是所有的请求后台都收；
2、主机管理，PC端agent注册后会将自身信息注册上来；
3、新增rtsp的时候可以直接绑定对应主机，就能完成鼠标监听的数据同步

## 采用的技术方案

### 后端技术方案：

- 数据库采用postgresql
- orm采用gorm
- web框架采用goa
- 转码工具采用ffmpeg

### 前端技术方案：

- angular 8.0
- ant design

### 系统托盘

Linux下需要先安装gtk3.0，否则程序无法编译

```
sudo apt-get install libgtk-3-dev libappindicator3-dev -y
```

使用的开源库如下：
https://github.com/getlantern/systray

文档可以参照：

https://pkg.go.dev/github.com/getlantern/systray?tab=doc

### 监听键盘事件的库

https://github.com/go-vgo/robotgo

### 图形化采用

https://github.com/fyne-io/fyne

### 客户端交叉编译

GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui

## 部署方案

### docker方式

#### 启动postgres数据库

首先执行：
``
docker volume create pgdata
``
正常启动：
```
docker run --restart=always --name postgres -d \
-v pgdata:/var/lib/postgresql/data  \
-e POSTGRES_USER=Jermine \
-e POSTGRES_PASSWORD=123456 \  
-e POSTGRES_DB=mean \
-e POSTGRES_INITDB_ARGS="--data-checksums"  \
-p 5432:5432 \
postgres:alpine
```
以上-e参数为变量，可以自行设置。

#### 启动Themis server端
```
sudo docker run --name themis --restart always -d \
-e DB_TYPE=postgres \
-e DB_CON_STR="sslmode=disable host=192.168.1.250 port=5432 user=Jermine dbname=%v password=123456" \
-e DB_IS_UPGRADE=true \
-e APP_PORT=8080 \
-e DOMAIN="https://jermine.vdo.pub" \
-e TOKEN_TIMEOUT=5184000 \
-p 8080:8080 \
-p 8081:8081 \
registry.cn-hangzhou.aliyuncs.com/vdo/themis:release-v1.0.4
```

以上-e参数为变量，可以自行设置，DB_CON_STR是数据库连接字符串，必填。

#### 启动Themis UI端

```
sudo docker run --name themis-ui --restart always -d  \
-e BASE_API="192.168.1.250:8081" \
-p 8090:80 \
registry.cn-hangzhou.aliyuncs.com/vdo/themis-ui:v1.0.4
```
BASE_API是填写服务端访问地址，对外映射端口可以更换8090
#### windows下启动Themis agent端

安装包安装后，填写后台server地址和使用后台生成的token，如：server：192.168.1.250：8081 (可以不加http) token：56d7340d-7bbc-4847-b5de-2ec7c9ea3b4e
。如果有其他需要更改的配置，找到安装目录，编辑config.toml文件即可。

### 单机运行命令

```
DB_TYPE=postgres \
DB_CON_STR="sslmode=disable host=192.168.1.250 port=5432 user=Jermine dbname=%v password=123456" \
DB_IS_UPGRADE=true \
APP_PORT=8080 \
TOKEN_TIMEOUT=5184000 ./themis
```