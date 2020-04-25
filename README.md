# theangrydev-runner
Remote code execution service. No authentication over SSL/TLS required since this will run in the private cluster without a public IP. Gunicorn workers will act as proxies and route payloads to this service

![runner](logo.png)

# Okay, but why?

Sure, I can just expose the Docker daemon through a TCP/IP port and have it listen for requests over HTTP. It would still live in a private cluster so I wouldn't have to bother with securing it with TLS. Django will use the Python Docker SDK to ping the daemon

Basically there are two ways to run code
 * `docker run ...`
 * `docker exec ...`

 `run` is going to spawn a container every time, whereas `exec` will run the code against an active container. Even though spawning a container is way way cheaper than spawning a VM and there isn't a lot of overhead involved, I prefer the `exec` way    


But, I want more granular control over what container that `exec` command is run against based on container stats. I want to have a set of dedicated containers up and ready to serve requests for each runtime.  For eg: 3 containers for Python requests, 2 for Ruby, 5 for Rust etc etc

Once a request comes in, the runner will check which runtime the source code belongs to and pick the right container for that runtime based on CPU utilization, memory usage, network usage etc

So as you can see, we're doing a lot of back and forth with the Docker daemon and I believe doing this back and forth over TCP/unix is much more efficient than doing it over TCP/IP. We have to do the initial request over TCP/IP of course, but the rest of the communication should be done over `/var/run/docker.sock`

