#!/bin/bash

  curl -L -o p7zip_16.02_x86_linux_bin.tar.bz2 "https://drive.google.com/uc?export=download&id=${P7Z_ID}"
  tar xvf p7zip_16.02_x86_linux_bin.tar.bz2 
  curl -c ./cookie -s -L "https://drive.google.com/uc?export=download&id=${FILE_ID}" > /dev/null
  curl -Lb ./cookie "https://drive.google.com/uc?export=download&confirm=`awk '/download/ {print $NF}' ./cookie`&id=${FILE_ID}" -o /tmp/vcenter.iso
  mkdir mnt
  ./p7zip_16.02/bin/7z x /tmp/vcenter.iso -o ./mnt
  rm vcenter.iso
  pwd
  ls

  cp main.tf.phase1 main.tf
  terraform init
  terraform apply --auto-approve
  cp main.tf.phase2 main.tf
  terraform init
  terraform apply --auto-approve
  rm -Rf ./mnt
