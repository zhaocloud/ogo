// Ogo

package ogo

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "reflect"
    "runtime"

    "github.com/VividCortex/godaemon"
    "github.com/nightlyone/lockfile"
    "github.com/zhaocloud/ogo/libs/config"
    "github.com/zhaocloud/ogo/libs/logs"
)

// ogo daemon framework version.
const VERSION = "0.1.0"

type Context struct {
    Env     *Environment           //环境参数
    Cfg     config.ConfigContainer //配置信息
    Workers map[string]*Worker
    Logger  *logs.OLogger //日志记录
}

var (
    Ctx       *Context
    AppConfig config.ConfigContainer
    Debugger  *logs.OLogger
)

// Run ogo application.
func Run() {
    defer func() {
        if err := recover(); err != nil {
            WriteMsg("App crashed with error:", err)
            for i := 1; ; i++ {
                _, file, line, ok := runtime.Caller(i)
                if !ok {
                    break
                }
                WriteMsg(file, line)
            }
            //panic要输出到console
            fmt.Println("App crashed with error:", err)
        }
    }()
    if Env.Daemonize {
        godaemon.MakeDaemon(&godaemon.DaemonAttr{})
    }
    //check&write pidfile, added by odin
    dir := filepath.Dir(Env.PidFile)
    if _, err := os.Stat(dir); err != nil {
        if os.IsNotExist(err) {
            //mkdir
            if err := os.Mkdir(dir, 0755); err != nil {
                panic(err)
            }
        }
    }
    if l, err := lockfile.New(Env.PidFile); err == nil {
        if le := l.TryLock(); le != nil {
            panic(le)
        }
    } else {
        panic(err)
    }

    var mainErr error
    //spawn worker,worker名字从配置文件来
    if workerName := AppConfig.String("Worker"); workerName != "" {
        Debugger.Debug("will run worker: %v", workerName)
        if worker, ok := Ctx.Workers[workerName]; ok {
            vw := reflect.New(worker.WorkerType)
            execWorker, ok := vw.Interface().(WorkerInterface)
            if !ok {
                panic("worker is not WorkerInterface")
            }

            //Init
            execWorker.Init(Ctx, workerName)

            //Main
            mainErr = execWorker.Main()
        }
    } else {
        mainErr = errors.New("not defined worker in cfg")
    }

    if mainErr != nil {
        Debugger.Critical("Main error: ", mainErr)
    }
}
