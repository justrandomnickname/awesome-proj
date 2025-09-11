package npc

import (
	"math/rand"
)

// TraitType представляет тип трейта
type TraitType string

const (
	TraitExperiencedWarrior TraitType = "experienced_warrior"
	TraitBookworm           TraitType = "bookworm"
	TraitNatureLover        TraitType = "nature_lover"
)

// Trait представляет трейт персонажа
type Trait struct {
	ID          string    `json:"id"`
	Type        TraitType `json:"type"`
	Name        string    `json:"name"`        // Человекочитаемое название
	Description string    `json:"description"` // Описание эффекта
}

// TraitEffect описывает влияние трейта на характеристики
type TraitEffect struct {
	PrudenceModifier       int `json:"prudence_modifier"`
	EmotionalityModifier   int `json:"emotionality_modifier"`
	IndependenceModifier   int `json:"independence_modifier"`
	OptimismModifier       int `json:"optimism_modifier"`
	FlexibilityModifier    int `json:"flexibility_modifier"`
	AggressivenessModifier int `json:"aggressiveness_modifier"`
}

// WeightCalculatorFunc - функция для расчёта веса трейта для конкретного NPC
type WeightCalculatorFunc func(npc *NPC) int

// TraitDefinition определяет трейт и его свойства
type TraitDefinition struct {
	Type             TraitType            `json:"type"`
	Name             string               `json:"name"`
	Description      string               `json:"description"`
	BaseWeight       int                  `json:"base_weight"` // Базовый вес (шанс выпадения из 100)
	Effect           TraitEffect          `json:"effect"`      // Влияние на характеристики
	WeightCalculator WeightCalculatorFunc `json:"-"`           // Функция расчёта веса
}

// TraitSystem управляет системой трейтов
type TraitSystem struct {
	traitDefinitions map[TraitType]*TraitDefinition
}

// NewTraitSystem создаёт новую систему трейтов
func NewTraitSystem() *TraitSystem {
	ts := &TraitSystem{
		traitDefinitions: make(map[TraitType]*TraitDefinition),
	}

	// Регистрируем трейты
	ts.registerTraits()

	return ts
}

// registerTraits регистрирует все доступные трейты
func (ts *TraitSystem) registerTraits() {
	// 1. Опытный Воин
	ts.traitDefinitions[TraitExperiencedWarrior] = &TraitDefinition{
		Type:        TraitExperiencedWarrior,
		Name:        "Опытный Воин",
		Description: "Многолетний опыт в боях повышает агрессивность и снижает эмоциональность",
		BaseWeight:  5, // 5% базовый шанс
		Effect: TraitEffect{
			AggressivenessModifier: +2, // Более агрессивный
			EmotionalityModifier:   -1, // Более хладнокровный
			PrudenceModifier:       +1, // Более осторожный в бою
		},
		WeightCalculator: func(npc *NPC) int {
			weight := 5 // Базовый вес

			// Влияние расы
			if npc.Race == "human" {
				weight += 3 // Люди чаще становятся воинами
			}

			// Влияние текущих черт характера
			if npc.TemperTraits.Aggressiveness >= 3 {
				weight += 5 // Агрессивные персонажи чаще становятся воинами
			}
			if npc.TemperTraits.Prudence >= 2 {
				weight += 2 // Осторожные тоже могут быть хорошими воинами
			}
			if npc.TemperTraits.Emotionality <= -2 {
				weight += 3 // Хладнокровные склонны к военному делу
			}

			return weight
		},
	}

	// 2. Книжный Червь
	ts.traitDefinitions[TraitBookworm] = &TraitDefinition{
		Type:        TraitBookworm,
		Name:        "Книжный Червь",
		Description: "Любовь к знаниям повышает осторожность и снижает агрессивность",
		BaseWeight:  8, // 8% базовый шанс
		Effect: TraitEffect{
			PrudenceModifier:       +2, // Более осторожный
			AggressivenessModifier: -2, // Менее агрессивный
			IndependenceModifier:   +1, // Более независимый в мышлении
		},
		WeightCalculator: func(npc *NPC) int {
			weight := 8 // Базовый вес

			// Влияние черт характера
			if npc.TemperTraits.Prudence >= 2 {
				weight += 4 // Осторожные люди чаще читают
			}
			if npc.TemperTraits.Aggressiveness <= -1 {
				weight += 3 // Неагрессивные предпочитают книги боям
			}
			if npc.TemperTraits.Independence >= 3 {
				weight += 2 // Независимые мыслители любят учиться
			}
			if npc.TemperTraits.Emotionality <= 0 {
				weight += 2 // Менее эмоциональные склонны к интеллектуальной деятельности
			}

			return weight
		},
	}

	// 3. Любитель Природы
	ts.traitDefinitions[TraitNatureLover] = &TraitDefinition{
		Type:        TraitNatureLover,
		Name:        "Любитель Природы",
		Description: "Глубокая связь с природой повышает оптимизм и гибкость",
		BaseWeight:  10, // 10% базовый шанс
		Effect: TraitEffect{
			OptimismModifier:       +2, // Более оптимистичный
			FlexibilityModifier:    +1, // Более гибкий
			AggressivenessModifier: -1, // Менее агрессивный
		},
		WeightCalculator: func(npc *NPC) int {
			weight := 10 // Базовый вес

			// Влияние черт характера
			if npc.TemperTraits.Optimism >= 1 {
				weight += 4 // Оптимисты тянутся к природе
			}
			if npc.TemperTraits.Flexibility >= 2 {
				weight += 3 // Гибкие люди легче адаптируются к природе
			}
			if npc.TemperTraits.Aggressiveness <= -2 {
				weight += 3 // Миролюбивые находят покой в природе
			}
			if npc.TemperTraits.Emotionality >= 2 {
				weight += 2 // Эмоциональные люди сильнее чувствуют связь с природой
			}

			return weight
		},
	}
}

// GenerateTraitsForNPC генерирует трейты для конкретного NPC
func (ts *TraitSystem) GenerateTraitsForNPC(npc *NPC, rng *rand.Rand) []Trait {
	var traits []Trait

	// Максимум 2 трейта на персонажа (можно настроить)
	maxTraits := 2

	// Для каждого возможного трейта рассчитываем вес
	weightedTraits := make([]TraitType, 0)
	for traitType, definition := range ts.traitDefinitions {
		weight := definition.WeightCalculator(npc)

		// Добавляем трейт в список с учётом веса
		for i := 0; i < weight; i++ {
			weightedTraits = append(weightedTraits, traitType)
		}
	}

	// Если нет доступных трейтов, возвращаем пустой список
	if len(weightedTraits) == 0 {
		return traits
	}

	// Генерируем трейты
	usedTraits := make(map[TraitType]bool)
	for len(traits) < maxTraits {
		// Случайно выбираем трейт из взвешенного списка
		selectedTraitType := weightedTraits[rng.Intn(len(weightedTraits))]

		// Проверяем, что трейт ещё не использован
		if usedTraits[selectedTraitType] {
			continue
		}

		// Дополнительная проверка на шанс получения трейта (можно настроить)
		if rng.Intn(100) < 30 { // 30% шанс получить трейт даже при наличии веса
			definition := ts.traitDefinitions[selectedTraitType]
			trait := Trait{
				ID:          generateTraitID(npc.ID, selectedTraitType),
				Type:        selectedTraitType,
				Name:        definition.Name,
				Description: definition.Description,
			}

			traits = append(traits, trait)
			usedTraits[selectedTraitType] = true
		}

		// Защита от бесконечного цикла
		if len(usedTraits) >= len(ts.traitDefinitions) {
			break
		}
	}

	return traits
}

// ApplyTraitEffects применяет эффекты трейтов к характеристикам NPC
func (ts *TraitSystem) ApplyTraitEffects(npc *NPC, traits []Trait) {
	for _, trait := range traits {
		if definition, exists := ts.traitDefinitions[trait.Type]; exists {
			effect := definition.Effect

			// Применяем модификаторы (с ограничениями -10...+10)
			npc.TemperTraits.Prudence = clampTrait(npc.TemperTraits.Prudence + effect.PrudenceModifier)
			npc.TemperTraits.Emotionality = clampTrait(npc.TemperTraits.Emotionality + effect.EmotionalityModifier)
			npc.TemperTraits.Independence = clampTrait(npc.TemperTraits.Independence + effect.IndependenceModifier)
			npc.TemperTraits.Optimism = clampTrait(npc.TemperTraits.Optimism + effect.OptimismModifier)
			npc.TemperTraits.Flexibility = clampTrait(npc.TemperTraits.Flexibility + effect.FlexibilityModifier)
			npc.TemperTraits.Aggressiveness = clampTrait(npc.TemperTraits.Aggressiveness + effect.AggressivenessModifier)
		}
	}
}

// GetTraitDefinition возвращает определение трейта
func (ts *TraitSystem) GetTraitDefinition(traitType TraitType) *TraitDefinition {
	return ts.traitDefinitions[traitType]
}

// GetTraitNamesByIDs возвращает названия трейтов по их ID
func (ts *TraitSystem) GetTraitNamesByIDs(traitIDs []string) []string {
	names := make([]string, 0, len(traitIDs))
	for _, id := range traitIDs {
		// Ищем "_trait_" в ID и извлекаем тип
		traitPrefix := "_trait_"
		traitIndex := -1
		for i := 0; i <= len(id)-len(traitPrefix); i++ {
			if id[i:i+len(traitPrefix)] == traitPrefix {
				traitIndex = i + len(traitPrefix)
				break
			}
		}

		if traitIndex != -1 && traitIndex < len(id) {
			traitTypeStr := id[traitIndex:]
			traitType := TraitType(traitTypeStr)
			if definition := ts.GetTraitDefinition(traitType); definition != nil {
				names = append(names, definition.Name)
			}
		}
	}
	return names
}

// clampTrait ограничивает значение черты характера диапазоном -10...+10
func clampTrait(value int) int {
	if value > 10 {
		return 10
	}
	if value < -10 {
		return -10
	}
	return value
}

// generateTraitID генерирует уникальный ID для трейта
func generateTraitID(npcID string, traitType TraitType) string {
	return npcID + "_trait_" + string(traitType)
}
