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
      .append($('<div class="create_date static date" />').text(data.create_date.replace('T', ' ').replace(/(\.\d+)?Z/, '')))
      .append($('<div class="create_date live" disabled />').text(data.create_date))
      .append($('<div class="vendor static" />').text(data.vendor.name))
      .append($('<select class="vendor live" />')
        .append(vendors)
        .data('vendor_uuid', data.vendor.id)
        .val(data.vendor.id))
  }

  $table
    .on('reinit', e => {
      vendors = []
      $.ajax({
        url: '/vendors',
        method: 'GET',
        async: false,
        success: (result, status, xhr) => {
          result.forEach(r => {
            vendors.push($(`<option value="${r.id}">${r.name}</option>`))
          })
        },
        error: console.log,
      })

      $(e.currentTarget).trigger('refresh', { newRow: newRow })
    })
    .on('click', '>.row', e => {
      $attributes.trigger('refresh', $(e.currentTarget))
    })

  function newAttributeRow(data) {
    data ||= {}
    return $('<div class="row hover" />')
      .append($('<div class=uuid />').text(data.id))
      .append($('<div class="name static" />').text(data.name))
      .append($('<input class="name live" />').val(data.name))
      .append($('<div class="value static" />').html(data.value))
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
        $selected.find('>.create_date.static').text(data.create_date)
        $selected
          .find('>.vendor.static')
          .text($selected.find('>.vendor.live>option:selected').text())
      },
      error: (xhr, status, error) => { $table.trigger('remove-selected') },
      buttonbar: $buttonbar
    })
  })

  $buttonbar.find('>.remove').on('click', e => {
    if ($(e.currentTarget).hasClass('active')) {
      $table.trigger('delete')
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
      ok: e => {
        var $row = $attributes.find('.selected')
        $.ajax({
          url: `/strain/${$table.find('.selected>.uuid').text()}/attribute`,
          contentType: 'application/json',
          method: 'PATCH',
          dataType: 'json',
          data: JSON.stringify({
            name: $row.find('>.name.live').val(),
            value: $row.find('>.value.live').val()
          }),
          async: true,
          success: (data, status, xhr) => {
            var $row = $attributes.find('.selected')
            $row.find('>.name.static').text($row.find('>.name.live').val())
            $row.find('>.value.static').text($row.find('>.value.live').val())
          },
          error: console.log,
        })
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
      ok: e => {
        var $row = $attributes.find('.selected')
        $.ajax({
          url: `/strain/${$table.find('.selected>.uuid').text()}/attribute`,
          contentType: 'application/json',
          method: 'POST',
          dataType: 'json',
          data: JSON.stringify({
            name: $row.find('>.name.live').val(),
            value: $row.find('>.value.live').val()
          }),
          async: true,
          success: (data, status, xhr) => {
            var $row = $attributes.find('.selected')
            $row.find('>.uuid').text(data.id)
            $row.find('>.name.static').text($row.find('>.name.live').val())
            $row.find('>.value.static').text($row.find('>.value.live').val())
          },
          error: (xhr, status, error) => { $attributes.trigger('remove-selected') },
        })
      },
      buttonbar: $attributebar
    })
  })

  $attributebar.find('>.remove').on('click', e => {
    if ($(e.currentTarget).hasClass('active')) {
      $attributes.trigger('delete',
        `/strain/${$table.find('.selected>.uuid').text()}/attribute/${$attributes.find('.selected>.uuid').text()}`
      )
    }
  })

  $attributebar.find('>.refresh').on('click', e => {
    $attributes.trigger('reinit')
  })

  $strain.on('activate', e => {
    $strain.addClass('active')
    $table.trigger('reinit')
  })
})