#!/bin/bash

# 本脚本从头构建一个区块链网络
# 请确保 cryptogen 和 configtxgen 这两个可执行文件已经被正确安装
# 创建一个通道 vmichannel

echo "1. Environmental cleaning"
mkdir -p config
mkdir -p crypto-config
rm -fr config/*
rm -fr crypto-config/*

echo "Clean up"

echo "2. Generate certificate and start block information"
cryptogen generate --config=./crypto-config.yaml
configtxgen -profile OneOrgOrdererGenesis -outputBlock ./config/genesis.block

echo "3. Blockchain: start"
docker-compose up -d        # 按照docker-compose.yaml的配置启动区块链网络并在后台运行

# 四、生成通道(这个动作会创建一个创世交易，也是该通道的创世交易)
echo "4. Generate TX file for channel"
configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./config/vmichannel.tx -channelID vmichannel

# 五、在区块链上按照刚刚生成的TX文件去创建通道
# 该操作和上面操作不一样的是，这个操作会写入区块链
echo "5. Create channels on the blockchain according to TX files"
docker exec cli peer channel create -o orderer.vmi.com:7050 -c vmichannel -f /etc/hyperledger/config/vmichannel.tx

# 六、让节点去加入到通道
echo "6. Let the node join the channel"
docker exec cli peer channel join -b vmichannel.block

# 七、链码安装
echo "7. Installation chain code"
docker exec cli peer chaincode install -n vmicc -v 1.0.0 -l golang -p github.com/vendor-manage-inventory/chaincode

#八、实例化链码
#-n 对应前文安装链码的名字 其实就是composer network start bna名字
#-v 为版本号，相当于composer network start bna名字@版本号
#-C 是通道，在fabric的世界，一个通道就是一条不同的链，composer并没有很多提现这点，composer提现channel也就在于多组织时候的数据隔离和沟通使用
#-c 为传参，传入init参数
echo "8. Instantiated chain code"
docker exec cli peer chaincode instantiate -o orderer.vmi.com:7050 -C vmichannel -n vmicc -l golang -v 1.0.0 -c '{"Args":["init"]}'


# 九、链码交互
# 规则： docker exec cli peer chaincode invoke -C 通道名 -n 前文安装链码的名字 -c 参数
# 见 invoke.sh

