package controller

import (
	log "github.com/sirupsen/logrus"

	. "github.com/flant/shell-operator/pkg/hook/binding_context"
	. "github.com/flant/shell-operator/pkg/hook/types"
	. "github.com/flant/shell-operator/pkg/webhook/mutating/types"

	"github.com/flant/shell-operator/pkg/webhook/mutating"
)

// A link between a hook and a kube monitor
type MutatingBindingToWebhookLink struct {
	BindingName     string
	ConfigurationId string
	WebhookId       string
	// Useful fields to create a BindingContext
	IncludeSnapshots []string
	Group            string
}

// ScheduleBindingsController handles schedule bindings for one hook.
type MutatingBindingsController interface {
	WithMutatingBindings([]MutatingConfig)
	WithWebhookManager(*mutating.WebhookManager)
	EnableMutatingBindings()
	DisableMutatingBindings()
	CanHandleEvent(event MutatingEvent) bool
	HandleEvent(event MutatingEvent) BindingExecutionInfo
}

type mutatingBindingsController struct {
	// Controller holds mutating bindings from one hook. Hook always belongs to one configurationId.
	ConfigurationId string
	// WebhookId -> link
	MutatingLinks map[string]*MutatingBindingToWebhookLink

	MutatingBindings []MutatingConfig

	webhookManager *mutating.WebhookManager
}

var _ MutatingBindingsController = &mutatingBindingsController{}

// NewKubernetesHooksController returns an implementation of KubernetesHooksController
var NewMutatingBindingsController = func() *mutatingBindingsController {
	return &mutatingBindingsController{
		MutatingLinks: make(map[string]*MutatingBindingToWebhookLink),
	}
}

func (c *mutatingBindingsController) WithMutatingBindings(bindings []MutatingConfig) {
	c.MutatingBindings = bindings
}

func (c *mutatingBindingsController) WithWebhookManager(mgr *mutating.WebhookManager) {
	c.webhookManager = mgr
}

func (c *mutatingBindingsController) EnableMutatingBindings() {
	confId := ""
	for _, config := range c.MutatingBindings {
		if config.Webhook.Metadata.ConfigurationId == "" && confId == "" {
			continue
		}
		if config.Webhook.Metadata.ConfigurationId != "" && confId == "" {
			confId = config.Webhook.Metadata.ConfigurationId
			continue
		}
		if config.Webhook.Metadata.ConfigurationId != confId {
			log.Errorf("Possible bug!!! kubernetesMutating has non-unique configurationIds: '%s' '%s'", config.Webhook.Metadata.ConfigurationId, confId)
		}
	}
	c.ConfigurationId = confId

	for _, config := range c.MutatingBindings {
		c.MutatingLinks[config.Webhook.Metadata.WebhookId] = &MutatingBindingToWebhookLink{
			BindingName:      config.BindingName,
			ConfigurationId:  c.ConfigurationId,
			WebhookId:        config.Webhook.Metadata.WebhookId,
			IncludeSnapshots: config.IncludeSnapshotsFrom,
			Group:            config.Group,
		}
		c.webhookManager.AddWebhook(config.Webhook)
	}
}

func (c *mutatingBindingsController) DisableMutatingBindings() {
	// TODO dynamic enable/disable mutating webhooks.
}

func (c *mutatingBindingsController) CanHandleEvent(event MutatingEvent) bool {
	if c.ConfigurationId != event.ConfigurationId {
		return false
	}
	_, has := c.MutatingLinks[event.WebhookId]
	return has
}

func (c *mutatingBindingsController) HandleEvent(event MutatingEvent) BindingExecutionInfo {
	if c.ConfigurationId != event.ConfigurationId {
		log.Errorf("Possible bug!!! Unknown mutating event: no binding for configurationId '%s' (webhookId '%s')", event.ConfigurationId, event.WebhookId)
		return BindingExecutionInfo{
			BindingContext: []BindingContext{},
			AllowFailure:   false,
		}
	}

	link, hasKey := c.MutatingLinks[event.WebhookId]
	if !hasKey {
		log.Errorf("Possible bug!!! Unknown mutating event: no binding for configurationId '%s', webhookId '%s'", event.ConfigurationId, event.WebhookId)
		return BindingExecutionInfo{
			BindingContext: []BindingContext{},
			AllowFailure:   false,
		}
	}

	bc := BindingContext{
		Binding:         link.BindingName,
		AdmissionReview: event.Review,
	}
	bc.Metadata.BindingType = KubernetesMutating
	bc.Metadata.IncludeSnapshots = link.IncludeSnapshots
	bc.Metadata.Group = link.Group

	return BindingExecutionInfo{
		BindingContext:   []BindingContext{bc},
		Binding:          link.BindingName,
		IncludeSnapshots: link.IncludeSnapshots,
		Group:            link.Group,
	}
}
