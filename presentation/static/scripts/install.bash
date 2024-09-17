#!/bin/bash 

if [[ "$OSTYPE" == "linux-gnu"* ]]; then
  echo "Linux operating system detected"
  curl -fsSL https://yer-app.fly.dev/scripts/yer_linux_amd64 --output year-end-recap
elif [[ "$OSTYPE" == "darwin"* ]]; then
  echo "Mac operating system detected"
  CPU_ARCH=$(uname -m)

  if [[ "$CPU_ARCH" == "arm64" ]]; then
    curl -fsSL https://yer-app.fly.dev/scripts/yer_darwin_arm64 --output year-end-recap
  elif [[ "$CPU_ARCH" == "x86_64" ]]; then
    curl -fsSL https://yer-app.fly.dev/scripts/yer_darwin_amd64 --output year-end-recap
  else
    echo "Unsupported Mac CPU architecture: $CPU_ARCH"
    exit 1
  fi
else
  echo "Unsupported Operating System: $OSTYPE"
  exit 1
fi

chmod +x year-end-recap
echo
echo "Year End Recap has been installed!"
echo

./year-end-recap -a