$(function () {
  let $generation = $('body>.main>.workspace>.generation')
  let $gentable = $generation.find('>.table.generation>.rows[name="editor"]')
  let $tablebar = $generation.find('>.table.generation>.buttonbar')
  let $sources = $gentable.find('>.row.sources')
  let $source = $sources.children().first()

  let $notes = $('body>.main>.workspace>.template>.notes')
    .clone(true, true)
    .insertAfter($gentable)

  $tablebar.trigger('subscribe', {
    clazz: 'notes',
    clicker: e => {
      e.stopPropagation()

      if ($generation
        .find('>.table.generation')
        .removeClass('selecting')
        .toggleClass('noting')
        .hasClass('noting')
      ) {
        $notes.trigger('refresh', $gentable.attr('id'))
      }
    },
  })

  $.ajax({
    url: '/substrates',
    method: 'GET',
    async: true,
    success: (result, status, xhr) => {
      var substrates = { bulk: [], grain: [], liquid: [], agar: [] }
      result.forEach(r => {
        substrates[r.type.toLowerCase()]
          .push($('<option>')
            .val(r.id)
            .text(`${r.name} | Vendor: ${r.vendor.name}`)
            .data('substrate', r))
      })
      $gentable
        .find('>.row.plating>select.value')
        .empty()
        .append(substrates.agar)

      $gentable
        .find('.row.liquid>select.value')
        .empty()
        .append(substrates.liquid)
    },
    error: console.log,
  })

  let newNdxRow = (data) => {
    return $('<div>')
      .addClass('row hover')
      .attr('id', data.id)
      .append($('<div class="lineage static" />').text(
        (sources => {
          let result = []
          sources.forEach(s => {
            result.push(s.strain.name)
          })
          return `${result.join(' + ')} - ${new Date(data.ctime).toLocaleString()}`
        })(data.sources)))
  }

  let $ndx = $generation
    .find('>.table.generation>.rows[name="generation"]')
    .on('click', '>.row', e => {
      if (e.isPropagationStopped()) {
        return false
      }
      var $row = $(e.currentTarget)
      $.ajax({
        url: '/generation/' + $row.attr('id'),
        method: 'GET',
        async: true,
        success: (result, status, xhr) => {
          $gentable
            .trigger('send', result)
            .data('record', result)
            .parents('.table.generation')
            .first()
            .find('>.columns>.column.full')
            .text($row.find('.lineage').text())
            .trigger('click')
        },
        error: console.log,
      })
    })

  let $events = $('.main>.workspace>.template>.table.events')
    .clone(true, true)
    .appendTo($generation)
    .find('>div.rows')
    .on('send', (e, ...ev) => {
      // $tablebar.find('.remove')[((ev ||= []).length !== 0 ? "removeClass" : "addClass")]('active')
    })
    .trigger('initialize', {
      parent: 'generation',
      $owner: $gentable,
      // $tablebar: $tablebar,
    })

  $gentable.off('refresh').off('add').off('edit').off('click').off('send')
    .on('send', (e, data = { plating_substrate: {}, liquid_substrate: {}, mtime: 'Now', ctime: 'Now' }) => {
      let $g = $(e.currentTarget)

      $g.attr('id', data.id)
      $g.find('>.row.plating>select').val(data.plating_substrate.id)
      $g.find('>.row.liquid>select').val(data.liquid_substrate.id)
      $g.find('>.row.mtime>.static.date').trigger('set', data.mtime)
      $g.find('>.row.ctime>.static.date').trigger('set', data.ctime)

      $sources.trigger('send', data.sources)
      $events.trigger('send', data.events)
    })
    .on('refresh', (e, data) => {
      let gid = data
      $.ajax({
        url: `/generation/${gid}`,
        method: 'GET',
        async: true,
        success: (result, status, xhr) => { $gentable.trigger('send', result) },
        error: console.log,
      })
    })
    .on('reset', e => {
      $(e.currentTarget)
        .removeClass('editing adding')
        .trigger('send', $(e.currentTarget).data('generation'))
        .find('>.row>select')
        .attr('disabled', true)
    })
    .parents('.table.generation')
    .on('click', '>.columns>.column.full', e => {
      $(e.delegateTarget)
        .removeClass('noting')
        .toggleClass('selecting')
    })

  $sources.find('.add-source').on('click', e => {
    let $curr = $source
      .clone(true, true)
      .prependTo($sources)
      .addClass('removable editing adding')

    $curr.find('>select').attr('disabled', false)
    $curr.find('>select[name="type"]').val('spore').trigger('change')
    $curr.find('>select[name="origin"]').val('strain').trigger('change')
  })

  $sources
    .on('send', (e, ...data) => {
      e.stopPropagation()

      $(e.currentTarget).find('>.source.removable').remove()

      data.forEach(s => {
        $source
          .clone(true, true)
          .prependTo($sources)
          .addClass('removable')
          .data('source', s)
          .trigger('send', s)
      })
    })
    .find('>.source')
    .on('send', (e, data = {}) => {
      e.stopPropagation()

      let $s = $(e.currentTarget)
      let strainsource = typeof data.lifecycle === 'undefined'

      // $s.attr('id', data.id)

      $s.find('>select[name="type"]')
        .val(data.type.toLowerCase())
        .trigger('change')

      $s.find('>select[name="origin"]')
        .val(strainsource ? 'strain' : 'event')
        .trigger('change')

      if (!strainsource) {
        $s.find('>select[name="parent"]')
          .val((data.lifecycle || {}).id)
          .trigger('change')
      }

      setTimeout(complete => {
        let $p = $s.find('>select[name="progenitor"]')
        if (strainsource) {
          $p.val(data.strain.id)
        } else {
          $p.val(data.lifecycle.events[0].id)
        }
        $p.trigger('change')
      }, 5)
    })
    .on('reset', e => {
      $(e.currentTarget)
        .trigger('send', $(e.currentTarget).data('source'))
        .find('select')
        .attr('disabled', true)
    })

    .on('change', '>select[name="progenitor"]', e => {
      let $p = $(e.currentTarget)

      $p
        .parents('.sources')
        .first()
        .find('>select[name="type"]')
        .val($p.text().match(/^Clone/) ? 'clone' : 'spore')
        .trigger('change')
    })
    .on('events', '>select[name="progenitor"]', (e, p) => {
      let $e = $(e.currentTarget).empty()

      $.ajax({
        url: `/lifecycle/${p}`,
        method: 'GET',
        async: true,
        success: (result = [], status, xhr) => {
          result.events.forEach(v => {
            if (!v.event_type.name.match(/^(Spore|Clone)/)) {
              return
            }
            $e.append($(`<option>`)
              .val(v.id)
              .text(`${v.event_type.name} | ${v.mtime.replace('T', ' ').replace(/:\d{1,2}(\..+)?Z.*/, '')}`)
              .data('send', v))
          })
          // TODO: figure out what's selected
        },
        error: console.log,
      })
    })
    .on('strains', '>select[name="progenitor"]', e => {
      let opts = []
      $generation.data('strains').forEach(v => {
        opts.push($(`<option>`)
          .val(v.id)
          .text(`${v.name || v.id} | ${(v.ctime || '').replace('T', ' ').replace(/:\d{1,2}(\..+)?Z.*/, '')}`)
          .data('send', { strain: v }))
      })
      $(e.currentTarget)
        .empty()
        .append(opts)
    })

    .on('change', '>select[name="parent"]', e => {
      $(e.currentTarget)
        .parents('.source')
        .find('>[name="progenitor"]')
        .trigger('events', $(e.currentTarget).val())
    })
    .on('refresh', '>select[name="parent"]', e => {
      e.stopPropagation()

      $(e.currentTarget)
        .empty()
        .append($generation.data('lifecycles'))
        .children()
        .first()
        .attr('selected', true)
    })

    .on('change', '>select[name="origin"]', e => {
      let origin = $(e.currentTarget).val()
      let $s = $(e.delegateTarget)
        .attr("source-origin", origin)

      if (origin === 'strain') {
        $s
          .find('>select[name="progenitor"]')
          .trigger('strains')
          .children()
          .first()
          .attr('selected', true)
        $s.trigger('change')
      } else {
        $s.find('>select[name="parent"]').trigger('refresh')
      }
    })

    .on('change', '>select[name="type"]', e => {
      $(e.currentTarget)
        .parents('.sources')
        .first()
        .attr('source-type', $(e.currentTarget).val())
    })

    .on('click', '>div.action', e => {
      let $src = $(e.delegateTarget)
      if ($src.toggleClass('editing').hasClass('adding')) {
        $src.remove()
        return
      } else if ($src.hasClass('editing')) {
        $src.find('select').attr('disabled', false)
      } else {
        $src.trigger('reset')
      }
    })

    .on('click', '>div.commit', e => {
      let $src = $(e.delegateTarget)
      let type = $src.find('>select[name="type"]').val()
      let method = 'POST'
      let url = `/generation/${$gentable.attr('id')}/sources`
      let other = {}

      $src.find('select').attr('disabled', true)

      if ($src.hasClass('adding')) {
        url = `${url}/${$src.attr('source-origin')}`
        other.data = JSON.stringify({
          ...$src
            .find('>select[name="progenitor"]>:selected')
            .data('send'),
          type: type[0].toUpperCase().concat(type.slice(1)) // mea culpa
        })
      } else if ($src.hasClass('editing')) {
        let temp = {
          id: $src.attr('id'),
          type: type[0].toUpperCase().concat(type.slice(1)),
          strain: {},
        }
        if ($src.attr('source-origin') !== 'strain') {
          method = 'PATCH'
          temp.lifecycle = {
            ...$src
              .find('>select[name="parent"]>:selected')
              .data('record'),
            events: [
              $src
                .find('>select[name="progenitor"]>:selected')
                .data('send'),
            ]
          }
        }
        other.data = JSON.stringify(temp)
      } else {
        url = `${url}/${$src.attr('id')}`
        method = 'DELETE'
        other.success = (result, status, xhr) => {
          $src.remove()
          if ($sources.find('.removable').length === 0) {
            $sources.attr('source-type', 'spore')
          }
        }
      }

      $.ajax({
        ...{
          url: url,
          method: method,
          async: true,
          success: (result, status, xhr) => { $src.removeClass('editing adding') },
          error: (xhr, status, err) => {
            console.log(status, err, xhr)
            $src.trigger('reset')
          },
        },
        ...other,
      })
    })

  $tablebar
    .on('click', '.edit', e => {
      if (!$(e.currentTarget).hasClass('active')) {
        return
      }

      $gentable
        .addClass('editing')
        .find('>.row>select')
        .attr('disabled', false)

      $tablebar.trigger('set', {
        "target": $gentable,
        "handlers": {
          "cancel": e => { $gentable.trigger('reset') },
          "ok": e => {
            let id = $gentable.attr('id')
            $.ajax({
              url: `/generation/${id}`,
              method: 'PATCH',
              async: true,
              data: JSON.stringify({
                id: id,
                plating_substrate: { id: $gentable.find('>.plating>select').val() },
                liquid_substrate: { id: $gentable.find('>.liquid>select').val() },
              }),
              success: (result, status, xhr) => {
                $gentable.removeClass('editing').data('generation', result)
              },
              error: (xhr, status, err) => {
                $gentable.trigger('reset')
                console.log(jqXhr, status, err)
              },
            })
          }
        }
      })
    })
    .on('click', '.add', e => {
      if (!$(e.currentTarget).hasClass('active')) {
        return
      }

      $gentable
        .addClass('adding')
        .trigger('send')
        .find('>div>select[disabled]')
        .attr('disabled', false)

      $tablebar.trigger('set', {
        "target": $gentable,
        "handlers": {
          "cancel": e => { $gentable.trigger('reset') },
          "ok": e => {
            $.ajax({
              url: '/generation',
              method: 'POST',
              async: true,
              data: JSON.stringify({
                plating_substrate: { id: $gentable.find('>.plating>select').val() },
                liquid_substrate: { id: $gentable.find('>.liquid>select').val() },
              }),
              success: (result, status, xhr) => {
                $gentable
                  .trigger('send', result)
                  .find('>div>select')
                  .attr('disabled', true)
                $sources.attr('source-type', 'spore')
              },
              error: console.log,
              complete: (jqXhr, status) => { $gentable.removeClass('editing adding') }
            })
          }
        }
      })
    })
    .on('click', '.remove', e => {
      if (!$(e.currentTarget).hasClass('active')) {
        return
      }

      $.ajax({
        url: `/generation/${$gentable.attr('id')}`,
        method: 'DELETE',
        async: true,
        success: (result, status, xhr) => {
          $ndx.trigger('refresh', { newRow: newNdxRow, buttonbar: $tablebar })
        },
        error: console.log,
      })
    })
    .on('click', '.refresh', e => {
      $ndx.trigger('refresh', { newRow: newNdxRow, buttonbar: $tablebar })
    })

  $generation.on('activate', e => {
    let $g = $(e.currentTarget).addClass('active')

    $.ajax({
      url: '/strains',
      method: 'GET',
      async: true,
      success: (result = [{ id: "No Strains found" }], status, xhr) => {
        $g.data('strains', result)
      },
      error: console.log,
    })

    $.ajax({
      url: `/lifecycles`,
      method: 'GET',
      async: true,
      success: (result = [{ id: 'No Lifecycles found' }], status, xhr) => {
        let lcopts = []
        result.forEach(v => {
          lcopts.push($('<option>')
            .val(v.id)
            .text(`${v.location || v.id} | ${(v.mtime || '').replace('T', ' ').replace(/:\d{1,2}(\..+)?Z.*/, '')}`)
            .data('lifecycle', v))
        })
        $g.data('lifecycles', lcopts)
      },
      error: console.log,
    })

    $ndx.trigger('refresh', { newRow: newNdxRow, buttonbar: $tablebar })
  })
})
