package usecase

import (
	"fmt"
	"github.com/Coflnet/auction-stats/internal/model"
	"github.com/Coflnet/auction-stats/internal/mongo"
	"github.com/Coflnet/auction-stats/internal/redis"
	"github.com/rs/zerolog/log"
	"sort"
	"time"
)

func ListUserNotifiers(userId int) ([]*model.Notifier, error) {
	notifiers, err := mongo.NotifiersForUser(userId)

	if err != nil {
		log.Error().Err(err).Msgf("error listing notifiers for user %d", userId)
		return nil, err
	}

	return notifiers, nil
}

func CreateNotifier(notifier *model.Notifier) error {

	// set last evaluation if it is not set aka 1.1.1970
	if notifier.LastEvaluation.IsZero() {
		notifier.LastEvaluation = time.Now()
	}

	// same as above
	if notifier.NextEvaluation.IsZero() {
		notifier.NextEvaluation = time.Now()
	}

	// init states
	if notifier.NotifierStates == nil {
		notifier.NotifierStates = make([]*model.NotifierState, 0)
	}

	return mongo.InsertNotifier(notifier)
}

func UpdateNotifier(notifier *model.Notifier) error {
	return mongo.ReplaceNotifier(notifier)
}

func DeleteNotifier(notifier *model.Notifier) error {
	return mongo.DeleteNotifier(notifier)
}

func StartNotifierSchedule() {
	go func() {
		for range time.Tick(time.Second * 10) {
			err := CheckNotifiers()
			if err != nil {
				log.Error().Err(err).Msg("error checking notifiers")
			}
		}
	}()
}

func CheckNotifiers() error {
	notifiers, err := mongo.NotifierToEvaluate()
	if err != nil {
		log.Error().Err(err).Msgf("error listing notifiers to evaluate")
		return err
	}

	for notifier := range notifiers {
		go func(n *model.Notifier) {
			if !n.Active {
				return
			}

			err := CheckNotifier(n)
			if err != nil {
				log.Error().Err(err).Msgf("error checking notifier, %v", n)
			}
		}(notifier)
	}

	return nil
}

func CheckNotifier(notifier *model.Notifier) error {

	alertTriggered, err := evaluateTemplates(notifier)
	if err != nil {
		log.Error().Err(err).Msgf("error evaluating templates for notifier %v", notifier.ID)
		return err
	}

	stateChanged, newState, err := newStateOfNotifier(notifier, alertTriggered)

	if err != nil {
		log.Error().Err(err).Msgf("error getting state for notifier, %v", notifier)
		return err
	}

	if !stateChanged {
		log.Info().Msgf("state for notifier %v is unchanged at %v do nothing", notifier.ID, newState)
		return nil
	}

	state := &model.NotifierState{
		State:     newState,
		Timestamp: time.Now(),
	}

	notifier.NotifierStates = append(notifier.NotifierStates, state)
	notifier.LastEvaluation = time.Now()
	notifier.NextEvaluation = time.Now().Add(notifier.TimeUntilTrigger)

	sort.Slice(notifier.NotifierStates, func(i, j int) bool {
		return notifier.NotifierStates[i].Timestamp.After(notifier.NotifierStates[j].Timestamp)
	})

	notifier.NotifierStates = notifier.NotifierStates[:100]

	err = mongo.ReplaceNotifier(notifier)
	if err != nil {
		log.Error().Err(err).Msgf("error replacing notifier, %v", notifier)
		return err
	}

	log.Warn().Msgf("logic to send the update is not implemented yet")

	return nil
}

func newStateOfNotifier(notifier *model.Notifier, alertTriggered bool) (bool, string, error) {

	if len(notifier.NotifierStates) == 0 {
		return true, model.NotifierStateOK, nil
	}

	lastState := notifier.NotifierStates[0]
	for _, state := range notifier.NotifierStates {
		if state.Timestamp.After(lastState.Timestamp) {
			lastState = state
		}
	}

	if lastState.State == model.NotifierStateOK {
		if alertTriggered {
			log.Info().Msgf("notifier %v goes from ok to pending", notifier.ID)
			return true, model.NotifierStatePending, nil
		}

		return false, model.NotifierStateOK, nil
	}

	if lastState.State == model.NotifierStateAlerting {
		if alertTriggered {
			return false, model.NotifierStateAlerting, nil
		}

		log.Info().Msgf("notifier %v goes from alerting to ok", notifier.ID)
		return true, model.NotifierStateOK, nil
	}

	if lastState.State == model.NotifierStatePending {
		if !alertTriggered {
			log.Info().Msgf("notifier %v goes from pending to ok", notifier.ID)
			return true, model.NotifierStateOK, nil
		}

		pendingDuration := time.Now().Sub(lastState.Timestamp)
		if pendingDuration > notifier.TimeUntilTrigger {
			log.Info().Msgf("notifier %v is now triggering", notifier.ID)
			return true, model.NotifierStateAlerting, nil
		}

		return false, model.NotifierStatePending, nil
	}

	return false, "", fmt.Errorf("unknown state for notifier %v", notifier.ID)
}

// evaluateTemplates checks if the notifier is triggered by evaluating the templates
func evaluateTemplates(notifier *model.Notifier) (bool, error) {

	if len(notifier.NotifierTemplates) == 0 {
		return false, nil
	}

	tempResults := make([]bool, 0)

	for _, template := range notifier.NotifierTemplates {
		var val int64 = 0
		switch template.Key {
		case model.NotifierTemplateKeyAuction:
			var err error
			val, err = redis.ReceiveAuctionCount(time.Now(), template.DurationToCheck)
			if err != nil {
				log.Error().Err(err).Msgf("error getting auction count for notifier, %v", notifier)
				return false, err
			}
			break
		case model.NotifierTemplateKeyFlip:
			return false, fmt.Errorf("flip notifiers are not implemented yet %s", template.Key)
		default:
			return false, fmt.Errorf("unknown notifier template key %s", template.Key)
		}

		r := false
		if template.Operator == model.NotifierTemplateOperatorBigger {
			r = val > template.Value
		} else if template.Operator == model.NotifierTemplateOperatorSmaller {
			r = val < template.Value
		} else {
			return false, fmt.Errorf("unknown notifier template operator %s", template.Operator)
		}

		tempResults = append(tempResults, r)
	}

	result := false
	if notifier.TemplateOperator == model.NotifierOperatorAND {
		result = true
		for _, tempResut := range tempResults {
			if !tempResut {
				result = false
				break
			}
		}
	} else if notifier.TemplateOperator == model.NotifierOperatorOR {
		result = false
		for _, tempResut := range tempResults {
			if tempResut {
				result = true
				break
			}
		}
	} else {
		return false, fmt.Errorf("unknown notifier operator %s", notifier.TemplateOperator)
	}

	return result, nil
}
