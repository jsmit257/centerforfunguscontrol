$(function () {
  let $strain = $('body>.main>.workspace>.strain')
    .on('activate', (e, selected) => {
      if (!selected || selected.constructor.prototype !== String.prototype) {
        selected = $table.find('>.row.selected').attr('id')
      }

      $table.trigger('reinit')

      setTimeout(_ => {
        $table.find(`>.row[id=${selected}]`).click()
      }, 20)
    })

  let $table = $strain

    .find('>.table.strain')
    .on('defer', (e, id) => {
      setTimeout(_ => {
        $(e.currentTarget)
          .find(`.rows>.row[id=${id}]:not(.selected)`)
          .click()
      }, 20)
    })

    .find('>.rows')
    .on('reinit', e => {
      vendors = []
      $.ajax({
        url: '/vendors',
        method: 'GET',
        async: true,
        success: (result, status, xhr) => {
          result.forEach(r => {
            vendors.push($('<option>')
              .val(r.id)
              .text(r.name))
          })
        },
        error: console.log,
      })

      $(e.currentTarget).trigger('refresh', {
        newRow: newRow,
        buttonbar: $buttonbar
      })
    })
    .on('click', '>.row', e => {
      if (e.isPropagationStopped()) {
        return
      }

      $.ajax({
        url: '/strainattributenames',
        method: 'GET',
        async: true,
        success: (result, status, xhr) => {
          let $attrlist = $strain.find('datalist[id=known-strain-names]').empty()
          result.forEach(r => { $attrlist.append($(`<option />`).val(r)) })
        },
        error: console.log,
      })

      $attributes.trigger('refresh', $(e.currentTarget))

      $buttonbar.find('>.optional>.genlink')[$(e.currentTarget).attr('gid') ? "addClass" : "removeClass"]('active')
    })

  let $buttonbar = $strain.find('>.table.strain>.buttonbar')
    .on('click', '>.edit.active', e => {
      let $row = $table.find('.selected')
      let $vendor = $row.find('>select.vendor').append(vendors)

      $vendor.val($vendor.data('vendor_uuid'))

      $table.trigger('edit', {
        url: `/strain/${$row.attr('id')}`,
        data: $selected => {
          return JSON.stringify({
            "name": $selected.find('>.name.live').val(),
            "species": $selected.find('>.species.live').val().trim(),
            "vendor": {
              "id": $selected.find('>.vendor.live').val()
            }
          })
        },
        success: (data, status, xhr) => {
          let $selected = $table.find('.selected')
          $selected.find('>.name.static').text($selected.find('>.name.live').val())
          $selected.find('>.species.static').html($selected.find('>.species.live').val() || "&nbsp")
          $selected
            .find('>.vendor.static')
            .text($selected.find('>.vendor.live>option:selected').text())
        },
        buttonbar: $(e.delegateTarget)
      })
    })
    .on('click', '>.add.active', e => {
      $attributes
        .data('attributes', $attributes.find('.row'))
        .empty()

      $table.trigger('add', {
        newRow: newRow,
        data: $selected => {
          return JSON.stringify({
            "name": $selected.find('>.name.live').val(),
            "species": $selected.find('>.species.live').val().trim(),
            "vendor": {
              "id": $selected.find('>.vendor.live').val()
            }
          })
        },
        success: (data, status, xhr) => {
          let $selected = $table.find('.selected')
          $selected.attr('id', data.id)
          $selected.find('>.name.static').text($selected.find('>.name.live').val())
          $selected.find('>.species.static').html($selected.find('>.species.live').val() || "&nbsp")
          $selected.find('>.ctime.static').text(data.ctime.replace('T', ' ').replace(/(\..+)?Z.*/, ''))
          $selected
            .find('>.vendor.static')
            .text($selected.find('>.vendor.live>option:selected').text())
        },
        error: (xhr, status, error) => {
          $table.trigger('remove-selected')
          $attributes.append($attributes.data('attributes'))
        },
        buttonbar: $(e.delegateTarget)
      })
    })
    .on('click', '.remove.active', e => {
      $table.trigger('delete', {
        url: `/strain/${$table.find('>.selected').attr('id')}`,
        buttonbar: $(e.delegateTarget),
      })
    })
    .on('click', '>.refresh.active', e => {
      $table.trigger('reinit')
    })
    .trigger('subscribe', {
      clazz: 'photos',
      clicker: e => {
        e.stopPropagation()
        if ($strain
          .find('>.table.strain')
          .toggleClass('photoing')
          .hasClass('photoing')
        ) {
          $genphotos
            .trigger('refresh', $table.find('>.selected').attr('id'))
            .addClass('gallery')
            .removeClass('singleton')
        }
      },
    })
    .trigger('subscribe', {
      clazz: 'genlink',
      clicker: e => {
        e.stopPropagation()

        $('body>.main>.header>.menuitem[name=generation]')
          .trigger('click', $table.find('>.row.selected').attr('gid'))
      },
    })
    .trigger('subscribe', {
      clazz: 'report',
      clicker: e => {
        e.stopPropagation()
        $('body>.main>.header>.menu-scroll')
          .trigger('select', ['reporting', 'history.str', $table.find('>.selected').attr('id')])
      }
    })


  let $attributes = $strain.find('>.table.strainattribute>.rows')
    .off('refresh').on('refresh', (e, row) => {
      $.ajax({
        url: `/strain/${$(row).attr('id')}`,
        method: 'GET',
        async: true,
        success: (result, sc, xhr) => {
          $attributes.empty()
          result.attributes ||= []
          result.attributes.forEach(a => { $attributes.append(newAttributeRow(a)) })
          $attributes.find('.row').first().click()
          $buttonbar.find('.remove')[$attributes.children().length > 0 ? "removeClass" : "addClass"]("active")
          $attributebar.find('.remove, .edit')[$attributes.children().length === 0 ? "removeClass" : "addClass"]("active")
        },
        error: console.log
      })
    })

  let $attributebar = $strain.find('>.table.strainattribute>.buttonbar')
    .on('click', '>.edit.active', e => {
      $attributes.trigger('edit', {
        url: `/strain/${$table.find('.selected').attr('id')}/attribute`,
        data: $selected => {
          return JSON.stringify({
            id: $selected.attr('id'),
            name: $selected.find('>.name.live').val(),
            value: $selected.find('>.value.live').val()
          })
        },
        success: (data, status, xhr) => {
          let $row = $attributes.find('.selected')
          $row.find('>.name.static').text($row.find('>.name.live').val())
          $row.find('>.value.static').text($row.find('>.value.live').val())
        },
        buttonbar: $(e.delegateTarget)
      })
    })
    .on('click', '>.add.active', e => {
      $attributes.trigger('add', {
        newRow: newAttributeRow,
        url: `/strain/${$table.find('.selected').attr('id')}/attribute`,
        data: $selected => {
          return JSON.stringify({
            name: $selected.find('>.name.live').val(),
            value: $selected.find('>.value.live').val()
          })
        },
        success: (data, status, xhr) => {
          let $row = $attributes.find('.selected')
          $row.attr('id', data.id)
          $row.find('>.name.static').text(data.name)
          $row.find('>.value.static').text(data.value)
          $buttonbar.find('.remove')[$attributes.children().length > 0 ? "removeClass" : "addClass"]("active")
        },
        error: (xhr, status, error) => { $attributes.trigger('remove-selected') },
        buttonbar: $(e.delegateTarget)
      })
    })
    .on('click', '.remove.active', e => {
      $attributes.trigger('delete', {
        url: `/strain/${$table.find('.selected').attr('id')}/attribute/${$attributes.find('.selected').attr('id')}`,
        buttonbar: $(e.delegateTarget)
      })
      if ($attributes.children().length === 0) {
        $buttonbar.find('.remove').addClass('active')
      }
    })
    .on('click', '>.refresh', e => {
      $table.trigger('reinit')
    })

  let vendors = []

  let $genphotos = $('body>.template>.photos')
    .clone(true, true)
    .insertAfter($('body>.main>.workspace>.strain>.table.strain>.rows'))

  function newRow(data = { vendor: {} }) {
    let attrs = { id: data.id }
    if (data.generation !== undefined) {
      attrs.gid = data.generation.id
    }
    return $('<div>')
      .addClass('row hover')
      .attr(attrs)
      .append($('<div class="name static" />').text(data.name))
      .append($('<input class="name live" />').val(data.name))
      .append($('<div class="species static" />').html(data.species || "&nbsp"))
      .append($('<input class="species live" />').val(data.species))
      .append($('<div class="vendor static" />').text(data.vendor.name))
      .append($('<select class="vendor live" />')
        .append(vendors)
        .data('vendor_uuid', data.vendor.id)
        .val(data.vendor.id))
      .append($('<div class="ctime static const date" />')
        .data('value', data.ctime)
        .text((data.ctime || "Now").replace('T', ' ').replace(/:\d{2}(\.\d+)?Z/, '')))
  }

  function newAttributeRow(data) {
    data ||= {}
    return $('<div>')
      .addClass('row hover')
      .attr('id', data.id)
      .append($('<div class="name static" />').text(data.name))
      .append($('<input list="known-strain-names" class="name live" />').val(data.name))
      .append($('<div class="value static" />').text(data.value))
      .append($('<input class="value live" />').val(data.value))
  }
})
