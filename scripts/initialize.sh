#!/bin/bash

# Esperar a que el contenedor arranque bien
sleep 2

# Iniciar mitmproxy con log permanente
mitmdump -s /app/intercept-master/inspect_requests.py \
         --listen-port 8080 \
         --set app=ALL \
         --set block_global=false \
         > /app/intercept-master/mitm.log 2>&1 &

# Iniciar control-master
./control-master traffic 4000 &

# Iniciar el executor
python3 ./executor-master/queue_receive.py



