package manager

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Lesterpig/board/config"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {
	supportedResources map[string]func(ctx context.Context) <-chan []Ingress
	timeout            time.Duration
	namespace          string
	kubeClientSet      *kubernetes.Clientset
}

func NewKubeClient(cfg *config.KubeClientConfig) (*KubeClient, error) {
	kubeClientSet, err := getKubeClientSet(&cfg.Kubeconfig, &cfg.Kubecontext)
	if err != nil {
		return nil, err
	}

	client := &KubeClient{
		kubeClientSet: kubeClientSet,
		timeout:       cfg.Timeout,
		namespace:     cfg.Namespace,
	}

	supResources := make(map[string]func(ctx context.Context) <-chan []Ingress)
	supResources["ingress"] = client.fetchIngressResources

	client.supportedResources = supResources

	return client, nil
}

// Get a Kubernetes clientset
// First, it checks if the program is running inside a Kubernetes cluster
// If not, it tries to load a kubeconfig file
// If no kubeconfig file is specified, it tries to load the default kubeconfig file
func getKubeClientSet(kubeconfig, kubecontext *string) (*kubernetes.Clientset, error) {
	clientConfig, err := rest.InClusterConfig()
	if err != nil {
		if errors.Is(err, rest.ErrNotInCluster) {
			loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
			if *kubeconfig != "" {
				loadingRules.ExplicitPath = *kubeconfig
			}
			overrides := clientcmd.ConfigOverrides{
				CurrentContext: *kubecontext,
			}

			clientConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &overrides).ClientConfig()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return kubernetes.NewForConfig(clientConfig)
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

func (k *KubeClient) fetchIngressResources(ctx context.Context) <-chan []Ingress {
	ch := make(chan []Ingress)
	go func() {
		defer close(ch)
		select {
		case <-ctx.Done():
			return
		default:
			t := int64(k.timeout.Seconds())
			ingresses, err := k.kubeClientSet.NetworkingV1().Ingresses(k.namespace).List(context.TODO(), meta.ListOptions{TimeoutSeconds: &t})
			if err != nil {
				return
			}

			results := make([]Ingress, 0)
			for _, ingress := range ingresses.Items {
				tlsHosts := make([]string, 0)
				for _, tls := range ingress.Spec.TLS {
					tlsHosts = append(tlsHosts, tls.Hosts...)
				}

				for _, rule := range ingress.Spec.Rules {

					if rule.Host == "" {
						// If no host is specified, the rule applies to all inbound HTTP traffic through the ingress controller external IP
						// TODO This can be retrieved from the status of the Ingress resource
						continue
					}
					for _, path := range rule.HTTP.Paths {
						results = append(results, Ingress{
							name: ingress.Name,
							host: rule.Host,
							path: path.Path,
							tls:  contains(tlsHosts, rule.Host),
						})
					}
				}
			}

			ch <- results
		}

	}()

	return ch

}

// Check if a string is in an array of strings
func contains(array []string, contain string) bool {
	for _, s := range array {
		if s == contain {
			return true
		}
	}
	return false
}
