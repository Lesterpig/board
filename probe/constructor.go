package probe

type Config struct {
	Name     string
	Config   ProberConfig
	Category string
	Type     string
}

var ProbeConstructors = map[string](func() Prober){
	"dns":  func() Prober { return &DNS{} },
	"http": func() Prober { return &HTTP{} },
	"port": func() Prober { return &Port{} },
	"smtp": func() Prober { return &SMTP{} },
}
