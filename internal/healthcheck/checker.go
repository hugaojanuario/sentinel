package healthcheck

import (
	"context"
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
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.checkAll(ctx)
		}
	}
}

func (c *Checker) checkAll(ctx context.Context) {
	containers, err := docker.ListAllContainers()
	if err != nil {
		return
	}

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
			}
		}(cont)
	}

	wg.Wait()
}
