#!/usr/bin/bash

echo "Installing and configuring Yogit" 
echo "Building executable..."
#$(cd ~/test_clone_ci && go build main.go)
# $(cd ~/.config/yogit && go build main.go)
$(go build main.go)

windows="main.exe"
linux="main"

if [ -f "$linux"]; then
  echo "Building for linux systems..."
  BuildLinux()
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
#cd ~/test_clone_ci
cd ~/.config/yogit
echo ""

#echo "changing directory"

#echo "cd ~/test_clone_ci"

#exec_path=$HOME/test_clone_ci/main.exe
echo "Building execution path"
exec_path=$HOME/.config/yogit/main

config="alias yogit=$exec_path"
echo "---Writing to .bashrc---"
echo $config

echo $config >> ~/.bashrc
echo "Sourcing .bashrc"
source ~/.bashrc

#echo "check .bashrc file to confirm alias is successfully configured as yogit2"
echo "---Confirm setup process was successful---"
cat -n ~/.bashrc

