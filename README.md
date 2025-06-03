# Regnerdetioslo.no

This repository contains the sourcecode for the website [regnerdetioslo.no](https://regnerdetioslo.no/)
 
A simple webpage that shows if it is currently raining in Oslo, or not.

## Running locally

```bash
docker run -it --rm -p 8080:8080 $(docker build -q .)
```