## [Center For Fungus Control](https://github.com/jsmit257/centerforfunguscontrol)
A front-end to drive an instance of a [huautla](https://jsmit257.github.io/huautla) database.

### Requirements

#### Standalone Deployment
- docker/docker-compose (required): if you can get a docker daemon running and find a way to run `docker-compose` commands, either in the CLI, scheduled actions, `systemd` services, or whatever, then you can run this application standalone with persistence between restarts. Sounds really simple, and it is, and if that's what you're looking for you can skip ahead to the dedicated [standalone page](./standalone/README.md).
- postgresql (optional): this is really only useful if you already have one, or have always wanted a reason to have one, and want to use it to manage your own backups, failover, etc. The default datasource will still start, unless you turn it off as a dependency in the `cffc-standalone:` service. More details about the data store can be found in [deployment](#deployment)

#### Development
- [Standalone Deployment Requirements](#standalone-deployment)
- golang development environment: 1.22.1 for now
- `git`: repo is public, but you'll need an account to push from to submit a pull-request
- `make`: you could work around it, but `Makefile` is a good place to write things down
- postgresql (optional) maybe easier if you have one, or at least the `postgres-client` tools (`psql`, `pg_dump`, `pg_restore`, etc)

### Docker
Docker images are the only build artifacts from this project, hosted on [dockerhub](https://hub.docker.com/repository/docker/jsmit257/cffc/tags?page=1&ordering=last_updated). Proper semantic or other versioning is still in committee, for now the `:lkg` tag LGTM.

Production images use the statically-compiled application from [./ingress/http/main.go] as `entrypoint`, and a copy of [./www] to host the UI static resources.

### Deployment
There are 3 ways to configure the database:
- Docker-only: means all data is lost when the cuntainer stops running
- Host-only: means all data is sourced from a non-Docker server, somewhere; data is persisted after shutdown
- Docker/Host-storage: the server process runs in a container, but the data is written to a volume mapped from the host's filesystem, so data is preserved after shutdown

And, 3 ways to run the HTTP server:
- Local (i.e. `go run ...`)
- Docker+build
- Docker:latest

Now you have to choose which of the 9 configurations best suits your current need. Sorry if that doesn't seem helpful, but there are just a few canned configurations, one of which should be sufficient for your case:
- the standalone solution described [here](./standalone/README.md) is arguably the best for a high-availability server with persistent storage, even if you're actively developing the project.
- `make run-local` all the database and http server configurations can be set on the commandline, as needed, and web-resources are published in real-time, so this usually best for development.
- `make tests`: do this before pushing new tags to dockerhub; the server stays up for user testing

Of course, you can run more than one instance for different reasons, even connected to the same database, but be especially careful never to let two server processes access the same postgres `data/` directory.

#### Database
`jsmit257/huautla:lkg` is the blessed docker image with the huautla database installed and pre-populated including a few values useful for metrics. A couple test-suites can populate a database with sample data so you don't have to enter a bunch yourself, but everything is lost when the service stops.

See [local huautla](https://github.com/jsmit257/huautla/blob/master/README.md#local-database) for tips on installing using native postgres-client tools.

For persistent storage w/o installing postgres on the host, a simple solution is to tweak the `DEST_*` environment variables in the `migration:` service in [standalone](./standalone/docker-compose.yml), and then start it. Be sure to revert your local changes before committing anything.

#### HTTP Server
- `docker-compose up --remove-orphans -d run-docker`: starts the `:latest` version of the server in the background, with *no* rebuild

- `[HUAUTLA_*=...] [HTTP_*=...] make run-local`: this is the easiest way to support development. All the configs are on the commandline. Web resources are read per-request, so changes are available w/o a server restart. And for `.go` changes, the server restarts faster this way than with docker `--build`.

- `make tests`: does a few things besides just starting the server
  - runs unit tests
  - starts a vanilla instance of jsmit257/huautla:lkg for storage
  - builds the standalone docker image as `cffc:latest` including compiling the server application and copying static web resources into the container
  - starts the server container
  - runs the suite of system tests (currently, just seeds some test data)
  - tags `cffc:latest` as `jsmit257/cffc:lkg`

### API
Crap! Didn't think about documentation much

URL base paths for vendor, stage, eventtype, substrate, ingredient, strain and lifecycle resources all follow the same basic pattern:

```
GET /<url-base-path>s        the plural of the path returns all rows
GET /<url-base-path>/$id     returns the row uniquely identified by $id
POST /<url-base-path>        adds a new record to the url-base-path resource
PATCH /<url-base-path>/$id   modifies the record uniquely identified by $id
DELETE /<url-base-path>/$id  deletes the record uniquely identified by $id
```

Request and response bodies are defined by the [data type](https://github.com/jsmit257/huautla/blob/master/types/data.go) corresponding to the resource name - i.e. `http://host:port/vendor` sends and recieves JSON marked up in `type Vendor struct{...}`. Normally, responses are the complete tree, the exception being `/lifecycles`, which only includes a few fields and *none* of the strain/substrate/event data.

Status codes for the preceeding resources are:
- Sussessful `POST` returns `201 Created` and the new resource with a newly-generated `id` attribute
- Successful `DELETE` returns `204 No Content`
- Successful anything else returns `200 OK`
- Errors are either `400 Bad Request` or `500 Internal Server Error`. Error responses also contain a `cid:` header which is the unique correlation ID generated by the server for each request and written on each log message in the call-stack. Error response bodies are a work in progress - right now they're not very informative.

A few 'irregular verbs' are defined for managing resources that are lists of attributes for their parent resource. In the cases of `substrate.ingredients`, `strain.attributes` and `lifecycle.events`, the general URL pattern is:

```
POST|PATCH /<parent-resource>/${parent_id}/<child-resource>
DELETE /<parent-resource>/${parent_id}/<child-resource>/${child_id}
```

`POST` and `PATCH` request bodies only contain the [data type](https://github.com/jsmit257/huautla/blob/master/types/data.go) of the child resource, but all methods including `DELETE` return a type of parent resource with all child resources populated.

There are no `GET`s because children only make sense in the contexts of their parents, and parents always eager-fetch all their children.

Status codes for these parent/child methods are the same the regular routes, except `DELETE` returns `200 OK` and a response body as described above

Only `PATCH /substrate/$substrate_id/ingredients/$current_ingredient_id` totally breaks the mold as it doesn't use surrogate keys to identify the record that needs changing, so '$current_ingredient_id' is what's in the database now, and the request body contains the new ingredient identifier. Keep in mind that the child resource type in this case is `SubstrateIngredient`, *not* `Ingredient`

An additional convenience method is provided to get a list of unique strain attribute names:

```
GET /strainattributenames  # returns an array of strings
```

### Contributing
License forthcoming, maybe creative commons or MIT, something with attribution. Don't let that stop you from contributing. Add issues, submit pull requests, etc.
