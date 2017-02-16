# Dynamic Routing of K8S Containers

This uses the net/http router to create servers
from environment variables.
Expects environment variables like: OUTSCORE_DEPLOYMENT_PORT_8080_TCP_ADDR=10.95.249.177

When running in K8S these variables will be created by default.

For example, the environment variable `OUTSCORE_DEPLOYMENT_PORT_8080_TCP_ADDR=10.95.249.177` will
create a reverse proxy on port 8080 to 10.95.249.177
