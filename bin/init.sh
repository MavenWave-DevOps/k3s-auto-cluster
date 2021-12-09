#!/bin/bash
file=/boot/cmdline.txt
data=(cgroup_enable=cpuset cgroup_memory=1 cgroup_enable=memory)
autok3sRelease=v0.0.2
service_account=scott-testing123@airwaves-testing-01.iam.gserviceaccount.com
key_path=/home/pi/sa.json
project_id=airwaves-testing-01
template_repo=k3s-deployment-templates
export BASE_NET=$(sudo hostname -I)
export NODE_QUANTITY=2

sudo apt-get update && sudo apt-get install google-cloud-sdk

for item in ${data[@]}; do
  if grep -Fxq "cgroup_enable=cpuset" $file; then
    echo -e -n " $item" >> $file;
    sudo reboot;
  fi;
done;

wget https://github.com/tony-mw/k3s-auto-cluster/releases/download/v1.0.0/autok3s
chmod 777 autok3s
sudo nohup ./autok3s -baseNet=$BASE_NET -nodeQuantity=$NODE_QUANTITY &>/dev/null &

echo "Sleeping for 7 minutes.."
sleep 420

kubectl get nodes

gcloud auth activate-service-account $service_account --key-file=$key_path --project=$project_id
gcloud source repos clone $template_repo --project=$project_id
cd $template_repo
git checkout main

kubectl create secret docker-registry gcr-json-key \
  --docker-server=gcr.io \
  --docker-username=_json_key \
  --docker-password="$(cat $key_path)" \
  --docker-email=antonio.prestifilippo@mavenwave.com
for i in $(ls); do
  kubectl apply -f $i;
done;




