package main

import (
	"fmt"
	"time"

	"github.com/echo9et/alerting/internal/agent/client"
	"github.com/echo9et/alerting/internal/entities"
)

// Используй флаги сборки
// go build -ldflags "-X main.buildVersion=1.0.0"
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	config, status := GetConfig()
	if !status {
		panic("Не верно проинцелизирован конфиг файл")
	}

	a := client.NewAgent(config.AddrServer, config.SelfIP, config.UseGRPC)
	r := time.Duration(config.ReportTimeout) * time.Second
	p := time.Duration(config.PollTimeout) * time.Second

	if config.CryptoKey != "" {
		pub, err := entities.GetPubKey(config.CryptoKey)
		if err != nil {
			panic(err)
		}
		a.UpdateMetrics(r, p, config.SecretKey, config.RateLimit, pub)
	}

	a.UpdateMetrics(r, p, config.SecretKey, config.RateLimit, nil)
}

// package main

// import (
// 	// ...
// 	"bytes"
// 	"context"
// 	"crypto/rand"
// 	"crypto/rsa"
// 	"crypto/sha256"
// 	"encoding/gob"
// 	"fmt"
// 	"log"

// 	"github.com/echo9et/alerting/internal/entities"
// 	pb "github.com/echo9et/alerting/proto"

// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// 	_ "google.golang.org/grpc/encoding/gzip"
// )

// // Используй флаги сборки
// // go build -ldflags "-X main.buildVersion=1.0.0"
// var (
// 	buildVersion string = "N/A"
// 	buildDate    string = "N/A"
// 	buildCommit  string = "N/A"
// )

// func main() {

// 	fmt.Printf("Build version: %s\n", buildVersion)
// 	fmt.Printf("Build date: %s\n", buildDate)
// 	fmt.Printf("Build commit: %s\n", buildCommit)

// 	config, status := GetConfig()
// 	if !status {
// 		panic("Не верно проинцелизирован конфиг файл")
// 	}

// 	// устанавливаем соединение с сервером
// 	conn, err := grpc.NewClient(":3200",
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 		grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer conn.Close()
// 	c := pb.NewMetricsClient(conn)
// 	pub, err := entities.GetPubKey(config.CryptoKey)
// 	if err != nil {
// 		fmt.Println(config.CryptoKey)
// 		panic(err)
// 	}
// 	TestGRPC(c, pub)
// }

// func TestGRPC(c pb.MetricsClient, pubKey *rsa.PublicKey) {
// 	// набор тестовых данных
// 	metrics := []*pb.Metric{
// 		{Id: "1255111", Type: pb.Metric_GAUGE, Value: 0.1},
// 		{Id: "reqests", Type: pb.Metric_GOUNTER, Delta: 1},
// 	}
// 	for _, metric := range metrics {
// 		fmt.Println(metric.Id)
// 		// добавляем пользователей
// 		resp, err := c.UpdateMetric(context.Background(), &pb.UpdateMetricRequest{
// 			Metric: metric,
// 		})
// 		if err != nil {
// 			log.Fatal("FATAL === ", err)
// 		}
// 		if resp.Error != "" {
// 			fmt.Println(resp.Error)
// 		}
// 	}
// 	{
// 		resp, err := c.UpdateMetrics(context.Background(), &pb.UpdateMetricsRequest{
// 			Metrics: metrics,
// 		})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if resp.Error != "" {
// 			fmt.Println("resp.Error", resp.Error)
// 		}
// 	}

// 	if pubKey != nil {
// 		var data bytes.Buffer // Stand-in for a network connection
// 		enc := gob.NewEncoder(&data)
// 		err := enc.Encode(metrics)
// 		if err != nil {
// 			log.Fatal("encode error:", err)
// 		}
// 		cd, err := rsa.EncryptOAEP(
// 			sha256.New(),
// 			rand.Reader,
// 			pubKey,
// 			data.Bytes(),
// 			nil,
// 		)
// 		if err != nil {
// 			log.Fatal("encode error:", err)
// 		}
// 		resp, err := c.UpdateEncrypteMetrics(context.Background(), &pb.UpdateEncrypteMetricsRequest{
// 			Data: cd,
// 		})
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		if resp.Error != "" {
// 			fmt.Println("resp.Error", resp.Error)
// 		}
// 	}
// }
