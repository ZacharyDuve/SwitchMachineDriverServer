#!/bin/sh

source ./service-common.sh

useradd_smds_user() {
    if id -u $smds_user_name >/dev/null; then
        echo $smds_user_name was already found
    else 
        echo $smds_user_name was not found. Going to add

        useradd $smds_user_name
    fi
}

place_executable() {
    echo copying "${executable_full_name}" to "${executable_dst_dir_path}/${executable_full_name}"
    cp "${executable_full_name}" "${executable_dst_dir_path}/${executable_full_name}"
    if [ $? -eq 0 ]; then
        echo 'Executalbe was successfully placed'
    else
        echo 'Error occured during copy so halting install'
        exit
    fi
    echo done copying executable to correct location
}

enable_serivce() {
    echo placing service into correct location
    cp $service_file_name $service_dst_dir/$service_file_name
    if [ $? -eq 0 ]; then
        echo Service was successfully placed
    else
        echo Error occured placing service stopping
        exit
    fi


    if systemctl enable $service_name; then
        echo Service was enabled
    else
        echo Error occured enabling service
        exit
    fi

    if systemctl start $service_name; then
        echo Service was started
    else
        echo Error starting the service
        exit
    fi
}

ensure_executable_dst_dir_exists() {
    if [ ! -d $executable_dst_dir_path ]; then
        mkdir -p $executable_dst_dir_path
    fi
}

#Need to check that we are running as root as it is needed for install
check_running_as_root

#Need to add the user that will be running the application
useradd_smds_user

#Need to lookup the name of the executable for the copy command
get_executable_full_name

#Need to make sure that the directory that we want to place the executable into exists so we can copy to it
ensure_executable_dst_dir_exists

#Put the application in the correct location
place_executable

#Need to place the service file and enable and start it
enable_serivce