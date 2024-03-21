$(function () {
  var $strain = $('body>.main>.workspace>.strain')
  var $table = $strain.find('>.table.strain>.rows')
  var $buttonbar = $strain.find('>.table.strain>.buttonbar')
  var $attributes = $strain.find('>.table.strainattribute>.rows')
  var $attributebar = $strain.find('>.table.strainattribute>.buttonbar')

  var vendors = []

  function newRow(data) {
    data ||= { vendor: {} }
    return $('<div class="row hover" />')
      .append($('<div class=uuid />').text(data.id))
      .append($('<div class="name static" />').text(data.name))
      .append($('<input class="name live" />').val(data.name))
      .append($('<div class="species static" />').html(data.species || "&nbsp"))
      .append($('<input class="species live" />').val(data.species))
      .append($('<div class="vendor static" />').text(data.vendor.name))
      .append($('<select class="vendor live" />')
        .append(vendors)
        .data('vendor_uuid', data.vendor.id)
        .val(data.vendor.id))
      .append($('<div class="create_date static const date" />')
        .data('value', data.create_date)
        .text((data.create_date || "Now").replace('T', ' ').replace(/:\d{2}(\.\d+)?Z/, '')))
  }

  $table
    .on('reinit', e => {
      vendors = []
      $.ajax({
        url: '/vendors',
        method: 'GET',
        async: true,
        success: (result, status, xhr) => {
          result.forEach(r => {
            vendors.push($(`<option value="${r.id}">${r.name}</option>`))
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
          var $attrlist = $strain.find('datalist[id=known-strain-names]').empty()
          result.forEach(r => { $attrlist.append($(`<option />`).val(r)) })
        },
        error: console.log,
      })

      $attributes.trigger('refresh', $(e.currentTarget))
    })

  function newAttributeRow(data) {
    data ||= {}
    return $('<div class="row hover" />')
      .append($('<div class=uuid />').text(data.id))
      .append($('<div class="name static" />').text(data.name))
      .append($('<input list="known-strain-names" class="name live" />').val(data.name))
      .append($('<div class="value static" />').text(data.value))
      .append($('<input class="value live" />').val(data.value))
  }

  $attributes
    .off('refresh')
    .on('refresh', (e, row) => {
      $.ajax({
        url: `/strain/${$(row).find('>.uuid').text()}`,
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

  $buttonbar.find('>.edit').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }

    $table.find('.selected>select.vendor')
      .append(vendors)
      .val($table.find('.selected>select.vendor').data('vendor_uuid'))

    $table.trigger('edit', {
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
        var $selected = $table.find('.selected')
        $selected.find('>.name.static').text($selected.find('>.name.live').val())
        $selected.find('>.species.static').html($selected.find('>.species.live').val() || "&nbsp")
        $selected
          .find('>.vendor.static')
          .text($selected.find('>.vendor.live>option:selected').text())
      },
      buttonbar: $buttonbar
    })
  })

  $buttonbar.find('>.add').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }

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
        var $selected = $table.find('.selected')
        $selected.find('>.uuid').text(data.id)
        $selected.find('>.name.static').text($selected.find('>.name.live').val())
        $selected.find('>.species.static').html($selected.find('>.species.live').val() || "&nbsp")
        $selected.find('>.create_date.static').text(data.create_date.replace('T', ' ').replace(/(\..+)?Z.*/, ''))
        $selected
          .find('>.vendor.static')
          .text($selected.find('>.vendor.live>option:selected').text())
      },
      error: (xhr, status, error) => {
        $table.trigger('remove-selected')
        $attributes.append($attributes.data('attributes'))
      },
      buttonbar: $buttonbar
    })
  })

  $buttonbar.find('>.remove').on('click', e => {
    if ($(e.currentTarget).hasClass('active')) {
      $table.trigger('delete', { buttonbar: $buttonbar })
    }
  })

  $buttonbar.find('>.refresh').on('click', e => {
    $table.trigger('reinit')
  })

  $attributebar.find('>.edit').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }
    $attributes.trigger('edit', {
      url: `/strain/${$table.find('.selected>.uuid').text()}/attribute`,
      data: $selected => {
        return JSON.stringify({
          name: $selected.find('>.name.live').val(),
          value: $selected.find('>.value.live').val()
        })
      },
      success: (data, status, xhr) => {
        var $row = $attributes.find('.selected')
        $row.find('>.name.static').text($row.find('>.name.live').val())
        $row.find('>.value.static').text($row.find('>.value.live').val())
      },
      buttonbar: $attributebar
    })
  })

  $attributebar.find('>.add').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }
    $attributes.trigger('add', {
      newRow: newAttributeRow,
      url: `/strain/${$table.find('.selected>.uuid').text()}/attribute`,
      data: $selected => {
        return JSON.stringify({
          name: $selected.find('>.name.live').val(),
          value: $selected.find('>.value.live').val()
        })
      },
      success: (data, status, xhr) => {
        var $row = $attributes.find('.selected')
        $row.find('>.uuid').text(data.id)
        $row.find('>.name.static').text(data.name)
        $row.find('>.value.static').text(data.value)
        $buttonbar.find('.remove')[$attributes.children().length > 0 ? "removeClass" : "addClass"]("active")
      },
      error: (xhr, status, error) => { $attributes.trigger('remove-selected') },
      buttonbar: $attributebar
    })
  })

  $attributebar.find('>.remove').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }
    $attributes.trigger('delete', {
      url: `/strain/${$table.find('.selected>.uuid').text()}/attribute/${$attributes.find('.selected>.uuid').text()}`,
      buttonbar: $attributebar
    })
    if ($attributes.children().length === 0) {
      $buttonbar.find('.remove').addClass('active')
    }
  })

  $attributebar.find('>.refresh').on('click', e => {
    $table.trigger('reinit')
  })

  $strain.on('activate', e => {
    $strain.addClass('active')
    $table.trigger('reinit')
  })
})