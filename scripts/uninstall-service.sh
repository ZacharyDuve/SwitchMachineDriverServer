#!/bin/sh

source ./service-common.sh

stop_remove_service() {
    systemctl stop $service_name

    systemctl disable $service_name

    rm $service_dst_dir/$service_file_name
}

# Need to check that we are root as we will need to run systemctl commands
check_running_as_root

stop_remove_service

