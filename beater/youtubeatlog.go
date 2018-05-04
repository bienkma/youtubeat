package beater

import (
	"fmt"
	"net/http"
	"time"
	"encoding/json"

	"github.com/bienkma/youtubeat/config"
	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/publisher/bc/publisher"
	"github.com/bienkma/youtubeat/handler"

	"github.com/Shopify/sarama"
	kafkaCluster "gopkg.in/bsm/sarama-cluster.v2"
	"google.golang.org/api/youtube/v3"
	"google.golang.org/api/googleapi/transport"
)

type YoutubeatLog struct {
	done     chan struct{}
	config   config.Config
	client   publisher.Client
	consumer *kafkaCluster.Consumer
	service  *youtube.Service
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	kafkaConfig := kafkaCluster.NewConfig()
	kafkaConfig.Consumer.Return.Errors = true
	kafkaConfig.Group.Return.Notifications = true
	kafkaConfig.Config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client := &http.Client{
		Transport: &transport.APIKey{Key: c.YoutubeApiKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		fmt.Errorf("Error creating new YouTube client: %v", err)
	}

	consumer, err := kafkaCluster.NewConsumer(c.Brokers, c.Group, c.Topics, kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("Error consumer topic: %s, erorrs is: %s", c.Topics, err)
	}

	bt := &YoutubeatLog{
		done:     make(chan struct{}),
		config:   c,
		consumer: consumer,
		service:  service,
	}
	return bt, nil
}

func (bt *YoutubeatLog) Run(b *beat.Beat) error {
	logp.Info("youtubeat is running! Hit CTRL-C to stop it.")

	bt.client = b.Publisher.Connect()

	for {
		select {
		case <-bt.done:
			bt.consumer.Close()
			return nil
		case ev := <-bt.consumer.Messages():
			tmp := handler.SearchKeyWord(bt.service, string(ev.Value))
			tmp2, _ := json.Marshal(tmp)
			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       b.Info.Name,
				"message":    string(tmp2),
			}
			if bt.client.PublishEvent(event, publisher.Guaranteed, publisher.Sync) {
				bt.consumer.MarkOffset(ev, "")
			}
		case notification := <-bt.consumer.Notifications():
			logp.Info("Rebalanced: %+v", notification)
		case err := <-bt.consumer.Errors():
			logp.Err("Error in Kafka consumer: %s", err.Error())
		}
	}
}

func (bt *YoutubeatLog) Stop() {
	bt.client.Close()
	close(bt.done)
}
