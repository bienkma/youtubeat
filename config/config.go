// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

type Config struct {
	Brokers []string `config:"brokers"`
	Topics []string `config:"topics"`
	Group string `config:"group"`
	YoutubeApiKey string `config:"key"`
}

var DefaultConfig = Config{
	Brokers: []string{"127.0.0.1:9092"},
	Topics: []string{"youtube"},
	Group: "youtubeat",
	YoutubeApiKey:"AIzaSyDRmG2dKkjEHTt9R5h3XAS2OPE4ZjaUivQ",
}
