package alert

// Alerter is the interface all alerters should implement.
type Alerter interface {
	Alert(title, body string)
}
