sidekick is a library that will update a key in etcd periodically.  

If sidekick.New() returns without an error, there will be a goroutine keeping the key refreshed in etcd.  

[![GoDoc](https://godoc.org/github.com/billhathahway/sidekick?status.png)](https://godoc.org/github.com/billhathaway/sidekick)

Example: 
 
    hostname,err := os.Hostname()  
	// comma separated list of server URLs  
    etcdServers := "http://localhost:4001"  
    etcdPath := "/servers/web/" + hostname  
    value := "running"  
    sk, _ := sidekick.New(etcdServers,etcdpath,value)  


  To change the value of the key, use:  

    sk.Value("newValue")  
	
  To stop updates to etcd, use:  

    sk.Stop()  
