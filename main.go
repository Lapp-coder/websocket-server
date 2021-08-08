package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"html/template"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	indexTemplate = template.Must(template.ParseFiles("index.html"))
	hub := newHub()
	handler := newHandler(hub)
	handler.initRoutes()
	go hub.listen()

	host := os.Getenv("CHAT_HOST")
	port := os.Getenv("CHAT_PORT")
	if host == "" || port == "" {
		logrus.Fatalf("incorrect address to start the server")
	}

	// Запуск сервера в go-рутине для его плавной остановки
	srv := newServer(host + ":" + port)
	go func() {
		if err := srv.start(); err != nil {
			logrus.Errorf("failed to start server: %s", err.Error())
		}
	}()

	logrus.Info("server started")

	// Ожидание на получение одного из системных сигналов (SIGINT, SIGTERM) для продолжение выполнения функции main
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	logrus.Info("server shutdown")

	// Плавная остановка сервера
	if err := srv.shutdown(context.Background()); err != nil {
		logrus.Errorf("failed to graceful shutdown server: %s", err.Error())
	}
}
