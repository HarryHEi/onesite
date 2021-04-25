package worker

import (
	"encoding/json"

	"onesite/core/dao"
	"onesite/core/worker/queue"
	"onesite/core/worker/tasks"
)

type Consumer interface {
	Run()
	AddHandler(func(string))
}

type DefaultConsumer struct {
	Topic    string
	Q        queue.Queue
	handlers []func(string)
}

func NewConsumer(topic string, q queue.Queue) *DefaultConsumer {
	return &DefaultConsumer{
		Topic: topic,
		Q:     q,
	}
}

func (t *DefaultConsumer) Run() {
	for {
		item := t.Q.GetTopic(t.Topic)
		for index := range t.handlers {
			t.handlers[index](item.(string))
		}
	}
}

func (t *DefaultConsumer) AddHandler(f func(string)) {
	t.handlers = append(t.handlers, f)
}

type Producer interface {
	ProduceTopic(topic string, v interface{})
}

type DefaultProducer struct {
	Q queue.Queue
}

func NewProducer(q queue.Queue) *DefaultProducer {
	return &DefaultProducer{
		q,
	}
}

func (p *DefaultProducer) ProduceTopic(topic string, v interface{}) {
	vStr, _ := json.Marshal(v)
	p.Q.PutTopic(topic, vStr)
}

type Worker struct {
	Q         queue.Queue
	Producer  Producer
	Consumers map[string]Consumer
}

func NewWorker(d *dao.Dao) *Worker {
	q := queue.NewRedisQueue(d.Redis.Db)
	return &Worker{
		Q:         q,
		Producer:  NewProducer(q),
		Consumers: make(map[string]Consumer),
	}
}

// Run 启动事件消费者
func (w *Worker) Run() {
	w.InitTasks()
	select {}
}

func (w *Worker) NewConsumer(topic string) Consumer {
	consumer, exists := w.Consumers[topic]
	if !exists {
		w.Consumers[topic] = NewConsumer(topic, w.Q)
		go func() {
			w.Consumers[topic].Run()
		}()
		return w.Consumers[topic]
	}
	return consumer
}

func (w *Worker) ProduceTopic(topic string, v string) {
	w.Producer.ProduceTopic(topic, v)
}

// InitTasks Create Consumer and register Handler
func (w *Worker) InitTasks() {
	// demo topic
	consumer := w.NewConsumer(tasks.DemoTopic)
	consumer.AddHandler(tasks.MakeGreeting)
}
