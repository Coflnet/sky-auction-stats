package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	NotifierOperatorAND = "and"
	NotifierOperatorOR  = "or"

	NotifierStateOK       = "ok"
	NotifierStatePending  = "pending"
	NotifierStateAlerting = "alerting"

	NotifierTemplateOperatorBigger  = ">"
	NotifierTemplateOperatorSmaller = "<"

	NotifierTemplateKeyAuction = "auction"
	NotifierTemplateKeyFlip    = "flip"
)

type Notifier struct {
	ID                 primitive.ObjectID  `json:"id" bson:"_id"`
	Active             bool                `json:"active" bson:"active"`
	UserId             int                 `json:"userId" bson:"user_id"`
	Name               string              `json:"name" bson:"name,omitempty"`
	Description        string              `json:"description" bson:"description,omitempty"`
	NotifierStates     []*NotifierState    `json:"notifierStates" bson:"notifier_states"`
	NotifierTemplates  []*NotifierTemplate `json:"notifierTemplates" bson:"notifier_templates"`
	TemplateOperator   string              `json:"templateOperator" bson:"template_operator"`
	AlertText          string              `json:"alertText" bson:"alert_text,omitempty"`
	LastEvaluation     time.Time           `json:"lastEvaluation" bson:"last_evaluation"`
	NextEvaluation     time.Time           `json:"nextEvaluation" bson:"next_evaluation"`
	EvaluationInterval int                 `json:"evaluationInterval" bson:"evaluation_interval"`

	// TimeUntilTrigger like the grafana pending state time
	TimeUntilTrigger time.Duration `json:"timeUntilTrigger" bson:"time_until_trigger"`
}

type NotifierTemplate struct {
	Operator        string        `json:"operator" bson:"operator"`
	Value           int64         `json:"value" bson:"value"`
	DurationToCheck time.Duration `json:"durationToCheck" bson:"duration_to_check"`
	Key             string        `json:"key" bson:"key"`
}

type NotifierState struct {
	State     string    `json:"state" bson:"state"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
