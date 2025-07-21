#!/bin/bash

# Script to run tests with coverage and check against target threshold

set -e

# Configuration
COVERAGE_TARGET=80
COVERAGE_FILE="coverage.out"
COVERAGE_HTML="coverage.html"

echo "🧪 Running tests with coverage..."

# Run tests with coverage
go test -v -race -covermode=atomic -coverprofile=$COVERAGE_FILE ./...

# Generate HTML coverage report
go tool cover -html=$COVERAGE_FILE -o $COVERAGE_HTML
echo "📊 Coverage report generated: $COVERAGE_HTML"

# Get total coverage percentage
COVERAGE=$(go tool cover -func=$COVERAGE_FILE | tail -1 | awk '{print $3}' | sed 's/%//')

echo "📈 Current coverage: $COVERAGE%"
echo "🎯 Target coverage: $COVERAGE_TARGET%"

# Check if coverage meets target
if (( $(echo "$COVERAGE >= $COVERAGE_TARGET" | bc -l) )); then
    echo "✅ Coverage target met! ($COVERAGE% >= $COVERAGE_TARGET%)"
    echo ""
    echo "🎉 Test suite summary:"
    echo "  - All tests passing"
    echo "  - Coverage above $COVERAGE_TARGET% threshold"
    echo "  - Ready for production!"
    exit 0
else
    echo "❌ Coverage below target ($COVERAGE% < $COVERAGE_TARGET%)"
    echo ""
    echo "📋 Coverage breakdown by package:"
    go tool cover -func=$COVERAGE_FILE | grep -v "total:" | sort -k3 -nr
    echo ""
    echo "💡 To improve coverage, add tests for uncovered functions."
    exit 1
fi 