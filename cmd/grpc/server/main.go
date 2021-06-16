package main

import (
	"fmt"
	"github.com/aleale2121/DSP_LAB/Music_Service/constant"
	"github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/module/artist/artist_server_service"
	"github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/module/auth"
	"github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/module/song/song_server_service"
	artistProto "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/artist"
	songProto "github.com/aleale2121/DSP_LAB/Music_Service/grpc/server/services/song"
	pgArtist "github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/artist"
	pgSong "github.com/aleale2121/DSP_LAB/Music_Service/storage/db_store/postgres/song"
	"github.com/aleale2121/DSP_LAB/Music_Service/storage/file_store"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
)

func main() {
	if err := RunServer(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func RunServer() error {
	connStr, dialect, err := constant.GetGormDatabaseConnectionString()
	if err != nil {
		panic(err)
	}

	dbConn, err := gorm.Open(dialect,
		connStr)
	if dbConn != nil {
		defer dbConn.Close()
	}
	if err != nil {
		panic(err)
	}
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("cannot get base path: %v", err)
	}
	store, err := file_store.NewStorage(basePath)
	if err != nil {
		log.Fatalf("cannot create storage: %v", err)
	}

	postgresArtist := pgArtist.NewArtistGormRepo(dbConn)
	artistService := artist_server_service.NewGrpcArtistAServer(postgresArtist)


	postgresSong := pgSong.NewSongGormRepo(dbConn)
	songService := song_server_service.NewGrpcSongServer(postgresSong, postgresArtist,*store)

	interceptor := auth.NewAuthInterceptor()
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	}

	gs := grpc.NewServer(serverOptions...)


	artistProto.RegisterArtistServiceServer(gs, artistService)
	songProto.RegisterSongServiceServer(gs, songService)

	reflection.Register(gs)

	grpcConnStr, err := constant.GetGrpcConnectionString()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", grpcConnStr)
	if err != nil {
		return err
	}

	err = gs.Serve(l)
	if err != nil {
		return err
	}
	return nil
}
