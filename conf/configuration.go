package conf

type KafkaConf struct {
	Address  []string `ini:"address"`
	ChanSize int      `init:"chanSize"`
}

type EtcdConf struct {
	Address     []string `ini:"address"`
	DialTimeout int      `ini:"dialTimeout"`
	Key         string   `ini:"key"`
	ChanSize    int      `ini:"chanSize"`
	TaskNum     int      `ini:"taskNum"`
}
