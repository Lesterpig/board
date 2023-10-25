package manager

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type KubeClient struct {
	supportedResources map[string]func(ctx context.Context) <-chan []Ingress
}

func NewKubeClient() *KubeClient {
	supResources := make(map[string]func(ctx context.Context) <-chan []Ingress)
	supResources["ingress"] = fetchIngressResources

	return &KubeClient{
		supportedResources: supResources,
	}
}

// Ingress Resource with only the used fields
type Ingress struct {
	name string
	host string
	path string
	tls  bool
}

func (k *KubeClient) Fetch(kubeResource string) (func(ctx context.Context) <-chan []Ingress, error) {
	// Look up the resource in support resources
	fetcher, ok := k.supportedResources[strings.ToLower(kubeResource)]
	if !ok {
		return nil, fmt.Errorf("unsupported kubernetes resource")
	}
	return fetcher, nil
}

func fetchIngressResources(ctx context.Context) <-chan []Ingress {
	ch := make(chan []Ingress)
	fmt.Printf("fetchIngressResources Called")
	go func() {
		defer close(ch)
		select {
		case <-ctx.Done():
			return
		default:
			// simulate http request
			time.Sleep(1 * time.Second)

			// Mapping
			ingresses := make([]Ingress, 0)
			ingresses = append(ingresses, Ingress{
				name: "Amazon Web services",
				host: "aws.com",
				path: "/",
				tls:  true,
			})
			ch <- ingresses
		}

	}()

	return ch

}
