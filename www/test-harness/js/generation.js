$(function () {
  let newNdxRow = (data) => {
    return $('<div>')
      .addClass('row hover')
      .attr('id', data.id)
      .append($('<div>')
        .addClass('lineage static')
        .text(
          (sources => {
            let result = []
            sources.forEach(s => {
              result.push(s.strain.name)
            })
            return `${result.join(' + ')} - ${new Date(data.ctime).toLocaleString()}`
          })(data.sources)))
  }

  let $generation = $('body>.main>.workspace>.generation')
    .on('activate', (e, selected) => {
      if (!selected || selected.constructor.prototype !== String.prototype) {
        selected = $ndx.find('>.selected').attr('id')
      }

      $.ajax({
        url: '/strains',
        method: 'GET',
        async: true,
        success: (result = [{ id: "No Strains found" }], status, xhr) => {
          let $progeny = $generation
            .find('select.progeny.strain-list')
            .empty()
            .append($('<option>')
              .addClass('null-pointer')
              .text('-- None --'))

          let $strains = $generation
            .find('select[name=strains].strain-list')
            .empty()

          result.forEach(s => {
            let $o = $('<option>')
              .data('record', s)
              .val(s.id)

            let gid = {}
            if (s.generation !== undefined) {
              gid.gid = s.generation.id
            }

            $o
              .clone(true, true)
              .attr(gid)
              .appendTo($progeny)
              .trigger('long-format', s)
            $o.clone(true, true).appendTo($strains).trigger('long-format', s)
          })
        },
        error: console.log,
      })

      $.ajax({
        url: `/lifecycles`,
        method: 'GET',
        async: true,
        success: (result = [{ id: 'No Lifecycles found' }], status, xhr) => {
          let $lcs = $source.find('select[name="lcs"]').empty()
          result.forEach(v => {
            $lcs
              .append($('<option>')
                .val(v.id)
                .text(`${v.location || v.id} | ${(v.mtime || '').replace('T', ' ').replace(/:\d{1,2}(\..+)?Z.*/, '')}`)
                .data('record', v))
          })
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
              .push($('<option>')
                .val(r.id)
                .text(`${r.name} | Vendor: ${r.vendor.name}`)
                .data('record', r))
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

      $source.find('>.src-ctl').addClass('active') // wtaf?!

      $ndx.trigger('refresh', { newRow: newNdxRow, buttonbar: $tablebar })

      setTimeout(_ => {
        $ndx.find(`>.row[id="${selected}"]:not(.selected)`).click()
      }, 50)
    })
    .on('select', (e, gid) => { $ndx.find(`.row[id=${gid}]`).click() })


  let $gentable = $generation.find('>.table.generation>.rows[name="editor"]')
    .off('refresh').off('add').off('edit').off('click').off('send')
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
    .on('refresh', (e, gid) => {
      $.ajax({
        url: `/generation/${gid}`,
        method: 'GET',
        async: true,
        success: (result, status, xhr) => { $gentable.trigger('send', result) },
        error: console.log,
      })
    })
    .on('resend', e => {
      $(e.currentTarget)
        .removeClass('editing adding')
        .trigger('send', $(e.currentTarget).data('record'))
        .find('>.row>select:not(.progeny)')
        .attr('disabled', true)
    })
    .on('change', '>.row.strain>.progeny', e => {
      let $p = $(e.currentTarget)
      let curr = $p.attr('curr-id')
      let val = $p.val()
      let url = `/strain/${val}/generation`

      if (curr) {
        $.ajax({
          url: `/strain/${curr}/generation`,
          method: 'DELETE',
          async: false,
          success: _ => {
            $p
              .removeAttr('curr-id')
              .find(`option[value=${curr}][disabled]`)
              .attr('disabled', false)
              .removeAttr('gid')
          },
          error: console.log,
        })
      }

      let $row = $p.parents('.row.strain').first()
      if (val === $p.children().first().text()) {
        $row.removeClass('link')
        return
      }

      $row.addClass('link')
      url = `${url}/${$row.parents('.rows').first().attr('id')}`

      $.ajax({
        url: url,
        method: 'PATCH',
        async: true,
        success: _ => {
          $p
            .attr('curr-id', val)
            .find(`option[value=${val}]`)
            .attr({
              disabled: true,
              gid: $row.parents('.rows').first().attr('id'),
            })
        },
        error: console.log,
      })
    })
    .on('activate', '>.row.strain>.progeny', (e, gid) => {
      e.stopPropagation()

      $.ajax({
        url: `/strain/${gid}/generation`,
        method: 'GET',
        async: true,
        success: (strain = { id: '' }, status, xhr) => {
          let $link = $(e.currentTarget)
            .removeAttr('curr-id')
            .val(strain.id || $(e.currentTarget).children().first().text())
            .parents('.row.strain')
            .removeClass('link')

          if (strain.id) {
            $(e.currentTarget).attr('curr-id', strain.id)
            $link.addClass('link')
          }
        },
        error: (xhr, status, err) => { console.log('error', xhr, status, err) },
      })
    })
    .on('click', '>.row.strain.link>.label', e => {
      e.stopPropagation()

      $('body>.main>.header>.menuitem[name=strain]')
        .trigger('click', $(e.delegateTarget)
          .find('>.row.strain.link>select.progeny')
          .attr('curr-id'))
    })

  let $tablebar = $generation.find('>.table.generation>.buttonbar')
    .on('click', '.edit.active', e => {
      $gentable
        .addClass('editing')
        .find('>.row>select:not(.progeny)')
        .attr('disabled', false)

      $tablebar.trigger('set', {
        "target": $gentable,
        "handlers": {
          "cancel": e => { $gentable.trigger('resend') },
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
                ctime: $gentable.find('>.mtime>.static.date').data('value'),
              }),
              success: (result, status, xhr) => {
                $gentable
                  .removeClass('editing')
                  .trigger('send', {
                    ...$gentable.data('record'),
                    ...result,
                  })
              },
              error: (xhr, status, err) => {
                console.log(xhr, status, err)
                $gentable.trigger('resend')
              },
            })
          }
        }
      })
    })
    .on('click', '.add.active', e => {
      $gentable
        .addClass('adding')
        .trigger('send')
        .find('>div>select[disabled]')
        .attr('disabled', false)

      $tablebar.trigger('set', {
        "target": $gentable,
        "handlers": {
          "cancel": e => { $gentable.trigger('resend') },
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
              error: (xhr, status, err) => {
                console.log(xhr, status, err)
                $gentable.trigger('resend')
              },
              complete: (xhr, status) => {
                $gentable.removeClass('adding')
              }
            })
          }
        }
      })
    })
    .on('click', '.remove.active', e => {
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
      $gentable.removeClass('editing')

      $ndx
        .find('>.row.selected')
        .removeClass('selected')
        .click()
    })
    .trigger('subscribe', {
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
    .trigger('subscribe', {
      clazz: 'report',
      clicker: e => {
        e.stopPropagation()
        $('body>.main>.header>.menu-scroll')
          .trigger('select', ['reporting', 'history.gen', $ndx.find('>.selected').attr('id')])
      }
    })

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
            .addClass('selecting') // gets removed on column.click
            .find('>.columns>.column.full')
            .text($row.find('.lineage').text())
            .trigger('click')

          $gentable
            .find('>.row.strain>.progeny')
            .trigger('activate', $row.attr('id'))
        },
        error: console.log,
      })
    })

  let $events = $('body>.template>.table.events')
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

  let $notes = $('body>.template>.notes')
    .clone(true, true)
    .insertAfter($gentable)

  let $sources = $gentable.find('>.row.sources')
  let $source = $sources.children().first()

  $gentable.parents('.table.generation')
    .on('click', '>.columns>.column.full', e => {
      $(e.delegateTarget)
        .removeClass('noting')
        .toggleClass('selecting')
    })

  $sources
    .on('send', (e, ...data) => {
      e.stopPropagation()

      $(e.currentTarget).find('>.source.removable').remove()

      data.forEach(s => {
        let $s = $source
          .clone(true, true)
          .insertBefore($(e.delegateTarget).find('.add-source'))
          .toggleClass('removable template')
          .data('record', s)
          .trigger('send', s)
      })
    })
    .on('click', '.add-source', e => {
      let $curr = $source
        .clone(true, true)
        .insertBefore($(e.delegateTarget).find('.add-source'))
        .toggleClass('removable editing adding template')

      $curr.find('>select').attr('disabled', false)
      $curr.find('>select[name="type"]').val('spore').trigger('change')
      $curr.find('>select[name="origin"]').val('strain').trigger('change')
    })

    .find('>.source')
    .on('send', (e, data = {}) => {
      e.stopPropagation()

      let $s = $(e.currentTarget).attr('id', data.id)
      let strainsource = typeof data.lifecycle === 'undefined'

      $s.find('>select[name="origin"]')
        .val(strainsource ? 'strain' : 'event')
        .trigger('change')

      $s.find('>select[name="strains"]')
        .val(data.strain.id)

      $s.find('>select[name="type"]')
        .val(data.type.toLowerCase())
        .trigger('change')

      if (!strainsource) {
        $s.find('>select[name="lcs"]')
          .val(data.lifecycle.id)
          .trigger('change')
      }
    })
    .on('resend', e => {
      e.stopPropagation()

      $(e.currentTarget)
        .trigger('send', $(e.currentTarget).data('record'))
        .find('select')
        .attr('disabled', true)
    })

    .on('change', '>select[name="origin"]', e => {
      $(e.delegateTarget).attr("source-origin", $(e.currentTarget).val())
    })
    .on('change', '>select[name="lcs"]', e => {
      $(e.currentTarget)
        .parents('.source')
        .first()
        .find('>.progenitor>select[name="progenitor"]')
        .trigger('events', $(e.currentTarget).val())
    })
    .on('events', '>.progenitor>select[name="progenitor"]', (e, p) => {
      let $e = $(e.currentTarget).empty()

      $.ajax({
        url: `/lifecycle/${p}`,
        method: 'GET',
        async: true,
        success: (result = [], status, xhr) => {
          result.events.forEach(v => {
            $e.append($(`<option>`)
              .attr(
                'disabled',
                !v.event_type.name
                  .toLowerCase()
                  .startsWith($(e.delegateTarget)
                    .parents('div[source-type]')
                    .first()
                    .attr('source-type')))
              .val(v.id)
              .text(`${v.event_type.name} | ${v.mtime.replace('T', ' ').replace(/:\d{1,2}(\..+)?Z.*/, '')}`)
              .data('record', v))
          })
          $e.val($(e.delegateTarget).data('record').lifecycle.events[0].id)
        },
        error: console.log,
      })
    })
    .on('change', '>select[name="type"]', e => {
      $(e.delegateTarget)
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
        $src.trigger('resend')
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
            .find('>.progenitor>select[name="progenitor"]>option:selected')
            .data('record'),
          strain: $src
            .find('>select[name="strains"]>option:selected')
            .data('record'),
          type: type[0].toUpperCase().concat(type.slice(1)) // mea culpa
        })
      } else if ($src.hasClass('editing')) {
        method = 'PATCH'
        let temp = {
          id: $src.attr('id'),
          type: type[0].toUpperCase().concat(type.slice(1)),
          strain: $src
            .find('>select[name="strains"]>option:selected')
            .data('record'),
        }
        if ($src.attr('source-origin') !== 'strain') {
          temp.lifecycle = {
            ...$src
              .find('>select[name="lcs"]>option:selected')
              .data('record'),
            events: [
              $src
                .find('>.progenitor>select[name="progenitor"]>option:selected')
                .data('record'),
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
          success: console.log,
          error: (xhr, status, err) => {
            console.log(status, err, xhr)
            $src.trigger('resend')
          },
          complete: _ => { $src.removeClass('adding editing') }
        },
        ...other,
      })
    })
})
