package services

import (
	"github.com/moleculer-go/cqrs"
	"github.com/moleculer-go/moleculer"
	"github.com/moleculer-go/store"
	"github.com/moleculer-go/store/sqlite"
)

// storeFactory high order func that returns a cqrs.StoreFactory function :)
// and merges the fields passed to this function, with the fields received by the cqrs.StoreFactory func.
func storeFactory(fields ...map[string]interface{}) cqrs.StoreFactory {
	return func(name string, cqrsFields, settings map[string]interface{}) store.Adapter {
		fields = append(fields, cqrsFields)
		return &sqlite.Adapter{
			URI:     "file:memory:?mode=memory",
			Table:   name,
			Columns: cqrs.FieldsToSQLiteColumns(fields...),
		}
	}
}

var events = cqrs.EventStore("userMediaEventStore", storeFactory())

//userMediaAggregate store pictures by user
var userMediaAggregate = cqrs.Aggregate(
	"userMediaAggregate",
	storeFactory(map[string]interface{}{
		"userId":   "string",
		"fileId":   "string",
		"picHash":  "string",
		"metadata": "map[string]string",
	}),
	cqrs.NoSnapshot,
).Snapshot("userMediaEventStore")

var allMediaAggregate = cqrs.Aggregate(
	"allMediaAggregate",
	storeFactory(map[string]interface{}{
		"picHash":       "string",
		"userPictureId": "string",
		"metadata":      "map[string]string",
	}),
	cqrs.NoSnapshot,
).Snapshot("userMediaEventStore")

var UserMediaService = moleculer.ServiceSchema{
	Name:   "userMedia",
	Mixins: []moleculer.Mixin{events.Mixin(), allMediaAggregate.Mixin(), userMediaAggregate.Mixin()},
	Actions: []moleculer.Action{
		{
			Name:    "create",
			Handler: events.PersistEvent("userMedia.created"),
		},
		{
			Name:    "transformUserMedia",
			Handler: transformUserMedia,
		},
		{
			Name:    "transformAllMedia",
			Handler: transformAllMedia,
		},
	},
	Events: []moleculer.Event{
		userMediaAggregate.On("userMedia.created").Create("userMedia.transformUserMedia"),
		allMediaAggregate.On("userMedia.created").Update("userMedia.transformAllMedia"),
	},
}

// transformUserMedia transform the event to be stored in the user media aggregate
func transformUserMedia(context moleculer.Context, event moleculer.Payload) interface{} {
	p := event.Get("payload")
	userId := p.Get("user")
	return p.Remove("user").Add("userId", userId)
}

// transformAllMedia transform the event to be stored in the all media aggregate
func transformAllMedia(context moleculer.Context, event moleculer.Payload) interface{} {
	p := event.Get("payload")
	userPictureId := p.Get("userId").String()
	p = p.Remove("userId", "fileId").Add("userPictureId", userPictureId)
	return p
}

// emitAll utility function to invoke all events.
func emitAll(eventHandlers ...moleculer.Event) moleculer.EventHandler {
	return func(context moleculer.Context, event moleculer.Payload) {
		for _, evtHandler := range eventHandlers {
			evtHandler.Handler(context, event)
		}
	}
}
