package ogo

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"

    "github.com/zhaocloud/ogo/libs/config"
    "github.com/zhaocloud/ogo/libs/logs"
    "github.com/zhaocloud/ogo/utils"
)

type Environment struct {
    WorkPath      string // working path(abs)
    AppPath       string // application path
    ProcName      string // proc name
    AppConfigPath string // config file path
    RunMode       string // run mode, "dev" or "prod"
    Daemonize     bool   // daemonize or not
    DebugLevel    int    // debug level
    PidFile       string // pidfile abs path
}

var (
    Env *Environment
)

func init() { //初始化环境变量,配置信息
    workPath, _ := os.Getwd()
    Env.WorkPath, _ = filepath.Abs(workPath)
    Env.AppPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
    Env.ProcName = filepath.Base(os.Args[0])

    //默认配置文件是 conf/{ProcName}.conf
    Env.AppConfigPath = filepath.Join(Env.AppPath, "conf", Env.ProcName+".conf")
    if !utils.FileExists(Env.AppConfigPath) {
        //不存在时指定为app.conf
        Env.AppConfigPath = filepath.Join(Env.AppPath, "conf", "app.conf")
    }

    if Env.WorkPath != Env.AppPath {
        if utils.FileExists(Env.AppConfigPath) {
            os.Chdir(Env.AppPath)
        } else {
            //在当前目录找配置文件
            Env.AppConfigPath = filepath.Join(Env.WorkPath, "conf", Env.ProcName+".conf")
            if !utils.FileExists(Env.AppConfigPath) {
                //不存在时指定为app.conf
                Env.AppConfigPath = filepath.Join(Env.WorkPath, "conf", "app.conf")
            }
        }
    }

    Env.RunMode = "dev" //default runmod

    //added by odin
    Env.Daemonize = false
    Env.PidFile = filepath.Join(Env.AppPath, "run", Env.ProcName+".pid")
    Env.DebugLevel = logs.LevelTrace //默认debug等级

    runtime.GOMAXPROCS(runtime.NumCPU())

    configErr := ParseConfig()

    // init Debugger
    Debugger = logs.NewLogger(2046)
    var err error
    if Env.Daemonize {
        err = Debugger.SetLogger("file", `{"filename":"logs/debug.log"}`)
    } else {
        err = Debugger.SetLogger("console", "")
    }
    if err != nil {
        fmt.Println("init logger error:", err)
    }
    Debugger.EnableFuncCallDepth(true)
    Debugger.SetLogFuncCallDepth(3)
    Debugger.SetLevel(Env.DebugLevel)

    if configErr != nil && !os.IsNotExist(configErr) {
        //放在这里才能使用Logger函数
        Debugger.Log("Info", "%v", err)
    }

    //Ctx
    Ctx.Env = Env
    Ctx.Cfg = AppConfig
    Ctx.Logger = Debugger
}

// ParseConfig parsed default config file.
func ParseConfig() (err error) {
    AppConfig, err = config.NewConfig("ini", Env.AppConfigPath)
    if err != nil {
        AppConfig = config.NewFakeConfig()
        return err
    } else {

        if runmode := AppConfig.String("RunMode"); runmode != "" {
            Env.RunMode = runmode
        }

        //added by odin
        if daemonize, err := AppConfig.Bool("Daemonize"); err == nil {
            Env.Daemonize = daemonize
        }
        if pidfile := AppConfig.String("PidFile"); pidfile != "" {
            // make sure pidfile is abs path
            if filepath.IsAbs(pidfile) {
                Env.PidFile = pidfile
            } else {
                Env.PidFile = filepath.Join(Env.AppPath, pidfile)
            }
        }
        if level, err := AppConfig.Int("DebugLevel"); err == nil {
            Env.DebugLevel = level
        }

    }
    return nil
}
