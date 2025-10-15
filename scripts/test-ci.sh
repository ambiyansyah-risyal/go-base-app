#!/bin/bash
echo "🔧 Testing CI/CD Pipeline Locally"
echo "================================="

# Test YAML syntax
echo "1. Checking YAML syntax..."
if python3 -c "import yaml; yaml.safe_load(open('.github/workflows/ci.yml'))" 2>/dev/null; then
    echo "✅ YAML syntax is valid"
else
    echo "❌ YAML syntax error"
    exit 1
fi

# Test workflow structure
echo ""
echo "2. Testing workflow structure..."
if ./bin/act --list > /dev/null 2>&1; then
    echo "✅ Workflow structure is valid"
    ./bin/act --list
else
    echo "❌ Workflow structure error"
    ./bin/act --list
    exit 1
fi

echo ""
echo "3. Checking Go compatibility..."
if go version | grep -q "go1\.2[2-4]"; then
    echo "✅ Go version is compatible"
else
    echo "⚠️  Go version may not match CI (1.22 expected)"
fi

echo ""
echo "4. Testing local builds..."
if go build ./cmd/... > /dev/null 2>&1; then
    echo "✅ Local builds successful"
else
    echo "❌ Local builds failed"
    go build ./cmd/...
    exit 1
fi

echo ""
echo "✅ CI/CD pipeline validation complete!"
echo "Ready for commit to GitHub."
