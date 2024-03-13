package data

type sqlMap map[string]map[string]string

var psqls = sqlMap{

	"event": {
		"all-by-lifecycle": `
    select e.uuid,
           e.temperature,
           e.humidity,
           e.mtime at time zone 'utc',
           e.ctime at time zone 'utc',
           et.uuid as eventtype_uuid,
           et.name as eventtype_name,
           et.severity as eventtype_severity,
           s.uuid as stage_uuid,
           s.name as stage_name
      from events e
      join event_types et
        on e.eventtype_uuid = et.uuid
      join stages s
        on et.stage_uuid = s.uuid
     where lifecycle_uuid = $1
     order
        by mtime desc`,
		"all-by-eventtype": `
    select e.uuid,
           e.temperature,
           e.humidity,
           e.mtime,
           e.ctime,
           et.uuid as eventtype_uuid,
           et.name as eventtype_name,
           et.severity as eventtype_severity,
           s.uuid as stage_uuid,
           s.name as stage_name
      from events e
      join event_types et
        on e.eventtype_uuid = et.uuid
      join stages s
        on et.stage_uuid = s.uuid
     where et.uuid = $1`,
		"select": `
    select e.temperature,
           e.humidity,
           e.mtime,
           e.ctime,
           et.uuid as eventtype_uuid,
           et.name as eventtype_name,
           et.severity as eventtype_severity,
           s.uuid as stage_uuid,
           s.name as stage_name
      from events e
      join event_types et
        on e.eventtype_uuid = et.uuid
      join stages s
        on et.stage_uuid = s.uuid
     where e.uuid = $1`,
		"add": `
    insert
      into events(uuid, temperature, humidity, mtime, ctime, lifecycle_uuid, eventtype_uuid)
    select $1, $2, $3, $4, $5, lc.uuid, et.uuid
      from lifecycles lc,
           event_types et
     where lc.uuid = $6
       and et.uuid = $7`,
		"change": `
    update events e
       set temperature = $1,
           humidity = $2,
           mtime = $3,
           eventtype_uuid = et.uuid
      from event_types et
     where e.uuid = $4
       and et.uuid = $5`,
		"remove": `delete from events where uuid = $1`,
	},

	"eventtype": {
		"select-all": `
    select e.uuid,
           e.name,
           e.severity,
           s.uuid as stage_uuid,
           s.name as stage_name
      from event_types e
      join stages s
        on e.stage_uuid = s.uuid
     order
        by s.name, e.name`,
		"select": `
    select e.name,
           e.severity,
           s.uuid as stage_uuid,
           s.name as stage_name
      from event_types e
      join stages s
        on e.stage_uuid = s.uuid
     where e.uuid = $1`,
		"insert": `
    insert
      into event_types(uuid, name, severity, stage_uuid)
    select $1,
           $2,
           $3,           
           s.uuid
      from stages s
     where s.uuid = $4`,
		"update": `
      update  event_types 
         set  name = $1,
              severity = $2,
              stage_uuid = $3
       where  uuid = $4`,
		"delete": `delete from event_types where uuid = $1`,
	},

	"ingredient": {
		"select-all": `select uuid, name from ingredients order by name`,
		"select":     `select name from ingredients where uuid = $1`,
		"insert":     `insert into ingredients(uuid, name) values($1, $2)`,
		"update":     `update ingredients set name = $1 where uuid = $2`,
		"delete":     `delete from ingredients where uuid = $1`,
	},

	"lifecycle": {
		// it's an ugly, bad precedent, except that it saves a lot of hits to the db
		"index": `
      select uuid,
             location,
             ctime
       from  lifecycles
      order
         by  ctime desc`,
		"select": `
      select  lc.location,
              lc.strain_cost,
              lc.grain_cost,
              lc.bulk_cost,
              lc.yield,
              lc.headcount,
              lc.gross,
              lc.mtime at time zone 'utc',
              lc.ctime at time zone 'utc',
              s.uuid as strain_uuid,
              s.species as strain_species,
              s.name as strain_name,
              s.ctime as strain_ctime,
              sv.uuid as strain_vendor_uuid,
              sv.name as strain_vendor_name,
              sv.website as strain_vendor_website,
              gs.uuid as grain_substrate_uuid,
              gs.name as grain_substrate_name,
              gs.type as grain_substrate_type,
              gv.uuid as grain_vendor_uuid,
              gv.name as grain_vendor_name,
              gv.website as grain_vendor_website,
              bs.uuid as bulk_substrate_uuid,
              bs.name as bulk_substrate_name,
              bs.type as bulk_substrate_type,
              bv.uuid as bulk_vendor_uuid,
              bv.name as bulk_vendor_name,
              bv.website as bulk_vendor_website
        from  lifecycles lc
        join  strains s
          on  lc.strain_uuid = s.uuid
        join  vendors sv
          on  s.vendor_uuid = sv.uuid
        join  substrates gs
          on  lc.grainsubstrate_uuid = gs.uuid
        join  vendors gv
          on  gs.vendor_uuid = gv.uuid
        join  substrates bs
          on  lc.bulksubstrate_uuid = bs.uuid 
        join  vendors bv
          on  bs.vendor_uuid = bv.uuid
       where  lc.uuid = $1`,
		"insert": `
      insert
        into lifecycles(
             uuid,
             location,
             strain_cost,
             grain_cost,
             bulk_cost,
             yield,
             headcount,
             gross,
             mtime,
             ctime,
             strain_uuid,
             grainsubstrate_uuid,
             bulksubstrate_uuid)
      select $1,
             $2,
             $3,
             $4,
             $5,
             $6,
             $7,
             $8,
             $9,
             $10,
             s.uuid,
             gs.uuid,
             bs.uuid
       from  strains s,
             substrates gs,
             substrates bs
      where  s.uuid = $11
        and  gs.uuid = $12
        and  gs.type = 'Grain'
        and  bs.uuid = $13
        and  bs.type = 'Bulk'`,
		"update": `
      update lifecycles
        set location = $1,
            strain_cost = $2,
            grain_cost = $3,
            bulk_cost = $4,
            yield = $5,
            headcount = $6,
            gross = $7,
            mtime = $8,
            strain_uuid = s.uuid,
            grainsubstrate_uuid = gs.uuid,
            bulksubstrate_uuid = bs.uuid
        from strains s,
            substrates gs,
            substrates bs
      where s.uuid = $9
        and gs.uuid = $10
        and gs.type = 'Grain'
        and bs.uuid = $11
        and bs.type = 'Bulk'
        and lifecycles.uuid = $12`,
		"modified": `update lifecycles set mtime = $1 where uuid = $2`,
		"delete":   `delete from lifecycles where uuid = $1`,
	},

	"stage": {
		"select-all": `select uuid, name from stages order by name`,
		"select":     `select name from stages where uuid = $1`,
		"insert":     `insert into stages(uuid, name) values($1, $2)`,
		"update":     `update stages set name = $1 where uuid = $2`,
		"delete":     `delete from stages where uuid = $1`,
	},

	"strain": {
		"select-all": `
      select  s.uuid,
              s.species,
              s.name,
              s.ctime,
              v.uuid as vendor_uuid,
              v.name as vendor_name,
              v.website as vendor_website
        from  strains s
        join  vendors v
          on  s.vendor_uuid = v.uuid
       order
          by  s.name`,
		"select": `
      select  s.species,
              s.name,
              s.ctime,
              v.uuid as vendor_uuid,
              v.name as vendor_name,
              v.website as vendor_website
        from  strains s
        join  vendors v
          on  s.vendor_uuid = v.uuid
       where  s.uuid = $1`,
		"insert": `
      insert
        into  strains(uuid, species, name, ctime, vendor_uuid)
      select  $1, $2, $3, $4, v.uuid
        from  vendors v
       where  v.uuid = $5`,
		"update": `
      update  strains 
         set  species = $1, 
              name = $2,
              vendor_uuid = $3
       where  uuid = $4`,
		"delete": `delete from strains where uuid = $1`,
	},

	"strainattribute": {
		"get-unique-names": `select distinct name from strain_attributes order by name`,
		"all": `
    select uuid, name, value
      from strain_attributes sa
     where strain_uuid = $1`,
		"add": `
    insert
      into strain_attributes (uuid, name, value, strain_uuid)
    select $1, $2, $3, s.uuid
      from strains s
     where s.uuid = $4`,
		"change": `
    update strain_attributes sa
       set value = $1
      from strains s
     where sa.name = $2
       and s.uuid = $3
       and sa.strain_uuid = s.uuid`,
		"remove": `delete from strain_attributes where uuid = $1`,
	},

	"substrate-ingredient": {
		"all": `
    select i.uuid,
           i.name
      from ingredients i
      join substrate_ingredients si
        on i.uuid = si.ingredient_uuid
     where si.substrate_uuid = $1`,
		"add": `
    insert
      into substrate_ingredients (uuid, substrate_uuid, ingredient_uuid)
    select $1, s.uuid, i.uuid
      from substrates s
      join ingredients i
        on s.uuid = $2
       and i.uuid = $3`,
		"change": `
    update substrate_ingredients
       set ingredient_uuid = $1
     where substrate_uuid = $2
       and ingredient_uuid = $3`,
		"remove": `
    delete
      from substrate_ingredients
     where substrate_uuid = $1
       and ingredient_uuid = $2`,
	},

	"substrate": {
		"select-all": `
      select  s.uuid,
              s.name,
              s.type,
              v.uuid as vendor_uuid,
              v.name as vendor_name,
              v.website as vendor_website
        from  substrates s
        join  vendors v
          on  s.vendor_uuid = v.uuid
       order
          by  s.name`,
		"select": `
      select  s.name,
              s.type,
              v.uuid as vendor_uuid,
              v.name as vendor_name,
              v.website as vendor_website
        from  substrates s
        join  vendors v
          on  s.vendor_uuid = v.uuid
       where  s.uuid = $1`,
		"insert": `
      insert
        into substrates(uuid, name, type, vendor_uuid)
      select $1, $2, $3, v.uuid
        from vendors v
      where v.uuid = $4`,
		"update": `
      update  substrates s
         set  name = $1,
              type = $2,
              vendor_uuid = v.uuid
        from  vendors v 
       where  v.uuid = $3
         and  s.uuid = $4`,
		"delete": `delete from substrates where uuid = $1`,
	},

	"vendor": {
		"select-all": `select uuid, name, website from vendors order by name`,
		"select":     `select name, website from vendors where uuid = $1`,
		"insert":     `insert into vendors(uuid, name, website) values($1, $2, $3)`,
		"update":     `update vendors set name = $1, website = $2 where uuid = $3`,
		"delete":     `delete from vendors where uuid = $1`,
	},
}
