# Upgrade
[reference](https://github.com/starkandwayne/om-tiler/blob/master/examples/README.md)

## requirements.
- [tile-config-generator](https://github.com/pivotalservices/tile-config-generator/releases).
- go 1.12 =>
- a running PCF Ops Manager

## generate/update tile configs.
get a legacy pivnet token [here](https://network.pivotal.io/users/dashboard/edit-profile).

```
tile-config-generator generate \
    --token=YOUR_LEGACY_TOKEN \
    --product-slug=elastic-runtime \
    --product-version=2.5.1 \
    --product-glob='srt*.pivotal' \
    --include-errands \
    --do-not-include-product-version \
    --base-directory=templates/assets/tiles/srt
```

## update pattern.yml.
edit/review `templates/assets/deployment.yml`

## generate template.
the yml templates in the assets directory should be embedded using go generate:

```
go generate src/omg-cli/templates/templates.go
```

# testing
follow general [readme](../../READM.md)
and run deploy_pcf.sh
