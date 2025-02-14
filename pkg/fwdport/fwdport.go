//go:generate ${gopath}/bin/mockgen -source=fwdport.go -destination=mock_fwdport.go -package=fwdport
package fwdport

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/c6o/kubefwd/pkg/fwdhosts"
	"github.com/c6o/kubefwd/pkg/fwdpub"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// ServiceFWD PodSyncer interface is used to represent a
// fwdservice.ServiceFWD reference, which cannot be used directly
// due to circular imports.  It's a reference from a pod to it's
// parent service.
type ServiceFWD interface {
	String() string
	SyncPodForwards(bool)
	ListServicePodNames() []string
	AddServicePod(pfo *PortForwardOpts)
	GetServicePodPortForwards(servicePodName string) []*PortForwardOpts
	RemoveServicePod(servicePodName string, stop bool)
	RemoveServicePodByPort(servicePodName string, podPort string, stop bool)
}

type PortForwardHelper interface {
	GetPortForwardRequest(pfo *PortForwardOpts) *restclient.Request
	NewOnAddresses(dialer httpstream.Dialer, addresses []string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}, out, errOut io.Writer) (*portforward.PortForwarder, error)
	ForwardPorts(forwarder *portforward.PortForwarder) error

	RoundTripperFor(config *restclient.Config) (http.RoundTripper, spdy.Upgrader, error)
	NewDialer(upgrader spdy.Upgrader, client *http.Client, method string, pfRequest *restclient.Request) httpstream.Dialer
}

// Todo: fix tests to not use this, use Host Modifier.
type HostsOperator interface {
	AddHosts()
	RemoveHosts()
	RemoveInterfaceAlias()
}

type PortForwardHelperImpl struct{}

// Todo: fix tests to not use this, use Host Modifier.
type PortForwardOptsHostsOperator struct {
	Pfo *PortForwardOpts
}

// HostsParams
type HostsParams struct {
	localServiceName string
	nsServiceName    string
	fullServiceName  string
	svcServiceName   string
}

// PortForwardOpts
type PortForwardOpts struct {
	Out        *fwdpub.Publisher
	Config     restclient.Config
	ClientSet  kubernetes.Clientset
	RESTClient restclient.Interface

	Service    string
	ServiceFwd ServiceFWD
	PodName    string
	PodPort    string
	LocalIp    net.IP
	LocalPort  string
	HostFile   *fwdhosts.HostFileWithLock

	// Context is a unique key (string) in kubectl config representing
	// a user/cluster combination. Kubefwd uses context as the
	// cluster name when forwarding to more than one cluster.
	Context string

	// Namespace is the current Kubernetes Namespace to locate services
	// and the pods that back them for port-forwarding
	Namespace string

	// ClusterN is the ordinal index of the cluster (from configuration)
	// cluster 0 is considered local while > 0 is remote
	ClusterN int

	// NamespaceN is the ordinal index of the namespace from the
	// perspective of the user. Namespace 0 is considered local
	// while > 0 is an external namespace
	NamespaceN int

	Domain         string
	HostsParams    *HostsParams
	Hosts          []string
	ManualStopChan chan struct{} // Send a signal on this to stop the portforwarding
	DoneChan       chan struct{} // Listen on this channel for when the shutdown is completed.

	StateWaiter       PodStateWaiter
	PortForwardHelper PortForwardHelper
	HostModifier fwdhosts.HostModifierOpts

	Headless bool
}

type pingingDialer struct {
	wrappedDialer     httpstream.Dialer
	pingPeriod        time.Duration
	pingStopChan      chan struct{}
	pingTargetPodName string
}

func (p pingingDialer) Dial(protocols ...string) (httpstream.Connection, string, error) {
	streamConn, streamProtocolVersion, dialErr := p.wrappedDialer.Dial(protocols...)
	if dialErr != nil {
		log.Warnf("Ping process will not be performed for %s, cannot dial", p.pingTargetPodName)
	}
	go func(streamConnection httpstream.Connection) {
		if streamConnection == nil || dialErr != nil {
			return
		}
		for {
			select {
			case <-time.After(p.pingPeriod):
				if pingStream, err := streamConnection.CreateStream(nil); err == nil {
					_ = pingStream.Reset()
				}
			case <-p.pingStopChan:
				log.Debugf("Ping process stopped for %s", p.pingTargetPodName)
				return
			}
		}
	}(streamConn)

	return streamConn, streamProtocolVersion, dialErr
}

// PortForward does the port-forward for a single pod.
// It is a blocking call and will return when an error occurred
// or after a cancellation signal has been received.
func PortForward(pfo *PortForwardOpts) error {
	defer close(pfo.DoneChan)

	transport, upgrader, err := pfo.PortForwardHelper.RoundTripperFor(&pfo.Config)
	if err != nil {
		return err
	}

	// check that pod port can be strconv.ParseUint
	_, err = strconv.ParseUint(pfo.PodPort, 10, 32)
	if err != nil {
		pfo.PodPort = pfo.LocalPort
	}

	fwdPorts := []string{fmt.Sprintf("%s:%s", pfo.LocalPort, pfo.PodPort)}
	req := pfo.PortForwardHelper.GetPortForwardRequest(pfo)

	pfStopChannel := make(chan struct{}, 1)      // Signal that k8s forwarding takes as input for us to signal when to stop
	downstreamStopChannel := make(chan struct{}) // @TODO: can this be the same as pfStopChannel?
	podRestartStopChannel := make(chan struct{}) // kubefwd hangs on exit if a separate stop channel is not used.

	localNamedEndPoint := fmt.Sprintf("%s:%s", pfo.Service, pfo.LocalPort)

	if pfo.Headless || pfo.getBrothersInPodsAmount() == 1 {
		pfo.HostModifier.AddHosts()
	}

	// Close created downstream channels if there are stop signal from above
	go func() {
		<-pfo.ManualStopChan
		close(downstreamStopChannel)
		close(podRestartStopChannel)
		close(pfStopChannel)
	}()

	// Waiting until the pod is running
	pod, err := pfo.StateWaiter.WaitUntilPodRunning(downstreamStopChannel)
	if err != nil {
		pfo.stopAndShutdown()
		return err
	} else if pod == nil {
		// if err is not nil but pod is nil
		// mean service deleted but pod is not runnning.
		// No error, just return
		pfo.stopAndShutdown()
		return nil
	}

	// Listen for pod is deleted
	go pfo.StateWaiter.ListenUntilPodDeleted(podRestartStopChannel, pod)

	p := pfo.Out.MakeProducer(localNamedEndPoint)

	dialer := pfo.PortForwardHelper.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, req)

	dialerWithPing := pingingDialer{
		wrappedDialer:     dialer,
		pingPeriod:        time.Second * 30,
		pingStopChan:      make(chan struct{}),
		pingTargetPodName: pfo.PodName,
	}

	var address []string
	if pfo.LocalIp != nil {
		address = []string{pfo.LocalIp.To4().String(), pfo.LocalIp.To16().String()}
	} else {
		address = []string{"localhost"}
	}

	fw, err := pfo.PortForwardHelper.NewOnAddresses(dialerWithPing, address, fwdPorts, pfStopChannel, make(chan struct{}), &p, &p)
	if err != nil {
		pfo.stopAndShutdown()
		return err
	}

	// Blocking call
	if err = pfo.PortForwardHelper.ForwardPorts(fw); err != nil {
		log.Errorf("ForwardPorts error: %s", err.Error())
		pfo.shutdown()
		return err
	} else {
		pfo.shutdown()
		//pfo.Stop() // Don't shut down, this gives connected clients time to move to a new pod.
	}

	return nil
}

// shutdown removes port-forward from ServiceFwd and removes hosts entries if it's necessary
func (pfo PortForwardOpts) shutdown() {
	pfo.ServiceFwd.RemoveServicePodByPort(pfo.String(), pfo.PodPort, true)
	if pfo.Headless || pfo.getBrothersInPodsAmount() == 0 {
		pfo.HostModifier.RemoveHosts()
	}
	pfo.HostModifier.RemoveInterfaceAlias()
}

// stopAndShutdown is shortcut for closing all downstream channels and shutdown
func (pfo PortForwardOpts) stopAndShutdown() {
	pfo.Stop()
	pfo.shutdown()
}

//// BuildHostsParams constructs the basic hostnames for the service
//// based on the PortForwardOpts configuration
//func (pfo *PortForwardOpts) BuildHostsParams() {
//
//	localServiceName := pfo.Service
//	nsServiceName := pfo.Service + "." + pfo.Namespace
//	fullServiceName := fmt.Sprintf("%s.%s.svc.cluster.local", pfo.Service, pfo.Namespace)
//	svcServiceName := fmt.Sprintf("%s.%s.svc", pfo.Service, pfo.Namespace)
//
//	// check if this is an additional cluster (remote from the
//	// perspective of the user / argument order)
//	if pfo.ClusterN > 0 {
//		fullServiceName = fmt.Sprintf("%s.%s.svc.cluster.%s", pfo.Service, pfo.Namespace, pfo.Context)
//	}
//	pfo.HostsParams.localServiceName = localServiceName
//	pfo.HostsParams.nsServiceName = nsServiceName
//	pfo.HostsParams.fullServiceName = fullServiceName
//	pfo.HostsParams.svcServiceName = svcServiceName
//}

// getBrothersInPodsAmount returns amount of port-forwards that proceeds on different ports under same pod
func (pfo *PortForwardOpts) getBrothersInPodsAmount() int {
	return len(pfo.ServiceFwd.GetServicePodPortForwards(pfo.String()))
}

// WaitUntilPodRunning Waiting for the pod running
func (waiter *PodStateWaiterImpl) WaitUntilPodRunning(stopChannel <-chan struct{}) (*v1.Pod, error) {
	pod, err := waiter.ClientSet.CoreV1().Pods(waiter.Namespace).Get(context.TODO(), waiter.PodName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	if pod.Status.Phase == v1.PodRunning {
		return pod, nil
	}

	watcher, err := waiter.ClientSet.CoreV1().Pods(waiter.Namespace).Watch(
		context.TODO(), metav1.SingleObject(metav1.ObjectMeta{Name: pod.ObjectMeta.Name, Namespace: pod.ObjectMeta.Namespace}))
	if err != nil {
		return nil, err
	}

	// if the os.signal (we enter the Ctrl+C)
	// or ManualStop (service delete or some thing wrong)
	// or RunningChannel channel (the watch for pod runnings is done)
	// or timeout after 300s
	// we'll stop the watcher
	// TODO: change the 300s timeout to custom settings.
	go func() {
		defer watcher.Stop()
		select {
		case <-stopChannel:
		case <-time.After(time.Second * 300):
		}
	}()

	// watcher until the pod status is running
	for {
		event, ok := <-watcher.ResultChan()
		if !ok {
			break
		}
		if event.Object != nil && event.Type == "MODIFIED" {
			changedPod := event.Object.(*v1.Pod)
			if changedPod.Status.Phase == v1.PodRunning {
				return changedPod, nil
			}
		}
	}
	return nil, nil
}

// listen for forwarded pod modification or deletion
// todo: Anticipate a dying pod and sync before it stops serving traffic.
func (waiter *PodStateWaiterImpl) ListenUntilPodDeleted(stopChannel <-chan struct{}, pod *v1.Pod) {

	watcher, err := waiter.ClientSet.CoreV1().Pods(waiter.Namespace).Watch(
		context.TODO(), metav1.SingleObject(metav1.ObjectMeta{Name: pod.ObjectMeta.Name, Namespace: pod.ObjectMeta.Namespace}))
	if err != nil {
		return
	}

	// Listen for stop signal from above
	go func() {
		<-stopChannel
		watcher.Stop()
	}()

	// watcher until the pod is deleted, then trigger a syncpodforwards
	for {
		event, ok := <-watcher.ResultChan()
		if !ok {
			return
		}
		switch event.Type {
		case watch.Deleted:
			log.Warnf("Pod %s deleted, resyncing the %s service pods.", pod.ObjectMeta.Name, waiter.ServiceFwd)
			go waiter.ServiceFwd.SyncPodForwards(true)
			break
		}
	}
}

// Stop sends the shutdown signal to the port-forwarding process.
// In case the shutdown signal was already given before, this is a no-op.
func (pfo *PortForwardOpts) Stop() {
	select {
	case <-pfo.DoneChan:
		return
	case <-pfo.ManualStopChan:
		return
	default:
	}
	close(pfo.ManualStopChan)
}

func (pfo *PortForwardOpts) String() string {
	return pfo.PodName
}

type PodStateWaiter interface {
	WaitUntilPodRunning(stopChannel <-chan struct{}) (*v1.Pod, error)
	ListenUntilPodDeleted(stopChannel <-chan struct{}, pod *v1.Pod)
}

type PodStateWaiterImpl struct {
	Namespace  string
	PodName    string
	ClientSet  kubernetes.Clientset
	ServiceFwd ServiceFWD
}

func (p PortForwardHelperImpl) GetPortForwardRequest(pfo *PortForwardOpts) *restclient.Request {
	// if need to set timeout, set it here.
	// restClient.Client.Timeout = 32
	return pfo.RESTClient.Post().
		Resource("pods").
		Namespace(pfo.Namespace).
		Name(pfo.PodName).
		SubResource("portforward")
}

func (p PortForwardHelperImpl) NewOnAddresses(dialer httpstream.Dialer, addresses []string, ports []string, stopChan <-chan struct{}, readyChan chan struct{}, out, errOut io.Writer) (*portforward.PortForwarder, error) {
	return portforward.NewOnAddresses(dialer, addresses, ports, stopChan, readyChan, out, errOut)
}

func (p PortForwardHelperImpl) RoundTripperFor(config *restclient.Config) (http.RoundTripper, spdy.Upgrader, error) {
	return spdy.RoundTripperFor(config)
}

func (p PortForwardHelperImpl) NewDialer(upgrader spdy.Upgrader, client *http.Client, method string, pfRequest *restclient.Request) httpstream.Dialer {
	return spdy.NewDialer(upgrader, client, method, pfRequest.URL())
}

func (p PortForwardHelperImpl) ForwardPorts(forwarder *portforward.PortForwarder) error {
	return forwarder.ForwardPorts()
}
