# Nftserver
## 1. Edit the configuration file (app.conf)
### Configure mysql related parameters, for example:
        #Database parameter configuration  
        #"username:userpassword@tcp(192.168.1.238:3306)/"
        #username  
        dbusername = *****
        #user password
        dbuserpassword = *******
        #Database server address
        dbserverip = 192.168.56.128
        #Database server port
        dbserverport = 3306
        #Name database
        dbname = nftdb
### Configure contract-related parameters, for example:
        #Contract event node access point
        WormholesNode = https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161
        EthersWsNode = wss://rinkeby.infura.io/ws/v3/97cb2119c79842b7818a7a37df749b2b

## 2. start the nftserver service
    setsid ./nftserver > log

