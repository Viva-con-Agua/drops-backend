#!/bin/bash
source .env
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

generate_model()
{
    cat <<EOF
{
    "uuid":"${1}",
    "name":"${2}",
    "service_name":"${3}",
    "type":"${4}"
}
EOF
}
init_drops() {
    echo -e "${GREEN}START INIT DROPS-BACKEND ${NC}"
    echo -e "${GREEN}STEP_1: ${NC}initial model for drops-backend service"
    uuid=$(uuidgen)
    echo $API_IP
    data=$(generate_model ${uuid} drops-backend drops-backend service)
    echo $data
    curl -v -d "${data}" -H "Content-Type: application/json" POST http://${API_IP}:1323/admin/models
}

init_drops
