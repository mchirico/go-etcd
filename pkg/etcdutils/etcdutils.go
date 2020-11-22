package etcdutils

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

type ETC struct {
	CertsDir string
}

func NewETC(certsDir ...string) ETC {
	e := ETC{}
	if certsDir == nil {
		e.CertsDir = "/certs"
	} else {
		e.CertsDir = certsDir[0]
	}

	return e
}

// "../../certs/client.pem"
func (e ETC) EtcdRun() string {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	cert, err := tls.LoadX509KeyPair(e.CertsDir+"/client.pem", e.CertsDir+"/client-key.pem")
	caCert, err := ioutil.ReadFile(e.CertsDir + "/ca.pem")
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

		TLS: tlsConfig,
	})
	defer cli.Close()
	kv := clientv3.NewKV(cli)

	s := GetSingleValueDemo(ctx, kv)
	GetSingleValueDemo2(ctx, kv)
	LeaseDemo(ctx, cli, kv)
	GetMultipleValuesWithPaginationDemo(ctx, kv)

	Watch(ctx, cli)

	for i := 0; i < 9; i++ {
		Txn(ctx, kv)
	}
	return s
}

func GetSingleValueDemo(ctx context.Context, kv clientv3.KV) string {
	fmt.Println("*** GetSingleValueDemo()")

	// Insert a key value
	pr, _ := kv.Put(ctx, "slop", "bob")
	rev := pr.Header.Revision
	fmt.Println("Revision:", rev)

	gr, _ := kv.Get(ctx, "slop")
	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)
	s := fmt.Sprintln("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

	// Modify the value of an existing key (create new revision)
	kv.Put(ctx, "slop", "555")

	gr, _ = kv.Get(ctx, "slop")
	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)
	s += fmt.Sprintln("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

	// Get the value of the previous revision
	gr, _ = kv.Get(ctx, "slop", clientv3.WithRev(rev))
	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

	return s
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
	time.Sleep(2 * time.Second)

	gr, _ = kv.Get(ctx, "key")
	if len(gr.Kvs) == 0 {
		fmt.Println("No more 'key'")
	}
}

func Txn(ctx context.Context, kv clientv3.KV) {

	tx := kv.Txn(ctx)

	txresp, err := tx.If(
		clientv3.Compare(clientv3.Value("foo"), "=", "bar"),
	).Then(
		clientv3.OpPut("foo", "sanfoo"), clientv3.OpPut("newfoo", "newbar"),
	).Else(
		clientv3.OpPut("foo", "bar"), clientv3.OpDelete("newfoo"),
	).Commit()
	fmt.Println(txresp, err)

	gr, _ := kv.Get(ctx, "foo")
	fmt.Println("\n\nTxn:\nValue: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)

}

func Watch(ctx context.Context, cli *clientv3.Client) {

	rch := cli.Watch(ctx, "foo", clientv3.WithPrefix())

	go func(chn clientv3.WatchChan) {
		for wresp := range chn {
			for _, ev := range wresp.Events {
				fmt.Printf("WATCH!!")
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}(rch)
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
