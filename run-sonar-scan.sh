#!/bin/bash
echo "🧪 Running Go tests..."
go test -coverprofile=coverage.out -json ./... > test-report.json

echo "📊 Coverage summary:"
go tool cover -func=coverage.out | grep total

echo "🔍 Running SonarQube scan..."
sonar-scanner -Dsonar.login=$SONAR_TOKEN

echo "✅ Scan completed! Check: http://localhost:9000"