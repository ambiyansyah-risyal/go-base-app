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
if which act > /dev/null 2>&1; then
    if act --list > /dev/null 2>&1; then
        echo "✅ Workflow structure is valid"
        act --list
    else
        echo "❌ Workflow structure error"
        act --list
        exit 1
    fi
else
    echo "⚠️  Act tool not found, skipping workflow validation"
    echo "   Install with: curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash"
fi

echo ""
echo "3. Checking Go compatibility..."
if go version | grep -q "go1\.2[2-4]"; then
    echo "✅ Go version is compatible"
else
    echo "⚠️  Go version may not match CI (1.22 expected)"
fi

echo ""
echo "4. Testing docs generation..."
if make docs > /dev/null 2>&1; then
    echo "✅ Documentation generation successful"
else
    echo "❌ Documentation generation failed"
    make docs
    exit 1
fi

echo ""
echo "5. Testing local builds..."
if make build > /dev/null 2>&1; then
    echo "✅ Local builds successful"
else
    echo "❌ Local builds failed"
    make build
    exit 1
fi

echo ""
echo "✅ CI/CD pipeline validation complete!"
echo "Ready for commit to GitHub."
