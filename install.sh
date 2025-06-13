#!/usr/bin/bash

echo "Installing and configuring Yogit" 
echo "Building executable..."
$(go build main.go)

windows="main.exe"
linux="main"

if [ -f "$linux" ]; then
  echo "Building for linux systems..."
  BuildLinux
else 
  echo "Building for windows system"
  curr_dir=pwd
  exec_path=$curr_dir/$linux
  echo "fatal: Build failed, please refer to docs at https://github.com/persona-mp3/yogit.git"
  echo "You can directly install the build from the repo"
fi


BuildLinux(){
  echo "Building execution Path"
  curr_dir=pwd
  exec_path=$curr_dir/$linux

  echo "Writing to bashrc..."
  config="alias yogit=$exec_path"
  echo $config

  echo $config >> ~/.bashrc
  echo "--- Confirm config was successful"
  echo "Run `yogit` to confirm or check .bashrc for $config"
}

