/*******************************************************
 *  File        :   main.go
 *  Author      :   nieaowei
 *  Date        :   2020/1/25 10:10 下午
 *  Notes       :
 *******************************************************/
package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"io/ioutil"
	_ "log-agent-demo/config"
	"log-agent-demo/logs/instance"
	"log-agent-demo/logs/kafka"
	"log-agent-demo/logs/tailf"
	"path/filepath"
	"strings"
)

func main() {
	inst, err := kafka.NewKafkaInst(g.Cfg().GetString("kafka.Server"), kafka.DefaultConfig())
	if err != nil {
		g.Log().Fatal("kafka failed.")
		return
	}
	instance.InstMgr.Register(inst)
	tailfInit(inst)
	instance.InstMgr.RunAll()
	for {
		option := ' '
		fmt.Scanf("%c", &option)
		if option == 'q' {
			return
		}
	}
}

func kafkaInit() {

}

func tailfInit(inst *kafka.KafkaInst) {
	path := g.Cfg().GetString("tailf.path")
	suffix := g.Cfg().GetString("tailf.suffix")

	files, err := ioutil.ReadDir(path)
	if err != nil {
		g.Log().Fatal("scan log file failed.")
	}

	for _, file := range files {
		if file.IsDir() {
			files1, err := ioutil.ReadDir(filepath.Join(path, "/", file.Name()))
			if err != nil {
				continue
			}
			for _, file1 := range files1 {
				if strings.HasSuffix(file1.Name(), suffix) {
					tailfInst, err := tailf.NewTailfInst(filepath.Join(path, file.Name(), file1.Name()), tailf.DefaultConfig())
					if err != nil {
						g.Log().Warning("New inst faile.")
						continue
					}
					tailfInst.Name = file.Name()
					tailfInst.BindChan(inst.GetMsgChan())
					instance.InstMgr.Register(tailfInst)
				}
			}
		}
	}
}
