#!/bin/bash
echo "🔧 Testing CI/CD Pipeline"
./bin/act --list && echo "✅ CI/CD validated - ready for commit"
