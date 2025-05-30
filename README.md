# RWA Backend (Golang)

Backend service for Real World Asset (RWA) tokenization platform built with Gin, GORM, and SQLite.

## 🚀 Quick Start

```bash
# 1. Install dependencies
go mod tidy

# 2. Run the server
go run main.go
```

The server will start on `http://localhost:8080`

## 📁 Project Structure

```
backend-lisk-go/
├── main.go                 # Entry point and server setup
├── models/
│   └── asset.go           # GORM models (Asset, MintRequest)
├── handlers/
│   └── asset_handler.go   # HTTP handlers
├── database/
│   └── database.go        # Database connection and migration
├── middleware/
│   └── cors.go           # CORS configuration
├── utils/
│   └── response.go       # API response helpers
├── go.mod                # Dependencies
└── README.md             # This file
```

## 🔌 API Endpoints

### Health Check
- `GET /` - Basic API info
- `GET /api/health` - Health check with database status

### Assets
- `GET /api/assets` - Get all assets
- `POST /api/assets/mint` - Add new minted asset
- `GET /api/assets/:id` - Get specific asset

## 📝 API Usage Examples

### Add Asset After Minting
```bash
curl -X POST http://localhost:8080/api/assets/mint \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Luxury Villa Bali",
    "symbol": "LVB",
    "institutionName": "PropTech Indonesia",
    "institutionAddress": "Jl. Sudirman 123, Jakarta",
    "description": "Luxury beachfront villa in Canggu",
    "totalSupply": "1000000",
    "pricePerRWA": "100",
    "contractAddress": "0x123...",
    "txHash": "0xabc...",
    "documentsURI": "ipfs://...",
    "imageURI": "ipfs://..."
  }'
```

### Get All Assets
```bash
curl http://localhost:8080/api/assets
```

### Response Format
All responses follow this standard format:
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... }
}
```

## 🗄️ Database Schema

### Asset Model
- `id` - Unique identifier
- `name` - Asset name
- `symbol` - Token symbol  
- `type` - Asset type (default: "Real Estate")
- `institution` - Institution name
- `institution_address` - Institution address
- `description` - Asset description
- `total_supply` - Total token supply
- `staked_amount` - Currently staked amount (default: 0)
- `price_usd` - Price per token in USD
- `annual_yield` - Expected annual yield (default: 8.5%)
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp
- `blockchain` - Blockchain name (default: "Lisk")
- `contract_address` - Smart contract address
- `tx_hash` - Creation transaction hash
- `documents_uri` - IPFS URI for documents
- `image_uri` - IPFS URI for images

## ⚙️ Configuration

### Environment Variables
Copy `env.example` to `.env` and adjust:

```bash
cp env.example .env
```

Available variables:
- `PORT` - Server port (default: 8080)
- `GIN_MODE` - Gin mode (debug/release)
- `DB_PATH` - SQLite database path

### CORS Settings
Currently allows:
- `http://localhost:3000` (development frontend)
- `https://yourdomain.com` (production domain)

## 🧪 Testing

### Manual Testing
```bash
# Health check
curl http://localhost:8080/api/health

# Test empty assets list
curl http://localhost:8080/api/assets

# Test asset creation
curl -X POST http://localhost:8080/api/assets/mint \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Asset","symbol":"TST","institutionName":"Test Corp","totalSupply":"1000","pricePerRWA":"10"}'
```

## 🔄 Integration with Frontend

### Mint Page Integration
After successful contract execution, call:
```javascript
const response = await fetch('http://localhost:8080/api/assets/mint', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    name: formData.name,
    symbol: formData.symbol,
    institutionName: formData.institutionName,
    // ... other form data
    contractAddress: transactionResult.contractAddress,
    txHash: transactionResult.hash
  })
})
```

### Explore Page Integration
```javascript
const response = await fetch('http://localhost:8080/api/assets')
const data = await response.json()
const assets = data.data.assets
```

## 🚀 Deployment

### Development
```bash
go run main.go
```

### Production Build
```bash
go build -o rwa-backend main.go
./rwa-backend
```

### Docker (Future)
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o main .
CMD ["./main"]
```

## 🔮 Future Enhancements

- [ ] Blockchain event listener
- [ ] PostgreSQL support
- [ ] JWT authentication
- [ ] Rate limiting
- [ ] Pagination and filtering
- [ ] Asset search functionality
- [ ] Real-time updates with WebSockets
- [ ] Docker deployment
- [ ] Comprehensive test suite

## 📊 Current Features

✅ SQLite database with GORM  
✅ RESTful API with Gin  
✅ CORS support for frontend  
✅ Standardized API responses  
✅ Asset CRUD operations  
✅ Health check endpoints  
✅ Auto database migration  
✅ Environment configuration  

---

**Ready to mint and explore RWA tokens! 🎯** 