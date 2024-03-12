$(function () {
  var $lifecycle = $('body>.main>.workspace>.lifecycle')
  var $ndx = $lifecycle.find('.table.ndx>.rows')
  var $table = $lifecycle.find('.table.lifecycle>.rows')
  var $tablebar = $lifecycle.find('.table.lifecycle>.buttonbar')
  var $events = $lifecycle.find('.table.events>.rows')
  var $eventbar = $lifecycle.find('.table.events>.buttonbar')

  var fields = ['id', 'location', 'strain_cost', 'grain_cost', 'bulk_cost', 'yield', 'count', 'gross', 'modified_date', 'create_date']

  function newNdxRow(data) {
    data ||= {}
    return $('<div class="row hover" />')
      .append($('<div class=uuid />').text(data.id))
      .append($('<div class="created_at static" />').text(data.create_date.replace('T', ' ').replace(/(\..+)?Z.*/, '')))
      .append($('<div class="location static" />').text(data.location))
  }

  function newEventRow(data) {
    data ||= {}
    return $('<div class="row hover" />')
      .append($('<div class=uuid />').text(data.id))
      .append($('<div class="temperature static" />').text(data.temperature))
      .append($('<input class="temperature live" />').val(data.temperature))
      .append($('<div class="humidity static" />').text(data.humidity))
      .append($('<input class="humidity live" />').val(data.humidity))
      .append($('<div class="eventtype static" />').text(data.event_type.id))
      .append($('<select class="eventtype live" />').val(data.event_type.id))
      .append($('<div class="modified_date static" />').text(data.modified_date.replace('T', ' ').replace(/(\..+)?Z.*/, '')))
      .append($('<div class="create_date static" />').text(data.create_date.replace('T', ' ').replace(/(\..+)?Z.*/, '')))
  }

  $ndx
    .on('click', '>.row', e => {
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

  $table
    .off('click')
    .off('refresh')
    .on('send', (e, lc) => {

      $.ajax({
        url: '/strains',
        method: 'GET',
        async: false,
        success: (result, status, xhr) => {
          var strains = []
          result.forEach(r => {
            strains.push($(`<option value="${r.id}">${r.name} | ${r.vendor.name}</option>`))
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
        async: false,
        success: (result, status, xhr) => {
          var substrates = { bulk: [], grain: [] }
          result.forEach(r => {
            substrates[r.type.toLowerCase()]
              .push($(`<option value="${r.id}">${r.name} | ${r.vendor.name}</option>`))
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

      fields.forEach(n => {
        $table.find(`.row.${n}>.static`).text(lc[n])
        $table.find(`.row.${n}>.live`).val(lc[n])
      })
      $table
        .find(`.row.strain>.static`)
        .text(`${lc.strain.name} | Species: ${lc.strain.species || "unknown"} | Vendor: ${lc.strain.vendor.name}`)
        .data('strain_uuid', lc.strain.id)
      $table.find(`.row>.strain>select`).val(lc.strain.id)
      $table
        .find(`.row.grain_substrate>.static`)
        .text(`${lc.grain_substrate.name} | Vendor: ${lc.grain_substrate.vendor.name}`)
        .data('grainsubstrate_uuid', lc.grain_substrate.id)
      $table.find(`.row>.grain_substrate>select`).val(lc.grain_substrate.id)
      $table
        .find(`.row.bulk_substrate>.static`)
        .text(`${lc.bulk_substrate.name} | Vendor: ${lc.bulk_substrate.vendor.name}`)
        .data('bulksubstrate_uuid', lc.bulk_substrate.id)
      $table.find(`.row>.bulk_substrate>select`).val(lc.bulk_substrate.id)

      $events.trigger('send', lc.events)
    })

  $events.on('send', (e, ...ev) => {
    ev ||= []

    $.ajax({
      url: '/eventtypes',
      method: 'GET',
      async: false,
      success: (result, status, xhr) => {
        var eventtypes = []
        result.forEach(r => {
          eventtypes.push($(`<option value="${r.id}">${r.name}</option>`))
        })
        $events
          .find('>.row.eventtypes>select.live')
          .empty()
          .append(eventtypes)
      },
      error: console.log,
    })

    $events.empty()
    ev.forEach(evt => { $events.append(newEventRow(evt)) })
    $events.find('.row').first().click()
    if (ev.length === 0) {
      $tablebar.find('.remove').addClass('active')
      $eventbar.find('.remove, .edit').removeClass('active')
    } else {
      $tablebar.find('.remove').removeClass('active')
      $eventbar.find('.remove, .edit').addClass('active')
    }
  })

  $tablebar.find('.remove').on('click', e => {
    if ($tablebar.find('.remove.active').length === 0) {
      return
    }
    $table.trigger('delete', `/lifecycle/${$table.find('.row>.uuid.static').text()}`)
    $ndx.trigger('remove-selected')
  })

  $tablebar.find('.refresh').on('click', e => {
    $ndx.find('.selected').removeClass('selected').click()
  })

  $eventbar.find('.remove').on('click', e => {
    if ($eventbar.find('.remove.active').length === 0) {
      return
    }
    $events.trigger('delete', `/lifecycle/${$table.find('.row>.uuid.static').text()}/events/${$events.find('.selected>.uuid').text()}`)
    if ($events.find('.row').length === 1) { // is this shady?
      $tablebar.find('.remove').addClass('active')
      $eventbar.find('.remove, .edit').removeClass('active')
    }
  })

  $eventbar.find('.refresh').remove()

  $lifecycle.on('activate', e => {
    $lifecycle
      .addClass('active')
      .find('>.ndx>.rows')
      .trigger('refresh', { newRow: newNdxRow })
  })
})