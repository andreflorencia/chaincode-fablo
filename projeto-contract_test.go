/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	
	"github.com/stretchr/testify/mock"
)

//const getStateError = "world state get error"

type MockStub struct {
	shim.ChaincodeStubInterface
	mock.Mock
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
	args := ms.Called(key)

	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
	args := ms.Called(key, value)

	return args.Error(0)
}

func (ms *MockStub) DelState(key string) error {
	args := ms.Called(key)

	return args.Error(0)
}

type MockContext struct {
	contractapi.TransactionContextInterface
	mock.Mock
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	args := mc.Called()

	return args.Get(0).(*MockStub)
}
/* 
func configureStub() (*MockContext, *MockStub) {
	var nilBytes []byte

	testProjeto := new(Projeto)
	testProjeto.Value = "set value"
	projetoBytes, _ := json.Marshal(testProjeto)

	ms := new(MockStub)
	ms.On("GetState", "statebad").Return(nilBytes, errors.New(getStateError))
	ms.On("GetState", "missingkey").Return(nilBytes, nil)
	ms.On("GetState", "existingkey").Return([]byte("some value"), nil)
	ms.On("GetState", "projetokey").Return(projetoBytes, nil)
	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

	mc := new(MockContext)
	mc.On("GetStub").Return(ms)

	return mc, ms
}

func TestProjetoExists(t *testing.T) {
	var exists bool
	var err error

	ctx, _ := configureStub()
	c := new(ProjetoContract)

	exists, err = c.ProjetoExists(ctx, "statebad")
	assert.EqualError(t, err, getStateError)
	assert.False(t, exists, "should return false on error")

	exists, err = c.ProjetoExists(ctx, "missingkey")
	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
	assert.False(t, exists, "should return false when no value for key in world state")

	exists, err = c.ProjetoExists(ctx, "existingkey")
	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
	assert.True(t, exists, "should return true when value for key in world state")
}

func TestCreateProjeto(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(ProjetoContract)

	err = c.CreateProjeto(ctx, "statebad", "some value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.CreateProjeto(ctx, "existingkey", "some value")
	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

	stub.AssertCalled(t, "PutState", "missingkey", []byte("{\"value\":\"some value\"}"))
}

func TestReadProjeto(t *testing.T) {
	var projeto *Projeto
	var err error

	ctx, _ := configureStub()
	c := new(ProjetoContract)

	projeto, err = c.ReadProjeto(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
	assert.Nil(t, projeto, "should not return Projeto when exists errors when reading")

	projeto, err = c.ReadProjeto(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
	assert.Nil(t, projeto, "should not return Projeto when key does not exist in world state when reading")

	projeto, err = c.ReadProjeto(ctx, "existingkey")
	assert.EqualError(t, err, "Could not unmarshal world state data to type Projeto", "should error when data in key is not Projeto")
	assert.Nil(t, projeto, "should not return Projeto when data in key is not of type Projeto")

	projeto, err = c.ReadProjeto(ctx, "projetokey")
	expectedProjeto := new(Projeto)
	expectedProjeto.Value = "set value"
	assert.Nil(t, err, "should not return error when Projeto exists in world state when reading")
	assert.Equal(t, expectedProjeto, projeto, "should return deserialized Projeto from world state")
}

func TestUpdateProjeto(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(ProjetoContract)

	err = c.UpdateProjeto(ctx, "statebad", "new value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when updating")

	err = c.UpdateProjeto(ctx, "missingkey", "new value")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when updating")

	err = c.UpdateProjeto(ctx, "projetokey", "new value")
	expectedProjeto := new(Projeto)
	expectedProjeto.Value = "new value"
	expectedProjetoBytes, _ := json.Marshal(expectedProjeto)
	assert.Nil(t, err, "should not return error when Projeto exists in world state when updating")
	stub.AssertCalled(t, "PutState", "projetokey", expectedProjetoBytes)
}

func TestDeleteProjeto(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(ProjetoContract)

	err = c.DeleteProjeto(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.DeleteProjeto(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when deleting")

	err = c.DeleteProjeto(ctx, "projetokey")
	assert.Nil(t, err, "should not return error when Projeto exists in world state when deleting")
	stub.AssertCalled(t, "DelState", "projetokey")
}
 */