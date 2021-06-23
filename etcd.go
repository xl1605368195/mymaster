package main

//
//func  InitEtcdClient(logger *logger.Writer) error {
//	// etcd客户端配置
//	etcdClientConfig := clientv3.Config{
//		Endpoints:   m.config.Endpoints,
//		DialTimeout: 5 * time.Second,
//	}
//	etcdClient, err := clientv3.New(etcdClientConfig)
//	// etcd建立连接
//	if err != nil {
//		m.logger.Err(constants.ETCD_ERROR, "Etcd clientv3.New error:%+v", err)
//		return err
//	}
//	m.EtcdClient = etcdClient
//	return nil
//}
