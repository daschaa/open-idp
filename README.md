# OpenIdP

OpenIdP is an open-source Identity Provider (IdP) server built to explore and demonstrate the OAuth 2.0 framework. 
It provides a hands-on implementation of various OAuth 2.0 grant types, so that I understand how they work and how they can be implemented in a real-world scenario.

## Roadmap

- [x] Client credentials grant
- [ ] Authorization code grant
- [ ] Implicit grant
- [ ] Resource owner password credentials grant

## Project Overview

**Language**: Written in **Go**, designed for simplicity and performance.
**Deployment**: Deployable as an AWS Lambda function.
**Infrastructure as Code**: Uses **AWS CDK** for infrastructure management. Infrastructure code is located in the _infrastructure directory.

## Getting Started

### Client Credentials Flow

The [client credentials flow](https://datatracker.ietf.org/doc/html/rfc6749#section-4.4) is designed for machine-to-machine authentication. 
This flow allows an application to use its client ID and client secret to request an access token, which can then be used to access protected resources.

#### Example Use Case
A backend service authenticates itself to access an API without requiring user involvement.

#### How It Works

1. The client sends a POST request to the /token endpoint with its client ID, client secret, and grant type.
2. If the credentials are valid, the server responds with an access token.
3. The client uses the token to authenticate API requests.

