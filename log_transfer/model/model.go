package model

type Config struct {
	KafkaConfig KafkaConfig `ini:"kafka"`
	ESConfig    ESConfig    `ini:"es"`
}
type KafkaConfig struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}
type ESConfig struct {
	Address string `ini:"address"`
	Index   string `ini:"index"`
	MaxSize int    `ini:"max_chan_size"`
	GoroutineNum int `ini:"goroutine_num"`
}
