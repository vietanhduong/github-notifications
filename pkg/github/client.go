package github

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v64/github"
	"github.com/samber/lo"
	"github.com/vietanhduong/github-notifications/pkg/logging"
)

var log = logging.WithField("pkg", "pkg/github")

type Interface interface {
	FetchNotifications(ctx context.Context, opt FetchNotificationsOptions) ([]*Notification, error)
	MarkNotificationAs(ctx context.Context, id string, done bool) error
}

type Client struct {
	token string
	gh    *github.Client
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

type FetchNotificationsOptions struct {
	// Only show results that were last updated after the given time
	Since time.Time
	// Only show notifications updated before the given time.
	Before time.Time
	// If true, show notifications marked as read.
	All bool
}

func (c *Client) FetchNotifications(ctx context.Context, opt FetchNotificationsOptions) ([]*Notification, error) {
	notifications, _, err := c.gh.Activity.ListNotifications(ctx, &github.NotificationListOptions{
		Since:  opt.Since,
		Before: opt.Before,
		All:    opt.All,
	})
	if err != nil {
		log.WithError(err).Debug("Failed to fetch notifications")
		return nil, fmt.Errorf("fetch notifications: %w", err)
	}
	return lo.Map(notifications, func(n *github.Notification, _ int) *Notification {
		return &Notification{
			Id:         n.GetID(),
			Reason:     n.GetReason(),
			Repository: n.GetRepository().GetFullName(),
			Subject:    n.GetSubject().GetTitle(),
			Unread:     n.GetUnread(),
			UpdatedAt:  n.GetUpdatedAt().Time,
		}
	}), nil
}

// MarkNotificationAs marks a notification as read or done. Default is read.
func (c *Client) MarkNotificationAs(ctx context.Context, id string, done bool) error {
	method := "PATCH"
	if done {
		method = "DELETE"
	}

	u := fmt.Sprintf("notifications/threads/%v", id)

	req, err := c.gh.NewRequest(method, u, nil)
	if err != nil {
		log.WithError(err).Debug("Failed to create new request")
		return fmt.Errorf("github new request: %w", err)
	}

	if _, err = c.gh.Do(ctx, req, nil); err != nil {
		log.WithError(err).Debug("Failed to mark notification")
		return fmt.Errorf("mark notification: %w", err)
	}
	return nil
}

func (c *Client) CurrentUser(ctx context.Context) (*User, error) {
	user, _, err := c.gh.Users.Get(ctx, "")
	if err != nil {
		log.WithError(err).Debug("Failed to get current user")
		return nil, fmt.Errorf("get current user: %w", err)
	}
	resp := &User{Login: user.GetLogin()}
	return resp, nil
}
