PROJECT_PATH=$(pwd)

closeProc() {
    server=$1
    ps -fe|grep ${server}|grep -v grep
    if [ $? -ne 0 ]
    then
        echo "${server} not existed"
    else
        echo "Closing old ${server}"
        killall ${server}
    fi
}

start_log() {
    cd ${PROJECT_PATH}
    mkdir -p logs
    chmod 777 logs
}

start_server() {
    server=$1

    cd ${PROJECT_PATH}/${server}
    echo "start to build ${server}"
    go build .
    echo "starting ${server}"
    nohup ./${server} -c ../conf/$server.yaml > ${PROJECT_PATH}/logs/${server}.log 2>&1 &
}

start_log
closeProc "tcpserver"
start_server "tcpserver"
closeProc "httpserver"
start_server "httpserver"
