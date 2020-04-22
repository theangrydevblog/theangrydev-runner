# theangrydev-runner
Remote code execution service. No authentication over SSL/TLS required since this will run in the private cluster without a public IP. Gunicorn workers will act as proxies and route payloads to this service

![runner](logo.png)
