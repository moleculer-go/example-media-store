package services

import (
	"github.com/moleculer-go/cqrs"
	"github.com/moleculer-go/moleculer"
	"github.com/moleculer-go/store"
	"github.com/moleculer-go/store/sqlite"
)

var events = cqrs.EventStore("userMediaEvents", storeFactory())

//userMediaStore store pictures by user
var userMediaStore = cqrs.Aggregate(
	"userMedia",
	storeFactory(map[string]interface{}{
		"userId":   "string",
		"fileId":   "string",
		"picHash":  "string",
		"metadata": "map",
	}),
	cqrs.NoSnapshot,
).Snapshot("userMediaEvents")

var allMediaStore = cqrs.Aggregate(
	"allMedia",
	storeFactory(map[string]interface{}{
		"picHash":       "string",
		"userPictureId": "string",
		"metadata":      "map",
	}),
	cqrs.NoSnapshot,
).Snapshot("userMediaEvents")

var MediaService = moleculer.ServiceSchema{
	Name: "media",
	Mixins: []moleculer.Mixin{
		events.Mixin(),
		allMediaStore.Mixin(),
		userMediaStore.Mixin()},
	Actions: []moleculer.Action{
		events.MapAction("create", "media.created"),
		{
			Name:    "toUserMedia",
			Handler: toUserMedia,
		},
		{
			Name:    "toAllMedia",
			Handler: toAllMedia,
		},
	},
	Events: []moleculer.Event{
		userMediaStore.On("media.created").Create("media.toUserMedia"),
		allMediaStore.On("media.created").Update("media.toAllMedia"),
	},
}

// toUserMedia transform the event to be stored in the user media aggregate
func toUserMedia(context moleculer.Context, event moleculer.Payload) interface{} {
	p := event.Get("payload")
	userId := p.Get("user").String()
	userMedia := p.Remove("user", "eventId").Add("userId", userId)
	return userMedia
}

// toAllMedia transform the event to be stored in the all media aggregate
func toAllMedia(context moleculer.Context, event moleculer.Payload) interface{} {
	p := event.Get("payload")
	userPictureId := p.Get("user").String()
	p = p.Remove("userId", "fileId", "eventId").Add("userPictureId", userPictureId)
	return p
}

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
