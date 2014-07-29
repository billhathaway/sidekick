Sidekick is a library that will update a key in etcd periodically.

Example:  
    hostname,err := os.Hostname()  
    etcdServers := "http://localhost:4001" // comma separated list of server URLs
    path := "/servers/web/" + hostname  
    value := "running"
    sidekick, _ := Sidekick.New(etcdServers,path,value)  

  If you want to stop the updates, use:  
    sidekick.Stop()  
