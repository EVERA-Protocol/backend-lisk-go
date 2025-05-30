package handlers

import (
	"fmt"
	"net/http"
	"rwa-backend/database"
	"rwa-backend/models"
	"rwa-backend/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// MintAsset handles POST /api/assets/mint
func MintAsset(c *gin.Context) {
	var req models.MintRequest

	// Bind and validate request
	if err := c.BindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Convert string numbers to proper types
	totalSupply, err := strconv.ParseInt(req.TotalSupply, 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid total supply", err.Error())
		return
	}

	expectedYield, err := strconv.ParseFloat(req.ExpectedYield, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid expected yield", err.Error())
		return
	}

	// Handle optional PricePerRWA - default to 1.0 if not provided
	pricePerRWA := 1.0
	if req.PricePerRWA != "" {
		parsedPrice, err := strconv.ParseFloat(req.PricePerRWA, 64)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid price per RWA", err.Error())
			return
		}
		pricePerRWA = parsedPrice
	}

	// Set contract address to "pending" if not provided
	contractAddress := req.ContractAddress
	if contractAddress == "" {
		contractAddress = "pending"
	}

	// Create asset model
	asset := models.Asset{
		ID:                 fmt.Sprintf("asset_%d", time.Now().Unix()),
		Name:               req.Name,
		Symbol:             req.Symbol,
		Type:               "Real Estate", // Default for MVP
		Institution:        req.InstitutionName,
		InstitutionAddress: req.InstitutionAddress,
		Description:        req.Description,
		TotalSupply:        totalSupply,
		StakedAmount:       0,           // Default
		PriceUsd:           pricePerRWA, // Use the parsed or default price
		AnnualYield:        expectedYield,
		Blockchain:         "Lisk", // Default
		ContractAddress:    contractAddress,
		TxHash:             req.TxHash,
		DocumentsURI:       req.DocumentsURI,
		ImageURI:           req.ImageURI,
	}

	// Save to database
	db := database.GetDB()
	if err := db.Create(&asset).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save asset", err.Error())
		return
	}

	// Return success response with computed fields
	assetResponse := asset.ToAssetResponse()
	utils.SuccessResponse(c, http.StatusCreated, "Asset minted successfully", gin.H{
		"id":    asset.ID,
		"asset": assetResponse,
	})
}

// GetAllAssets handles GET /api/assets
func GetAllAssets(c *gin.Context) {
	var assets []models.Asset

	db := database.GetDB()

	// Fetch all assets ordered by creation date (newest first)
	if err := db.Order("created_at DESC").Find(&assets).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch assets", err.Error())
		return
	}

	// Convert to AssetResponse with computed fields
	var assetResponses []models.AssetResponse
	for _, asset := range assets {
		assetResponses = append(assetResponses, asset.ToAssetResponse())
	}

	// Return assets with computed fields
	utils.SuccessResponse(c, http.StatusOK, "Assets fetched successfully", gin.H{
		"assets": assetResponses,
		"total":  len(assetResponses),
	})
}

// GetAssetByID handles GET /api/assets/:id
func GetAssetByID(c *gin.Context) {
	id := c.Param("id")

	var asset models.Asset
	db := database.GetDB()

	if err := db.First(&asset, "id = ?", id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Asset not found", err.Error())
		return
	}

	// Return with computed fields
	assetResponse := asset.ToAssetResponse()
	utils.SuccessResponse(c, http.StatusOK, "Asset fetched successfully", assetResponse)
}

// UpdateContractAddress handles PATCH /api/assets/:id/contract
func UpdateContractAddress(c *gin.Context) {
	id := c.Param("id")
	var req models.UpdateContractRequest

	// Bind and validate request
	if err := c.BindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Find asset
	var asset models.Asset
	db := database.GetDB()
	if err := db.First(&asset, "id = ?", id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Asset not found", err.Error())
		return
	}

	// Update contract address and tx hash
	asset.ContractAddress = req.ContractAddress
	if req.TxHash != "" {
		asset.TxHash = req.TxHash
	}
	asset.UpdatedAt = time.Now()

	// Save changes
	if err := db.Save(&asset).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update asset", err.Error())
		return
	}

	// Return updated asset with computed fields
	assetResponse := asset.ToAssetResponse()
	utils.SuccessResponse(c, http.StatusOK, "Contract address updated successfully", assetResponse)
}

// GetAssetStats handles GET /api/assets/:id/stats
func GetAssetStats(c *gin.Context) {
	id := c.Param("id")

	var asset models.Asset
	db := database.GetDB()

	if err := db.First(&asset, "id = ?", id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Asset not found", err.Error())
		return
	}

	// For MVP, return calculated stats (in production, fetch from blockchain)
	stats := models.TokenStats{
		TotalSupply:       asset.TotalSupply,
		CirculatingSupply: asset.TotalSupply, // All tokens are circulating for RWA
		HolderCount:       3,                 // Mock holder count
		Price:             asset.PriceUsd,
		MarketCap:         float64(asset.TotalSupply) * asset.PriceUsd,
	}

	utils.SuccessResponse(c, http.StatusOK, "Asset stats fetched successfully", stats)
}

// UpdateAssetStaking handles PATCH /api/assets/:id/staking (for testing staking progress)
func UpdateAssetStaking(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		StakedAmount int64 `json:"stakedAmount" binding:"required"`
	}

	// Bind and validate request
	if err := c.BindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	// Find asset
	var asset models.Asset
	db := database.GetDB()
	if err := db.First(&asset, "id = ?", id).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Asset not found", err.Error())
		return
	}

	// Validate staked amount doesn't exceed total supply
	if req.StakedAmount > asset.TotalSupply {
		utils.ErrorResponse(c, http.StatusBadRequest, "Staked amount cannot exceed total supply", "")
		return
	}

	// Update staked amount
	asset.StakedAmount = req.StakedAmount
	asset.UpdatedAt = time.Now()

	// Save changes
	if err := db.Save(&asset).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update staking", err.Error())
		return
	}

	// Return updated asset with computed fields
	assetResponse := asset.ToAssetResponse()
	utils.SuccessResponse(c, http.StatusOK, "Staking updated successfully", assetResponse)
}

// HealthCheck handles GET /api/health
func HealthCheck(c *gin.Context) {
	// Check database connection
	db := database.GetDB()
	sqlDB, err := db.DB()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Database connection error", err.Error())
		return
	}

	if err := sqlDB.Ping(); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Database ping failed", err.Error())
		return
	}

	// Get asset count for health info
	var count int64
	db.Model(&models.Asset{}).Count(&count)

	utils.SuccessResponse(c, http.StatusOK, "API is healthy", gin.H{
		"status":       "ok",
		"timestamp":    time.Now(),
		"database":     "connected",
		"total_assets": count,
	})
}
