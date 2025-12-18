package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/pur1fying/GO_BAAS/internal/config"
	"github.com/pur1fying/GO_BAAS/internal/database"
	"github.com/pur1fying/GO_BAAS/internal/mail"

	"github.com/pur1fying/GO_BAAS/internal/global_info"
	"github.com/pur1fying/GO_BAAS/internal/logger"
)

func ResponseWithJson(w http.ResponseWriter, statusCode int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		logger.BAASError(fmt.Sprintf("JSON 序列化失败: %v", err))
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	ResponseWithJson(w, 200, map[string]string{"status": "ready"})
}

func main() {
	global_info.InitGlobalInfo()
	logger.InitGlobalLogger()

	// Global Config Load
	err := config.Load("")
	if err != nil {
		logger.BAASCritical("Failed to load config:", err.Error())
	}

	// DB
	err = database.Init()
	if err != nil {
		logger.BAASCritical("Failed to init database:", err.Error())
	}

	// Mail Smtp Init
	err = mail.InitMail()
	if err != nil {
		logger.BAASCritical("Failed to init mail:", err.Error())
	}

	message := &mail.Message{
		From:    mail.Address{Email: "2274916027@qq.com", DisplayName: "测试发送者1"},
		ReplyTo: []mail.Address{},
		TO:      []mail.Address{mail.Address{Email: "2274916027@qq.com", DisplayName: "测试接受者2"}, mail.Address{Email: "283829027@qq.com", DisplayName: "1sy"}},
		CC:      []mail.Address{mail.Address{Email: "2022280486@email.szu.edu.cn", DisplayName: "szdx测试接收者"}},
		BCC:     []string{},

		Subject: "测试邮件",
		//Text:    "woshixuqihao",
		Html: "<h1>这是一封测试邮件。</h1>",
	}
	err = mail.SendMail(message)
	if err != nil {
		logger.BAASCritical("Failed to send mail:", err.Error())
	}

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.HandleFunc("/readiness", handlerReadiness)

	svr := &http.Server{
		Handler: router,
		Addr:    ":8080",
	}
	logger.Flush()
	logger.HighLight("Starting Server on :8080")
	svr.ListenAndServe()
	return
}
