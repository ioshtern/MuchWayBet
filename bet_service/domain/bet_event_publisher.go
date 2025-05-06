package domain

type BetEventPublisher interface {
	PublishBetCreated(bet *Bet) error
	PublishBetUpdated(bet *Bet) error
	PublishBetDeleted(bet *Bet) error
}
