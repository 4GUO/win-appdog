## win-appdog

### 用途
    用来远程管理window进程的开闭，无论是局域网还是公网都可以
### 配置说明
    MainProcessName : zhgAppDog.exe # 主程序进程名
    ActMqttTopic : zhg/appdog/act # 监听的动作指令主题
    StatusMqttTopic : zhg/appdog/status # 状态上报的主题
    ProcessStatusCheckRate : 60 # 状态上报的速度
    Mqtt: # MQTT配置 依赖MQTT协议
        Host : mqtt://127.0.0.1:19007
        ClientId : zhg-app-dog
        Username : zhg1
        Password : Zy3K6PkSpGf43
        CleanSession : true
    Apps:
      -
        Name : sublime 测试程序 
        ProcessName : Code.exe # 进程名称
        ShortcutName : D:\Coding\project\appDog\shortcut\tc3 # 快捷方式路径
        UniqueId : 1.0001 # 进程对于本程序的逻辑Id
### 用法
    主题 zhg/appdog/act 发送json {"UniqueId": "1.0002","Act":1} 打开进程
    主题 zhg/appdog/act 发送json {"UniqueId": "1.0002","Act":0} 关闭进程
    主题 zhg/appdog/status {"UniqueId":"1.0002","Status":true} 进程状态
    
### 打包
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-H windowsgui  -w -s"

### 技术点
    chatgpt 70%
    golang
    mqtt


        