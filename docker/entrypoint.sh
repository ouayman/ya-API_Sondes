#!/bin/sh
set -e

echo "Spring Cloud Config URI: $APP_SPRING_CONFIG_URI"
echo "app profiles: $APP_PROFILES"

if [ -n "$APP_SPRING_CONFIG_URI" ] && [ -n "$APP_PROFILES" ]; then
    echo "Run with Spring Cloud Config"
	exec "$@"
else
	echo 'Expecting a spring cloud config in environement variable $APP_CONFIG_FILE and $APP_PROFILES'
fi
