# First Time Setup Guide

## Goal

The goal of this doc is to store what I have done to set up the development environment for the first time.



## Create namespace controller

```sh
kubebuilder create api --group core --version v1 --kind Namespace --controller=true --resource=false
```

## Create athenz-domain poller

We want the operator to poll Athenz for changes to the domain, so we create a poller over time.

Maybe not the best practice, but we are giving a shot


