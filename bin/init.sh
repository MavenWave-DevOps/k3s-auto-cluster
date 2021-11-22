#!/bin/bash
file=/boot/cmdline.txt
data=(cgroup_enable=cpuset cgroup_memory=1 cgroup_enable=memory)

for item in ${data[@]}; do
  echo -e -n " $item" >> $file;
done

