package controller

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-set/v3"
	"github.com/vietanhduong/github-notifications/pkg/github"
	"github.com/vietanhduong/github-notifications/pkg/logging"
)

// DefaultFetchInterval is the default interval to fetch notifications from GitHub.
// By default, GitHub set the rate limit for a GitHub token is 5000 requests per hour.
// 5 seconds / request = 720 requests per hour
const DefaultFetchInterval = 5 * time.Second

var supportedReasons = set.From[string]([]string{
	"approval_requested",
	"assign",
	"author",
	"comment",
	"ci_activity",
	"invitation",
	"manual",
	"member_feature_requested",
	"mention",
	"review_requested",
	"security_alert",
	"security_advisory_credit",
	"state_change",
	"subscribed",
	"team_mention",
})

var log = logging.WithField("pkg", "github_notifications/controller")

type Controller struct {
	gh            github.Interface
	fetchInterval time.Duration
	listenReasons *set.Set[string]
	client        *http.Client
	lastFetch     time.Time
}

func New(gh github.Interface, opt ...Option) *Controller {
	c := &Controller{
		gh:            gh,
		fetchInterval: 5 * time.Minute,
		listenReasons: set.From[string](nil),
		client:        http.DefaultClient,
	}

	for _, o := range opt {
		o(c)
	}

	if c.listenReasons.Empty() {
		c.listenReasons = supportedReasons.Copy()
	}
	return c
}

func (c *Controller) Run(stop <-chan struct{}) error {
	log.WithField("fetch_interval", c.fetchInterval).Trace("Start fetching notifications")
	ticker := time.NewTicker(c.fetchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			log.Info("Stop fetching notifications")
			return nil
		case <-ticker.C:

		}
	}
}
