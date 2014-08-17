## ogo

Odin's go daemon framework

## Original Features

* Daemonize, 可在<appname>.conf中用 Daemonize={bool}配置, pidfile默认写到程序目录的run/<appname>.pid
* DebugLevel配置,<appname>.conf中 DebugLevel={int}配置,数字越高级别越高

## 使用
* main.go如下:

```
package main

import (
    "github.com/zhaocloud/ogo"
    _ "<gopath_to_workers>/workers"
)

func main() {
    ogo.Run()
}
```

* workers目录下,文件名自定,以'test'为例,代码如下:

```
package workers

import (
    "github.com/zhaocloud/ogo"
)

type TestWorker struct {
    ogo.Worker
}

func init() {
    ogo.AddWorker(&TestWorker{})
}

func (w *TestWorker) Main() error {
    return nil
}

```

主要看worker的名字, 框架会自动调用Main()

