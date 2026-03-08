package app

import "github.com/phrara/mallive/inventory/app/query"

type Application struct {
	Commands Commands
	Queries Queries
}

type Commands struct {}

type Queries struct {
	GetItems query.GetItemsHandler
	CheckIfItemsInInventory query.CheckIfItemsInInventoryHandler
}