/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ProjetoContract contract for managing CRUD for Projeto
type ProjetoContract struct {
	contractapi.Contract
}

// ProjetoExists returns true when asset with given ID exists in world state
func (c *ProjetoContract) ProjetoExists(ctx contractapi.TransactionContextInterface, projetoID string) (bool, error) {
	data, err := ctx.GetStub().GetState(projetoID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateProjeto creates a new instance of Projeto
func (c *ProjetoContract) CreateProjeto(ctx contractapi.TransactionContextInterface, projetoID, name, description, owner string, members []string, tasks []Task) (string, error) {
	exists, err := c.ProjetoExists(ctx, projetoID)
	if err != nil {
		return "", fmt.Errorf("could not read from world state. %s", err)
	} else if exists {
		return "", fmt.Errorf("the project %s already exists", projetoID)
	}

	projeto := &Projeto{
		DocType:     "Projeto",
		Id:          projetoID,
		Name:        name,
		Description: description,
		Owner:       owner,
		Members:     members,
		Tasks:       tasks,
	}

	bytes, err := json.Marshal(projeto)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(projetoID, bytes)
	if err != nil {
		return "", fmt.Errorf("failed to put to world state: %v", err)
	}

	return fmt.Sprintf("Project %s created successfully", projetoID), nil
}

func (c *ProjetoContract) ReadAllProjetos(ctx contractapi.TransactionContextInterface) ([]*Projeto, error) {
	iter, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("could not retrieve projects. %s", err)
	}
	defer iter.Close()

	var projetos []*Projeto
	for iter.HasNext() {
		keyValue, err := iter.Next()
		if err != nil {
			return nil, fmt.Errorf("could not retrieve next project. %s", err)
		}

		projeto := new(Projeto)
		err = json.Unmarshal(keyValue.Value, projeto)
		if err != nil {
			return nil, fmt.Errorf("could not unmarshal project data. %s", err)
		}

		projetos = append(projetos, projeto)
	}

	return projetos, nil
}


// ReadProjeto retrieves an instance of Projeto from the world state
func (c *ProjetoContract) ReadProjeto(ctx contractapi.TransactionContextInterface, projetoID string) (*Projeto, error) {
	exists, err := c.ProjetoExists(ctx, projetoID)
	if err != nil {
		return nil, fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("the project %s does not exist", projetoID)
	}

	bytes, err := ctx.GetStub().GetState(projetoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get state: %v", err)
	}

	projeto := new(Projeto)
	err = json.Unmarshal(bytes, projeto)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal data to type Projeto: %v", err)
	}

	return projeto, nil
}




func (c *ProjetoContract) UpdateProjeto(ctx contractapi.TransactionContextInterface, projetoID, name, description, owner string, members []string, tasks []Task) (string, error) {
    exists, err := c.ProjetoExists(ctx, projetoID)
    if err != nil {
        return "", fmt.Errorf("could not read from world state. %s", err)
    } else if !exists {
        return "", fmt.Errorf("the project %s does not exists", projetoID)
    }

    projeto := &Projeto{
        DocType:     "Projeto",
        Id:          projetoID,
        Name:        name,
        Description: description,
        Owner:       owner,
        Members:     members,
        Tasks:       tasks,
    }

    bytes, err := json.Marshal(projeto)
    if err != nil {
        return "", err
    }

    err = ctx.GetStub().PutState(projetoID, bytes)
    if err != nil {
        return "", fmt.Errorf("failed to put to world state: %v", err)
    }

    return fmt.Sprintf("Project %s updated successfully", projetoID), nil
}



// DeleteProjeto deletes an instance of Projeto from the world state
func (c *ProjetoContract) DeleteProjeto(ctx contractapi.TransactionContextInterface, projetoID string) error {
	exists, err := c.ProjetoExists(ctx, projetoID)
	if err != nil {
		return fmt.Errorf("could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("the asset %s does not exist", projetoID)
	}

	return ctx.GetStub().DelState(projetoID)
}

func (c *ProjetoContract) DeleteTaskFromProjeto(ctx contractapi.TransactionContextInterface, projetoID, taskID string) (string, error) {
    exists, err := c.ProjetoExists(ctx, projetoID)
    if err != nil {
        return "", fmt.Errorf("could not read from world state. %s", err)
    } else if !exists {
        return "", fmt.Errorf("the project %s does not exist", projetoID)
    }

    projetoBytes, err := ctx.GetStub().GetState(projetoID)
    if err != nil {
        return "", fmt.Errorf("failed to get project from world state: %v", err)
    }

    var projeto Projeto
    err = json.Unmarshal(projetoBytes, &projeto)
    if err != nil {
        return "", fmt.Errorf("failed to unmarshal project: %v", err)
    }

    // Check if the task exists in the project
    taskIndex := -1
    for i, task := range projeto.Tasks {
        if task.Id == taskID {
            taskIndex = i
            break
        }
    }
    if taskIndex == -1 {
        return "", fmt.Errorf("the task %s does not exist in the project %s", taskID, projetoID)
    }

    // Delete the task from the project
    projeto.Tasks = append(projeto.Tasks[:taskIndex], projeto.Tasks[taskIndex+1:]...)

    projetoBytes, err = json.Marshal(projeto)
    if err != nil {
        return "", err
    }

    err = ctx.GetStub().PutState(projetoID, projetoBytes)
    if err != nil {
        return "", fmt.Errorf("failed to put project to world state: %v", err)
    }

    return fmt.Sprintf("Task %s deleted from project %s", taskID, projetoID), nil
}






