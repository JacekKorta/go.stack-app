# go.stack-questions
<a href="https://www.repostatus.org/#wip"><img src="https://www.repostatus.org/badges/latest/wip.svg" alt="Project Status: WIP – Initial development is in progress, but there has not yet been a stable, usable release suitable for the public." /></a><br>

This service is a part of the repo: [microservices-training-ground](https://github.com/JacekKorta/microservices-training-ground)<br>
The main purpose of this service is to subscribe questions from the stackoverflow site. Each instance of the service can subscribe to one tag. If you need to subscribe more tags please run more instances of the service. 

### How to run?

You should run this service via docker compose in main repo [microservices-training-ground](https://github.com/JacekKorta/microservices-training-ground)

Create env file:

```bash
scripts/create_env_file.sh <your tag as arg>
```

### Additional info about env variables

APP_URL - stackoverflow api address<br>
FILTER - filters are responsible for the SO rest api response shape. You shouldn’t change it. <br>
TAGGED - Your tag <br>
REQEST_LIMIT_PER_SEC= request limit per second. should be lower than 30 <br>
DELAY_BETWEEN_CHECKS=time (minutes) between checking for the new questions. Should be greater or equal to 1. <br>
