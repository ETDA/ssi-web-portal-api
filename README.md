
<h1 align="center">
    Web Portal API 
</h1>

<p align="center">
  <a href="#about">About</a> •
  <a href="#development">Development</a> 
</p>

## About

This repository contains the source code of the Self-Sovereign Identity for Web Portal API.

The SSI Web Portal API are for interact with another service and communicate with Frontend and mobile application. The Web Portal API required another services to run which provided in <a href = "##Prerequisites">Prerequisites</a>.

## Development

### Prerequisites

- [core-api](https://github.com/ETDA/ssi-core-api)
- [key-repository-api](https://github.com/ETDA/ssi-vc-schema-repository-api)
- [cloud-wallet-backup-api](https://github.com/ETDA/ssi-cloud-wallet-backup-api)
- [vc-schema-api](https://github.com/ETDA/ssi-vc-schema-repository-api)
- [mobile-e-wallet-api](https://github.com/ETDA/ssi-mobile-e-wallet-api)

#### Start Service
    
- Copy file `.env.sample` to `.env`
- run `docker-compose up -d`
- you can access the service via `http://localhost:8085`

