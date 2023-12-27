#!/bin/sh

source ./app-description.sh

smds_user_name='smdsuser'

executable_dst_dir_path="/opt/com/zmanhobbies/smds"

service_name=smds
service_file_name=$service_name.service
service_dst_dir=/etc/systemd/system

get_executable_full_name() {
    executable_full_name=$(ls . | grep "${service_name}-")
}

check_running_as_root() {
    if [ "$EUID" -ne 0 ]
    then echo "Please run as root"
    exit
    fi
}