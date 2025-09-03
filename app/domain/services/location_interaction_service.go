package services

import (
	"awesome-proj/app/domain/entities"
	"fmt"
	"math/rand"
	"time"
)

type LocationInteractionService struct {
}

func NewLocationInteractionService() *LocationInteractionService {
	return &LocationInteractionService{}
}

func (lis *LocationInteractionService) GenerateResponseToAction(action entities.Interaction, locationName string, npcCount int) entities.Interaction {
	var responses = []string{
		"Окружающие с интересом наблюдают за вашими действиями.",
		"Ваши действия привлекли внимание местных жителей.",
		"Некоторые персонажи переглядываются, обсуждая произошедшее.",
		"Атмосфера в локации слегка изменилась после ваших действий.",
		"Местные жители реагируют на происходящее по-разному.",
		"Ваше поведение не осталось незамеченным.",
	}
	response := responses[rand.Intn(len(responses))]
	return entities.Interaction{
		ID:         generateInteractionID(),
		Type:       entities.InteractionTypeLocationResponse,
		Content:    response,
		LocationID: action.LocationID,
		Timestamp:  time.Now(),
	}
}

func generateInteractionID() string {
	return fmt.Sprintf("interaction_%d", time.Now().UnixNano())
}
