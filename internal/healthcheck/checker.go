package healthcheck

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/hugaojanuario/sentinel/internal/docker"
)

type Event struct {
	ContainerID   string
	ContainerName string
	Status        string
	OOMKilled     bool
	ExitCode      int
	ExitError     string
	FinishedAt    time.Time
	CheckedAt     time.Time
}

type Checker struct {
	interval time.Duration
	events   chan Event
}

func NewChecker(interval time.Duration) *Checker {
	return &Checker{
		interval: interval,
		events:   make(chan Event, 100),
	}
}

func (c *Checker) Events() <-chan Event {
	return c.events
}

func (c *Checker) Run(ctx context.Context) {
	log.Printf("[checker] iniciado com intervalo=%s", c.interval)
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("[checker] encerrando")
			return
		case <-ticker.C:
			c.checkAll(ctx)
		}
	}
}

func (c *Checker) checkAll(ctx context.Context) {
	containers, err := docker.ListAllContainers()
	if err != nil {
		log.Printf("[checker] erro ao listar containers: %v", err)
		return
	}
	log.Printf("[checker] verificando %d containers", len(containers))

	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for _, cont := range containers {
		wg.Add(1)
		go func(info docker.ContainerInfo) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			case sem <- struct{}{}:
			}
			defer func() { <-sem }()

			state, err := docker.InspectContainer(info.ID)
			if err != nil {
				log.Printf("[checker] erro ao inspecionar container %s: %v", info.ID[:12], err)
				return
			}

			event := Event{
				ContainerID:   state.ID,
				ContainerName: state.Name,
				Status:        state.Status,
				OOMKilled:     state.OOMKilled,
				ExitCode:      state.ExitCode,
				ExitError:     state.ExitError,
				FinishedAt:    state.FinishedAt,
				CheckedAt:     time.Now(),
			}

			select {
			case c.events <- event:
			default:
				log.Printf("[checker] canal cheio, evento de container=%s descartado", event.ContainerName)
			}
		}(cont)
	}

	wg.Wait()
}
