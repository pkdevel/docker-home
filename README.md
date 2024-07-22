# docker-home
A simple docker home page

## Usage

1. docker
   ```bash
   docker run -d \
     -v /var/run/docker.sock:/var/run/docker.sock \
     -p 8080:8080 \
     --name=docker-home \
     pkdevel/docker-home
   ```

1. docker-compose
   ```yaml
   docker-home:
     image: pkdevel/docker-home:latest
     container_name: docker-home
     restart: unless-stopped
     ports:
       - 9080:8080
     volumes:
       - /var/run/docker.sock:/var/run/docker.sock
   ```

## License

MIT

## Repobeats
![Alt](https://repobeats.axiom.co/api/embed/29317c1e3f336aea42d4fd0edd5c3deb9c60ed38.svg "Repobeats analytics image")
