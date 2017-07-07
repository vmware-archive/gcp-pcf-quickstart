def GenerateConfig(ctx):
  """Reads SSL certificate and key from a file."""
  ssl = {'name': 'ssl-cert',
         'type': 'compute.v1.sslCertificate',
         'properties': {
             'certificate': '\n'.join([ctx.imports[ctx.properties['sslCertificatePath']]]),
             'privateKey': ctx.imports[ctx.properties['sslPrivateKeyPath']]}}
  return {'resources': [ssl],
          'outputs': [{
             'name': 'link',
             'value': '$(ref.ssl-cert.selfLink)'
          }]}
