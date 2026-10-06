package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	matchmaking "npln.nintendo.net/npln-practice/proto/matchmaking/v1"
	auth "npln.nintendo.net/npln-practice/proto/auth/v1"
)

const (
	grpcAddr = "127.0.0.1:443"

	healthAddr = "127.0.0.1:8080"

	tenantID = "t-e1c218b5-lp1"

	certFile = "certs/server.crt"
	keyFile  = "certs/server.key"
)

// logUnaryRPC registra cada chamada unary gRPC recebida pelo laboratório.
func logUnaryRPC(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	log.Printf("[RPC] RECEBIDO: %s", info.FullMethod)

	resp, err := handler(ctx, req)

	if err != nil {
		log.Printf("[RPC] ERRO: %s -> %v", info.FullMethod, err)
		return resp, err
	}

	log.Printf("[RPC] SUCESSO: %s", info.FullMethod)

	return resp, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("==============================================")
	log.Println(" MHRise NPLN Lab v0.2")
	log.Println("==============================================")

	// ------------------------------------------------------------
	// 1. Verificar certificado
	// ------------------------------------------------------------

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Printf("[TLS] Certificado não encontrado: %s", certFile)
		log.Println("[TLS] Gerando certificado do laboratório...")

		if err := GenerateCertificate(); err != nil {
			log.Fatalf("[TLS] Falha ao gerar certificado: %v", err)
		}

		log.Println("[TLS] Certificado gerado com sucesso")

	} else if err != nil {
		log.Fatalf("[TLS] Erro verificando certificado: %v", err)
	}

	// ------------------------------------------------------------
	// 2. Carregar certificado TLS
	// ------------------------------------------------------------

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("[TLS] Falha ao carregar certificado: %v", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{
			cert,
		},

		// gRPC utiliza HTTP/2.
		NextProtos: []string{"h2"},
	}

	// ------------------------------------------------------------
	// 3. Criar servidor gRPC
	// ------------------------------------------------------------

	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),

		// Registrar todas as chamadas unary.
		grpc.UnaryInterceptor(logUnaryRPC),
	)

	// ------------------------------------------------------------
	// 4. Registrar serviços
	// ------------------------------------------------------------

	// IMPORTANTE:
	// Mantemos exatamente os construtores que já existem
	// no nosso projeto.

	matchmaking.RegisterMatchmakerServer(
		grpcServer,
		NewMatchmakerServer(),
	)

	matchmaking.RegisterGameSessionServiceServer(
		grpcServer,
		NewGameSessionServer(),
	)

	auth.RegisterAuthServer(grpcServer, NewAuthServer())

	// ------------------------------------------------------------
	// 5. Abrir porta TCP 443
	// ------------------------------------------------------------

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("[gRPC] Falha ao abrir %s: %v", grpcAddr, err)
	}

	// ------------------------------------------------------------
	// 6. Health check HTTP
	// ------------------------------------------------------------

	go func() {
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK\n"))
		})

		log.Printf("[HTTP] health: http://%s/healthz", healthAddr)

		if err := http.ListenAndServe(healthAddr, nil); err != nil {
			log.Printf("[HTTP] health server encerrado: %v", err)
		}
	}()

	// ------------------------------------------------------------
	// 7. Informações do laboratório
	// ------------------------------------------------------------

	log.Println("----------------------------------------------")
	log.Printf("[LAB] gRPC/TLS : %s", grpcAddr)
	log.Printf("[LAB] Health    : %s", healthAddr)
	log.Printf("[LAB] Tenant    : %s", tenantID)
	log.Printf("[LAB] TLS ALPN  : h2")
	log.Println("----------------------------------------------")

	// ------------------------------------------------------------
	// 8. Iniciar servidor
	// ------------------------------------------------------------

	log.Printf("[gRPC] Servidor aguardando conexões em %s", grpcAddr)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("[gRPC] Servidor encerrado: %v", err)
	}
}