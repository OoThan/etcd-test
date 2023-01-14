package main

import (
	"context"
	"ercd-test/cmd/user/server"
	"ercd-test/internal/conf"
	"ercd-test/internal/logger"
	"ercd-test/internal/pb"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//port := flag.String("port", "8070", "UserRPC Default port is 8070.")
	//flag.Parse()

	//addr := fmt.Sprintf("0.0.0.0:%s", *port)

	conf.InitYaml()

	//  grpc start
	gs := grpc.NewServer()

	// create user service instance
	userSvc := server.NewUserService()
	pb.RegisterUserServiceServer(gs, userSvc)

	// listening
	//addr := fmt.Sprintf("0.0.0.0:%s", *port)
	//l, err := net.Listen("tcp", addr)
	//if err != nil {
	//	logger.Logrus.Fatalf("Failed to Listen User RPC server: %v", err)
	//}
	//defer l.Close()

	//// service register with etcd
	//cli, err := clientv3.New(clientv3.Config{
	//	Endpoints: []string{"127.0.0.1:2379"},
	//})
	//if err != nil {
	//	logger.Logrus.Fatalf("Failed to connet with etcd: %v", err)
	//}
	//
	//em, err := endpoints.NewManager(cli, "userRPC")
	//if err != nil {
	//	logger.Logrus.Panic(err)
	//}
	//
	//if err := em.AddEndpoint(context.Background(), "userRPC/"+addr, endpoints.Endpoint{
	//	Addr: addr,
	//}); err != nil {
	//	logger.Logrus.Panic(err)
	//}

	// register endpoint etcd
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: conf.Etcd().Hosts,
	})
	if err != nil {
		logger.Logrus.Error(err)
		panic(err)
	}

	em, err := endpoints.NewManager(cli, conf.RPCSvc().UserRPC.Name)
	if err != nil {
		logger.Logrus.Error(err)
		panic(err)
	}

	for _, host := range conf.RPCSvc().UserRPC.Hosts {

		//addr := fmt.Sprintf("0.0.0.0:%s", host)
		l, err := net.Listen("tcp", host)
		if err != nil {
			logger.Logrus.Fatalf("Failed to Listen User RPC server: %v", err)
		}

		// start grpc server
		go func() {
			logger.Logrus.Println("User RPC server is starting on: ", host)
			if err := gs.Serve(l); err != nil {
				logger.Logrus.Fatalf("Failed to start server: %v", err)
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		key := fmt.Sprintf("%s/%s", conf.RPCSvc().UserRPC.Name, host)
		err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
			Addr: host,
		})
		if err != nil {
			logger.Logrus.Error(err)
			panic(err)
		}
	}

	userSvc.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c

	cli.Close()
	gs.Stop()
}
