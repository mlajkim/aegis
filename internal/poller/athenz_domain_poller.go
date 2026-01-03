package poller

import (
	"context"
	"time"

	"github.com/mlajkim/aegis/internal/syncer"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// RolePoller runs periodically to sync all roles in all namespaces.
// It implements controller-runtime's manager.Runnable interface.
type RolePoller struct {
	Syncer   *syncer.Syncer
	Interval time.Duration
}

// New creates a new RolePoller instance.
func New(s *syncer.Syncer, d time.Duration) *RolePoller {
	return &RolePoller{
		Syncer:   s,
		Interval: d,
	}
}

// Start is called by the Manager when it starts.
// It creates a ticker and runs the sync logic periodically.
func (p *RolePoller) Start(ctx context.Context) error {
	logger := log.FromContext(ctx).WithName("poller").WithName("role")
	logger.Info("Starting Role Poller", "interval", p.Interval)

	ticker := time.NewTicker(p.Interval)
	defer ticker.Stop()

	p.tick(ctx) // Run immediately on start

	for {
		select {
		case <-ticker.C:
			p.tick(ctx)
		case <-ctx.Done():
			logger.Info("Stopping Role Poller")
			return nil
		}
	}
}

// tick executes the business logic via Syncer.
func (p *RolePoller) tick(ctx context.Context) {
	logger := log.FromContext(ctx)

	// Delegate the heavy lifting to Syncer
	if err := p.Syncer.AthenzDomainIntoK8sRb(ctx); err != nil {
		logger.Error(err, "Failed to sync all namespaces")
	} else {
		logger.Info("Successfully polled target athenz domains into cluster", "poller", "athenz-domain-poller", "nextRunIn", p.Interval.String())
	}
}

// NeedLeaderElection ensures this poller only runs on the leader pod.
// This prevents multiple pods from spamming the API simultaneously.
func (p *RolePoller) NeedLeaderElection() bool {
	return true
}
