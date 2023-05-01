#!/bin/bash

cd ..
ls -la
exit
cd adapter

printf "\nCompiling code...\n"
npm run build

printf "Compilation done! Building image...\n\n"
DOCKER_BUILDKIT=1 docker build -f deployment/Dockerfile -t coffee-supply-chain .

printf "\nImage built! Starting all components...\n\n"
cd ../deployment
./start_all_component.sh

printf "\nDone!\n"
