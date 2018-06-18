
### create ingress

**Aliases**   :
  ingr, ingresses, ing
**Usage**     :
 Create ingress. Available options: TLS with LetsEncrypt and custom certs.
**Example**   :
  chkit create ingress [--force] [--filename ingress.json] [-n prettyNamespace]
**Flags**     :
  + force f : create ingress without confirmation
  + host  : ingress host (example: prettyblog.io), required
  + path  : path to endpoint (example: /content/pages), optional
  + port  : ingress endpoint port (example: 80, 443), optional
  + service  : ingress endpoint service, required
  + tls-cert  : TLS cert file, optional
  + tls-secret  : TLS secret string, optional
  
**Subcommand**:
  
