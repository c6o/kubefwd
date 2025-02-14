//go:generate ${gopath}/bin/mockgen -source=fwdhosts.go -destination=mock_fwdhosts.go -package=fwdhosts
package fwdhosts

import (
    "fmt"
    "github.com/c6o/kubefwd/pkg/fwdnet"
    log "github.com/sirupsen/logrus"
    "github.com/txn2/txeh"
    "net"
    "sync"
)

// HostFileWithLock
type HostFileWithLock struct {
    Hosts *txeh.Hosts
    sync.Mutex
}

type HostModifierOpts struct {
    NamespaceN int
    ClusterN int
    Service string
    Domain string
    Namespace string
    Context string
    HostFile *HostFileWithLock
    LocalIp net.IP
    Hosts []string
}

func (params *HostModifierOpts) AddHosts() {

    params.HostFile.Lock()
    // pfo.Service holds only the service name
    // start with the smallest allowable hostname

    // bare service name
    if params.ClusterN == 0 && params.NamespaceN == 0 {
        params.addHost(params.Service)
        if params.Domain != "" {
            params.addHost(fmt.Sprintf(
                "%s.%s",
                params.Service,
                params.Domain,
            ))
        }
    }

    // alternate cluster / first namespace
    if params.ClusterN > 0 && params.NamespaceN == 0 {
       params.addHost(fmt.Sprintf(
           "%s.%s",
           params.Service,
           params.Context,
       ))
    }

    // namespaced without cluster
    if params.ClusterN == 0 {
       params.addHost(fmt.Sprintf(
           "%s.%s",
           params.Service,
           params.Namespace,
       ))

       params.addHost(fmt.Sprintf(
           "%s.%s.svc",
           params.Service,
           params.Namespace,
       ))

       params.addHost(fmt.Sprintf(
           "%s.%s.svc.cluster.local",
           params.Service,
           params.Namespace,
       ))

       if params.Domain != "" {
           params.addHost(fmt.Sprintf(
               "%s.%s.svc.cluster.%s",
               params.Service,
               params.Namespace,
               params.Domain,
           ))
       }
    }

    params.addHost(fmt.Sprintf(
       "%s.%s.%s",
       params.Service,
       params.Namespace,
       params.Context,
    ))

    params.addHost(fmt.Sprintf(
       "%s.%s.svc.%s",
       params.Service,
       params.Namespace,
       params.Context,
    ))

    params.addHost(fmt.Sprintf(
       "%s.%s.svc.cluster.%s",
       params.Service,
       params.Namespace,
       params.Context,
    ))

    err := params.HostFile.Hosts.Save()
    if err != nil {
        log.Error("Error saving hosts file", err)
    }
    params.HostFile.Unlock()
}

// RemoveHosts removes hosts /etc/hosts  associated with a forwarded pod
func (params *HostModifierOpts) RemoveHosts() {
    // we should lock the pfo.HostFile here
    // because sometimes other goroutine write the *txeh.Hosts
    params.HostFile.Lock()
    // other applications or process may have written to /etc/hosts
    // since it was originally updated.
    err := params.HostFile.Hosts.Reload()
    if err != nil {
        log.Errorf("Unable to reload /etc/hosts: %s", err.Error())
        return
    }

    log.Debugf("Removing hostfile entries... %v", params.Hosts)
    // remove all hosts
    for _, host := range params.Hosts {
        log.Debugf("Removing host %s in namespace %s from context %s", host, params.Namespace, params.Context)
        params.HostFile.Hosts.RemoveHost(host)
    }

    err = params.HostFile.Hosts.Save()
    if err != nil {
        log.Errorf("Error saving /etc/hosts: %s\n", err.Error())
    }
    params.HostFile.Unlock()
}

func (params *HostModifierOpts) RemoveInterfaceAlias() {
    fwdnet.RemoveInterfaceAlias(params.LocalIp)
}

func (params *HostModifierOpts) addHost(host string) {
    // add to list of hostnames for this port-forward
    params.Hosts = append(params.Hosts, host)

    // remove host if it already exists in /etc/hosts
    params.HostFile.Hosts.RemoveHost(host)

    // add host to /etc/hosts
    params.HostFile.Hosts.AddHost(params.LocalIp.String(), host)
}
