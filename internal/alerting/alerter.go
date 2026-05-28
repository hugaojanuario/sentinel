package alerting

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hugaojanuario/sentinel/internal/healthcheck"
)

type Alerter struct {
	botToken string
	chatID   string
	events   <-chan healthcheck.Event
	states   map[string]string
	mu       sync.Mutex
}

func NewAlerter(botToken, chatID string, events <-chan healthcheck.Event) *Alerter {
	return &Alerter{
		botToken: botToken,
		chatID:   chatID,
		events:   events,
		states:   make(map[string]string),
	}
}

func (a *Alerter) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-a.events:
			a.handle(event)
		}
	}
}

func (a *Alerter) handle(event healthcheck.Event) {
	key := stateKey(event)

	a.mu.Lock()
	prev := a.states[event.ContainerID]
	if prev == key {
		a.mu.Unlock()
		return
	}
	a.states[event.ContainerID] = key
	a.mu.Unlock()

	msg := buildMessage(event, prev)
	if msg == "" {
		return
	}

	_ = a.sendTelegram(msg)
}

func stateKey(event healthcheck.Event) string {
	switch event.Status {
	case "exited", "dead":
		if event.OOMKilled {
			return fmt.Sprintf("oom:%d", event.ExitCode)
		}
		return fmt.Sprintf("exited:%d", event.ExitCode)
	case "restarting":
		return "restarting"
	case "running":
		return "running"
	default:
		return event.Status
	}
}

func buildMessage(event healthcheck.Event, prevKey string) string {
	ts := event.CheckedAt.Format("02/01/2006 15:04:05")

	switch event.Status {
	case "exited", "dead":
		if event.OOMKilled {
			return fmt.Sprintf(
				"🔴 <b>%s CAIU</b>\nMotivo: OOM Kill (sem memória)\nExit code: %d\nParou em: %s",
				event.ContainerName, event.ExitCode, formatTime(event.FinishedAt),
			)
		}
		reason := exitReason(event)
		return fmt.Sprintf(
			"🔴 <b>%s CAIU</b>\nMotivo: %s\nExit code: %d\nParou em: %s",
			event.ContainerName, reason, event.ExitCode, formatTime(event.FinishedAt),
		)
	case "restarting":
		return fmt.Sprintf(
			"⚠️ <b>%s está reiniciando</b>\nO container entrou em restart loop.\nÚltimo exit code: %d\nVerificado em: %s",
			event.ContainerName, event.ExitCode, ts,
		)
	case "running":
		if prevKey != "" && prevKey != "running" {
			return fmt.Sprintf(
				"✅ <b>%s voltou</b>\nO container se recuperou e está saudável novamente.\nEm: %s",
				event.ContainerName, ts,
			)
		}
		return ""
	default:
		return ""
	}
}

func exitReason(event healthcheck.Event) string {
	if event.ExitError != "" {
		return event.ExitError
	}
	return "saída inesperada"
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "desconhecido"
	}
	return t.Format("02/01/2006 15:04:05")
}

type telegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func (a *Alerter) sendTelegram(text string) error {
	if a.botToken == "" || a.chatID == "" {
		return nil
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", a.botToken)

	payload := telegramMessage{
		ChatID:    a.chatID,
		Text:      text,
		ParseMode: "HTML",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram respondeu com status %d", resp.StatusCode)
	}

	return nil
}
