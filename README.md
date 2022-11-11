# Version Notifier</br>[![Codacy Badge](https://app.codacy.com/project/badge/Grade/4704892fd733422bbb6dbec098c709be)](https://www.codacy.com/gh/yuvalpress/version-notifier/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=yuvalpress/version-notifier&amp;utm_campaign=Badge_Grade) [![Docker Build](https://github.com/yuvalpress/version-notifier/workflows/Docker%20Build/badge.svg)](https://github.com/yuvalpress/version-notifier/actions?query=workflow%3ADocker%20Build) [![Chart Release](https://github.com/yuvalpress/version-notifier/workflows/Chart%20Release/badge.svg)](https://github.com/yuvalpress/version-notifier/actions?query=workflow%3A%22Chart+Release%22) [![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
🕷 **Version Notifier** - your Friendly Neighborhood Spiderman, only geeker 🤓

Version Notifier is a modern solution for the "being notified" aspect of each Techy's day-to-day work.
</br>By using it, you'll be notified for any new global GitHub repository release you choose, directly to your Slack channel.</br></br>

## Getting Started 🏁
You can deploy the application in one of two ways:</br>
### HELM
Download the latest release and deploy it to your Kubernetes cluster </br>
  ```shell
    helm repo add vnotifier https://yuvalpress.github.io/version-notifier
    
    kubectl create ns notifier
    helm install version-notifier vnotifier/version-notifier -n notifier
  ```

### Docker Image
Create a dockerfile from the Version-Notifier base image and deploy it as a standalone container:
  ```dockerfile
    # Name this file Dockerfile
    FROM yuvalpress/version-notifier:latest
    
    # You MUST Set this environment variables for the application to send notification to slack
    ENV SLACK_CHANNEL {{ value }}
    ENV SLACK_TOKEN {{ value }}

    # Optional
    ENV NOTIFY {{ value }}
  ```
  
  Build and Deploy:
  ```shell
    # Run this command from the Dockerfile dir
    docker build -t {{ value }} .
    docker run --name {{ value }}
  ```
</br></br>
## Configuration Options 🕹
### NOTIFY
List represented as string with the following possible keywords: `major, minor, patch, all`
</br>This value can be set in both HELM values.yaml file under `application.notify` and as environment variable in your custom Dockerfile.
</br></br> Possible combinations:
  * "all" - `all` must be set alone
  * "major, patch" - only notify for `major` and `patch` version changes
  * "minor" - only notify about `minor` version changes

If not set, NOTIFY will be automatically set to `all`</br></br>

### config.yaml
The config.yaml file holds the repositories to be scraped by Version-Notifier.

Example repo template: `<github-user>: <repository>`
#### Edit with HELM:
![values.yaml](./docs/repos-helm.png)

#### Add to custom Dockerfile:
1. Create a file called config.yaml, place it under the same folder as the Dockerfile and populate it like such:
```yaml
repos:
    - confluentinc: terraform-provider-confluent
    - hashicorp: terraform-provider-aws
    - hashicorp: terraform-provider-google
```
2. Add it to your custom Dockerfile:
  ```dockerfile
    # Name this file Dockerfile
    FROM yuvalpress/version-notifier:latest

    # add custom config.yaml file
    COPY config.yaml ./config.yaml
    
    # You MUST Set this environment variables for the application to send notification to slack
    ENV SLACK_CHANNEL {{ value }}
    ENV SLACK_TOKEN {{ value }}
  ```
</br></br>
## Verification of Success 🎯
If the deployment was successful, you'll see the logs rolling out of your container:
### Using Docker
If you executed Version Notifier using Docker, you'll see the logs roll after you run the container.</br></br>
![Docker Run](docs/docker-run.gif)

### Watch logs with kubernetes
```shell
pod=$(kubectl get pods -n notifier -l app=version-notifier -o yaml | yq '.items[0].metadata.name') && kubectl logs $pod -n notifier -f
```
<br></br>
## Upcoming Features ✨
* Analyzing tags of Helm Charts released using the `helm/chart-releaser-action` GitHub Action.
* Support for more notification methods (currently Slack only). 

<br></br>
## Want to contribute? 💻
PR's are more than welcome!

#### Steps:
1. Open a branch in the following form: `feature/<feature_name>`.
2. Make sure to bump the Docker Image version by incrementing the version inside the docker.version file.
3. Open PR!