package commands

import (
	"encoding/json"

	"github.com/iotaledger/hive.go/ierrors"
	"github.com/iotaledger/hive.go/logger"

	"github.com/iotaledger/hive.go/web/subscriptionmanager"
	"github.com/iotaledger/hive.go/web/websockethub"
)

type CommandType string

type BaseCommand struct {
	Command CommandType `json:"command"`
}

type EventType string

type BaseEvent struct {
	Event EventType `json:"event"`
}

var (
	ErrFailedToDeserializeCommand = ierrors.New("failed to deserialize command")
	ErrFailedToValidateCommand    = ierrors.New("failed to validate command")
	ErrFailedToSendMessage        = ierrors.New("failed to send message")
)

type CommandHandler interface {
	SupportsCommand(commandType CommandType) bool
	HandleCommand(client *websockethub.Client, message []byte) error
}

type CommandManager struct {
	log                 *logger.Logger
	commands            []CommandHandler
	subscriptionManager *subscriptionmanager.SubscriptionManager[websockethub.ClientID, string]
}

func NewCommandHandler(log *logger.Logger, subscriptionManager *subscriptionmanager.SubscriptionManager[websockethub.ClientID, string]) *CommandManager {
	return &CommandManager{
		log: log,
		commands: []CommandHandler{
			// Register new commands here
			&SubscriptionCommandHandler{
				log:                 log,
				subscriptionManager: subscriptionManager,
			},
		},
		subscriptionManager: subscriptionManager,
	}
}

func (c *CommandManager) handleError(err error, client *websockethub.Client, commandType CommandType) error {
	unwrappedError := ierrors.Unwrap(err)

	switch {
	case ierrors.Is(err, ErrFailedToDeserializeCommand):
		c.log.Warnf("Failed to deserialize for client:[%d], command:[%s], err:[%v]", client.ID(), commandType, unwrappedError)
	case ierrors.Is(err, ErrFailedToSendMessage):
		c.log.Warnf("Failed to send event to client:[%d], command:[%s], err:[%v]", client.ID(), commandType, unwrappedError)
	case ierrors.Is(err, ErrFailedToValidateCommand):
		c.log.Warnf("Failed to validate received command from client:[%d], command:[%s], err:[%v]", client.ID(), commandType, unwrappedError)
	default:
		c.log.Warnf("Unhandled error in websocket command handler for client:[%d], command:[%s], err:[%v]", client.ID(), commandType, err)
	}

	return err
}

func (c *CommandManager) HandleNodeCommands(client *websockethub.Client, message []byte) error {
	var baseCommand BaseCommand
	if err := json.Unmarshal(message, &baseCommand); err != nil {
		return c.handleError(ierrors.Wrap(ErrFailedToDeserializeCommand, err.Error()), client, baseCommand.Command)
	}

	if baseCommand.Command == "" {
		return c.handleError(ierrors.Wrap(ErrFailedToValidateCommand, "Command is empty"), client, baseCommand.Command)
	}

	for _, commandHandler := range c.commands {
		if commandHandler.SupportsCommand(baseCommand.Command) {
			if err := commandHandler.HandleCommand(client, message); err != nil {
				return c.handleError(err, client, baseCommand.Command)
			}
		}
	}

	return nil
}
