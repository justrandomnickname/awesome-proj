package npc

type TemperTraits struct {
	ID             string `json:"id"`
	NPC_ID         string `json:"npc_id"`
	Prudence       int    `json:"prudence"`
	Emotionality   int    `json:"emotionality"`
	Independence   int    `json:"independence"`
	Optimism       int    `json:"optimism"`
	Flexibility    int    `json:"flexibility"`
	Aggressiveness int    `json:"aggressiveness"`
}
