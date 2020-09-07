package etcd

import (
	"fmt"
	"github.com/NullExceptionNull/go_structrue/conf"
	"go.etcd.io/etcd/clientv3"
	"gopkg.in/ini.v1"
	"time"
)

var etcdClient *clientv3.Client
var etcdConfig = new(conf.EtcdConf)

func init() {
	initProps()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdConfig.Address,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		fmt.Println("etcd 初始化失败")
	}
	etcdClient = cli
	defer cli.Close()
}

func initProps() *conf.EtcdConf {
	load, err := ini.Load("../conf/conf.ini")
	if err != nil {
		fmt.Println("init 文件加载失败")
	}
	err = load.Section("etcd").MapTo(etcdConfig)

	if err != nil {
		fmt.Println("init 文件解析失败")
	}
	return etcdConfig
}
