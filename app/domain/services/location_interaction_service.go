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

func (lis *LocationInteractionService) GenerateInitialLocationState(locationID string, locationName string, npcCount int) entities.Interaction {
	var stateDescriptions = []string{
		fmt.Sprintf("В %s царит обычная атмосфера. %d персонажей занимаются своими делами.", locationName, npcCount),
		fmt.Sprintf("Вы входите в %s. Здесь довольно тихо, %d местных жителей мирно проводят время.", locationName, npcCount),
		fmt.Sprintf("Атмосфера в %s спокойная. %d персонажей общаются между собой.", locationName, npcCount),
		fmt.Sprintf("В %s слышен тихий гул разговоров. %d персонажей заняты своими повседневными делами.", locationName, npcCount),
	}

	description := stateDescriptions[rand.Intn(len(stateDescriptions))]
	return entities.Interaction{
		ID:         generateInteractionID(),
		Type:       entities.InteractionTypeLocationState,
		Content:    description,
		LocationID: locationID,
		Timestamp:  time.Now(),
	}
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
