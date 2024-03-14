package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mrdan4es/test-go-containerregistry-proxy/server"
	"log/slog"

	"github.com/mrdan4es/test-go-containerregistry-proxy/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("parse config:", err)
	}
	msg, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		slog.Error("log service config:", err)
		return
	}

	slog.Info(fmt.Sprintf("config: %s", msg))

	remoteServer := server.NewSecureRemoteServer(cfg.Remote, "docker.io")
	remoteRepo := server.NewRemoteRepository(remoteServer)

	descriptor, err := remoteRepo.FetchReleaseDescriptor(context.Background())
	if err != nil {
		slog.Error(err.Error())
	}

	descriptorIndent, err := json.MarshalIndent(descriptor, "", "  ")
	if err != nil {
		slog.Error("log service config:", err)
		return
	}

	slog.Info(string(descriptorIndent))
}
