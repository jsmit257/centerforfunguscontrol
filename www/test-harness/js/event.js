$(_ => {
  let newEventRow = (data = { event_type: { stage: {} }, mtime: "Now", ctime: "Now" }) => {
    return $('<div>')
      .addClass('row hover')
      .attr('id', data.id)
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
          eventtypes.push($(new Option())
            .val(r.id)
            .text(`${r.name + "/" + r.stage.name}`))
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

  $('body>.template>.notes')
    .clone(true, true)
    .insertAfter($('body>.template>.table.events>.rows'))

  $('body>.template>.photos')
    .clone(true, true)
    .insertAfter($('body>.template>.table.events>.rows'))

  $('body>.main>.workspace .table.events>.rows, body>.template>.table>.rows').off('send')
    .on('send', (e, ...data) => {
      let $table = $(e.currentTarget).data('events', data)

      $table.empty()
      data.forEach(d => {
        $table.append(newEventRow(d))
      })
      $table.find('.row').first().click()
      $table
        .data('config')
        .$eventbar
        .find('>.remove')[data.length !== 0 ? "addClass" : "removeClass"]('active')
    })
    .on('initialize', (e, cfg) => {
      let $events = $(e.currentTarget)
      let $eventbar = $events
        .parents('div.table.events')
        .first()
        .find('>div.buttonbar')
      let $owner = cfg.$owner

      $events.data('config', { $eventbar: $eventbar })

      $eventbar
        .on('click', '>.edit', e => {
          if (!$(e.currentTarget).hasClass('active')) {
            return
          }
          getEventtypes($events)
            .val($events.find('.selected>.eventtype.live').data('eventtype_uuid'))

          var $modifiedDate = $events.find('.selected>.mtime.static').text("Now")

          $events.trigger('edit', {
            url: `${cfg.parent}/${$owner.data('record').id}/events`,
            data: $selected => {
              return JSON.stringify({
                "id": $selected.attr('id'),
                "temperature": parseFloat($selected.find('>.temperature.live').val()) || 0,
                "humidity": parseFloat($selected.find('>.humidity.live').val()) || 0,
                "event_type": {
                  "id": $selected.find('>.eventtype.live').val()
                },
                "ctime": new Date($selected.find('>.ctime.static').data('value')).toISOString(),
              })
            },
            success: (data, status, xhr) => { $owner.trigger('send', data) },
            cancel: _ => { $modifiedDate.trigger('reset') },
            buttonbar: $eventbar
          })
        })
        .on('click', '>.add', e => {
          if (!$(e.currentTarget).hasClass('active')) {
            return
          }
          if ($events
            .trigger('add', {
              newRow: newEventRow,
              url: `${cfg.parent}/${$owner.data('record').id}/events`,
              data: $selected => {
                return JSON.stringify({
                  "temperature": parseFloat($selected.find('>.temperature.live').val()) || 0,
                  "humidity": parseFloat($selected.find('>.humidity.live').val().trim()) || 0,
                  "event_type": {
                    "id": $selected.find('>.eventtype.live').val()
                  },
                })
              },
              success: (data, status, xhr) => { $owner.trigger('send', data) },
              error: _ => { $events.trigger('remove-selected') },
              buttonbar: $eventbar
            })
            .children()
            .length === 0
          ) {
            cfg.$tablebar.find('.remove').addClass('active')
          }
          getEventtypes($events)
        })
        .on('click', '>.remove', e => {
          if (!$(e.currentTarget).hasClass('active')) {
            return
          }
          if ($events
            .trigger('delete', {
              url: `${cfg.parent}/${$owner.data('record').id}/events/${$events.find('.selected').attr('id')}`,
              buttonbar: $eventbar
            })
            .children()
            .length === 0
          ) {
            cfg.$tablebar.find('.remove').addClass('active')
          }
        })
        .trigger('subscribe', {
          clazz: 'notes',
          clicker: e => {
            e.stopPropagation()

            let $p = $(e.delegateTarget)
              .parents('.table.events')
              .first()

            if ($p.toggleClass('noting').hasClass('noting')
            ) {
              $p.find('div.notes').trigger('refresh', $events.find('.selected').attr('id'))
            }
          },
        })
        .trigger('subscribe', {
          clazz: 'photos',
          clicker: e => {
            e.stopPropagation()

            let $p = $(e.delegateTarget)
              .parents('.table.events')
              .first()

            if ($p.toggleClass('photoing').hasClass('photoing')
            ) {
              $p
                .find('div.photos')
                .addClass('gallery')
                .removeClass('singleton')
                .trigger('refresh', $events.find('.selected').attr('id'))
            }
          },
        })
        .find('.refresh').remove()
    })
})