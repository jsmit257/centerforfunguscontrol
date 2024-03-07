## Huautla
Named after the town in Jiménez that was the bridge between native traditions and what qualifies as a modern understanding of psychedelics. There's an interesting story of people and politics behind their noteriety, so like any memorial, this is a nod to what we consider to be a remarkable moment in history.

### Overview
The goal of this project is to track analytics observing the lifecycle of fungi. There's a lot more to consider before we could call this list of attributes comprehensive for many fungi, but for now we've captured the few events that cover most of our needs. Even this highly abstract overview of a colony over time is missing many species-specific attributes. We'll continue to evolve the analytics as we encounter more distinct charecteristics.

More details of a lifecycle are described in the [object model](#object-model), but at a higher level, lifecycles evolve through some of the following 4 stages:

* Gestation: These are the range of events between culturing spores, and having a live colony in stasis, so to speak, that's ready to be introduced to an environment that can support growth to the point of procreation. Simply and roughly put, this is the time between a spore-print, and life outside of sugar-water
* Colonization: Maybe 'adolescence' is a better name. Probably sealed in a bag of grain(s), woodchips or ???, this is when cosmopolitan organisms establish a colony, share resources, build strong contamination resistance - all in the safety of closed bubble.
* Majority: This is when the colony is established and strong and can grow to the excess needed to fruit new generations. An example is mixing colonized adolescents from the previous stage with porous, moist, nutritive bulk substrate in a bed. They can consume much more than the colony needs, so the extra will be used for procreation.
* Vacation: A sort of suspended animation when a fungus is ahead of schedule and needs to stop working. Fortunately, they chill well with minimal damage, and simply wake up when they get warm again. So time spent on `Vacation` within a lifecycle isn't fully counted into a normal temporal model for a particular strain, but it still needs to be predictable in terms of how long they can chill, and whether they end with fruits, or the colony dies.

### Requirements
To develop and test locally, you'll need at least:
* Git (you'll need an account if you plan to change anything)
* Golang environment
* docker/docker-compose (sometimes they're packaged separately, get both)
* postgres-client (optional) is useful if you want to poke around the db after a feiled test, otherwise it's not needed
* postgres-server (optional) if you want to work/test on a database that won't be removed when the tests finish
* make
* bash-compatible shell (bourne doesn't handle arrays nicely)

### Basics
???

### Object Model
These are the database tables described as a golang object tree; the sql perspective is [here](./sql/init.pgsql), and there are [cliff-notes](./sql/cliff-notes); no, this isn't really yaml

```yaml
vendor: &vendor
  uuid: surrogate key
  name: must be unique

strainattribute: &strainattribute
  uuid: surrogate key
  name: the key in the key/value pair, can't be changed after it's created
  value: arbitrary value - use it for anything

strain: &strain
  uuid: surrogate key
  name: just a name, without a vendor it doesn't mean much
  vendor: *vendor ### vendors give reputation to names, so (name+vendor) is unique in this table
  attributes: list of *strainattribute

ingredient: &ingredient
  uuid: surrogate key
  name: unique in this table, but can be shared by many substrates

substrate: &substrate
  uuid: surrogate key
  name: the name the vendor gave it
  type: constrained to ('Bulk', 'Grain', TBD)
  vendor: *vendor ### substrates are unique by (name+vendor) for the same reasons as strains
  ingredients: list of 0 or more *ingredient

grainsubstrate: &grainsubstrate *substrate
  ### these are just substrates that pass an API level check (not a check constraint on the database, yet) that only allows types of 'Grain'

bulksubstrate: &bulksubstrate
  ### like grainsubstrate, except the type check is for 'Bulk'

stage: &stage
  uuid: surrogate key
  name: unique to this table

eventtype: &eventtype
  uuid: surrogate key
  name: something short and descriptive
  severity: label for quick filtering of log noise
  stage: *stage ### unique with name, b/c some types wouldn't apply to some stages, and there may be dublicate names

event: &event
  uuid: surrogate key
  temperature:
  humidity:
  mtime: date/time when this record was last modified
  ctime: date/time when this record was created
  eventtype: *eventtype

lifecycle:
  uuid: surrogate key
  name: unique identifier ### probably going away or being otherwise refactored
  location: where the bin/bag/jar/agar is being stored
  grain_cost:
  bulk_cost:
  yield: in grams, how much dried or fresh product was shipped (see gross)
  count: how many caps were harvested?
  market_price: TBD - this should be the market value when the lifecycle begins; it may go up, it may go down, the only thing you know for sure is this is where it's at now; once entered, you cacn't change this, but that's not an invitation to speculate
  gross: product fresh weight, mostly to anticipate weight loss from dehydration
  mtime: date/time when this record was last modified
  ctime: date/time when this record was created
  strain: *strain
  grainsubstrate: *grainsubstrate
  bulksubstrate: *bulksubstrate
  events: list of *event
```
The root of anything interesting is the Lifecycle. It's generally a path from inception to completion, but it's a 'cycle' and not a 'span' because it doesn't always start at the beginning, and it will eventually support spawning one lifecycle to create new lifecycles, which makes them more tree-like than discrete instances (though ironically not cyclical). But that isn't really finished yet

### Configuring
Basically, fill out this form and submit it to `huautla.New()`, along with a logger.
```go
Config struct {
  PGHost string
  PGUser string
  PGPass string
  PGPort uint
  PGSSL  string
}
```
There's a workable reference implementation in the system test [init()](./tests/system/main_test.go) function.

FWIW, that test implementation expects a return type of `types.DB`. As you can see from [the api](./types/api.go), this is all the methods for all the entities in the database. That's a miserable test-harness to have to maintain in a properly modular client. Instead, consider this pattern:

```go
db, err := huautla.New(...)
if err != nil { ... }

huautlaAdaptors := struct{
  ClientEventer
  ClientEventTyper
  ...
  ClientVendorer
}{
  ClientEventer: ClientEventer {
    config: ...,
    Eventer: db,
  },
  ClientEventTyper: { ..., EventTyper: db },
  ...
  ClientVendorer: : { ..., Vendorer: db },
}
```
Where each `Client*` type is a candidate to be a receiver, encapsulating its own context/configuration/etc, with its reference to `types.DB` constrained to just the functionality for a single relation, while identifying this database unambigiously from others used by this client.

There are also all the anonymous interface shennanigans you can use in function signatures and other typecasts, but worry about that when you really need to.

### Installing
This just means installing the database, since the API is installed by vendoring. The `install` target in [the makefile](./Makefile), along with the `install` service in [docker-compose](./docker-compose.yml) will connect with a database instance and run the create statements necessary to build a complete `huautla` database with users, tables, etc. By default, the `install` service connects to the docker-compose `postgres` service, which is great if you're going to run tests, but everything gets lost when the containers go down. Overriding the pg* environment vars in `docker-compose` allows to connect to a durable instance (e.g. not a container), and the created database will be suitable for hosting any application that uses this project. It will *not* clobber an existing database - you've got to do that manually, for legal reasons.

### Using
Comments are inlined in the [api](./types/api.go) source, rather than have to maintain that file and a README separately.

### Testing
TBD: tag docker postgres+huautla images with code releases

Until then, the [docker-compose](./docker-compose.yml) should be pretty portable to client environments, but you'll need to change the location of the [sql](./sql/) and [bin](./bin) script sources, unless you're cloning this repo standalone (really not a bad option). You only need enough tweaking to run the `install` service mentioned [above](#installing).

### Contributing
