$(function () {
  var $lifecycle = $('body>.main>.workspace>.lifecycle')
  var $ndx = $lifecycle.find('.table.ndx>.rows')
  var $table = $lifecycle.find('.table.lifecycle>.rows')
  var $tablebar = $lifecycle.find('.table.lifecycle>.buttonbar')
  var $events = $lifecycle.find('.table.events>.rows')
  var $eventbar = $lifecycle.find('.table.events>.buttonbar')

  var fields = ['id', 'location', 'strain_cost', 'grain_cost', 'bulk_cost', 'yield', 'count', 'gross']

  let newNdxRow = (data) => {
    return $('<div class="row hover" />')
      .append($('<div class=uuid />').text(data.id))
      .append($('<div class="created_at static date" />').text(data.ctime.replace('T', ' ').replace(/:\d{2}(\..+)?Z.*/, '')))
      .append($('<div class="location static" />').text(data.location))
  }

  let selectLists = _ => {
    $.ajax({
      url: '/strains',
      method: 'GET',
      async: true,
      success: (result, status, xhr) => {
        var strains = []
        result.forEach(r => {
          strains.push($(`<option value="${r.id}">${r.name} | Species: ${r.species || "unknown"} | Vendor: ${r.vendor.name}</option>`)
            .data('strain', r))
        })
        $table
          .find('>.row.strain>select.live')
          .empty()
          .append(strains)
      },
      error: console.log,
    })

    $.ajax({
      url: '/substrates',
      method: 'GET',
      async: true,
      success: (result, status, xhr) => {
        var substrates = { bulk: [], grain: [], liquid: [], agar: [] }
        result.forEach(r => {
          substrates[r.type.toLowerCase()]
            .push($(`<option value="${r.id}">${r.name} | Vendor: ${r.vendor.name}</option>`)
              .data('substrate', r))
        })
        $table
          .find('>.row.grain_substrate>select.live')
          .empty()
          .append(substrates.grain)
        $table
          .find('.row.bulk_substrate>select.live')
          .empty()
          .append(substrates.bulk)
      },
      error: console.log,
    })
  }

  $ndx.on('click', '>.row', e => {
    if (e.isPropagationStopped()) {
      return false
    }
    var $row = $(e.currentTarget)
    $.ajax({
      url: '/lifecycle/' + $row.find('div.uuid').text(),
      method: 'GET',
      async: true,
      success: (result, status, xhr) => { $table.trigger('send', result) },
      error: console.log,
    })
  })

  $table.off('click').off('refresh').off('edit').off('add')
    .on('send', (e, lc) => {

      $table
        .data('lifecycle', lc)
        .parent()
        .find('.columns>.column>span.title')
        .text(`${lc.strain.name} - ${new Date(lc.ctime).toLocaleString()}`)

      fields.forEach(n => {
        $table.find(`.row.${n}>.static`).text(lc[n] || 0)
        $table.find(`.row.${n}>.live`).val(lc[n] || 0)
      })
      $table.find(`.row.mtime>.static`).trigger('set', lc.mtime)
      $table.find(`.row.ctime>.static`).trigger('set', lc.ctime)

      $table
        .find(`.row.strain>.static`)
        .text(`${lc.strain.name} | Species: ${lc.strain.species || "unknown"} | Vendor: ${lc.strain.vendor.name}`)
      $table.find(`.row.strain>select`).val(lc.strain.id)
      $table
        .find(`.row.grain_substrate>.static`)
        .text(`${lc.grain_substrate.name} | Vendor: ${lc.grain_substrate.vendor.name}`)
      $table.find(`.row.grain_substrate>select`).val(lc.grain_substrate.id)
      $table
        .find(`.row.bulk_substrate>.static`)
        .text(`${lc.bulk_substrate.name} | Vendor: ${lc.bulk_substrate.vendor.name}`)
      $table.find(`.row.bulk_substrate>select`).val(lc.bulk_substrate.id)

      $events.trigger('send', lc.events)
    })
    .on('edit', (e, args) => {
      selectLists()

      $table
        .trigger('set-editing', 'edit')
        .find('input, select')
        .first()
        .focus()

      var $modifiedDate = $table.find('.mtime>.static').text("Now")

      args.buttonbar.trigger('set', {
        target: $table,
        handlers: {
          cancel: (xhr, status, error) => {
            $table
              .removeClass('editing')
              .trigger('set-editing')
            $modifiedDate.trigger('reset')
          },
          ok: args.ok || (e => {
            $.ajax({
              url: args.url || `/lifecycle/${$table.find('.row.id>.uuid').text()}`,
              contentType: 'application/json',
              method: 'PATCH',
              dataType: 'json',
              data: args.data(),
              async: true,
              success: result => {
                args.success(result)
                $ndx.find('.selected>.location').text(result.location)
                $table.trigger('set-editing')
              },
              error: args.error || console.log,
            })
          })
        }
      })
    })
    .on('add', (e, args) => {
      selectLists()

      $table
        .trigger('set-editing', 'add')
        .find('input, select')
        .val("")
        .first()
        .focus()

      $table.find('.mtime>.static, .ctime>.static').text("Now")

      $events.empty()

      args.buttonbar.trigger('set', {
        target: $table,
        handlers: {
          cancel: (xhr, status, error) => {
            $table.trigger('send', $table.data('lifecycle'))
            $table.trigger('set-editing')
          },
          ok: args.ok || (e => {
            $.ajax({
              url: args.url || `/lifecycle`,
              contentType: 'application/json',
              method: 'POST',
              dataType: 'json',
              data: args.data(),
              async: false,
              success: result => {
                args.success(result)
                $table.trigger('set-editing')
                var $ndxRow = newNdxRow(result)
                  .trigger('click')
                  .addClass('selected')
                $ndx
                  .find('.selected')
                  .removeClass('selected')
                var $children = $ndx.children()
                if ($children.length === 0) {
                  $ndx.append($ndxRow)
                } else {
                  $ndxRow.insertBefore($children.first())
                }
              },
              error: args.error || console.log,
            })
          })
        }
      })
    })
    .on('set-editing', (e, status) => {
      $table[(!status) ? "removeClass" : "addClass"](`editing ${status || "add edit"}`)
        .parent()
        .find('.columns>.column>span.status')
        .text(status ? `(${status})` : '')
    })

  let newEventRow = (data = { event_type: { stage: {} }, mtime: "Now", ctime: "Now" }) => {
    return $('<div class="row hover" />')
      .append($('<div class=uuid />').text(data.id))
      .append($('<div class="eventtype static" />').text(data.event_type.name + "/" + data.event_type.stage.name))
      .append($('<select class="eventtype live" />')
        .data('eventtype_uuid', data.event_type.id)
        .val(data.event_type.id))
      .append($('<div class="temperature static" />').text(data.temperature))
      .append($('<input class="temperature live" min="0" type="number">').val(data.temperature))
      .append($('<div class="humidity static" />').text(data.humidity))
      .append($('<input class="humidity live" type="number" min="0" max="100">').val(data.humidity))
      .append($('<div class="mtime static const date" />')
        .data('value', data.mtime)
        .text(data.mtime.replace('T', ' ').replace(/:\d{1,2}(\..+)?Z.*/, '')))
      .append($('<div class="ctime static const date" />')
        .data('value', data.ctime)
        .text(data.ctime.replace('T', ' ').replace(/:\d{1,2}(\..+)?Z.*/, '')))
  }

  let getEventtypes = $events => {
    let $result

    $.ajax({
      url: '/eventtypes',
      method: 'GET',
      async: false,
      success: (result, status, xhr) => {
        let eventtypes = []
        result.forEach(r => {
          eventtypes.push($(`<option value="${r.id}">${r.name + "/" + r.stage.name}</option>`))
        })
        $result = $events
          .find('>.row.selected>select.eventtype.live')
          .empty()
          .append(eventtypes)
      },
      error: console.log,
    })

    return $result
  }

  $events.on('send', (e, ...ev) => {
    ev ||= []

    let $events = $(e.currentTarget).data('events', ev)

    $events.empty()
    ev.forEach(evt => { $events.append(newEventRow(evt)) })
    $events.find('.row').first().click()

    $tablebar.find('.remove')[(ev.length !== 0 ? "removeClass" : "addClass")]('active')
    $eventbar.find('.remove, .edit')[(ev.length !== 0 ? "addClass" : "removeClass")]('active')
  })

  $tablebar.find('.edit').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }

    $table.trigger('edit', {
      data: _ => {
        return JSON.stringify({
          "id": $table.find('.row>.uuid.static').text(),
          "location": $table.find('.row.location>.live').val(),
          "strain_cost": parseFloat($table.find('.row.strain_cost>.live').val()) || 0,
          "grain_cost": parseFloat($table.find('.row.grain_cost>.live').val()) || 0,
          "bulk_cost": parseFloat($table.find('.row.bulk_cost>.live').val()) || 0,
          "yield": parseFloat($table.find('.row.yield>.live').val()) || 0,
          "count": parseFloat($table.find('.row.count>.live').val()) || 0,
          "gross": parseFloat($table.find('.row.gross>.live').val()) || 0,
          "strain": {
            "id": $table.find('.row.strain>.live').val(),
          },
          "grain_substrate": {
            "id": $table.find('.row.grain_substrate>.live').val(),
          },
          "bulk_substrate": {
            "id": $table.find('.row.bulk_substrate>.live').val(),
          },
          "ctime": new Date($table.find('.row.ctime>.static').text()).toISOString()
        })
      },
      success: (result, status, xhr) => {
        $table
          .trigger('send', {
            ...result,
            ...{
              strain: $table.find('.row.strain>.live>option:selected').data('strain'),
              grain_substrate: $table.find('.row.grain_substrate>.live>option:selected').data('substrate'),
              bulk_substrate: $table.find('.row.bulk_substrate>.live>option:selected').data('substrate'),
              events: $events.data('events')
            }
          })
      },
      error: _ => { $table.removeClass('editing') },
      buttonbar: $tablebar
    })
  })

  $tablebar.find('.add').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }
    $table.trigger('add', {
      data: _ => {
        return JSON.stringify({
          "location": $table.find('.row.location>.live').val(),
          "strain_cost": parseFloat($table.find('.row.strain_cost>.live').val()) || 0,
          "grain_cost": parseFloat($table.find('.row.grain_cost>.live').val()) || 0,
          "bulk_cost": parseFloat($table.find('.row.bulk_cost>.live').val()) || 0,
          "yield": parseFloat($table.find('.row.yield>.live').val()) || 0,
          "count": parseFloat($table.find('.row.count>.live').val()) || 0,
          "gross": parseFloat($table.find('.row.gross>.live').val()) || 0,
          "strain": {
            "id": $table.find('.row.strain>.live').val(),
          },
          "grain_substrate": {
            "id": $table.find('.row.grain_substrate>.live').val(),
          },
          "bulk_substrate": {
            "id": $table.find('.row.bulk_substrate>.live').val(),
          },
        })
      },
      success: (result, status, xhr) => {
        $table
          .trigger('send', {
            ...result,
            ...{
              strain: $table.find('.row.strain>.live>option:selected').data('strain'),
              grain_substrate: $table.find('.row.grain_substrate>.live>option:selected').data('substrate'),
              bulk_substrate: $table.find('.row.bulk_substrate>.live>option:selected').data('substrate'),
            }
          })
      },
      error: _ => { $table.removeClass('editing') },
      buttonbar: $tablebar
    })
  })

  $tablebar.find('.remove').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }
    $ndx.trigger('delete', {
      url: `/lifecycle/${$table.find('.row>.uuid.static').text()}`,
      buttonbar: $tablebar
    })
  })

  $tablebar.find('.refresh').on('click', e => {
    $ndx.find('.selected').removeClass('selected').click()
  })

  $eventbar.find('>.edit').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }

    getEventtypes($events)
      .val($events.find('.selected>.eventtype.live').data('eventtype_uuid'))

    var $modifiedDate = $events.find('.selected>.mtime.static').text("Now")

    $events.trigger('edit', {
      url: `/lifecycle/${$table.find('.row.id>.uuid').text()}/events`,
      data: $selected => {
        return JSON.stringify({
          "id": $selected.find('>.uuid').text(),
          "temperature": parseFloat($selected.find('>.temperature.live').val()) || 0,
          "humidity": parseFloat($selected.find('>.humidity.live').val().trim()) || 0,
          "event_type": {
            "id": $selected.find('>.eventtype.live').val()
          },
          "ctime": new Date($selected.find('>.ctime.static').data('value')).toISOString(),
        })
      },
      success: (data, status, xhr) => { $table.trigger('send', data) },
      cancel: _ => {
        $modifiedDate.trigger('reset')
      },
      buttonbar: $eventbar
    })
  })

  $eventbar.find('>.add').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }

    $events.trigger('add', {
      newRow: newEventRow,
      url: `/lifecycle/${$table.find('.row.id>.uuid').text()}/events`,
      data: $selected => {
        return JSON.stringify({
          "temperature": parseFloat($selected.find('>.temperature.live').val()) || 0,
          "humidity": parseFloat($selected.find('>.humidity.live').val().trim()) || 0,
          "event_type": {
            "id": $selected.find('>.eventtype.live').val()
          },
        })
      },
      success: (data, status, xhr) => {
        $table.trigger('send', data)
        $eventbar.find('.remove, .edit').removeClass('active')
      },
      error: (xhr, status, error) => { $events.trigger('remove-selected') },
      buttonbar: $eventbar
    })

    getEventtypes($events)
  })

  $eventbar.find('.remove').on('click', e => {
    if ($eventbar.find('.remove.active').length === 0) {
      return
    }
    $events.trigger('delete', {
      url: `/lifecycle/${$table.find('.row>.uuid.static').text()}/events/${$events.find('.selected>.uuid').text()}`,
      buttonbar: $eventbar
    })
    if ($events.children().length === 0) {
      $tablebar.find('.remove').addClass('active')
    }
  })

  $eventbar.find('.refresh').remove()

  // XXX: consider changing this to put delete/refresh on the ndx and remove them 
  //      from (or leave them in) the main table; and/or, add filters for 
  //      sunset/unmodified, etc
  $lifecycle.find('.table.ndx>.buttonbar')
    .empty()
    .append($('<div />')
      .html('&nbsp'))

  $lifecycle.on('activate', e => {
    $lifecycle
      .addClass('active')
      .find('>.ndx>.rows')
      .trigger('refresh', { newRow: newNdxRow, buttonbar: $tablebar })
  })
})