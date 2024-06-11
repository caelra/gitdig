# GITDIG

## Setup

### Prerequisites

Docker installed on your machine. You can download and install [Docker](https://www.docker.com/products/docker-desktop/)

### Environment variables

```bash
GITHUB_TOKEN="github_token"
GOOGLE_MAIL="some@gmail.com"
GOOGLE_APP_PASS="XXXXXXXXXXXXXXXXX"
```

### Building Docker image

To build the Docker image, run the following command:

```bash
docker compose build
```

After building the Docker image, you can run the container using:

```bash
docker run --rm -it gitdig-cli
```
