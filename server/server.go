package server

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/mrdan4es/test-go-containerregistry-proxy/config"
	"net/http"
	"net/url"
)

type RemoteServer interface {
	ServerName() string
	Transport() http.RoundTripper
}

type RemoteRepository struct {
	remoteServer RemoteServer
}

func NewRemoteRepository(remoteServer RemoteServer) *RemoteRepository {
	return &RemoteRepository{
		remoteServer: remoteServer,
	}
}

func (rr *RemoteRepository) RemoteServerURL() string {
	return rr.remoteServer.ServerName()
}

func (rr *RemoteRepository) Ping(ctx context.Context) error {
	_, err := rr.FetchReleaseDescriptor(ctx)
	return err
}

func (rr *RemoteRepository) FetchReleaseDescriptor(ctx context.Context) (*v1.Descriptor, error) {
	releaseRef, err := rr.ParseReference(rr.RemoteServerURL() + "/nginx:latest")
	if err != nil {
		return nil, fmt.Errorf("failed to parse release description image reference: %w", err)
	}

	releaseImgDesc, err := remote.Head(releaseRef, rr.remoteOpts(ctx)...)
	if err != nil {
		return nil, fmt.Errorf("failed to read release descriptor: %w", err)
	}

	return releaseImgDesc, err
}

func (rr *RemoteRepository) ParseReference(s string) (name.Reference, error) {
	return name.ParseReference(s)
}

func (rr *RemoteRepository) remoteOpts(ctx context.Context) []remote.Option {
	return []remote.Option{
		remote.WithContext(ctx),
		remote.WithTransport(rr.remoteServer.Transport()),
	}
}

type SecureRemoteServer struct {
	serverName   string
	roundTripper http.RoundTripper
}

func NewSecureRemoteServer(cfg config.Remote, serverName string) *SecureRemoteServer {
	rt := http.DefaultTransport.(*http.Transport).Clone()

	if cfg.Proxy.Url != "" {
		rt.Proxy = func(*http.Request) (*url.URL, error) {
			return url.Parse(cfg.Proxy.Url)
		}
	}

	return &SecureRemoteServer{
		serverName:   serverName,
		roundTripper: rt,
	}
}

func (rs *SecureRemoteServer) ServerName() string {
	return rs.serverName
}

func (rs *SecureRemoteServer) Transport() http.RoundTripper {
	return rs.roundTripper
}
