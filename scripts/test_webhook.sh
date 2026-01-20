#!/bin/bash

echo "╔════════════════════════════════════════════════════════╗"
echo "║                                                        ║"
echo "║           🧪 TESTE DE WEBHOOK - SIMULADO 🧪             ║"
echo "║                                                        ║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""
echo "📨 Enviando webhook simulado de invoice pago..."
echo ""

curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "event": {
      "subscription": "invoice",
      "log": {
        "type": "credited",
        "invoice": {
          "id": "test-invoice-12345",
          "amount": 10000,
          "fee": 0,
          "name": "Test User",
          "taxId": "012.345.678-90"
        }
      }
    }
  }'

echo ""
echo ""
echo "✅ Webhook enviado!"
echo "📊 Verifique os logs do servidor (make run) para ver o resultado"
echo ""
