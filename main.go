package main

import (  

    "context"
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "log"

    "github.com/etcd-io/etcd/clientv3"
    "time"

)

var (  
    dialTimeout    = 2 * time.Second
    requestTimeout = 10 * time.Second
)

func main() {  
    ctx, _ := context.WithTimeout(context.Background(), requestTimeout)
    cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client-key.pem")
    caCert, err := ioutil.ReadFile("certs/ca.pem")
    caCertPool := x509.NewCertPool()

    if err != nil {
        log.Fatalf("ERR: %v\n",err)
    }

    caCertPool.AppendCertsFromPEM(caCert)


    // Setup HTTPS client
    tlsConfig := &tls.Config{
        Certificates: []tls.Certificate{cert},
        RootCAs:      caCertPool,
    }
    //tlsConfig.BuildNameToCertificate()


    cli, _ := clientv3.New(clientv3.Config{
        DialTimeout: dialTimeout,
        Endpoints: []string{"etcd.cwxstat.io:2379"},
        TLS: tlsConfig,
    })
    defer cli.Close()
    kv := clientv3.NewKV(cli)

    GetSingleValueDemo(ctx, kv)

}

func GetSingleValueDemo(ctx context.Context, kv clientv3.KV) {
    fmt.Println("*** GetSingleValueDemo()")

    // Insert a key value
    pr, _ := kv.Put(ctx, "slop", "bob")
    rev := pr.Header.Revision
    fmt.Println("Revision:", rev)

}