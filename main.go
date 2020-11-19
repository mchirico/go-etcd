package main

import (
    "context"
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "github.com/etcd-io/etcd/clientv3"
    "io/ioutil"
    "log"
    "strconv"
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
		log.Fatalf("ERR: %v\n", err)
	}

	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	cli, _ := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   []string{"etcd.cwxstat.io:2379"},
		TLS:         tlsConfig,
	})
	defer cli.Close()
	kv := clientv3.NewKV(cli)

	GetSingleValueDemo(ctx, kv)
    GetSingleValueDemo2(ctx, kv)
    LeaseDemo(ctx, cli, kv)
    GetMultipleValuesWithPaginationDemo(ctx, kv)
}

func GetSingleValueDemo(ctx context.Context, kv clientv3.KV) {
	fmt.Println("*** GetSingleValueDemo()")

	// Insert a key value
	pr, _ := kv.Put(ctx, "slop", "bob")
	rev := pr.Header.Revision
	fmt.Println("Revision:", rev)


    gr, _ := kv.Get(ctx, "slop")
    fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

    // Modify the value of an existing key (create new revision)
    kv.Put(ctx, "slop", "555")

    gr, _ = kv.Get(ctx, "slop")
    fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

    // Get the value of the previous revision
    gr, _ = kv.Get(ctx, "slop", clientv3.WithRev(rev))
    fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)


}

func LeaseDemo(ctx context.Context, cli *clientv3.Client, kv clientv3.KV) {
    fmt.Println("*** LeaseDemo()")
    // Delete all keys
    kv.Delete(ctx, "key", clientv3.WithPrefix())

    gr, _ := kv.Get(ctx, "key")
    if len(gr.Kvs) == 0 {
        fmt.Println("No 'key'")
    }


    lease, _ := cli.Grant(ctx, 1)

    // Insert key with a lease of 1 second TTL
    kv.Put(ctx, "key", "value", clientv3.WithLease(lease.ID))

    gr, _ = kv.Get(ctx, "key")
    if len(gr.Kvs) == 1 {
        fmt.Println("Found 'key'")
    }

    // Let the TTL expire
    time.Sleep(3 * time.Second)

    gr, _ = kv.Get(ctx, "key")
    if len(gr.Kvs) == 0 {
        fmt.Println("No more 'key'")
    }
}

func GetMultipleValuesWithPaginationDemo(ctx context.Context, kv clientv3.KV) {
    fmt.Println("*** GetMultipleValuesWithPaginationDemo()")
    // Delete all keys
    kv.Delete(ctx, "key", clientv3.WithPrefix())

    // Insert 20 keys
    for i := 0; i < 20; i++ {
        k := fmt.Sprintf("key_%02d", i)
        kv.Put(ctx, k, strconv.Itoa(i))
    }

    opts := []clientv3.OpOption{
        clientv3.WithPrefix(),
        clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
        clientv3.WithLimit(3),
    }

    gr, _ := kv.Get(ctx, "key", opts...)

    fmt.Println("--- First page ---")
    for _, item := range gr.Kvs {
        fmt.Println(string(item.Key), string(item.Value))
    }

    lastKey := string(gr.Kvs[len(gr.Kvs)-1].Key)

    fmt.Println("--- Second page ---")
    opts = append(opts, clientv3.WithFromKey())
    gr, _ = kv.Get(ctx, lastKey, opts...)

    // Skipping the first item, which the last item from from the previous Get
    for _, item := range gr.Kvs[1:] {
        fmt.Println(string(item.Key), string(item.Value))
    }
}


func GetSingleValueDemo2(ctx context.Context, kv clientv3.KV) {
    fmt.Println("*** GetSingleValueDemo()")
    // Delete all keys
    kv.Delete(ctx, "key", clientv3.WithPrefix())

    // Insert a key value
    pr, _ := kv.Put(ctx, "key", "444")
    rev := pr.Header.Revision
    fmt.Println("Revision:", rev)

    gr, _ := kv.Get(ctx, "key")
    fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

    // Modify the value of an existing key (create new revision)
    kv.Put(ctx, "key", "555")

    gr, _ = kv.Get(ctx, "key")
    fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

    // Get the value of the previous revision
    gr, _ = kv.Get(ctx, "key", clientv3.WithRev(rev))
    fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)
}