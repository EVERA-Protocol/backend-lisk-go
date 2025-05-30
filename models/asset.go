package models

import (
	"time"

	"gorm.io/gorm"
)

// Asset represents a tokenized real-world asset
type Asset struct {
	ID                 string    `json:"id" gorm:"primaryKey"`
	Name               string    `json:"name" gorm:"not null"`
	Symbol             string    `json:"symbol" gorm:"not null"`
	Type               string    `json:"type" gorm:"default:'Real Estate'"`
	Institution        string    `json:"institution" gorm:"not null"`
	InstitutionAddress string    `json:"institutionAddress"`
	Description        string    `json:"description"`
	TotalSupply        int64     `json:"totalSupply"`
	StakedAmount       int64     `json:"stakedAmount" gorm:"default:0"`
	PriceUsd           float64   `json:"priceUsd"`
	AnnualYield        float64   `json:"annualYield" gorm:"default:8.5"`
	CreatedAt          time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	Blockchain         string    `json:"blockchain" gorm:"default:'Lisk'"`
	ContractAddress    string    `json:"contractAddress"`
	TxHash             string    `json:"txHash"`
	DocumentsURI       string    `json:"documentsURI"`
	ImageURI           string    `json:"imageURI"`

	// Computed fields for frontend compatibility (not stored in DB)
	Documents  []Document `json:"documents" gorm:"-"`
	TopStakers []Staker   `json:"topStakers" gorm:"-"`
}

// AssetResponse represents the enhanced response with computed fields
type AssetResponse struct {
	Asset
	// Computed financial metrics
	AvailableSupply  int64   `json:"availableSupply"`
	MarketCap        float64 `json:"marketCap"`
	MinInvestment    float64 `json:"minInvestment"`
	MaxInvestment    float64 `json:"maxInvestment"`
	StakingProgress  float64 `json:"stakingProgress"`
	IsContractActive bool    `json:"isContractActive"`
	TotalValue       float64 `json:"totalValue"`
	StakedValue      float64 `json:"stakedValue"`
	AvailableValue   float64 `json:"availableValue"`
}

// Document represents supporting documentation
type Document struct {
	Name string `json:"name"`
	Date string `json:"date"`
	URL  string `json:"url"`
}

// Staker represents a token holder
type Staker struct {
	Address    string  `json:"address"`
	Amount     float64 `json:"amount"`
	Percentage float64 `json:"percentage"`
}

// MintRequest represents the request payload for minting new assets
type MintRequest struct {
	Name               string `json:"name" binding:"required"`
	Symbol             string `json:"symbol" binding:"required"`
	InstitutionName    string `json:"institutionName" binding:"required"`
	InstitutionAddress string `json:"institutionAddress"`
	Description        string `json:"description"`
	TotalSupply        string `json:"totalSupply" binding:"required"`
	ExpectedYield      string `json:"expectedYield" binding:"required"`
	PricePerRWA        string `json:"pricePerRWA"` // Optional field, defaults if not provided
	ContractAddress    string `json:"contractAddress"`
	TxHash             string `json:"txHash"`
	DocumentsURI       string `json:"documentsURI"`
	ImageURI           string `json:"imageURI"`
}

// UpdateContractRequest represents updating contract address after deployment
type UpdateContractRequest struct {
	ContractAddress string `json:"contractAddress" binding:"required"`
	TxHash          string `json:"txHash"`
}

// TokenStats represents on-chain token statistics
type TokenStats struct {
	TotalSupply       int64   `json:"totalSupply"`
	CirculatingSupply int64   `json:"circulatingSupply"`
	HolderCount       int64   `json:"holderCount"`
	Price             float64 `json:"price"`
	MarketCap         float64 `json:"marketCap"`
}

// ToAssetResponse converts Asset to AssetResponse with computed fields
func (a *Asset) ToAssetResponse() AssetResponse {
	// Calculate computed fields
	availableSupply := a.TotalSupply - a.StakedAmount
	if availableSupply < 0 {
		availableSupply = 0
	}

	marketCap := float64(a.TotalSupply) * a.PriceUsd
	totalValue := marketCap
	stakedValue := float64(a.StakedAmount) * a.PriceUsd
	availableValue := float64(availableSupply) * a.PriceUsd

	// Calculate staking progress percentage
	stakingProgress := 0.0
	if a.TotalSupply > 0 {
		stakingProgress = (float64(a.StakedAmount) / float64(a.TotalSupply)) * 100
	}

	// Determine if contract is active (not "pending")
	isContractActive := a.ContractAddress != "" && a.ContractAddress != "pending"

	// Set investment limits based on available supply
	minInvestment := a.PriceUsd     // Minimum 1 token
	maxInvestment := availableValue // Maximum all available tokens
	if maxInvestment < minInvestment {
		maxInvestment = minInvestment
	}

	return AssetResponse{
		Asset:            *a,
		AvailableSupply:  availableSupply,
		MarketCap:        marketCap,
		MinInvestment:    minInvestment,
		MaxInvestment:    maxInvestment,
		StakingProgress:  stakingProgress,
		IsContractActive: isContractActive,
		TotalValue:       totalValue,
		StakedValue:      stakedValue,
		AvailableValue:   availableValue,
	}
}

// BeforeCreate adds computed fields before returning to frontend
func (a *Asset) AfterFind(tx *gorm.DB) error {
	// Add default document if DocumentsURI exists
	if a.DocumentsURI != "" {
		a.Documents = []Document{
			{
				Name: "Supporting Documents",
				Date: a.CreatedAt.Format("2006-01-02"),
				URL:  a.DocumentsURI,
			},
		}
	} else {
		a.Documents = []Document{}
	}

	// Initialize mock stakers data for demonstration
	if a.StakedAmount > 0 {
		// Create mock staker data (in production, this would come from blockchain)
		totalStaked := float64(a.StakedAmount)
		a.TopStakers = []Staker{
			{
				Address:    "0x1234...5678",
				Amount:     totalStaked * 0.4,
				Percentage: 40.0,
			},
			{
				Address:    "0xabcd...efgh",
				Amount:     totalStaked * 0.35,
				Percentage: 35.0,
			},
			{
				Address:    "0x9876...5432",
				Amount:     totalStaked * 0.25,
				Percentage: 25.0,
			},
		}
	} else {
		a.TopStakers = []Staker{}
	}

	return nil
}
