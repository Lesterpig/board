package alert

type AlertConfig struct {
	Type    string
	Token   string
	Webhook string
	Channel string
}

var AlertConstructors = map[string](func(c AlertConfig) Alerter){
	"pushbullet": func(c AlertConfig) Alerter {
		return NewPushbullet(c.Token)
	},
	"slack": func(c AlertConfig) Alerter {
		return NewSlack(c.Webhook, c.Channel)
	},
}
