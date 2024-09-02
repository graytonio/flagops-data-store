# Configuration Variables

The flagops datastore is configured through environment variables. This table describes the available env and what they control.

| Key                                | Description                                                                                      | Default        |
| ---------------------------------- | ------------------------------------------------------------------------------------------------ | -------------- |
| FLAGOPS_SECRET_PROVIDER            | Select which provider to store secrets in                                                        | asm            |
| FLAGOPS_ASM_DELETION_RECOVERY      | Number of days to use for recovery window when deleting identities with the ASM secrets provider | 7              |
| FLAGOPS_FACT_PROVIDER              | Select which provider to store facts in                                                          | redis          |
| FLAGOPS_POSTGRES_DB_DSN            | DSN for postgres db to connect to for storing user and permissions data                          | ""             |
| FLAGOPS_USER_SESSION_SALT          | Random salt string used in securing user seesions. Recommended to set for production deployments | "flagops-salt" |
| FLAGOPS_REDIS_URI                  | URI for redis when using the redis facts provider                                                | ""             |
| FLAGOPS_OAUTH_PROVIDER             | Oauth2 login provider                                                                            | ""             |
| FLAGOPS_GITHUB_OAUTH_CLIENT_KEY    | Github oauth client key when using github oauth                                                  | ""             |
| FLAGOPS_GITHUB_OAUTH_CLIENT_SECRET | Github oauth client secret when using github oauth                                               | ""             |
| FLAGOPS_HOSTNAME                   | Domain of deployment used in Oauth2 redirections                                                 | ""             |
