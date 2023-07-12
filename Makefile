.PHONY: video_service
video_service:
	protoc --go_out=. --go_opt=paths=source_relative -I. --go-grpc_out=paths=source_relative:. --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative --plugin=protoc-gen-grpc-gateway=${GOPATH}/bin/protoc-gen-grpc-gateway api/VideoService.proto

.PHONY: playlist_service
playlist_service:
	protoc --go_out=. --go_opt=paths=source_relative -I. --go-grpc_out=paths=source_relative:. --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative --plugin=protoc-gen-grpc-gateway=${GOPATH}/bin/protoc-gen-grpc-gateway api/PlaylistService.proto

.PHONY: account_service
account_service:
	protoc --go_out=. --go_opt=paths=source_relative -I. --go-grpc_out=paths=source_relative:. --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative --plugin=protoc-gen-grpc-gateway=${GOPATH}/bin/protoc-gen-grpc-gateway api/AccountService.proto

.PHONY: subscription_service
subscription_service:
	protoc --go_out=. --go_opt=paths=source_relative -I. --go-grpc_out=paths=source_relative:. --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative --plugin=protoc-gen-grpc-gateway=${GOPATH}/bin/protoc-gen-grpc-gateway api/SubscriptionService.proto

.PHONY: channel_service
channel_service:
	protoc --go_out=. --go_opt=paths=source_relative -I. --go-grpc_out=paths=source_relative:. --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative --plugin=protoc-gen-grpc-gateway=${GOPATH}/bin/protoc-gen-grpc-gateway api/ChannelService.proto
