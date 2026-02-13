package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"os/signal"
	"syscall"

	"gin-alpine/src/pkg/utils"
	"gin-alpine/src/services/worker/internal/bootstrap"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"
)

type Server struct {
	Srv       *asynq.Server
	Bootstrap *bootstrap.Bootstrap
}

func Run() {
	srv := NewConsumerWorkerServer()
	srv.runConsumerWorkerServer()
}

func NewConsumerWorkerServer() *Server {
	b := bootstrap.MustGetBootstrapInstance()
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: b.Config.RedisURL,
		},
		asynq.Config{
			Concurrency: 10,
			Logger: &AsynqStdLogger{
				log.Default(),
			},
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
	return &Server{
		Srv:       srv,
		Bootstrap: b,
	}
}

func (s *Server) runConsumerWorkerServer() {
	mux := asynq.NewServeMux()
	mux.HandleFunc(utils.TaskSendingEmailType, s.EmailTaskHandler)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	log.Println("Starting consummer worker... ")
	go func() {
		if err := s.Srv.Run(mux); err != nil {
			log.Fatalf("asynq stopped: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down notifications worker...")
	s.Srv.Shutdown()
}

func (s *Server) EmailTaskHandler(ctx context.Context, t *asynq.Task) error {
	var p utils.EmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("fail to unmarshal email task %v: %v", err, asynq.SkipRetry)
	}
	auth := smtp.PlainAuth(
		"",
		s.Bootstrap.Config.EmailSender,
		s.Bootstrap.Config.EmailSenderPass,
		s.Bootstrap.Config.EmailSMTP,
	)
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
	msg := fmt.Sprintf("Subject: %s\n%s\n\n%s", p.Subject, headers, p.Body)
	err := smtp.SendMail(
		s.Bootstrap.Config.SMTPAddress,
		auth,
		s.Bootstrap.Config.EmailSender,
		*p.Addresses,
		[]byte(msg),
	)
	if err == nil {
		s.Bootstrap.Logger.Info(
			utils.TaskSendEmail,
			zap.String("to", (*p.Addresses)[0]),
			zap.String("subject", p.Subject),
		)
		log.Printf("email sended to: %s - subject: %s\n", (*p.Addresses)[0], p.Subject)
	} else {
		s.Bootstrap.Logger.Error(utils.TaskSendEmail, zap.String("error", err.Error()))
		log.Printf("error sending email: %v", err)
	}
	return err
}

type AsynqStdLogger struct {
	l *log.Logger
}

func (l *AsynqStdLogger) Debug(args ...any) { l.l.Println(args...) }
func (l *AsynqStdLogger) Info(args ...any)  { l.l.Println(args...) }
func (l *AsynqStdLogger) Warn(args ...any)  { l.l.Println(args...) }
func (l *AsynqStdLogger) Error(args ...any) { l.l.Println(args...) }
func (l *AsynqStdLogger) Fatal(args ...any) { l.l.Fatal(args...) }
