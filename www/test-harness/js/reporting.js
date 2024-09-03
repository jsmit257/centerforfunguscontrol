$(_ => {
  let labelmap = {
    'bulk_cost': 'Bulk cost',
    'bulk_substrate': 'Bulk',
    'count': 'Count/kg',
    'ctime': 'Created',
    'dtime': 'Deleted',
    'generation': 'Parent(s)',
    'grain_cost': 'Grain cost',
    'grain_substrate': 'Grain',
    'id': 'ID',
    'liquid_substrate': 'Liquid',
    'mtime': 'Modified',
    'plating_substrate': 'Plating',
    'strain_cost': 'Strain cost',
    'yield': 'yield (g)',
  }

  let pluralmap = {
    'attributes': 'attribute',
    'events': 'event',
    'generations': 'generation',
    'ingredients': 'ingredient',
    'lifecycles': 'lifecycle',
    'notes': 'note',
    'photos': 'photo',
    'sources': 'source',
  }

  let format = (p => {
    return (d) => {
      if (!d.replace) {
        return d
      }
      return d.replace(p, '$1 $2')
    }
  })(/^(\d{4}.\d\d.\d\d).(\d\d.\d\d.\d\d).*/)

  let cloneentity = ($tmpl => {
    $tmpl.find('>.ndx').remove()

    return name => $tmpl
      .clone(true, true)
      .attr('name', name)
      .attr('sort-key', name)
  })($('.workspace>.history>.entity').clone(true, true))

  let newentity = ((entityname, data, $parent) => {
    // return $parent.trigger('add-child', entityname)
    //   .find('>.list>div')
    //   .last()
    //   .trigger('send', data)
    //   .find('>.entity-name')
    //   .trigger('map', entityname)
    //   .parent()
    return cloneentity(entityname)
      .appendTo($parent)
      .trigger('send', data)
      .find('>.entity-name')
      .trigger('map', entityname)
      .parent()
  })

  let parsedata = (k, data, $parent) => {
    switch (data[k].constructor.prototype) {
      case Object.prototype:
        newentity(k, data[k], $parent.find('>.list'))
        break

      case Array.prototype:
        let $list = newentity(k, [], $parent.find('>.list'))
          .removeClass('collapsed')
          .find('>.list')

        $list.prev().text(`(${data[k].length})`)

        data[k].forEach(v => {
          newentity(pluralmap[k] || k, v, $list)
        })

        break

      default:
        let $row = $('<div>')
          .addClass(`scalar`)
          .attr('sort-key', k)
          .append($('<div>')
            .addClass('label'))
          .append($('<div>')
            .addClass('value')
            .html(format(data[k])))
          .appendTo($parent.find('>.list'))
          .find('>.label')
          .trigger('map', k)
          .parent()
    }
  }

  $('.main>.workspace>.history')
    .on('activate', (e, id) => {
      let entityname = $('.main>.header>.menuitem.selecting').attr('entity-name')

      $(e.currentTarget)
        .attr('name', entityname)
        .find('>.entity')
        .attr('name', entityname)
        .addClass('collapsed')
        .trigger('reinit')
        .find('>.ndx')
        .trigger('refresh', id)
        .parent()
        .find('>.entity-name')
        .trigger('map', entityname)
    })
    .on('refresh', '>.entity>.ndx', (e, id) => {
      e.stopPropagation()

      $.ajax({
        url: `/${$(e.delegateTarget).attr('name')}s`,
        method: 'GET',
        async: true,
        success: (data, status, xhr) => {
          $(e.currentTarget)
            .empty()
            .trigger('send', data)
            .find(`>.row#${id}`)
            .click()
        },
        error: console.log,
      })
    })
    .on('send', '>.entity[nme="eventtype"]>.ndx', (e, ...data) => {
      e.stopPropagation()

      data.forEach(ev => {
        $('<div>')
          .addClass('row hover')
          .attr('id', ev.id)
          .append($('<div>')
            .addClass('event-stage')
            .html(`${ev.name} &larr; ${ev.stage.name}`))
          .appendTo(e.currentTarget)
      })
    })
    .on('send', '>.entity[name="generation"]>.ndx', (e, ...data) => {
      e.stopPropagation()

      data.forEach(gen => {
        let strains = []
        gen.sources.forEach(v => {
          strains.push(v.strain.name)
        })

        $('<div>')
          .addClass('row hover')
          .attr('id', gen.id)
          .append($('<div>')
            .addClass('mtime static date'))
          .append($('<div>')
            .addClass('que')
            .html(strains.join(' & ')))
          .appendTo(e.currentTarget)
          .find('>.static.date')
          .trigger('set', gen.mtime)
      })
    })
    .on('send', '>.entity[name="lifecycle"]>.ndx', (e, ...data) => {
      e.stopPropagation()

      data.forEach(lc => {
        $('<div>')
          .addClass('row hover')
          .attr('id', lc.id)
          .append($('<div>')
            .addClass('mtime static date'))
          .append($('<div>')
            .addClass('location')
            .text(lc.location))
          .appendTo(e.currentTarget)
          .find('>.static.date')
          .trigger('set', lc.mtime)
      })
    })
    .on('send', '>.entity[name="strain"]>.ndx', (e, ...data) => {
      e.stopPropagation()

      data.forEach(strain => {
        $('<div>')
          .addClass('row hover')
          .attr('id', strain.id)
          .append($('<div>')
            .html(`${strain.name} &bull; <i>${strain.species}</i>`))
          .appendTo(e.currentTarget)
      })
    })
    .on('send', '>.entity[name="substrate"]>.ndx', (e, ...data) => {
      e.stopPropagation()

      data.forEach(sub => {
        $('<div>')
          .addClass('row hover')
          .attr({
            'id': sub.id,
            'param': `${sub.type}-id=${sub.id}`,
            'owner': ['grain', 'bulk'].indexOf(sub.type) !== -1 ? 'lifecycle' : 'generation',
          })
          .append($('<div>')
            .html(`${sub.name} &bull; <i>${sub.vendor.name}</i>`))
          .appendTo(e.currentTarget)
      })
    })
    .on('send', '>.entity[name="vendor"]>.ndx', (e, ...data) => {
      e.stopPropagation()

      data.forEach(sub => {
        $('<div>')
          .addClass('row hover')
          .attr({
            'id': sub.id,
          })
          .append($('<div>')
            .text(sub.name))
          .appendTo(e.currentTarget)
      })
    })
    .on('click', '>.entity>.ndx>.row', (e, parent) => {
      e.stopPropagation()

      let entityname = $(e.delegateTarget).attr('name')
      let $entity = $(e.currentTarget)
        .parents('.entity')
        .first()

      $entity.find('>.ndx>.selected').removeClass('selected')

      let $row = $(e.currentTarget).addClass('selected')

      $.ajax({
        url: $row.data('url') || `/${entityname}/${$row.attr('id')}`,
        method: 'GET',
        async: true,
        success: data => {
          $entity
            .removeClass('collapsed')
            .find('>.list')
            .empty()
            .parent()
            .trigger('send', data)
            .trigger('sort')
        },
        error: console.log,
      })
    })
    .on('send', '.entity[name="attribute"]', (e, data) => {
      e.stopPropagation()

      $(e.currentTarget).trigger('cliff-notes', [data.name, data.value])
    })
    .on('send', '.entity[name="event"]', (e, data) => {
      e.stopPropagation()

      $(e.currentTarget)
        .trigger('get-child', { id: data.id, entityname: 'notes' })
        .trigger('get-child', { id: data.id, entityname: 'photos' })
        .trigger('cliff-notes', [format(data.mtime), data.event_type.name])
    })
    .on('send', '.entity[name="event_type"]', (e, data) => {
      e.stopPropagation()

      $(e.currentTarget).find('>.cliff-notes').text(data.name)
    })
    .on('send', '.entity[name="generation"]', (e, data) => {
      e.stopPropagation()

      let strains = []
      data.sources.forEach(v => {
        strains.push(v.strain.name)
      })

      $(e.currentTarget)
        .trigger('get-child', { id: data.id, entityname: 'notes' })
        .trigger('cliff-notes', [format(data.mtime), strains.join(' & ')])
    })
    .on('send', '.entity[name="ingredient"]', (e, data) => {
      e.stopPropagation()

      $(e.currentTarget).trigger('cliff-notes', data.name)
    })
    .on('send', '.entity[name="lifecycle"]', (e, data) => {
      e.stopPropagation()

      let totalcost = ((data.strain.strain_cost = (data.strain_cost || 0))
        + (data.grain_substrate.grain_cost = (data.grain_cost || 0))
        + (data.bulk_substrate.bulk_cost = (data.bulk_cost || 0)))

      delete data.strain_cost
      delete data.grain_cost
      delete data.bulk_cost

      $(e.currentTarget)
        .trigger('get-child', { id: data.id, entityname: 'notes' })
        .trigger('cliff-notes', [data.strain.name, format(data.mtime), `$${totalcost || '~'}`])
    })
    .on('send', '.entity[name="note"]', (e, data) => {
      e.stopPropagation()

      $(e.currentTarget).trigger('cliff-notes', [format(data.mtime), `${data.note.slice(0, 25)}...`])
    })
    .on('send', '.entity[name="photo"]', (e, data) => {
      e.stopPropagation()

      $(e.currentTarget)
        .trigger('get-child', { id: data.id, entityname: 'notes' })
        .trigger('cliff-notes', format(data.mtime))

      data.image = `<a href=/album/${data.image} target=_lobby>${data.image}</a>`
    })
    .on('send', '.entity[name="source"]', (e, data) => {
      e.stopPropagation()

      $(e.currentTarget)
        .find('>.cliff-notes')
        .html(`${data.type} &rArr; ${data.strain.name}`)
    })
    .on('send', '.entity[name="strain"]', (e, data) => {
      e.stopPropagation()

      let $s = $(e.currentTarget)

      $s
        .trigger('get-child', { id: data.id, entityname: 'photos' })
        .trigger('cliff-notes', [
          data.name,
          data.species,
          data.vendor.name,
          `$${data.strain_cost || '~'}`,
        ])

      if (data.generation) {
        $s.trigger('get-child', { id: data.generation.id, entityname: 'generation' })
        delete data.generation
      }

      if ($s.parents('.entity[name="lifecycle"], .entity[name="generation"]').length === 0) {
        $s
          .trigger('get-child', {
            entityname: "lifecycles",
            url: `/reports/lifecycles?strain-id=${data.id}`,
          })
        $s
          .trigger('get-child', {
            entityname: "generations",
            url: `/reports/generations?strain-id=${data.id}`,
          })
      }

      // delete data.strain_cost
    })
    .on('send', '.entity[name="vendor"]', (e, data) => {
      e.stopPropagation()

      data.website = `<a href=${data.website} target=_macondo>${data.website}</a>`

      $(e.currentTarget).trigger('cliff-notes', data.name)
    })
    .on('send', ['.entity[name="plating', 'liquid', 'grain', 'bulk_substrate"]'].join('_substrate"], .entity[name="'), (e, data) => {
      e.stopPropagation()

      $(e.currentTarget)
        .trigger('cliff-notes', [data.name, data.vendor.name, `$${data.grain_cost || data.bulk_cost || '~'}`])

      // delete data.grain_cost
      // delete data.strain_cost
    })
    .on('send', '.entity[name="substrate"]', (e, data) => {
      let $s = $(e.currentTarget)
      let $selected = $s
        .parents('.history')
        .find('>.entity>.ndx>.row.selected')
      let owner = $selected.attr('owner')

      if ($s
        .trigger('cliff-notes', [data.name, data.vendor.name])
        .parents(`.entity[name = "${owner}"]`).length === 0
      ) {
        $s
          .trigger('get-child', {
            entityname: `${owner}s`,
            url: `/reports/${owner}s?${$selected.attr('param')}`,
          })
      }
    })
    .on('send', '.entity', (e, data) => {
      e.stopPropagation()

      $(e.currentTarget).attr('id', (data || {}).id)

      for (let el in data) {
        if (Object.prototype.hasOwnProperty.call(data, el)) {
          parsedata(el, data, $(e.currentTarget))
        }
      }

      $(e.currentTarget).trigger('sort')
    })
    .on('sort', '.entity', e => {
      e.stopPropagation()

      let ndx = sortindices[$(e.currentTarget).attr('name')]
      if (!ndx) {
        return
      }

      let $list = $(e.currentTarget).find('>.list');
      $list
        .append(...$list
          .children()
          .sort((a, b) => ndx.indexOf(a.getAttribute('sort-key')) - ndx.indexOf(b.getAttribute('sort-key'))))
    })
    .on('cliff-notes', '.entity', (e, ...data) => {
      e.stopPropagation()

      let $cliff = $(e.currentTarget)
        .find('>.cliff-notes')
        .empty()

      data.forEach(v => $cliff.append($('<span>').text(v)))
    })
    .on('reinit', '>.entity', e => {
      $(e.currentTarget)
        .find('>.list')
        .empty()
        .parent()
        .find('>.cliff-notes')
        .html('(choose ye)')
    })
    .on('add-child', '.entity', (e, name) => {
      $(e.currentTarget)
        .find('>.list')
        .append(cloneentity(name))
    })
    .on('get-child', '.entity', (e, data) => {
      e.stopPropagation()

      $.ajax({
        url: data.url || `/${data.entityname}/${data.id}`,
        method: 'GET',
        async: true,
        success: (result, status, xhr) => {
          if (!result) {
            return
          }

          let entity = {}
          entity[`${data.entityname}`] = result  // inlining doesn't work: { `${ data.entityname }` : result }

          parsedata(data.entityname, entity, $(e.currentTarget))

          $(e.currentTarget)
            .trigger('sort')
            .parents('.entity')
            .first()
            .trigger('sort')
        },
        error: console.log,
      })
    })
    .on('click', '.entity-name, .cliff-notes', (e, data) => {
      e.stopPropagation()

      let $e = $(e.currentTarget)
        .parents('.entity')
        .first()
        .toggleClass('collapsed')

      if (e.shiftKey) {
        $e.find('.entity')[$e.hasClass('collapsed') ? 'addClass' : 'removeClass']('collapsed')
      }
    })
    .on('map', '.label, .entity-name', (e, key) => {
      $(e.currentTarget).html(labelmap[key] || key)
    })

  let sortindices = (_ => {
    let result = {
      attribute: ['id', 'name', 'value'],
      event: ['id', 'temperature', 'humidity', 'event_type', 'notes', 'mtime', 'ctime'],
      event_type: ['id', 'name', 'severity', 'stage'],
      generation: ['id', 'sources', 'plating_substrate', 'liquid_substrate', 'events', 'notes', 'mtime', 'ctime'],
      ingredient: ['id', 'name'],
      lifecycle: ['id', 'location', 'yield', 'count', 'strain', 'grain_substrate', 'bulk_substrate', 'events', 'notes', 'mtime', 'ctime'],
      note: ['id', 'note', 'mtime', 'ctime'],
      photo: ['id', 'image', 'notes', 'mtime', 'ctime'],
      // source: [],
      stage: ['id', 'name'],
      strain: ['id', 'name', 'species', 'strain_cost', 'generation', 'attributes', 'photos', 'vendor', 'lifecycles', 'generations', 'ctime', 'dtime'],
      vendor: ['id', 'name', 'website'],
    }

    result.grain_substrate
      = result.bulk_substrate
      = result.plating_substrate
      = result.liquid_substrate
      = ['id', 'name', 'type', 'grain_cost', 'bulk_cost', 'ingredients', 'vendor']

    return result
  })()
})