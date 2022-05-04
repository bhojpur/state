package rpc

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/cors"

	"github.com/bhojpur/state/internal/pubsub"
	"github.com/bhojpur/state/internal/rpc/core"
	"github.com/bhojpur/state/internal/state"
	"github.com/bhojpur/state/internal/state/indexer"
	"github.com/bhojpur/state/pkg/config"
	"github.com/bhojpur/state/pkg/libs/log"
	"github.com/bhojpur/state/pkg/rpc/jsonrpc/server"
)

// Server defines parameters for running an Inspector rpc server.
type Server struct {
	Addr    string // TCP address to listen on, ":http" if empty
	Handler http.Handler
	Logger  log.Logger
	Config  *config.RPCConfig
}

type eventBusUnsubscriber interface {
	UnsubscribeAll(ctx context.Context, subscriber string) error
}

// Routes returns the set of routes used by the Inspector server.
func Routes(cfg config.RPCConfig, s state.Store, bs state.BlockStore, es []indexer.EventSink, logger log.Logger) core.RoutesMap {
	env := &core.Environment{
		Config:     cfg,
		EventSinks: es,
		StateStore: s,
		BlockStore: bs,
		Logger:     logger,
	}
	return core.RoutesMap{
		"blockchain":       server.NewRPCFunc(env.BlockchainInfo),
		"consensus_params": server.NewRPCFunc(env.ConsensusParams),
		"block":            server.NewRPCFunc(env.Block),
		"block_by_hash":    server.NewRPCFunc(env.BlockByHash),
		"block_results":    server.NewRPCFunc(env.BlockResults),
		"commit":           server.NewRPCFunc(env.Commit),
		"validators":       server.NewRPCFunc(env.Validators),
		"tx":               server.NewRPCFunc(env.Tx),
		"tx_search":        server.NewRPCFunc(env.TxSearch),
		"block_search":     server.NewRPCFunc(env.BlockSearch),
	}
}

// Handler returns the http.Handler configured for use with an Inspector server. Handler
// registers the routes on the http.Handler and also registers the websocket handler
// and the CORS handler if specified by the configuration options.
func Handler(rpcConfig *config.RPCConfig, routes core.RoutesMap, logger log.Logger) http.Handler {
	mux := http.NewServeMux()
	wmLogger := logger.With("protocol", "websocket")

	var eventBus eventBusUnsubscriber

	websocketDisconnectFn := func(remoteAddr string) {
		err := eventBus.UnsubscribeAll(context.Background(), remoteAddr)
		if err != nil && err != pubsub.ErrSubscriptionNotFound {
			wmLogger.Error("Failed to unsubscribe addr from events", "addr", remoteAddr, "err", err)
		}
	}
	wm := server.NewWebsocketManager(logger, routes,
		server.OnDisconnect(websocketDisconnectFn),
		server.ReadLimit(rpcConfig.MaxBodyBytes))
	mux.HandleFunc("/websocket", wm.WebsocketHandler)

	server.RegisterRPCFuncs(mux, routes, logger)
	var rootHandler http.Handler = mux
	if rpcConfig.IsCorsEnabled() {
		rootHandler = addCORSHandler(rpcConfig, mux)
	}
	return rootHandler
}

func addCORSHandler(rpcConfig *config.RPCConfig, h http.Handler) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: rpcConfig.CORSAllowedOrigins,
		AllowedMethods: rpcConfig.CORSAllowedMethods,
		AllowedHeaders: rpcConfig.CORSAllowedHeaders,
	})
	h = corsMiddleware.Handler(h)
	return h
}

// ListenAndServe listens on the address specified in srv.Addr and handles any
// incoming requests over HTTP using the Inspector rpc handler specified on the server.
func (srv *Server) ListenAndServe(ctx context.Context) error {
	listener, err := server.Listen(srv.Addr, srv.Config.MaxOpenConnections)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	return server.Serve(ctx, listener, srv.Handler, srv.Logger, serverRPCConfig(srv.Config))
}

// ListenAndServeTLS listens on the address specified in srv.Addr. ListenAndServeTLS handles
// incoming requests over HTTPS using the Inspector rpc handler specified on the server.
func (srv *Server) ListenAndServeTLS(ctx context.Context, certFile, keyFile string) error {
	listener, err := server.Listen(srv.Addr, srv.Config.MaxOpenConnections)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		listener.Close()
	}()
	return server.ServeTLS(ctx, listener, srv.Handler, certFile, keyFile, srv.Logger, serverRPCConfig(srv.Config))
}

func serverRPCConfig(r *config.RPCConfig) *server.Config {
	cfg := server.DefaultConfig()
	cfg.MaxBodyBytes = r.MaxBodyBytes
	cfg.MaxHeaderBytes = r.MaxHeaderBytes
	// If necessary adjust global WriteTimeout to ensure it's greater than
	// TimeoutBroadcastTxCommit.
	if cfg.WriteTimeout <= r.TimeoutBroadcastTxCommit {
		cfg.WriteTimeout = r.TimeoutBroadcastTxCommit + 1*time.Second
	}
	return cfg
}
