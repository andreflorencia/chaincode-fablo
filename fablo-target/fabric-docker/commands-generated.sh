#!/usr/bin/env bash

generateArtifacts() {
  printHeadline "Generating basic configs" "U1F913"

  printItalics "Generating crypto material for Orderer" "U1F512"
  certsGenerate "$FABLO_NETWORK_ROOT/fabric-config" "crypto-config-orderer.yaml" "peerOrganizations/orderer.example.com" "$FABLO_NETWORK_ROOT/fabric-config/crypto-config/"

  printItalics "Generating crypto material for Org1" "U1F512"
  certsGenerate "$FABLO_NETWORK_ROOT/fabric-config" "crypto-config-org1.yaml" "peerOrganizations/org1.example.com" "$FABLO_NETWORK_ROOT/fabric-config/crypto-config/"

  printItalics "Generating genesis block for group group1" "U1F3E0"
  genesisBlockCreate "$FABLO_NETWORK_ROOT/fabric-config" "$FABLO_NETWORK_ROOT/fabric-config/config" "Group1Genesis"

  # Create directory for chaincode packages to avoid permission errors on linux
  mkdir -p "$FABLO_NETWORK_ROOT/fabric-config/chaincode-packages"
}

startNetwork() {
  printHeadline "Starting network" "U1F680"
  (cd "$FABLO_NETWORK_ROOT"/fabric-docker && docker-compose up -d)
  sleep 4
}

generateChannelsArtifacts() {
  printHeadline "Generating config for 'my-channel1'" "U1F913"
  createChannelTx "my-channel1" "$FABLO_NETWORK_ROOT/fabric-config" "MyChannel1" "$FABLO_NETWORK_ROOT/fabric-config/config"
}

installChannels() {
  printHeadline "Creating 'my-channel1' on Org1/peer0" "U1F63B"
  docker exec -i cli.org1.example.com bash -c "source scripts/channel_fns.sh; createChannelAndJoin 'my-channel1' 'Org1MSP' 'peer0.org1.example.com:7041' 'crypto/users/Admin@org1.example.com/msp' 'orderer0.group1.orderer.example.com:7030';"

  printItalics "Joining 'my-channel1' on  Org1/peer1" "U1F638"
  docker exec -i cli.org1.example.com bash -c "source scripts/channel_fns.sh; fetchChannelAndJoin 'my-channel1' 'Org1MSP' 'peer1.org1.example.com:7042' 'crypto/users/Admin@org1.example.com/msp' 'orderer0.group1.orderer.example.com:7030';"
}

installChaincodes() {
  if [ -n "$(ls "$CHAINCODES_BASE_DIR/.")" ]; then
    local version="3"
    printHeadline "Packaging chaincode 'workshop'" "U1F60E"
    chaincodeBuild "workshop" "golang" "$CHAINCODES_BASE_DIR/."
    chaincodePackage "cli.org1.example.com" "peer0.org1.example.com:7041" "workshop" "$version" "golang" printHeadline "Installing 'workshop' for Org1" "U1F60E"
    chaincodeInstall "cli.org1.example.com" "peer0.org1.example.com:7041" "workshop" "$version" ""
    chaincodeInstall "cli.org1.example.com" "peer1.org1.example.com:7042" "workshop" "$version" ""
    chaincodeApprove "cli.org1.example.com" "peer0.org1.example.com:7041" "my-channel1" "workshop" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "" ""
    printItalics "Committing chaincode 'workshop' on channel 'my-channel1' as 'Org1'" "U1F618"
    chaincodeCommit "cli.org1.example.com" "peer0.org1.example.com:7041" "my-channel1" "workshop" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "" "peer0.org1.example.com:7041" "" ""
  else
    echo "Warning! Skipping chaincode 'workshop' installation. Chaincode directory is empty."
    echo "Looked in dir: '$CHAINCODES_BASE_DIR/.'"
  fi

}

notifyOrgsAboutChannels() {
  printHeadline "Creating new channel config blocks" "U1F537"
  createNewChannelUpdateTx "my-channel1" "Org1MSP" "MyChannel1" "$FABLO_NETWORK_ROOT/fabric-config" "$FABLO_NETWORK_ROOT/fabric-config/config"

  printHeadline "Notyfing orgs about channels" "U1F4E2"
  notifyOrgAboutNewChannel "my-channel1" "Org1MSP" "cli.org1.example.com" "peer0.org1.example.com" "orderer0.group1.orderer.example.com:7030"

  printHeadline "Deleting new channel config blocks" "U1F52A"
  deleteNewChannelUpdateTx "my-channel1" "Org1MSP" "cli.org1.example.com"
}

upgradeChaincode() {
  local chaincodeName="$1"
  if [ -z "$chaincodeName" ]; then
    echo "Error: chaincode name is not provided"
    exit 1
  fi

  local version="$2"
  if [ -z "$version" ]; then
    echo "Error: chaincode version is not provided"
    exit 1
  fi

  if [ "$chaincodeName" = "workshop" ]; then
    if [ -n "$(ls "$CHAINCODES_BASE_DIR/.")" ]; then
      printHeadline "Packaging chaincode 'workshop'" "U1F60E"
      chaincodeBuild "workshop" "golang" "$CHAINCODES_BASE_DIR/."
      chaincodePackage "cli.org1.example.com" "peer0.org1.example.com:7041" "workshop" "$version" "golang" printHeadline "Installing 'workshop' for Org1" "U1F60E"
      chaincodeInstall "cli.org1.example.com" "peer0.org1.example.com:7041" "workshop" "$version" ""
      chaincodeInstall "cli.org1.example.com" "peer1.org1.example.com:7042" "workshop" "$version" ""
      chaincodeApprove "cli.org1.example.com" "peer0.org1.example.com:7041" "my-channel1" "workshop" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "" ""
      printItalics "Committing chaincode 'workshop' on channel 'my-channel1' as 'Org1'" "U1F618"
      chaincodeCommit "cli.org1.example.com" "peer0.org1.example.com:7041" "my-channel1" "workshop" "$version" "orderer0.group1.orderer.example.com:7030" "" "false" "" "peer0.org1.example.com:7041" "" ""

    else
      echo "Warning! Skipping chaincode 'workshop' upgrade. Chaincode directory is empty."
      echo "Looked in dir: '$CHAINCODES_BASE_DIR/.'"
    fi
  fi
}

stopNetwork() {
  printHeadline "Stopping network" "U1F68F"
  (cd "$FABLO_NETWORK_ROOT"/fabric-docker && docker-compose stop)
  sleep 4
}

networkDown() {
  printHeadline "Destroying network" "U1F916"
  (cd "$FABLO_NETWORK_ROOT"/fabric-docker && docker-compose down)

  printf "\nRemoving chaincode containers & images... \U1F5D1 \n"
  docker rm -f $(docker ps -a | grep dev-peer0.org1.example.com-workshop-3-* | awk '{print $1}') || echo "docker rm failed, Check if all fabric dockers properly was deleted"
  docker rmi $(docker images dev-peer0.org1.example.com-workshop-3-* -q) || echo "docker rm failed, Check if all fabric dockers properly was deleted"
  docker rm -f $(docker ps -a | grep dev-peer1.org1.example.com-workshop-3-* | awk '{print $1}') || echo "docker rm failed, Check if all fabric dockers properly was deleted"
  docker rmi $(docker images dev-peer1.org1.example.com-workshop-3-* -q) || echo "docker rm failed, Check if all fabric dockers properly was deleted"

  printf "\nRemoving generated configs... \U1F5D1 \n"
  rm -rf "$FABLO_NETWORK_ROOT/fabric-config/config"
  rm -rf "$FABLO_NETWORK_ROOT/fabric-config/crypto-config"
  rm -rf "$FABLO_NETWORK_ROOT/fabric-config/chaincode-packages"

  printHeadline "Done! Network was purged" "U1F5D1"
}
