package main

import (
	"net/http"
	"net/url"
	"strings"

	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
	controllerruntime "sigs.k8s.io/controller-runtime"
)

func kubeproxy(w http.ResponseWriter, r *http.Request) {
	kconfig := controllerruntime.GetConfigOrDie()
	transport, err := rest.TransportFor(kconfig)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	proxyUrl := *r.URL
	u, err := url.Parse(kconfig.Host)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	proxyUrl.Host = u.Host
	proxyUrl.Scheme = u.Scheme
	proxyUrl.Path = strings.TrimPrefix(proxyUrl.EscapedPath(), "/kube-proxy")

	kubeProxy := proxy.NewUpgradeAwareHandler(&proxyUrl, transport, true, false, nil)
	kubeProxy.UpgradeTransport = proxy.NewUpgradeRequestRoundTripper(transport, transport)
	kubeProxy.ServeHTTP(w, r)
}

// 运行
// go run main.go
// 1. 查看所有可用接口
// curl http://localhost:8888/kube-proxy/
// 2. 查看所有命名空间
// curl http://localhost:8888/kube-proxy/api/v1/namespaces

func main() {
	http.HandleFunc("/kube-proxy/", kubeproxy)
	http.ListenAndServe(":8888", nil)
}
