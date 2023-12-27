#!/bin/sh

source ./app-description.sh

build_dir_path=../build

create_build_dir_not_exist() {
    #echo Going to create build directory if it doesnt exist

    if [ ! -d $build_dir_path ]; then 
        mkdir $build_dir_path
    fi
}

clean_build_dir() {
    if [ -d $build_dir_path ]; then 
        rm -r $build_dir_path
    fi
}

get_app_version_from_git() {
    app_version=$(git describe --tags --abbrev=0)
}

get_executable_full_versioned_name() {
    os_name=$(uname -s)
    if [ -n "${GOOS}" ]; then
        os_name=$GOOS
    fi

    arch_name=$(uname -m)
    if [ -n "${GOARCH}" ]; then
        if [ -n "${GOARM}" ]; then 
            arch_name=$GOARCH-$GOARM
        else
            arch_name=$GOARCH
        fi
    fi

    #Need to call this to set the variable holding the app version
    get_app_version_from_git

    executable_full_versioned_name=$executable_name-$os_name-$arch_name-$app_version
}

build_app() {
    

    create_build_dir_not_exist
    
    get_executable_full_versioned_name

    app_output_path=$build_dir_path/$executable_full_versioned_name

    echo building app, output to: $app_output_path

    go build -o $app_output_path $main_file_path
}