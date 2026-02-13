package tasks

import (
	"log"
	"testing"

	"gin-alpine/src/pkg/utils"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"
)

func TestEmailTask(t *testing.T) {
	description := "If the send email scenarios are working correctly"
	defer func() {
		log.Printf("Test: %s\n", description)
		log.Println("Deferred tearing down.")
	}()

	redisURL := "localhost:6379"

	t.Run("should send a welcome email", func(t *testing.T) {
		addressTestEmail := "jonas.w.martins@gmail.com"
		tmpl, err := GetEmailTemplate("welcome.html")
		assert.NoError(t, err)

		vars := utils.WelcomeEmailVars{
			Username:       "TestUser",
			LoginURL:       "/",
			UnsubscribeURL: "/",
			Password:       "qwerty-quwety",
			AssetsBaseURL:  "http://localhost:3000",
		}
		welcomePassRendered, err := RenderWelcomeTemplate(tmpl, &vars)
		assert.NoError(t, err)
		payload := utils.EmailPayload{
			Addresses: &[]string{addressTestEmail},
			Subject:   WelcomePasswordSubject,
			Body:      welcomePassRendered.String(),
		}
		task, err := NewEmailTask(&payload)
		assert.NoError(t, err)
		client := asynq.NewClient(asynq.RedisClientOpt{
			Addr: redisURL,
		})
		defer func() {
			if errClose := client.Close(); errClose != nil {
				t.Logf("erro closing asynq %v", errClose)
			}
		}()
		_, err = client.Enqueue(task)
		assert.NoError(t, err)
	})

}
