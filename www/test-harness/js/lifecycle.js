$(_ => {
  let $lifecycle = $('body>.main>.workspace>.lifecycle')
  let $ndx = $lifecycle.find('.table.lifecycle>.rows.ndx')
  let $table = $lifecycle.find('.table.lifecycle>.rows.lc')
  let $tablebar = $lifecycle.find('.table.lifecycle>.buttonbar')

  let fields = ['id', 'location', 'strain_cost', 'grain_cost', 'bulk_cost', 'yield', 'count', 'gross']

  let newNdxRow = (data) => {
    return $('<div>')
      .addClass('row hover')
      .attr('id', data.id)
      .append($('<div class="mtime static date" />')
        .data('value', data.mtime)
        .text(data.mtime.replace('T', ' ').replace(/:\d{2}(\..+)?Z.*/, '')))
      .append($('<div class="location static" />').text(data.location))
  }

  let $gennotes = $('body>.main>.workspace>.template>.notes')
    .clone(true, true)
    .insertAfter($table)

  $tablebar.trigger('subscribe', {
    clazz: 'notes',
    clicker: e => {
      e.stopPropagation()
      if ($lifecycle.find('>.table.lifecycle').toggleClass('noting').hasClass('noting')) {
        $gennotes.trigger('refresh', ($table.data('record') || {}).id)
      }
    },
  })

  $ndx.on('click', '>.row', e => {
    if (e.isPropagationStopped()) {
      return false
    }

    $.ajax({
      url: `/lifecycle/${$(e.currentTarget).attr('id')}`,
      method: 'GET',
      async: true,
      success: (result, status, xhr) => {
        $table
          .trigger('send', result)
          .parents('.table.lifecycle')
          .first()
          .removeClass('noting')

        $events
          .parents('.table.events')
          .first()
          .removeClass('noting')
      },
      error: console.log,
    })
  })

  let $events = $('.main>.workspace>.template>.table.events')
    .clone(true, true)
    .appendTo($lifecycle)
    .find('>div.rows')
    .on('send', (e, ...ev) => {
      $tablebar.find('.remove')[((ev ||= []).length !== 0 ? "removeClass" : "addClass")]('active')
    })
    .trigger('initialize', {
      parent: 'lifecycle',
      $owner: $table,
      $tablebar: $tablebar,
    })

  $table.off('click').off('refresh').off('edit').off('add').off('send')
    .on('send', (e, lc) => {
      $table
        .data('record', lc)
        .parent()
        .find('.columns>.column>span.title')
        .text(`${lc.strain.name} - ${new Date(lc.ctime).toLocaleString()}`)

      fields.forEach(n => {
        $table.find(`.row.${n}>.static`).text(lc[n] || 0)
        $table.find(`.row.${n}>.live`).val(lc[n] || 0)
      })
      $table.find(`.row.mtime>.static`).trigger('set', lc.mtime)
      $table.find(`.row.ctime>.static`).trigger('set', lc.ctime)

      $.each({
        "strain": `${lc.strain.name} | Species: ${lc.strain.species || "unknown"} | Vendor: ${lc.strain.vendor.name}`,
        "grain_substrate": `${lc.grain_substrate.name} | Vendor: ${lc.grain_substrate.vendor.name}`,
        "bulk_substrate": `${lc.bulk_substrate.name} | Vendor: ${lc.bulk_substrate.vendor.name}`,
      }, (k, v) => {
        $table.find(`.row.${k}>.static`).text(v)
        $table.find(`.row.${k}>select`).data('fkey', lc[k].id)
      })

      $events.trigger('send', lc.events)
    })
    .on('edit', (e, args) => {
      new Array("strain", "bulk_substrate", "grain_substrate").forEach(v => {
        ($t => {
          $t.val($t.data('fkey'))
        })($table.find(`.row.${v}>select`))
      })

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
              url: args.url || `/lifecycle/${$ndx.find('>.selected').attr('id')}`,
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
            $table.trigger('send', $table.data('record'))
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

  $tablebar
    .on('click', '.edit', e => {
      if (!$(e.currentTarget).hasClass('active')) {
        return
      }
      console.log('calling')

      $table.trigger('edit', {
        data: _ => {
          return JSON.stringify({
            "id": $ndx.find('>.selected').attr('id'),
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
    // .on('click', '.edit', e => {
    //   if (!$(e.currentTarget).hasClass('active')) {
    //     console.log('cancelling')
    //     e.stopPropagation()
    //     return
    //   }
    // })
    .on('click', '.add', e => {
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
    .on('click', '.remove', e => {
      if (!$(e.currentTarget).hasClass('active')) {
        return
      }
      $ndx.trigger('delete', {
        url: `/lifecycle/${$table.attr('id')}`,
        buttonbar: $tablebar
      })
    })

  $tablebar.find('.refresh').on('click', e => {
    $ndx.find('.selected').removeClass('selected').click()
  })

  $lifecycle
    .on('activate', e => {
      $.ajax({
        url: '/strains',
        method: 'GET',
        async: false,
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
        async: false,
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

      $(e.currentTarget)
        .addClass('active')
        .find('.rows.ndx')
        .trigger('refresh', { newRow: newNdxRow, buttonbar: $tablebar })
    })
})