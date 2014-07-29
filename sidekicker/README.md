sidekicker is a standalone etcd updating client

    go get github.com/billhathaway/sidekick/sidekicker

Usage of sidekicker:  
    -key="": key in etcd  
    -servers="http://localhost:4001": etcdServers in CSV list  
    -val="": value for key  
	
Example:  
  sidekicker -servers http://server1:4001,http://server2:4001,http://server3:4001 -key /web/servers/`` `uname -n` `` -val running  
  