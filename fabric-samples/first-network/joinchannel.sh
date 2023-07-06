export CORE_PEER_TLS_ENABLED=true
export FABRIC_CFG_PATH=$PWD
export CHANNEL_NAME=mychannel

setGlobalsForPeer0qtl() {
  export CORE_PEER_LOCALMSPID="qtlMSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/qtl.infotelconnect.com/peers/peer0.qtl.infotelconnect.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=/home/qtl/qtl.infotelconnect.com/users/Admin@qtl.infotelconnect.com/msp
  export CORE_PEER_ADDRESS=peer0.qtl.infotelconnect.com:7051
}

setGlobalsForPeer1qtl() {
  export CORE_PEER_LOCALMSPID="qtlMSP"
  export CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/vmipl.com/peers/peer0.vmipl.com/tls/ca.crt
  export CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/qtl.infotelconnect.com/users/Admin@qtl.infotelconnect.com/msp
  export CORE_PEER_ADDRESS=peer1.qtl.infotelconnect.com:8051
}

createChannel() {

  setGlobalsForPeer0qtl

  peer channel create -o orderer0.ucccpr.com:7050 -c $CHANNEL_NAME \
    --ordererTLSHostnameOverride orderer0.ucccpr.com \
    -f ./${CHANNEL_NAME}.tx --outputBlock ./${CHANNEL_NAME}.block \
    --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA

}

joinChannel() {
  setGlobalsForPeer0qtl
  peer channel join -b ./$CHANNEL_NAME.block

  setGlobalsForPeer1qtl
  peer channel join -b ./$CHANNEL_NAME.block
}

updateAnchorPeers() {
  setGlobalsForPeer0qtl
  peer channel update -o orderer0.ucccpr.com:7050 --ordererTLSHostnameOverride orderer0.ucccpr.com -c $CHANNEL_NAME -f ${CORE_PEER_LOCALMSPID}anchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA

}

createChannel
sleep 2
joinChannel
sleep 2
updateAnchorPeers