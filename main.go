/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	projetoContract := new(ProjetoContract)
	projetoContract.Info.Version = "0.0.1"
	projetoContract.Info.Description = "My Smart Contract"
	projetoContract.Info.License = new(metadata.LicenseMetadata)
	projetoContract.Info.License.Name = "Apache-2.0"
	projetoContract.Info.Contact = new(metadata.ContactMetadata)
	projetoContract.Info.Contact.Name = "John Doe"

	chaincode, err := contractapi.NewChaincode(projetoContract)
	chaincode.Info.Title = "projeto3 chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from ProjetoContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
