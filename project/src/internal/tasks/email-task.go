// Package tasks
package tasks

import (
	"bytes"
	"encoding/json"
	"html/template"

	"gin-alpine/src/pkg/utils"

	"github.com/hibiken/asynq"
)

func NewEmailTask(pl *utils.EmailPayload) (*asynq.Task, error) {
	payload, err := json.Marshal(pl)
	if err != nil {
		return nil, err
	}
	// NewTask must have the same tag used by the handler
	// ex: mux.HandleFunc(utils.TaskSendingEmailType, s.EmailTaskHandler)
	task := asynq.NewTask(utils.TaskSendingEmailType, payload)
	return task, nil
}

func GetEmailTemplate(templateName string) (*template.Template, error) {
	htmlTemplateFile, err := utils.GetFilePath(
		[]string{"src", "services", "notifications", "internal", "templates", templateName},
	)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.ParseFiles(htmlTemplateFile)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func RenderHelloTemplate(tmpl *template.Template, vars *utils.EmailHelloVars) (*bytes.Buffer, error) {
	var emailRendered bytes.Buffer
	if err := tmpl.Execute(&emailRendered, *vars); err != nil {
		return nil, err
	}
	return &emailRendered, nil
}

func RenderResetPasswordTemplate(tmpl *template.Template, vars *utils.ResetPasswordVars) (*bytes.Buffer, error) {
	var emailRendered bytes.Buffer
	if err := tmpl.Execute(&emailRendered, *vars); err != nil {
		return nil, err
	}
	return &emailRendered, nil
}

func RenderWelcomeTemplate(tmpl *template.Template, vars *utils.WelcomeEmailVars) (*bytes.Buffer, error) {
	var emailRendered bytes.Buffer
	if err := tmpl.Execute(&emailRendered, *vars); err != nil {
		return nil, err
	}
	return &emailRendered, nil
}
