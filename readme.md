# Nftserver
## 1. Edit the configuration file (app.conf)
### Configure mysql related parameters, for example:
        #Database parameter configuration  
        #"admin:user123456@tcp(192.168.1.238:3306)/"
        #username  
        dbusername = admin
        #user password
        dbuserpassword = user123456
        #Database server address
        dbserverip = 192.168.56.128
        #Database server port
        dbserverport = 3306
        #Name database
        dbname = nftdb
### Configure contract-related parameters, for example:
        #Exchange contract
        TradeAddr = 0xD8D5D49182d7Abf3cFc1694F8Ed17742886dDE82
        #1155 Contract
        NFT1155Addr = 0xA1e67a33e090Afe696D7317e05c506d7687Bb2E5
        #Admin List Contract
        AdminAddr = 0x56c971ebBC0cD7Ba1f977340140297C0B48b7955
        #Contract event node access point
        EthersNode = https://rinkeby.infura.io/v3/9aa3d95b3bc440fa88ea12eaa4456161
        EthersWsNode = wss://rinkeby.infura.io/ws/v3/97cb2119c79842b7818a7a37df749b2b

## 2. the signature configuration file
    Before signing, you need to delete all the data after the [time] tag (including the [time] tag) in the app.conf file.
    ./signappconf -f app.conf -key without private key starting with 0x
    Copy the signed app.conf file to the conf directory in the same directory as the nftserver execution file.

## 3. start the nftserver service
    setsid ./nftserver > log
