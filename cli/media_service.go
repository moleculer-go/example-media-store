package main

import (
	"github.com/moleculer-go/example-life-album/services"
	gateway "github.com/moleculer-go/gateway"
	"github.com/moleculer-go/gateway/websocket"
	"github.com/moleculer-go/moleculer"
	"github.com/moleculer-go/moleculer/broker"
	"github.com/moleculer-go/moleculer/cli"
	"github.com/spf13/cobra"
)

func getGatewayConfig(cmd *cobra.Command) map[string]interface{} {
	env, _ := cmd.Flags().GetString("env")
	if env == "dev" {
		return map[string]interface{}{
			"reverseProxy": map[string]interface{}{
				"target": "http://localhost:3000",
			},
		}
	}
	return map[string]interface{}{}
}

func main() {
	websocketMixin := &websocket.WebSocketMixin{
		Mixins: []websocket.SocketMixin{
			&websocket.EventsMixin{},
		},
	}

	cli.Start(
		&moleculer.Config{LogLevel: "debug"},
		func(broker *broker.ServiceBroker, cmd *cobra.Command) {
			gatewaySvc := &gateway.HttpService{
				Settings: getGatewayConfig(cmd),
				Mixins:   []gateway.GatewayMixin{websocketMixin},
			}
			broker.Publish(gatewaySvc, services.Upload, services.MediaService)
			broker.Start()
		})
}
