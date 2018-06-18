
### replace ingress

**Aliases**   :
  ingr, ingresses, ing
**Usage**     :
 Replace ingress with a new one, use --force flag to write one-liner command, omitted attributes are inherited from the previous ingress.
**Example**   :
  chkit replace ingress $INGRESS [--force] [--service $SERVICE] [--port 80] [--tls-secret letsencrypt]
**Flags**     :
  + force f : replace ingress without confirmation
  + host  : ingress host, optional
  + port  : ingress endpoint port, optional
  + service  : ingress endpoint service, optional
  + tls-secret  : ingress tls-secret, use 'letsencrypt' for automatic HTTPS, '-' to use HTTP, optional
  
**Subcommand**:
  
