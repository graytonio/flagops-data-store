# FlagOps Datastore

**NOTE:** Experimental not intended for production use

This server is use to support an open feature provider tailored for flagops usage. It allows supporting data to be stored and accessed through a unified sdk instead of several distinct ones. This data is then used as part of the ArgoCD rollout process to template out and ultimately deploy the k8s manifests to the cluster.

## Kinds of Data

For this project data has been split into 3 sections:

1. Facts - These are non-sensitive traits about a deployment such as the cluster it is indented to be deployed in, application name, owner of the application, etc.

2. Secrets - Sensitive data required for the application to run such as db connection strings, certificates, or cloud service credentials.

3. Feature Flags - These are special configuration flags that change how the deployment manifests are templated. They are managed separately from this server but use the facts stored here to determine their output.

## Expected Usecases

The intended usecase is using flagops setup in ArgoCD as the templating plugin for your manifests which use feature flags from one of the supported providers to determine their value.

The facts both from the ArgoCD application template as well as any facts configured in this server are passed to the feature flag provider as part of the evaluation context which then returns the state of feature flags.

This allows other mechanisms outside of the git state to update and fetch metadata about a deployment and update facts based on automated processes.

As an example, for a process each customer generates a single tenant deployment which is tied to their account. Based on what tier this customer signed up for they receive a max replica count of either 1, 3 or 5. The account level can change and is not stored in git but in some CMS system. This system can export the account level data through a webhook or api call to update the facts for customers as soon as it changes. This change is then applied to the feature flag rules which updates the max replicas accordingly.

## Roadmap

- [ ] API Keys
- [ ] Fact provider Postgres
- [ ] UI for visualizing
- [ ] Test OpenFeature Provider Golang
- [ ] OpenFeature Provider Python
- [ ] Metrics
