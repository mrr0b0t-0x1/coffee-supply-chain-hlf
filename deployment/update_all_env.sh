for DIR in *_api/; do echo updating env vars for $DIR; cd $DIR; ./deployment/set_env_vars.sh; cd ..; done
