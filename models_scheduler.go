package main

//region FeatureToInclude

// UserPage Page utilisateur
type UserPage string

const (
	// Preview Previews
	Preview UserPage = "preview"
	// Console Console de commande
	Console UserPage = "console"
	// Monitor Page de monitoring
	Monitor UserPage = "monitor"
)

//endregion FeatureToInclude

// Operator Classe opérateur
type Operator struct {
	Name            string `json:"name"`
	ChannelID       string `json:"channel_id"`
	ControlledDrone string `json:"controlled_drone"`
	IsAnonymous     bool   `json:"is_anonymous"`
}

// OperatorIdentifier Classe identification opérateur
type OperatorIdentifier struct {
	// Name Nom
	Name string `json:"name"`
}
