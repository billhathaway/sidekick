Sidekick is a library that will update a key in etcd periodically.

Example:  
    hostname,err := os.Hostname()  
	// comma separated list of server URLs  
    etcdServers := "http://localhost:4001"  
    path := "/servers/web/" + hostname  
    value := "running"  
	
    sidekick, _ := Sidekick.New(etcdServers,path,value)  


  To change the value of the key, use:  
    sidekick.Value("somethingnew")  
	
  To stop updates to etcd, use:  
    sidekick.Stop()  
