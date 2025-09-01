package services

import (
	"awesome-proj/app/domain/entities"
	"fmt"
	"math/rand"
	"time"
)

// LocationInteractionService handles location interactions and responses
type LocationInteractionService struct {
	// В будущем здесь будет AI client для генерации ответов
}

// NewLocationInteractionService creates a new interaction service
func NewLocationInteractionService() *LocationInteractionService {
	return &LocationInteractionService{}
}

// GenerateInitialLocationState creates the first interaction when player enters location
func (lis *LocationInteractionService) GenerateInitialLocationState(locationID string, locationName string, npcCount int) entities.Interaction {
	// Заглушка для начального состояния локации
	// В будущем здесь будет промпт в AI с описанием локации и NPC
	
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

// GenerateResponseToAction creates a response to player action
func (lis *LocationInteractionService) GenerateResponseToAction(action entities.Interaction, locationName string, npcCount int) entities.Interaction {
	// Заглушка для ответа на действие игрока
	// В будущем здесь будет промпт в AI с контекстом действия
	
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

// generateInteractionID creates a unique interaction ID
func generateInteractionID() string {
	return fmt.Sprintf("interaction_%d", time.Now().UnixNano())
}
