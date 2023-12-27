#!/bin/sh

source ./app-description.sh
source ./build-common.sh
source ./service-common.sh

deployable_dir_path=../deployable

create_deployable_dir_not_exist() {
    if [ ! -d $deployable_dir_path ]; then 
        mkdir $deployable_dir_path
    fi
}

clean_deployable_dir() {
    if [ -d $deployable_dir_path ]; then 
        rm -r $deployable_dir_path
    fi
}

create_specific_service_file() {
    #cp ../systemd/$service_file_name $deployable_dir_path/$service_file_name
    sed "s/ExecStart=/ExecStart=\/opt\/com\/zmanhobbies\/smds\/${executable_full_versioned_name}/" ../systemd/$service_file_name > $deployable_dir_path/$service_file_name 
}

cleanup_specific_service_file() {
    rm $deployable_dir_path/$service_file_name
}

create_deployable() {

    create_deployable_dir_not_exist

    get_executable_full_versioned_name

    create_specific_service_file

    zip -j "${deployable_dir_path}/${executable_full_versioned_name}.zip" "${build_dir_path}/${executable_full_versioned_name}" "./service-common.sh" "./install-service.sh" "./uninstall-service.sh" "./app-description.sh" $deployable_dir_path/$service_file_name

    cleanup_specific_service_file
}