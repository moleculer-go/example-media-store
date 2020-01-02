package services

import (
	"github.com/moleculer-go/moleculer"
)

var Upload = moleculer.ServiceSchema{
	Name: "upload",
	Actions: []moleculer.Action{
		{
			Name: "picture",
			Handler: func(ctx moleculer.Context, params moleculer.Payload) interface{} {
				user := params.Get("user").String()
				pic := params.Get("picture").Stream()

				return nil
			},
		},
	},
}
