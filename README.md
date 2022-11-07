# go-keyboard
基于 [robotogo](https://github.com/go-vgo/robotgo) 实现的Go语言循环按键脚本

## compiler for windows on macOS
```shell
# 安装
brew install mingw-w64

# 跨平台编译
GOOS=windows \
GOARCH=amd64 \
CGO_ENABLED=1 \
CC=x86_64-w64-mingw32-gcc \
CXX=x86_64-w64-mingw32-g++ \
go build -x -o ./go-keyboard.exe
```

## config.yaml配置说明

`keys_press`参考: https://github.com/go-vgo/robotgo/blob/master/docs/keys.md

```yaml
app:
  name: "go-keyboard"
  namespace: "test"
  debug_event: true # 是否仅调试跟踪模式，此模式仅会打印键盘、鼠标事件到日志输入
log:
  log_path: ./logs/go_keyboard.log
  debug: true # 日志debug若为true，则直接将信息输出到标准输出
keyboard:
  - name: "loop press number3(ctrl+shift+9/0)"
    process_name: "iterm"
    key_press: "num3" # 循环持续按键 3
    key_sleep: 100 # 模拟按键间隔，单位ms
    start_stop_key_diff: true # 如果启、停按键有差异，则需要分别配置start_key和stop_key，一致则仅需要配置start_key即可
    start_key: [ "9", "ctrl", "shift" ]
    stop_key: [ "0", "ctrl", "shift" ]
  - name: "F12"
    process_name: "goland"
    key_press: "num4"
    key_sleep: 500
    start_stop_key_diff: false 
    start_key: [ "+", "ctrl" ]
```

## 同步文件
可以通过windows和mac局域网共享，同步编译生成的二进制和代码文件测试: `访达-前往-连接服务器(cmd+k)`

## 运行
```shell
./go-keyboard
```

## 问题
1. 针对windows部分客户端无响应，通过开启debug后发现robotgo事件未响应，初步断点是客户端做了键盘、鼠标模拟（待进一步分析）
2. 目前功能比较有限，仅支持单一按键配置