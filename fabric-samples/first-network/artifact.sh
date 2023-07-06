../bin/cryptogen generate --config=./crypto-config.yaml

# channel name defaults to "mychannel"

export CHANNEL_NAME=mychannel

echo $CHANNEL_NAME

export FABRIC_CFG_PATH=$PWD

# Generate System Genesis block

../bin/configtxgen -profile TwoOrgsOrdererGenesis -channelID byfn-sys-channel -outputBlock ./channel-artifacts/genesis.block

# Generate channel configuration block

export CHANNEL_NAME=mychannel  && ../bin/configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME

echo "****   Generating anchor peer update for QtlMSP  ****"

../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/QTLMSPanchors.tx -channelID mychannel -asOrg qtlMSP

echo "****   Generating anchor peer update for VmiplMSP  ****"

../bin/configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/VMIPLMSPanchors.tx -channelID mychannel -asOrg vmiplMSP

sleep 2