#!/bin/bash
echo "🧪 Running Go tests..."
go test -coverprofile=coverage.out -json ./... > test-report.json

echo "📊 Coverage summary:"
go tool cover -func=coverage.out | grep total

echo "🔍 Running SonarQube scan..."
sonar-scanner -Dsonar.login="sqa_c335cddb5384589e7de1f3787bd9d3c31cc8e106"

echo "✅ Scan completed! Check: http://localhost:9000"