package manager

import (
	"context"
	"time"
)

type KubeClient struct {
}

func NewKubeClient() *KubeClient {
	return &KubeClient{}
}

// Ingress Resource with only the used fields
type Ingress struct {
	name string
	host string
	path string
	tls  bool
}

func (k *KubeClient) fetchIngress(ctx context.Context) <-chan []Ingress {
	ch := make(chan []Ingress)

	go func() {
		defer close(ch)
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(1 * time.Second) // simulate http request

			ingresses := make([]Ingress, 0)
			ingresses = append(ingresses, Ingress{
				name: "Amazon Weob services",
				host: "aws.com",
				path: "/",
				tls:  true,
			})
			ch <- ingresses
		}

	}()

	return ch

}
