# Computation Movement Enactor

The Computation Movement Enactor (CME) is the DITAS component in charge of coordinating the movement of VDCs across different clusters.
It is intended to be run inside a VDM and to be automatically deployed by the DITAS [Deployment Engine](https://github.com/DITAS-Project/deployment-engine)

## Build

The CME is written in golang and go version 1.13 is required. To facilitate building and executing a [Dockerfile](Dockerfile) is provided. It can be built with ```docker build . -t ditas/computation-movement-enactor```
A [Jenkins file](Jenkinsfile) is provided as well to integrate into the DITAS CI environment.

## Configuration

The CME expects is configured by a file named ```cme.properties``` that must be mounted as a volume in the container path ```/etc/ditas```. The DITAS Deployment Engine will do it automatically when deploying this component inside a VDM. The properties that can be configured are:
- ```port```: The port in which the CME will listen for requests. By default it's 8080
- ```tombstone.use_ssl```: Set this property to ```true``` to communicate with the Tombstone component by https. Otherwise, it will try to communicate by http. By default it's ```false```
- ```deployment_engine.url```: This **mandatory** property holds the base URL of the DITAS Deployment Engine
- ```blueprint.id```: This **mandatory** property holds the abstract blueprint identifier of the VDM that's holding the CME
- ```sign.key```: This **mandatory** property holds the shared secret key between the CME and the Tombstone component of each of the VDCs it manages.
