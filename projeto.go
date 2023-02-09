/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

// Projeto stores a value
type Projeto struct {
	DocType 	string   `json:"DocType"`
	Id          string   `json:"Id"`
	Name        string   `json:"Name"`
	Description string   `json:"Description"`
	Owner       string   `json:"Owner"`
	Members     []string `json:"Members"`
	Tasks       []Task   `json:"Tasks"`
}

// Task represents a task in the project management application
type Task struct {
	Id 			string `json:"Id"`
	Description string `json:"Description"`
	Assignee    string `json:"Assignee"`
	Status      string `json:"Status"`
}