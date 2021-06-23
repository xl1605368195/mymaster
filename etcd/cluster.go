package etcd

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/rand"
	"sync"
	"time"

	etcdclient "github.com/coreos/etcd/clientv3"

	"hids-manager/gourd/context/logging"
)

type ICluster interface {
	Register(value INodeInfo) error
	GetEtcdClient() *etcdclient.Client
	GetClusterNodes() []string
	GetConfigClusterNodes() []string

	//设置最大重试次数
	SetMaxRetry(uint)

	GetKey(value INodeInfo) (*etcdclient.GetResponse, error)
	Connect() error
	Revoke() error
	UploadConf(key, value string) error
	WatchKey(string) etcdclient.WatchChan
}

type EtcdConfig struct {
	Etcd_key_online_prefix      string //注册自己时的 key
	Etcd_key_config_host_prefix string //获取自己的配置 key
	//Etcd_key_status_prefix       string   //上报自身状态的 key // 20190516 废弃
	Etcd_dsn            []string //ETCD vip配置
	Etcd_dsn_backup     []string //ETCD 配置
	Etcd_ca_pem         string
	Etcd_client_pem     string
	Etcd_client_key_pem string
	Etcd_username       string //environ.conf中获取的用户名，用于server端role角色授权
	Etcd_passwd         string //environ.conf中获取的密码，用于server端role角色授权
}

type EtcdCluster struct {
	sync.Mutex
	logger         *logging.Writer
	etcdConfig     *EtcdConfig // ETCD地址列表
	kapi           *etcdclient.Client
	watchGroupChan etcdclient.WatchChan
	watchHostChan  etcdclient.WatchChan
	watchGroupKey  string
	watchHostKey   string
	vipFailed      bool // 已经尝试过连接VIP，但失败了
	maxRetry       uint
	curRetry       uint
	shutdown       chan interface{}
	nodeInfo       INodeInfo
	//leaseInfo     *etcdclient.LeaseGrantResponse
}

const (
	ETCD_CONNECT_TIMEOUT   = 5 * time.Second
	ETCD_TRANSPORT_TIMEOUT = 5 * time.Second

)

func NewEtcdCluster(etcdConfig *EtcdConfig, logger *logging.Writer) *EtcdCluster {

	logger.Debug("$$$$$$$ This is ETCD vip list from config: %v $$$$$$$", etcdConfig.Etcd_dsn)

	logger.Debug("$$$$$$$ This is ETCD group node list from config: %v $$$$$$$", etcdConfig.Etcd_dsn_backup)

	cluster := &EtcdCluster{
		etcdConfig: etcdConfig,
		logger:     logger,
	}
	return cluster
}

func (this *EtcdCluster) SetMaxRetry(n uint) {
	this.maxRetry = n
}

// 在这个函数里面，增加retry的判断逻辑，如果是vip连接出错，那么再次尝试的时候，尝试原生的ip。
// 但是这个函数里面无法知道retry的状况，所以，第一次连接不上vip，第二次就直接连接原生节点了。
func (this *EtcdCluster) Connect() error {
	// 也可以考虑按照尝试过几次后，再连接备用的。
	// 是否已经尝试过vip连接，且尝试过的次数已经达到限制
	if !this.vipFailed && this.curRetry >= this.maxRetry/5 {
		this.etcdConfig.Etcd_dsn = this.etcdConfig.Etcd_dsn_backup
		this.etcdConfig = randEtcdDsn(this.etcdConfig)
		this.vipFailed = true
	}

	// load cert
	//this.logger.Debug("etcdConfig.Etcd_client_pem:%s", this.etcdConfig.Etcd_client_pem)
	//this.logger.Debug("etcdConfig.Etcd_client_key_pem:%s", this.etcdConfig.Etcd_client_key_pem)
	cert, err := tls.X509KeyPair([]byte(this.etcdConfig.Etcd_client_pem), []byte(this.etcdConfig.Etcd_client_key_pem))
	if err != nil {
		return err
	}

	//this.logger.Debug("client key 、client pem 读取完成")
	// load root ca
	//caData, err := ioutil.ReadFile(this.etcdConfig.Etcd_ca_pem)
	//if err != nil {
	//	return  err
	//}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(this.etcdConfig.Etcd_ca_pem))

	_tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}

	cfg := etcdclient.Config{
		Endpoints:   this.etcdConfig.Etcd_dsn,
		TLS:         _tlsConfig,
		DialTimeout: ETCD_CONNECT_TIMEOUT,
		//AutoSyncInterval: time.Second * 300,
		//Username:         this.etcdConfig.Etcd_username,
		//Password:         this.etcdConfig.Etcd_passwd,
	}

	client, err := etcdclient.New(cfg)
	this.curRetry++
	if err != nil {
		return err
	}

	this.kapi = client
	return nil
}

func (this *EtcdCluster) GetEtcdClient() *etcdclient.Client {
	return this.kapi
}

func (this *EtcdCluster) GetKey(value INodeInfo) (*etcdclient.GetResponse, error) {
	keyName := fmt.Sprintf(this.etcdConfig.Etcd_key_online_prefix, value.GetUuid())
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_TRANSPORT_TIMEOUT)
	rs, err := this.kapi.Get(ctx, keyName)
	defer cancel() //修复后
	if err != nil {
		this.logger.Debug("## GetKey (key : %v) got Error : %v ", keyName, err)
		return nil, err
	}

	//cancel()	//修复前
	this.logger.Debug("*** Got key value from etcd, key : %s , response : %v", keyName, rs)

	return rs, nil
}

func (this *EtcdCluster) Revoke() error {
	keyName := fmt.Sprintf(this.etcdConfig.Etcd_key_online_prefix, this.nodeInfo.GetUuid())
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_TRANSPORT_TIMEOUT)
	//leaseResp,err := this.kapi.Revoke(ctx,this.leaseInfo.ID)
	_, err := this.kapi.Delete(ctx, keyName)
	cancel()
	if err != nil {
		return err
	}
	//this.logger.Emerg(leaseResp.Header.String())
	return nil
}

func (this *EtcdCluster) Register(value INodeInfo) error {
	keyName := fmt.Sprintf(this.etcdConfig.Etcd_key_online_prefix, value.GetUuid())

	this.logger.Info("*** 注册etcd，register key：%s", keyName)

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_TRANSPORT_TIMEOUT)
	//leaseResp, err := this.GetEtcdClient().Grant(ctx, 600)	//租约时间设定为600秒
	//cancel()
	//if err != nil {
	//	return err
	//}

	//ctx, cancel = context.WithTimeout(context.Background(), ETCD_TRANSPORT_TIMEOUT)
	//kvc := etcdclient.NewKV(this.kapi)
	//var txnResp *etcdclient.TxnResponse
	//txnResp, err = kvc.Txn(ctx).
	//	If(etcdclient.Compare(etcdclient.CreateRevision(keyName), "=", 0)).
	//	Then(etcdclient.OpPut(keyName, value.GetRegisterInfo(),etcdclient.WithLease(etcdclient.LeaseID(leaseResp.ID)))).
	//	Commit()
	_, err = this.kapi.Put(ctx, keyName, value.GetRegisterInfo())
	cancel()

	if err != nil {
		return err
	}

	//if !txnResp.OpResponse() {
	//	return fmt.Errorf("无法注册自己:%s ，该 Key 已经存在",keyName)
	//}

	//ctx, cancel = context.WithTimeout(context.Background(), ETCD_TRANSPORT_TIMEOUT)
	//_, err = this.GetEtcdClient().KeepAlive(context.TODO(), leaseResp.ID)
	//cancel()
	//if err != nil {
	//	return err
	//}

	this.nodeInfo = value
	//this.leaseInfo = leaseResp

	return nil
}

//put value to etcd group key
func (this *EtcdCluster) UploadConf(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_TRANSPORT_TIMEOUT)
	rs, err := this.kapi.Put(ctx, key, value)
	cancel()
	this.logger.Debug("update config to etcd : %s , response: %v", key, rs)
	return err
}

func (this *EtcdCluster) WatchKey(key string) etcdclient.WatchChan {
	// only watch the change of union_conf, won't need WithPrefix()
	resp := this.kapi.Watch(context.TODO(), key)

	return resp
}

func (this *EtcdCluster) GetConfigClusterNodes() []string {
	return this.etcdConfig.Etcd_dsn
}

func (this *EtcdCluster) GetClusterNodes() []string {
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_TRANSPORT_TIMEOUT)
	memberlist, err := this.kapi.Cluster.MemberList(ctx)
	cancel()
	if err != nil {
		return nil
	}

	var mem []string
	members := memberlist.Members

	for _, member := range members {
		clients := member.GetClientURLs()
		for _, client := range clients {
			mem = append(mem, client)
		}
	}

	return mem
}
