package main

import (
	"fmt"
	"github.com/pborman/getopt"
	"go.etcd.io/etcd/clientv3"
	"os"
	"path/filepath"
)

type Master struct {
	AppName    string
	config     MasterConfig
	logger     *Writer
	HostName   string
	RunDir     string
	SelfPid    int
	Watchkey   string
	EtcdClient *clientv3.Client
}

func NewMaster() *Master {
	return &Master{
		AppName:  "jrasp-master",
		HostName: GetHostname(),
		SelfPid:  os.Getpid(),
	}
}

func main() {

	var s = getopt.New()

	var help = false
	var logOutputStd = false
	var isProd = false

	s.BoolVarLong(&logOutputStd, "logoutput", 'l', "Print log to stdout")
	s.BoolVarLong(&help, "help", 'h', "Print this help")

	s.Parse(os.Args)

	if help {
		s.PrintUsage(os.Stderr)
		return
	}

	m := NewMaster()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		_, err = os.Stderr.WriteString(fmt.Sprintf("can't get pwd , %v", err))
		if err != nil {
			fmt.Printf("%v", err)
			return
		}
		return
	}

	// 日志初始化
	m.RunDir = dir
	m.logger = NewLog(dir, m.HostName, m.AppName, isProd, logOutputStd)

	// master配置初始化
	if err := m.InitConfig(); err != nil {
		return
	}

	// etcd客户端初始化
	if err := m.InitEtcdClient(); err != nil {
		return
	}

	//watch 的 key
	watchKey := fmt.Sprintf(JRASP_MASTER_WATCH_KEY, m.HostName, m.SelfPid)
	m.logger.Info(ETCD_INFO, "Current master watchKey=%s", watchKey)

	select {} //todo block
}
